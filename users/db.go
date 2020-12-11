package main

import "gopkg.in/mgo.v2"

type DbStruct struct {
	DbSession       *mgo.Session
	UsersCollection *mgo.Collection
	AdminCollection *mgo.Collection
}

func connectToDB() *mgo.Session {
	DB, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	return DB
}

func connectToUsersCollection() *mgo.Collection {
	return DB.DbSession.DB(mainDbName).C("users")
}

func connectToAdminCollection() *mgo.Collection {
	return DB.DbSession.DB(mainDbName).C("admin")
}
