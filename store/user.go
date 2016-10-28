package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/easylifewell/purifier-server/model"
)

// AddUser 添加User对象
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

// UpdateUser 更新User对象
func UpdateUser(p model.User) (string, error) {
	query := func(c *mgo.Collection) error {
		return c.Update(bson.M{"phone": p.Phone}, p)
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

func GetLastUser() *model.User {
	user := new(model.User)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"last_user": true}).One(&user)
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
