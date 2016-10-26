package sms

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

func SendSMS(ctx context.Context, phone, code string) error {
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
		return err
	}
	URL := fmt.Sprintf("https://api.ucpaas.com/2014-06-30/Accounts/%s/Messages/templateSMS?sig=%s", SID, sig)
	req, err := http.NewRequest("POST", URL, buf)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", auth)
	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	})

	return err
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
