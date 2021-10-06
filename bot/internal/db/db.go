package db

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"time"
)

type DbStruct struct {
	DbSession     *mgo.Session
	PicCollection *mgo.Collection
	Logger *zap.Logger
}

type obj map[string]interface{}

const (
	purpleLogoID  = "logo_purple"
	purpleRulesID = "rules_purple"
	yellowLogoID  = "logo_yellow"
	yellowRulesID = "rules_yellow"
)

func NewDbStruct(logger *zap.Logger, url, dbName, collectionName string) *DbStruct {
	db, err := mgo.Dial(url)
	if err != nil {
		logger.Fatal("error when connecting to db", zap.Error(err))
	}
	logger.Info("connected to db",
		zap.Any("url", url),
		zap.Any("db_name", dbName),
		zap.Any("collection", collectionName))
	return &DbStruct{
		DbSession:     db,
		PicCollection: db.DB(dbName).C(collectionName),
		Logger: logger,
	}
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
func (db DbStruct) GetLogoAndRules() (logo *structs.Pic, rules *structs.Pic, err error) {
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
