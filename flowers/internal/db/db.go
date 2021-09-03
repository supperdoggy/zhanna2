package db

import (
	"errors"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	"gopkg.in/mgo.v2"
	"gopkg.in/night-codes/types.v1"
	"log"
	"math/rand"
	"sync"
	"time"
)

type obj map[string]interface{}

type DbStruct struct {
	DbSession                *mgo.Session
	UserFlowerDataCollection *mgo.Collection
	FlowerCollection         *mgo.Collection
	mut                      sync.Mutex
	m                        []uint64
}

var DB DbStruct

func init() {
	s, err := mgo.Dial("")
	if err != nil {
		panic("Init error:" + err.Error())
	}

	DB = DbStruct{
		DbSession:                s,
		UserFlowerDataCollection: s.DB(cfg.DBName).C(cfg.UserFlowerDataCollection),
		FlowerCollection:         s.DB(cfg.DBName).C(cfg.FlowerCollection),
	}

	var allFlowerIDs []obj
	if err := DB.FlowerCollection.Find(nil).All(&allFlowerIDs); err != nil {
		panic(err.Error())
	}
	DB.mut.Lock()
	for _, v := range allFlowerIDs {
		DB.m = append(DB.m, types.Uint64(v["_id"]))
	}
	DB.mut.Unlock()
}

func (d *DbStruct) AddFlower(f structs.Flower) (err error) {
	err = d.FlowerCollection.Insert(f)
	if err != nil {
		return
	}
	d.mut.Lock()
	d.m = append(d.m, f.ID)
	d.mut.Unlock()
	return
}

func (d *DbStruct) RemoveFlower(id uint64) (err error) {
	err = d.FlowerCollection.Remove(obj{"_id": id})
	if err != nil {
		return err
	}
	err = d.removeIDFromCache(id)
	if err != nil {
		log.Println("error removing id from cache " + err.Error())
	}
	return err
}

func (d *DbStruct) EditFlower(id uint64, f structs.Flower) (err error) {
	return d.FlowerCollection.Update(obj{"_id": id}, obj{"$set": f})
}

func (d *DbStruct) GetFlower(id uint64, f structs.Flower) (result structs.Flower, err error) {
	err = d.FlowerCollection.Find(obj{"_id": id}).One(&f)
	return f, err
}

// getAllFlowers - returns all flower types
func (d *DbStruct) GetAllFlowers() (result []structs.Flower, err error) {
	err = d.FlowerCollection.Find(nil).All(&result)
	return
}

func (d *DbStruct) GetRandomFlower() (result structs.Flower, err error) {
	err = d.FlowerCollection.Find(obj{"_id": d.GetRandomID()}).One(&result)
	return
}

// returns growing user flower
func (d *DbStruct) GetUserCurrentFlower(owner int) (result structs.Flower, err error) {
	err = d.UserFlowerDataCollection.Find(obj{"owner": owner, "hp": obj{"$ne": 100}, "dead": false}).One(&result)
	return
}

func (d *DbStruct) CountFlowers(owner int) (total int, err error) {
	total, err = DB.UserFlowerDataCollection.Find(obj{"owner":owner}).Count()
	return
}

func (d *DbStruct) GetUserFlowerById(id uint64) (structs.Flower, error) {
	var f structs.Flower
	err := d.UserFlowerDataCollection.Find(obj{"id": id}).One(&f)
	return f, err
}

func (d *DbStruct) GetAllUserFlowers(owner int) ([]structs.Flower, error) {
	var result []structs.Flower
	err := d.UserFlowerDataCollection.Find(obj{"owner": owner, "hp": 100, "dead": false}).All(&result)
	return result, err
}

// edit user flower
func (d *DbStruct) EditUserFlower(id uint64, f structs.Flower) (err error) {
	return d.UserFlowerDataCollection.Update(obj{"_id": id}, obj{"$set": f})
}

func (d *DbStruct) GetRandomID() uint64 {
	rand.Seed(time.Now().UnixNano())
	d.mut.Lock()
	id := d.m[rand.Intn(len(d.m))]
	d.mut.Unlock()
	return id
}

func (d *DbStruct) removeIDFromCache(val uint64) error {
	var i int
	var ok bool
	d.mut.Lock()
	for k, v := range d.m {
		if val == v {
			i = k
			ok = true
			break
		}
	}
	if !ok {
		return errors.New("no such element in slice")
	}

	d.m[len(d.m)-1], d.m[i] = d.m[i], d.m[len(d.m)-1]
	d.m = d.m[:len(d.m)-1]
	d.mut.Unlock()
	return nil
}
