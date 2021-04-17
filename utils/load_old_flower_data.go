package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	ai "github.com/night-codes/mgo-ai"
	"gopkg.in/mgo.v2"
)

// Flower - struct of flower
type Flower struct {
	ID   uint64 `json:"id" bson:"_id"`
	HP   uint8  `json:"durability" bson:"hp,omitempty"`
	Name string `json:"name" bson:"name"`
	Icon string `json:"icon" bson:"icon"`
	Type string `json:"type" bson:"type"`
	// owner id
	Owner int   `json:"owner" bson:"owner"`
	Grew  uint8 `json:"grew" bson:"grew"`
	Dead  bool  `json:"dead" bson:"dead"`

	CreationTime time.Time `json:"creationTime" bson:"creationTime"`
	LastTimeGrow time.Time `json:"lastTimeGrow" bson:"lastTimeGrow"`
}

type Type struct {
	Fullname   string `json:"fullname"`
	Durability int    `json:"durability"`
	Icon       string `json:"icon"`
	Name       string `json:"name"`
}

type flowerData struct {
	ID        int      `json:"id"`
	Username  string   `json:"username"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Types     []Flower `json:"types"`
}

var DB = initDb()

type DbStruct struct {
	DbSession                *mgo.Session
	UserFlowerDataCollection *mgo.Collection
	FlowerCollection         *mgo.Collection
}

func initDb() DbStruct {
	s, err := mgo.Dial("")
	if err != nil {
		panic("Init error:" + err.Error())
	}
	return DbStruct{
		DbSession:                s,
		UserFlowerDataCollection: s.DB("Zhanna2").C("UserFlowerDataCollection"),
		FlowerCollection:         s.DB("Zhanna2").C("FlowerCollection"),
	}
}

func main1() {
	ai.Connect(DB.FlowerCollection)
	ai.Connect(DB.UserFlowerDataCollection)
	log.Println("Start...")
	files := getAllFilesPaths()
	if len(files) == 0 {
		log.Println("files are empty")
		return
	}

	parsed := getParsedData(files)
	if len(parsed) == 0 {
		log.Println("error getting parsed data all are empty")
		return
	}
	log.Printf("got %v users\n", len(parsed))
	flowers := createFlowers(parsed)
	if len(flowers) == 0 {
		log.Println("error creating flowers")
		return
	}
	writeFlowers(flowers)
	log.Println("Done!")
}

func writeFlowers(data []Flower) {
	log.Println("starting inserting flower...")
	for _, v := range data {
		v.ID = ai.Next(DB.UserFlowerDataCollection.Name)
		if err := DB.UserFlowerDataCollection.Insert(v); err != nil {
			go log.Println("...error inserting flower", v)
		}
		go log.Println("...inserted id", v.ID)
	}
	log.Println("finished inserting")
}

func createFlowers(fData []*flowerData) (result []Flower) {
	log.Println("started creating flowers...")
	for _, v := range fData {
		for _, t := range v.Types {
			t.Owner = v.ID
			t.CreationTime = time.Now()
			result = append(result, t)
		}
	}
	log.Println("finished created flowers with amount", len(result))
	return result
}

func getParsedData(paths []string) (result []*flowerData) {
	log.Println("started parsing data...")
	for _, v := range paths {
		data, err := ioutil.ReadFile(v)
		if err != nil {
			log.Println("error reading file", v)
			continue
		}
		var f flowerData
		if err := json.Unmarshal(data, &f); err != nil {
			log.Println("error unmarshal data", string(v), err.Error())
			continue
		}
		result = append(result, &f)
	}
	log.Println("parsed data with len", len(result))
	return
}

func getAllFilesPaths() (result []string) {
	log.Println("start getting files path")
	err := filepath.Walk("data_transfer/flower_data",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			result = append(result, path)
			// fmt.Println(path, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	log.Println("done reading files with len", len(result))
	return
}
