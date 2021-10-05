package db

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg2 "github.com/supperdoggy/superSecretDevelopement/structs/services/fortune"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/night-codes/types.v1"
	"math/rand"
	"time"
)

type obj map[string]interface{}

type DbStruct struct {
	DbSession        *mgo.Session
	Logger *zap.Logger
	cookieCollection *mgo.Collection
	m                []int32
}

var (
	DB = getDB()
)

func getDB() *DbStruct {
	logger, _ := zap.NewDevelopment()
	s, err := mgo.Dial("")
	if err != nil {
		logger.Fatal("error dialing with db", zap.Error(err))
	}
	DB := DbStruct{
		DbSession:        s,
		Logger: logger,
		cookieCollection: s.DB(cfg2.DBName).C(cfg2.FortuneCollection),
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
	if err := DB.cookieCollection.Find(obj{"_id": id}).One(&result); err != nil {
		return result, err
	}
	return result, nil
}
