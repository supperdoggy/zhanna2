package main

import "gopkg.in/mgo.v2"

var DB = initDb()

type DbStruct struct {
	DbSession        *mgo.Session
	CookieCollection *mgo.Collection
}

func initDb() DbStruct {
	s, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	return DbStruct{
		DbSession:        s,
		CookieCollection: s.DB("Zhanna2").C("Cookies"),
	}
}
