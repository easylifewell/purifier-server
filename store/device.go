package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/easylifewell/purifier-server/model"
)

// AddDevice 添加Device对象
func AddDevice(p model.Device) (string, error) {
	p.ID = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection("device", query)
	if err != nil {
		return "", err
	}
	return p.ID.Hex(), nil
}

// GetDeviceByID 根据设备ID获取设备的信息
func GetDeviceByID(id string) *model.Device {
	objid := bson.ObjectIdHex(id)
	device := new(model.Device)
	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&device)
	}
	witchCollection("device", query)
	return device
}

// GetDeviceByUserID 根据userID获取所有的设备
func GetDeviceByUserID(userID string) ([]model.Device, error) {
	var datas []model.Device
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"user_id": userID}).All(&datas)
	}
	err := witchCollection("device", query)
	if err != nil {
		return datas, err
	}
	return datas, nil
}

// UpdateDevice 更新device数据
func UpdateDevice(p model.Device) (string, error) {
	exop := func(c *mgo.Collection) error {
		return c.Update(bson.M{"device_id": p.DeviceID}, p)
	}
	err := witchCollection("device", exop)
	if err != nil {
		return "", err
	}
	return p.ID.Hex(), nil
}

// SearchDevice 执行查询，此方法可拆分做为公共方法
// [SearchDevice description]
// @param {[type]} collectionName string [description]
// @param {[type]} query          bson.M [description]
// @param {[type]} sort           bson.M [description]
// @param {[type]} fields         bson.M [description]
// @param {[type]} skip           int    [description]
// @param {[type]} limit          int)   (results      []interface{}, err error [description]
func SearchDevice(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	err = witchCollection("device", exop)
	return
}
