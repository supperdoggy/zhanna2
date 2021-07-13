package db

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/aneks"
	"gopkg.in/mgo.v2"
	"math/rand"
	"sync"
	"time"
)

type DbStruct struct {
	m          []int
	mut        sync.Mutex
	Collection *mgo.Collection
}

type obj map[string]interface{}

var DB DbStruct

func init() {
	DB.Collection = connectToDb()
	var o []structs.Anek
	if err := DB.Collection.Find(nil).All(&o); err != nil {
		panic(err.Error())
	}
	for _, v := range o {
		DB.m = append(DB.m, v.Id)
	}
}

func connectToDb() *mgo.Collection {
	b, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	return b.DB(cfg.MainDBName).C(cfg.CollectionName)
}

func (d *DbStruct) GetAnekById(id int) (result structs.Anek) {
	if err := d.Collection.Find(obj{"_id": id}).One(&result); err != nil {
		fmt.Println(err.Error())
	}
	return
}

func (d *DbStruct) GetRandomAnek() (result structs.Anek, err error) {
	rand.Seed(time.Now().UnixNano())
	return d.GetAnekById(d.m[rand.Intn(len(d.m)-1)]), err
}

// TODO MAKE REMOVE FROM CACHE
func (d *DbStruct) DeleteAnek(id int) (err error) {
	err = d.Collection.Remove(obj{"_id": id})
	return
}

func (d *DbStruct) AddAnek(text string) (err error) {
	id, err := d.Collection.Count()
	if err != nil {
		return
	}
	a := structs.Anek{
		Id:   id + 1,
		Text: text,
	}

	if err = d.Collection.Insert(a); err != nil {
		return
	}
	d.AddIdToCache(a.Id)
	return
}

func (db *DbStruct) AddIdToCache(id int) {
	db.mut.Lock()
	db.m = append(db.m, id)
	db.mut.Unlock()
}
