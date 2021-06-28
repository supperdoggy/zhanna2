package main

import (
	"math/rand"
	"time"

	"gopkg.in/mgo.v2"
)

type obj map[string]interface{}

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

func (d *DbStruct) addFlower(f Flower) (err error) {
	return d.FlowerCollection.Insert(f)
}

func (d *DbStruct) removeFlower(id uint64) (err error) {
	return d.FlowerCollection.Remove(obj{"_id": id})
}

func (d *DbStruct) editFlower(id uint64, f Flower) (err error) {
	return d.FlowerCollection.Update(obj{"_id": id}, obj{"$set": f})
}

func (d *DbStruct) getFlower(id uint64, f Flower) (result Flower, err error) {
	err = d.FlowerCollection.Find(obj{"_id": id}).One(&f)
	return f, err
}

// getAllFlowers - returns all flower types in aneks
func (d *DbStruct) getAllFlowers() (result []Flower, err error) {
	err = d.FlowerCollection.Find(nil).All(&result)
	return
}

func (d *DbStruct) getRandomFlower() (result Flower, err error) {
	rand.Seed(time.Now().Unix())
	count, err := d.FlowerCollection.Count()
	if err != nil {
		return result, err
	}
	id := rand.Intn(count + 1)
	if id == 0 {
		id++
	}
	err = d.FlowerCollection.Find(obj{"_id": id}).One(&result)
	return
}

// returns growing user flower
func (d *DbStruct) getUserFlower(owner int) (result Flower, err error) {
	err = d.UserFlowerDataCollection.Find(obj{"owner": owner, "hp": obj{"$ne": 100}, "dead": false}).One(&result)
	return
}

// returns map of flower name and count
func (d *DbStruct) getAllUserFlowersMap(owner int) (map[string]int, error) {
	resultMap := make(map[string]int)
	resultSlice := []Flower{}
	if err := d.UserFlowerDataCollection.Find(obj{"owner": owner, "hp": 100, "dead": false}).All(&resultSlice); err != nil {
		return nil, err
	}

	if len(resultSlice) == 0 {
		return resultMap, nil
	}

	for _, v := range resultSlice {
		resultMap[v.Icon+" "+v.Name]++
	}
	return resultMap, nil
}

func (d *DbStruct) countFlowers(owner int) (total int, err error) {
	flowers, err := DB.getAllUserFlowersMap(owner)
	if err != nil {
		return
	}
	for _, v := range flowers {
		total += v
	}
	return
}

func (d *DbStruct) getUserFlowerById(id uint64) (Flower, error) {
	var f Flower
	err := d.UserFlowerDataCollection.Find(obj{"id": id}).One(&f)
	return f, err
}

func (d *DbStruct) countUserFlowers(owner int) (int, error) {
	i, err := d.UserFlowerDataCollection.Find(obj{"owner": owner, "hp": 100, "dead": false}).Count()
	return i, err
}

func (d *DbStruct) getAllUserFlowers(owner int) ([]Flower, error) {
	var result []Flower
	err := d.UserFlowerDataCollection.Find(obj{"owner": owner, "hp": 100, "dead": false}).All(&result)
	return result, err
}

// edit user flower
func (d *DbStruct) editUserFlower(id uint64, f Flower) (err error) {
	return d.UserFlowerDataCollection.Update(obj{"_id": id}, obj{"$set": f})
}
