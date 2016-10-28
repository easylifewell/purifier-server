package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	ID            bson.ObjectId `bson:"_id"`
	SID           int64         `json:"sid" bson:"sid"`
	Phone         string        `json:"phone" bson:"phone"`
	NickName      string        `json:"nickname" bson:"nickname"`
	RealName      string        `json:"realname" bson:"realname"`
	SMSCode       string        `json:"sms_code" bson:"sms_code"`
	SMSSendDate   time.Time     `json:"sms_send_date" bson:"sms_send_date"`
	SMSChangeDate time.Time     `json:"sms_change_date" bson:"sms_change_date"`
	CreateDate    time.Time     `json:"create_date" bson:"create_date"`
	LoginDate     time.Time     `json:"login_date" bson:"login_date"`
	Email         string        `json:"email" bson:"email"`
	Avatar        string        `json:"avatar_url" bson:"avatar_url"`
	Password      string        `bson:"password"`
	LastUser      bool          `json:"last_user" bson:"last_user"`
	Devices       []string      `json:"devices" bson:"devices"`
}
