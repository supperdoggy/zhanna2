package utils

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"time"
)

func LoadPicsFromSRC() {
	db, err := mgo.Dial("")
	if err != nil {
		panic("error when connecting to db: " + err.Error())
	}
	collection := db.DB(cfg.DBName).C(cfg.PicCollectionName)

	files, err := IOReadDir("/Users/mmarchy/go/src/github.com/supperdoggy/zhanna2/bot/src")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(files)
	for _, v := range files {
		data, err := ioutil.ReadFile("/Users/mmarchy/go/src/github.com/supperdoggy/zhanna2/bot/src/" + v)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		p := structs.Pic{
			ID:          removejpg(v),
			Data:        data,
			TimeCreated: time.Now(),
		}
		err = collection.Insert(p)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func IOReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	files = files[1:]
	return files, nil
}

func removejpg(s string) string {
	return s[:len(s)-4]
}
