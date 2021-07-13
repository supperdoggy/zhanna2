package db

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/NHIE"
	"math/rand"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

type DBStruct struct {
	DbSession                *mgo.Session
	NeverHaveIEverCollection *mgo.Collection
	QuestionsCache           []structs.NHIE
	mut                      sync.Mutex
}

var (
	DB DBStruct
)

func init() {
	s, err := mgo.Dial("")
	if err != nil {
		panic("Init error:" + err.Error())
	}
	DB = DBStruct{
		DbSession:                s,
		NeverHaveIEverCollection: s.DB(cfg.DBNAME).C(cfg.NHIECollection),
	}
	// OK only if we have not that much data
	if err := DB.NeverHaveIEverCollection.Find(nil).All(&DB.QuestionsCache); err != nil {
		panic(err.Error())
	}
}

func (d *DBStruct) GetRandomNHIE() structs.NHIE {
	rand.Seed(time.Now().Unix())
	d.mut.Lock()
	result := d.QuestionsCache[rand.Intn(len(d.QuestionsCache))]
	d.mut.Unlock()
	return result
}
