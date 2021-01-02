package main

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/mgo.v2"
)

type DbStruct struct {
	DbSession      *mgo.Session
	TostCollection *mgo.Collection
}

type obj map[string]interface{}

func connectToDB() *mgo.Session {
	d, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	return d
}

func (d *DbStruct) connectToTostCollection() *mgo.Collection {
	return d.DbSession.DB(mainDbName).C("tostv")
}

func getRandomTost() (result Tost, err error) {
	rand.Seed(time.Now().UnixNano())
	size, err := db.TostCollection.Count()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return getTostById(rand.Intn(size - 1)), err
}

func getTostById(id int) (result Tost) {
	if err := db.TostCollection.Find(obj{"_id": id}).One(&result); err != nil {
		fmt.Println(err.Error(), id)
		fmt.Println(err.Error())
	}
	return
}

func deleteTost(id int) (err error) {
	err = db.TostCollection.Remove(obj{"_id": id})
	return
}

func addTost(text string) (err error) {
	id, err := db.TostCollection.Count()
	if err != nil {
		return
	}
	a := Tost{
		ID:   id + 1,
		Text: text,
	}

	if err = db.TostCollection.Insert(a); err != nil {
		return
	}
	return
}
