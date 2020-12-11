package main

import "gopkg.in/mgo.v2"

type DbStruct struct {
	DbSession          *mgo.Session
	MessagesCollection *mgo.Collection
	CfgCollection      *mgo.Collection
	PhrasesCollection  *mgo.Collection
}

func connectToDB() *mgo.Session {
	DB, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	return DB
}

func connectToMessagesCollection() *mgo.Collection {
	return DB.DbSession.DB(mainDbName).C("messages")
}

func connectToCfgCollection() *mgo.Collection {
	return DB.DbSession.DB(mainDbName).C("cfg")
}

func connectToPhrasesCollection() *mgo.Collection {
	return DB.DbSession.DB(mainDbName).C("phrases")
}
