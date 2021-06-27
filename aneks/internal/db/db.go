package db

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/db"
	"gopkg.in/mgo.v2"
	"math/rand"
	"time"
)

type DB struct {}

type obj map[string]interface{}

func connectToDb() *mgo.Collection {
	b, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	return b.DB(cfg.MainDBName).C(cfg.CollectionName)
}

var AneksCollection = connectToDb()

func (d *DB ) GetAnekById(id int) (result structs.Anek) {
	if err := AneksCollection.Find(obj{"_id": id}).One(&result); err != nil {
		fmt.Println(err.Error())
	}
	return
}

func (d *DB ) GetRandomAnek() (result structs.Anek, err error) {
	rand.Seed(time.Now().UnixNano())
	size, err := AneksCollection.Count()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return d.GetAnekById(rand.Intn(size - 1)), err
}

func (d *DB ) DeleteAnek(id int) (err error) {
	err = AneksCollection.Remove(obj{"_id": id})
	return
}

func (d *DB ) AddAnek(text string) (err error) {
	id, err := AneksCollection.Count()
	if err != nil {
		return
	}
	a := structs.Anek{
		Id:   id + 1,
		Text: text,
	}

	if err = AneksCollection.Insert(a); err != nil {
		return
	}
	return
}
