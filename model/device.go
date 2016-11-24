package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Device struct {
	ID       bson.ObjectId `bson:"_id"`
	DeviceID string        `json:"ID" bson:"device_id"`
	SPhone   string        `json:"sphone" bson:"sphone"`
	RPhone   string        `json:"rphone" bson:"rphone"`
	CarName  string        `json:"car_name" bson:"car_name"`
	Started  bool          `json:"started" bson:"started"`
	UserID   int64         `json:"user_id" bson:"user_id"`
}
