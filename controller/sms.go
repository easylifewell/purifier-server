package controller

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/easylifewell/purifier-server/model"
	"github.com/easylifewell/purifier-server/sms"
	"github.com/easylifewell/purifier-server/store"
	"github.com/gorilla/mux"
)

var (
	Phone = regexp.MustCompile("^1[3|4|5|7|8][0-9]{9}$")
)

const (
	KEYGEN = "SUSHENGYUAN"
)

type SMSController struct {
}

func NewSMSController() *SMSController {
	return &SMSController{}
}

func (dc SMSController) Login(w http.ResponseWriter, r *http.Request) {
	if phone, ok := isLogin(r); ok {
		logrus.WithFields(logrus.Fields{
			"user.phone": phone,
		}).Info("利用Cookie登录成功")
		Response200(w, "登录成功")
		return
	}

	vars := mux.Vars(r)
	phone := vars["phone"]
	password := vars["password"]

	if phone == "" || password == "" {
		Response400(w, "请求参数不全")
		return
	}
	if !Phone.MatchString(phone) {
		Response400(w, "无效的手机号码")
		return
	}

	if len(password) < 6 || len(password) > 20 {
		Response400(w, "无效的密码")
		return
	}

	user := store.GetUserByPhone(phone)
	var err error
	if user.Phone == "" {
		fmt.Printf("用户 %s 尚未注册，并登录\n", user.Phone)
		Response400(w, "您尚未注册，请先进行注册")
		return

	} else {
		fmt.Printf("用户(%s)登录\n", phone)
		user, err = store.Login(fmt.Sprintf("%d", user.SID), password)
		if err != nil {
			Response500(w, err.Error())
			return
		}
	}
	// 登陆成功，返回cookie
	createUIDC(w, r, store.GetUserByPhone(phone))
	return
}

func (dc SMSController) Register(w http.ResponseWriter, r *http.Request) {
	if phone, ok := isLogin(r); ok {
		logrus.WithFields(logrus.Fields{
			"user.phone": phone,
		}).Info("利用Cookie登录成功")
		Response200(w, "登录成功")
		return
	}

	vars := mux.Vars(r)
	phone := vars["phone"]
	smscode := vars["smscode"]
	password := vars["password"]

	if phone == "" || smscode == "" || password == "" {
		Response400(w, "请求参数不全")
		return
	}
	if !Phone.MatchString(phone) {
		Response400(w, "无效的手机号码")
		return
	}

	if len(password) < 6 || len(password) > 20 {
		Response400(w, "无效的密码")
		return
	}

	user := store.GetUserByPhone(phone)
	var err error
	if user.Phone == "" {
		fmt.Printf("用户 %s 注册，并登录\n", user.Phone)
		if user, err = store.Register(phone, smscode, password); err != nil {
			Response500(w, err.Error())
			return
		}

	} else {
		fmt.Printf("用户忘记密码(%s)，进行重置密码，并登录\n", phone)
		user, err = store.ForgetPassword(fmt.Sprintf("%d", user.SID), smscode, password)
		if err != nil {
			Response500(w, err.Error())
			return
		}
	}
	// 登陆成功，返回cookie
	createUIDC(w, r, store.GetUserByPhone(phone))
	return
}

func createUIDC(w http.ResponseWriter, r *http.Request, user *model.User) {
	fmt.Println("创建cookie")
	now := time.Now()
	expires := now.Add(time.Duration(time.Hour * 24 * 14))
	m := fmt.Sprintf("%X", md5.Sum(([]byte(KEYGEN + user.Phone + now.String()))))
	head := hex.EncodeToString([]byte(m))
	UIDC := base64.StdEncoding.EncodeToString([]byte(head + "_phone_" + user.Phone))

	// 更新User的token信息
	user.LoginDate = now
	user.Token = UIDC
	if _, err := store.UpdateUser(*user); err != nil {
		Response500(w, err.Error())
		return
	}

	// Set Cookies
	cookie := http.Cookie{
		Name:     "UIDC",
		Value:    UIDC,
		Path:     "/",
		Domain:   "easylifewell.com",
		Expires:  expires,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	Response200(w, "登陆成功")

}
func (dc SMSController) SendSMS(w http.ResponseWriter, r *http.Request) {
	if phone, ok := isLogin(r); ok {
		logrus.WithFields(logrus.Fields{
			"user.phone": phone,
		}).Info("利用Cookie登录成功")
		Response200(w, "登录成功")
		return
	}

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

	var isGuest bool
	user := store.GetUserByPhone(phone)
	if user.Phone == "" {
		isGuest = true
	} else {
		isGuest = false
	}

	// Get the random code
	code := getRandomCode()
	logrus.WithFields(logrus.Fields{
		"isGuest": isGuest,
		"code":    code,
	}).Info("发送验证码信息")
	respCode, err := sms.SendSMS(ctx, phone, code, isGuest)
	if err != nil {
		Response500(w, err.Error())
		return
	}

	if respCode != "000000" {
		Response500(w, "短信发送失败")
		return
	}
	logrus.WithFields(logrus.Fields{
		"isGuest": isGuest,
		"code":    code,
		"phone":   phone,
	}).Info("发送验证码成功")
	Response200(w, "验证码发送成功")
}
