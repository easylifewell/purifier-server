package controller

import (
	"github.com/easylifewell/purifier-server/sms"
	"github.com/gorilla/mux"
	"net/http"
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

	timeout, err := time.ParseDuration(req.FormValue("timeout"))
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
	code := vars["code"]
	err := sms.SendSMS(ctx, phone, code)
	if err != nil {
		log.Print(err)
		return
	}
}
