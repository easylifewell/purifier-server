package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Data struct {
	ID       bson.ObjectId `bson:"_id"`
	Version  string        `json:"Version" bson:"version"`
	DeviceID string        `json:"ID" bson:"device_id"`
	Temp     float64       `json:"TEMP" bson:"temperature"`
	Humi     float64       `json:"HUMI" bson:"humidity"`
	Form     float64       `json:"FORM" bson:"formaldehyde"`
	Pm25     float64       `json:"PM2.5" bson:"pm25"`
	Time     time.Time     `bson:"time"`
}
