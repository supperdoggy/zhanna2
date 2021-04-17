package main

import (
	"bufio"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

type Cookie struct {
	ID   int32  `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}

type DbStruct struct {
	DbSession        *mgo.Session
	CookieCollection *mgo.Collection
}

func initDb() DbStruct {
	s, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	return DbStruct{
		DbSession:        s,
		CookieCollection: s.DB("Zhanna2").C("Cookies"),
	}
}

var DB = initDb()

func main() {
	data := scanFile("/Users/maks/go/src/github.com/supperdoggy/superSecretDevelopement/utils/data.txt")
	for k, v := range data {
		DB.CookieCollection.Insert(Cookie{ID: int32(k) + 1, Text: v})
	}
}

func scanFile(filepath string) (data []string) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if len(sc.Text()) >= 2 {
			data = append(data, sc.Text())
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}
	return
}
