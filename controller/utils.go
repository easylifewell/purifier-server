package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/easylifewell/purifier-server/store"
)

type ReturnInfo struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func responseWithMessage(w http.ResponseWriter, message string, code int) {
	info := ReturnInfo{
		Message: message,
		Code:    code,
	}

	b, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
	}
	ResponseWithJSON(w, b, code)
}

// Response200 return OK status
func Response200(w http.ResponseWriter, message string) {
	responseWithMessage(w, message, http.StatusOK)
}

// Response400 return when the request is bad
func Response400(w http.ResponseWriter, message string) {
	responseWithMessage(w, message, http.StatusBadRequest)
}

// Response500 return internal server error
func Response500(w http.ResponseWriter, message string) {
	responseWithMessage(w, message, http.StatusInternalServerError)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

// getRandomCode 获取随机的验证吗
func getRandomCode() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return fmt.Sprintf("%d", r1.Int31())[2:6]
}

// isLogin  判断用户是否登录
func isLogin(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("UIDC")
	if err != nil {
		return "", false
	}
	// decode the UIDC
	UIDC, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", false
	}

	fileds := strings.Split(string(UIDC), "_")
	phone := fileds[len(fileds)-1]
	user := store.GetUserByPhone(phone)
	if user.Phone == phone && user.Token == cookie.Value {
		return phone, true
	}
	return "", false
}

func Contains(values []string, value string) bool {
	if len(value) == 0 {
		return false
	}

	for _, i := range values {
		if i == value {
			return true
		}
	}

	return false
}
