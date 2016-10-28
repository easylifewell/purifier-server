package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/easylifewell/purifier-server/model"
)

// AddGuest 添加Guest对象
func AddGuest(p model.Guest) (string, error) {
	p.ID = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection("guest", query)
	if err != nil {
		return "", err
	}
	return p.ID.Hex(), nil
}

// AddGuest 添加Guest对象
func UpdateGuest(p model.Guest) (string, error) {
	query := func(c *mgo.Collection) error {
		return c.Update(bson.M{"phone": p.Phone}, p)
	}
	err := witchCollection("guest", query)
	if err != nil {
		return "", err
	}
	return p.ID.Hex(), nil
}

// GetGuestByPhone 根据手机号码获取用户的信息
func GetGuestByPhone(phone string) *model.Guest {
	guest := new(model.Guest)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"phone": phone}).One(&guest)
	}
	witchCollection("guest", query)
	return guest
}
