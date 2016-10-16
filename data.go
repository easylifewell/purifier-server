package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	URL = "127.0.0.1:27017"
)

type Data struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Version  string        `json:"Version" bson:"version"`
	DeviceID string        `json:"ID" bson:"device_id"`
	Temp     float64       `json:"TEMP" bson:"temperature"`
	Humi     float64       `json:"HUMI" bson:"humidity"`
	Form     float64       `json:"FORM" bson:"formaldehyde"`
	Pm25     float64       `json:"PM2.5" bson:"pm25"`
	Time     time.Time     `bson:"time"`
}

var (
	mgoSession *mgo.Session
	dataBase   = "YuanView"
	collection = "yuanview"
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

// AddData 添加Data对象
func AddData(p Data) (string, error) {
	p.ID = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection(collection, query)
	if err != nil {
		return "", err
	}
	return p.ID.Hex(), nil
}

// GetDataByID 获取一条记录通过objectid
func GetDataByID(id string) *Data {
	objid := bson.ObjectIdHex(id)
	data := new(Data)
	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&data)
	}
	witchCollection(collection, query)
	return data
}

// GetDataByDeviceID 根据设备ID获取数据
func GetDataByDeviceID(deviceID string) ([]Data, error) {
	var datas []Data
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"device_id": deviceID}).All(&datas)
	}
	err := witchCollection(collection, query)
	if err != nil {
		return datas, err
	}
	return datas, nil
}

// GetDatas 获取所有的data数据
func GetDatas() ([]Data, error) {
	var datas []Data
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&datas)
	}
	err := witchCollection(collection, query)
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
	err := witchCollection(collection, exop)
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
	err = witchCollection(collectionName, exop)
	return
}
