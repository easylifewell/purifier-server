package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/easylifewell/purifier-server/model"
)

// AddData 添加Data对象
func AddUser(p model.User) (string, error) {
	p.ID = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection("user", query)
	if err != nil {
		return "", err
	}
	return p.ID.Hex(), nil
}

func GetUserBySID(sid string) *model.User {
	user := new(model.User)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"sid": sid}).One(&user)
	}
	witchCollection("user", query)
	return user
}

// GetUserByPhone 根据手机号码获取用户的信息
func GetUserByPhone(phone string) *model.User {
	user := new(model.User)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"phone": phone}).One(&user)
	}
	witchCollection("user", query)
	return user
}
