package db

import (
	"errors"
	ai "github.com/night-codes/mgo-ai"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/night-codes/types.v1"
	"math/rand"
	"sync"
	"time"
)

type obj map[string]interface{}

type DbStruct struct {
	DbSession                *mgo.Session
	Logger                   *zap.Logger
	userFlowerDataCollection *mgo.Collection
	flowerCollection         *mgo.Collection
	mut                      sync.Mutex
	m                        []uint64
}

var DB = getDB()

func getDB() *DbStruct {
	logger, _ := zap.NewDevelopment()
	s, err := mgo.Dial("")
	if err != nil {
		logger.Fatal("error dialing with db", zap.Error(err))
	}

	DB := DbStruct{
		DbSession:                s,
		userFlowerDataCollection: s.DB(cfg.DBName).C(cfg.UserFlowerDataCollection),
		flowerCollection:         s.DB(cfg.DBName).C(cfg.FlowerCollection),
		Logger:                   logger,
	}
	ai.Connect(DB.flowerCollection)
	ai.Connect(DB.userFlowerDataCollection)

	var allFlowerIDs []obj
	if err := DB.flowerCollection.Find(nil).All(&allFlowerIDs); err != nil {
		logger.Fatal("error finding all flower ids", zap.Error(err), zap.Any("db", DB))
	}
	DB.mut.Lock()
	for _, v := range allFlowerIDs {
		DB.m = append(DB.m, types.Uint64(v["_id"]))
	}
	DB.mut.Unlock()
	return &DB
}

func (d *DbStruct) AddFlower(f structs.Flower) (err error) {
	f.ID = ai.Next(d.flowerCollection.Name)
	err = d.flowerCollection.Insert(f)
	if err != nil {
		return
	}
	d.mut.Lock()
	d.m = append(d.m, f.ID)
	d.mut.Unlock()
	return
}

func (d *DbStruct) RemoveFlower(id uint64) (err error) {
	err = d.flowerCollection.Remove(obj{"_id": id})
	if err != nil {
		return err
	}
	err = d.removeIDFromCache(id)
	if err != nil {
		d.Logger.Error("error removing id from cache", zap.Error(err), zap.Any("id", id))
	}
	return err
}

func (d *DbStruct) EditFlower(id uint64, f structs.Flower) (err error) {
	return d.flowerCollection.Update(obj{"_id": id}, obj{"$set": f})
}

func (d *DbStruct) GetFlower(id uint64, f structs.Flower) (result structs.Flower, err error) {
	err = d.flowerCollection.Find(obj{"_id": id}).One(&f)
	return f, err
}

// getAllFlowers - returns all flower types
func (d *DbStruct) GetAllFlowers() (result []structs.Flower, err error) {
	err = d.flowerCollection.Find(nil).All(&result)
	return
}

func (d *DbStruct) GetRandomFlower() (result structs.Flower, err error) {
	err = d.flowerCollection.Find(obj{"_id": d.GetRandomID()}).One(&result)
	return
}

// returns growing user flower
func (d *DbStruct) GetUserCurrentFlower(owner int) (result structs.Flower, err error) {
	err = d.userFlowerDataCollection.Find(obj{"owner": owner, "hp": obj{"$ne": 100}, "dead": false}).One(&result)
	return
}

func (d *DbStruct) CountFlowers(owner int) (total int, err error) {
	total, err = DB.userFlowerDataCollection.Find(obj{"owner": owner}).Count()
	return
}

func (d *DbStruct) GetUserFlowerById(id uint64) (structs.Flower, error) {
	var f structs.Flower
	err := d.userFlowerDataCollection.Find(obj{"id": id}).One(&f)
	return f, err
}

func (d *DbStruct) GetAllUserFlowers(owner int) ([]structs.Flower, error) {
	var result []structs.Flower
	err := d.userFlowerDataCollection.Find(obj{"owner": owner, "hp": 100, "dead": false}).All(&result)
	return result, err
}

func (d *DbStruct) RemoveUserFlower(cryteria defaultCfg.Obj) error {
	return d.userFlowerDataCollection.Remove(cryteria)
}

// edit user flower
func (d *DbStruct) EditUserFlower(f structs.Flower) (err error) {
	if f.ID == 0 {
		f.ID = ai.Next(d.flowerCollection.Name)
	}
	_, err = d.userFlowerDataCollection.Upsert(obj{"_id": f.ID}, obj{"$set": f})
	return
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
	// kinda slow but I assume we wont have more than 300 elements in slice
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

	// removing element
	d.m[len(d.m)-1], d.m[i] = d.m[i], d.m[len(d.m)-1]
	d.m = d.m[:len(d.m)-1]
	d.mut.Unlock()
	return nil
}

func (d *DbStruct) UserFlowerSlice(ids []int) (result []structs.Flower, err error) {
	// building query for request to mongo
	query := make([]defaultCfg.Obj, 1)
	for _, v := range ids {
		query = append(query, defaultCfg.Obj{"owner": v})
	}
	err = d.userFlowerDataCollection.Find(obj{"$and": defaultCfg.Arr{obj{"$or": query}, obj{"dead": false}}}).Select(obj{"owner": 1, "hp": 1}).All(&result)
	return
}
