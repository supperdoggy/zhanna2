package db

import (
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	"gopkg.in/mgo.v2"
	"log"
)

type DbStruct struct {
	DbSession *mgo.Session
	PicCollection *mgo.Collection
}

var DB DbStruct

func init() {
	db, err := mgo.Dial("")
	if err != nil {
		panic("error when connecting to db: " + err.Error())
	}
	DB = DbStruct{
		DbSession:     db,
		PicCollection: db.DB(cfg.DBName).C(cfg.PicCollectionName),
	}
	log.Println("connected to db")
}
