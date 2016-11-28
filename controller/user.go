package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/easylifewell/purifier-server/store"
	"github.com/gorilla/mux"
)

type UserInfo struct {
	SID        int64     `json:"sid"`
	Phone      string    `json:"phone"`
	NickName   string    `json:"nickname"`
	RealName   string    `json:"realname"`
	CreateDate time.Time `json:"create_date"`
	Email      string    `json:"email"`
	Avatar     string    `json:"avatar_url"`
	Devices    []string  `json:"devices"`
}

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	phone, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	user := store.GetUserByPhone(phone)
	if user.Phone == "" {
		Response400(w, "用户不存在，请您先注册系统")
		return
	}

	u := new(UserInfo)
	u.SID = user.SID
	u.Phone = user.Phone
	u.NickName = user.NickName
	u.RealName = user.RealName
	u.Email = user.Email
	u.Avatar = user.Avatar
	u.Devices = user.Devices
	u.CreateDate = user.CreateDate
	res, err := json.Marshal(u)
	if err != nil {
		Response500(w, "获取用户信息失败")
		return
	}
	ResponseWithJSON(w, res, 200)
	return

}
func (uc UserController) SetNickName(w http.ResponseWriter, r *http.Request) {
	phone, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	vars := mux.Vars(r)
	nickname := vars["nickname"]
	if nickname == "" {
		Response400(w, "请输入请求参数nickname")
		return
	}

	user := store.GetUserByPhone(phone)
	user.NickName = nickname
	if _, err := store.UpdateUser(*user); err != nil {
		Response500(w, "更新NickName失败")
		return
	}
	Response200(w, "更新NickName成功")
	return
}

func (uc UserController) SetRealName(w http.ResponseWriter, r *http.Request) {
	phone, ok := isLogin(r)
	if !ok {
		Response400(w, "请您登录先登录系统")
		return
	}

	vars := mux.Vars(r)
	realname := vars["realname"]
	if realname == "" {
		Response400(w, "请输入请求参数realname")
		return
	}

	user := store.GetUserByPhone(phone)
	user.RealName = realname
	if _, err := store.UpdateUser(*user); err != nil {
		Response500(w, "更新RealName失败")
		return
	}
	Response500(w, "更新RealName成功")
	return
}
