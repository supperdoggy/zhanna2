package db

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	"gopkg.in/mgo.v2"
	"log"
)

type DbStruct struct {
	DbSession *mgo.Session
	PicCollection *mgo.Collection
}

type obj map[string]interface{}

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

func (db DbStruct) GetPicFromDB(id string) (*structs.Pic, error) {
	var p structs.Pic
	err := db.PicCollection.Find(obj{"_id": id}).One(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
