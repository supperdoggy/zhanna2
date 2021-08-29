package db

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	den4ikcfg "github.com/supperdoggy/superSecretDevelopement/structs/services/den4ik"
	"gopkg.in/mgo.v2"
	"log"
)

type obj map[string]interface{}
type arr []interface{}

type DbStruct struct {
	DbSession              *mgo.Session
	GameSessionsCollection *mgo.Collection
}

var DB DbStruct

func init() {
	s, err := mgo.Dial("")
	if err != nil {
		log.Fatal("Init error:", err.Error())
	}
	DB = DbStruct{
		DbSession:              s,
		GameSessionsCollection: s.DB(den4ikcfg.DBName).C(den4ikcfg.GameSessionsCollection),
	}
	log.Println("connected to db")
}

func (d *DbStruct) InsertGameSession(session structs.Session) error {
	err := d.GameSessionsCollection.Insert(session)
	return err
}

func (d *DbStruct) GetSession(id int) (session structs.Session, err error) {
	err = d.GameSessionsCollection.Find(obj{"_id": id}).One(&session)
	return
}

func (d *DbStruct) UpdateSession(id int, session structs.Session) error {
	err := d.GameSessionsCollection.Update(obj{"_id": id}, session)
	return err
}

func (d *DbStruct) DeleteSession(id int) error {
	err := d.GameSessionsCollection.Remove(obj{"_id": id})
	return err
}

