package store

import (
	"gopkg.in/mgo.v2"
)

const (
	URL = "127.0.0.1:27017"
)

var (
	mgoSession *mgo.Session
	dataBase   = "YuanView"
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
