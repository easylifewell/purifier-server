package model

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id"`
	NickName string        `json:"nickname" bson:"nickname"`
	RealName string        `json:"realname" bson:"realname"`
	Email    string        `json:"email" bson:"email"`
	Avatar   string        `json:"avatar_url" bson:"avatar_url"`
	Phone    string        `json:"phone" bson:"phone"`
	Password string        `bson:"password"`
	Code     string        `bson:"code"`
	Devices  []string      `json:"devices" bson:"devices"`
}
