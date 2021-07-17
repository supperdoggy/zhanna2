package db

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/tost"
	"math/rand"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

type DbStruct struct {
	DbSession      *mgo.Session
	TostCollection *mgo.Collection
	m              []int
	mut            sync.Mutex
}

type obj map[string]interface{}

var DB DbStruct

func init() {
	d, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	DB.DbSession = d
	DB.TostCollection = DB.DbSession.DB(cfg.DBName).C(cfg.CollectionName)

	var o []structs.Tost
	if err := DB.TostCollection.Find(nil).Select(obj{"_id": 1}).All(&o); err != nil {
		panic(err.Error())
	}
	for _, v := range o {
		DB.m = append(DB.m, v.ID)
	}
}

func (db *DbStruct) GetRandomTost() (result structs.Tost, err error) {
	rand.Seed(time.Now().UnixNano())
	result, err = db.GetTostById(DB.m[rand.Intn(len(DB.m)-1)])
	return
}

func (db *DbStruct) GetTostById(id int) (result structs.Tost, err error) {
	err = db.TostCollection.Find(obj{"_id": id}).One(&result)
	return
}

func (db *DbStruct) DeleteTost(id int) (err error) {
	err = db.TostCollection.Remove(obj{"_id": id})
	if err != nil {
		return
	}
	db.RemoveFromCacheByVal(id)
	return nil
}

func (db *DbStruct) AddTost(text string) (err error) {
	a := structs.Tost{
		ID:   db.m[len(db.m)]+1,
		Text: text,
	}

	if err = db.TostCollection.Insert(a); err != nil {
		return
	}
	db.AddIdToCache(a.ID)
	return
}

func (db *DbStruct) AddIdToCache(id int) {
	db.mut.Lock()
	db.m = append(db.m, id)
	db.mut.Unlock()
}

func (db *DbStruct) RemoveFromCacheByVal(val int) {
	db.mut.Lock()
	for k, v := range db.m {
		if v == val {
			db.m = append(db.m[:k], db.m[k+1:]...)
		}
	}
	db.mut.Unlock()
}

