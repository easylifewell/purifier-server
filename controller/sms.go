package controller

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/easylifewell/purifier-server/sms"
	"github.com/gorilla/mux"
)

var (
	Phone = regexp.MustCompile("^1[3|4|5|7|8][0-9]{9}$")
)

type SMSController struct {
}

func NewSMSController() *SMSController {
	return &SMSController{}
}

func (dc SMSController) SendSMS(w http.ResponseWriter, r *http.Request) {
	// ctx is the Context for this handler. Calling cancel closes the ctx.Done
	// channel, which is the cancellation signal for requests started by this handler.
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	timeout, err := time.ParseDuration(r.FormValue("timeout"))
	if err == nil {
		// The request has a timeout, so create a context that is
		// canceled automatically when the timeout expires
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel() // Cancel ctx as soon as handleSearch returns
	vars := mux.Vars(r)
	phone := vars["phone"]
	if !Phone.MatchString(phone) {
		Response400(w, "无效的手机号码")
		return
	}

	code := fmt.Sprintf("%d", rand.Int31())[2:6]
	respCode, err := sms.SendSMS(ctx, phone, code)
	if err != nil {
		Response500(w, err.Error())
		return
	}

	if respCode != "000000" {
		Response500(w, "短信发送失败")
		return
	}
	Response200(w, "短信发送成功")
}
