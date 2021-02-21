package main

import (
	"math/rand"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	DB        = initDb()
	Questions = struct {
		m []NHIE
		sync.Mutex
	}{m: []NHIE{}}
)

func init() {
	if err := DB.NeverHaveIEverCollection.Find(nil).All(&Questions.m); err != nil {
		panic(err.Error())
	}
}

type NHIE struct {
	ID   int    `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}

type DBStruct struct {
	DbSession                *mgo.Session
	NeverHaveIEverCollection *mgo.Collection
}

func initDb() DBStruct {
	s, err := mgo.Dial("")
	if err != nil {
		panic("Init error:" + err.Error())
	}
	return DBStruct{
		DbSession:                s,
		NeverHaveIEverCollection: s.DB("Zhanna2").C("NeverHaveIEver"),
	}
}

func getRandomNHIE() NHIE {
	rand.Seed(time.Now().Unix())
	return Questions.m[rand.Intn(len(Questions.m))]
}
