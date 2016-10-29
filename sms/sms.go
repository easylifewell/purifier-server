package sms

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/easylifewell/purifier-server/store"
)

const (
	APPID      = "61fd63a1d7a94f8c81c5ee5a01e96a01"
	SID        = "588379b7b6c442b674187e6de5ae8b9c"
	TOKEN      = "6c8cc0cd398a133c0000b6d64578ff22"
	TEMPLATEID = "30048"
)

type Sms struct {
	TemplateSMS struct {
		AppID      string `json:"appId"`
		Param      string `json:"param"`
		TemplateID string `json:"templateId"`
		To         string `json:"to"`
	} `json:"templateSMS"`
}

type SmsResp struct {
	Resp struct {
		RespCode    string `json:"respCode"`
		Failure     int    `json:"failure"`
		TemplateSMS struct {
			CreateDate int    `json:"createData"`
			SmsID      string `json:"smdId"`
		} `json:"templateSMS"`
	} `json:"resp"`
}

func SendSMS(ctx context.Context, phone, code string, isGuest bool) (string, error) {
	if isGuest {
		return sendSMSForGuest(ctx, phone, code)
	} else {
		return sendSMSForUser(ctx, phone, code)
	}
}

// sendSMSForGuest 首次登录发送验证码的函数
func sendSMSForGuest(ctx context.Context, phone, code string) (string, error) {
	guest := store.GetGuestByPhone(phone)

	now := time.Now()
	// 第一次发送
	if guest.Phone == "" {
		guest.SMSCode = code
		guest.SMSSendDate = now
		guest.Phone = phone
		guest.SMSChangeDate = now.Add(time.Duration(time.Minute * 5))
		store.AddGuest(*guest)
		return sendSMS(ctx, phone, code)
	}

	// 短信过期，再次发送
	if now.Sub(guest.SMSChangeDate) >= 0 {
		guest.SMSCode = code
		guest.SMSSendDate = now
		guest.Phone = phone
		guest.SMSChangeDate = now.Add(time.Duration(time.Minute * 5))
		store.UpdateGuest(*guest)
		return sendSMS(ctx, phone, code)
	}

	// 60s 内不允许再次发送短信
	sendover := now.Sub(guest.SMSSendDate)
	if sendover < time.Duration(time.Minute) {
		return "", errors.New(fmt.Sprintf("请%f s后再次发送", sendover.Seconds()))
	} else {
		// 60s后再次发送
		guest.SMSSendDate = now
		guest.SMSCode = code
		guest.Phone = phone
		guest.SMSChangeDate = now.Add(time.Duration(time.Minute * 5))
		store.UpdateGuest(*guest)
		return sendSMS(ctx, phone, code)
	}

}

// sendSMSForUser 非首次登录发送验证码的函数
func sendSMSForUser(ctx context.Context, phone, code string) (string, error) {
	user := store.GetUserByPhone(phone)

	now := time.Now()
	// 短信过期，再次发送
	if now.Sub(user.SMSChangeDate) >= 0 {
		fmt.Println("短信过期，再次发送")
		user.SMSCode = code
		user.SMSSendDate = now
		user.Phone = phone
		user.SMSChangeDate = now.Add(time.Duration(time.Minute * 5))
		if _, err := store.UpdateUser(*user); err != nil {
			return "", err
		}
		return sendSMS(ctx, phone, code)
	}

	// 60s 内不允许再次发送短信
	sendover := now.Sub(user.SMSSendDate)
	fmt.Println("now:", now)
	fmt.Println("sendate:", user.SMSSendDate)
	if sendover < time.Duration(time.Minute) {
		return "", errors.New(fmt.Sprintf("请%f s后再次发送", sendover.Seconds()))
	} else {
		fmt.Printf("60s后，再次发送 sendover = %v\n", sendover)
		// 60s后再次发送
		user.SMSSendDate = now
		user.SMSCode = code
		user.Phone = phone
		user.SMSChangeDate = now.Add(time.Duration(time.Minute * 5))
		if _, err := store.UpdateUser(*user); err != nil {
			return "", err
		}
		return sendSMS(ctx, phone, code)
	}
}

// SendSMS 发送短信，返回服务器的响应码
func sendSMS(ctx context.Context, phone, code string) (string, error) {

	// just for test
	if phone == "18901030365" || phone == "13641309915" {
		return "000000", nil
	}

	// http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1%E9%AA%8C%E8%AF%81:rest_yz
	t := time.Now()
	now := fmt.Sprintf("%d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	auth := base64.StdEncoding.EncodeToString([]byte(SID + ":" + now))
	sig := fmt.Sprintf("%X", md5.Sum([]byte(SID+TOKEN+now)))

	var sms Sms
	sms.TemplateSMS.AppID = APPID
	sms.TemplateSMS.Param = code
	sms.TemplateSMS.TemplateID = TEMPLATEID
	sms.TemplateSMS.To = phone

	var buf io.ReadWriter
	buf = new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(sms)
	if err != nil {
		return "", err
	}
	URL := fmt.Sprintf("https://api.ucpaas.com/2014-06-30/Accounts/%s/Messages/templateSMS?sig=%s", SID, sig)
	req, err := http.NewRequest("POST", URL, buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", auth)
	var respCode string
	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var data SmsResp
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		respCode = data.Resp.RespCode
		return nil
	})

	return respCode, err
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return
		return ctx.Err()
	case err := <-c:
		return err
	}

}
