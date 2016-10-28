package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Guest struct {
	ID            bson.ObjectId `bson:"_id"`
	Phone         string        `json:"phone" bson:"phone"`
	SMSCode       string        `json:"sms_code" bson:"sms_code"`
	SMSSendDate   time.Time     `json:"sms_send_date" bson:"sms_send_date"`
	SMSChangeDate time.Time     `json:"sms_change_date" bson:"sms_change_date"`
}
