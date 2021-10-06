package db

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/night-codes/types.v1"
	"math/rand"
	"time"
)

type obj map[string]interface{}

type (
	DbStruct struct {
		dbSession        *mgo.Session
		logger           *zap.Logger
		cookieCollection *mgo.Collection
		m                []int32
	}
	IDbStruct interface {
		GetRandomFortune() (structs.Cookie, error)
	}
)

func NewDB(logger *zap.Logger, url, dbName, collectionName string) *DbStruct {
	s, err := mgo.Dial(url)
	if err != nil {
		logger.Fatal("error dialing with db", zap.Error(err))
	}
	DB := DbStruct{
		dbSession:        s,
		logger:           logger,
		cookieCollection: s.DB(dbName).C(collectionName),
		m:                make([]int32, 0),
	}
	var o []obj
	if err := DB.cookieCollection.Find(obj{}).Select(obj{"_id": 1}).All(&o); err != nil {
		logger.Fatal("error finding cookies", zap.Error(err))
	}
	for _, v := range o {
		DB.m = append(DB.m, types.Int32(v["_id"]))
	}
	return &DB
}

func (d *DbStruct) GetRandomFortune() (structs.Cookie, error) {
	rand.Seed(time.Now().UnixNano())
	id := d.m[rand.Intn(len(d.m)-1)]
	var result structs.Cookie
	if err := d.cookieCollection.Find(obj{"_id": id}).One(&result); err != nil {
		return result, err
	}
	return result, nil
}
