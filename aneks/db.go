package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"math/rand"
	"time"
)

func connectToDb() *mgo.Collection{
	b, err := mgo.Dial("")
	if err != nil{
		panic(err.Error())
	}
	return b.DB("Zhanna2").C("Aneks")
}


var AneksCollection = connectToDb()

func getAnekById(id int) (result Anek){
	if err := AneksCollection.Find(obj{"_id":id}).One(&result); err != nil{
		fmt.Println(err.Error())
	}
	return
}

func getRandomAnek() (result Anek){
	rand.Seed(time.Now().UnixNano())
	size, err := AneksCollection.Count()
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	return getAnekById(rand.Intn(size-1))
}

