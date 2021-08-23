package db

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

type DbStruct struct {
	DbSession *mgo.Session
	PicCollection *mgo.Collection
}

type obj map[string]interface{}

var DB DbStruct

const (
	purpleLogoID = "logo_purple"
	purpleRulesID = "rules_purple"
	yellowLogoID = "logo_yellow"
	yellowRulesID = "rules_yellow"
)

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

// GetLogoAndRules - returns logo and rules according to time now
func (db DbStruct) GetLogoAndRules() (logo *structs.Pic, rules *structs.Pic, err error){
	hour := time.Now().Hour()
	if hour > 19 && hour < 10 {
		// then pick dark logo and rules
		logo, err := db.GetPicFromDB(purpleLogoID)
		if err != nil {
			return nil, nil, err
		}
		rules, err := db.GetPicFromDB(purpleRulesID)
		if err != nil {
			return nil, nil, err
		}
		return logo, rules, nil
	} else {
		// else pick light logo and rules
		logo, err := db.GetPicFromDB(yellowLogoID)
		if err != nil {
			return nil, nil, err
		}
		rules, err := db.GetPicFromDB(yellowRulesID)
		if err != nil {
			return nil, nil, err
		}
		return logo, rules, nil
	}
}
