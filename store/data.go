package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/easylifewell/purifier-server/model"
)

// AddData 添加Data对象
func AddData(p model.Data) (string, error) {
	p.ID = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection("data", query)
	if err != nil {
		return "", err
	}
	return p.ID.Hex(), nil
}

// GetDataByID 获取一条记录通过objectid
func GetDataByID(id string) *model.Data {
	objid := bson.ObjectIdHex(id)
	data := new(model.Data)
	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&data)
	}
	witchCollection("data", query)
	return data
}

// GetDataByDeviceID 根据设备ID获取数据
func GetDataByDeviceID(deviceID string) ([]model.Data, error) {
	var datas []model.Data
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"device_id": deviceID}).All(&datas)
	}
	err := witchCollection("data", query)
	if err != nil {
		return datas, err
	}
	return datas, nil
}

// GetDatas 获取所有的data数据
func GetDatas() ([]model.Data, error) {
	var datas []model.Data
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&datas)
	}
	err := witchCollection("data", query)
	if err != nil {
		return datas, err
	}
	return datas, nil
}

// UpdateData 更新data数据
func UpdateData(query bson.M, change bson.M) string {
	exop := func(c *mgo.Collection) error {
		return c.Update(query, change)
	}
	err := witchCollection("data", exop)
	if err != nil {
		return "true"
	}
	return "false"
}

// SearchData 执行查询，此方法可拆分做为公共方法
// [SearchData description]
// @param {[type]} collectionName string [description]
// @param {[type]} query          bson.M [description]
// @param {[type]} sort           bson.M [description]
// @param {[type]} fields         bson.M [description]
// @param {[type]} skip           int    [description]
// @param {[type]} limit          int)   (results      []interface{}, err error [description]
func SearchData(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	err = witchCollection("data", exop)
	return
}
