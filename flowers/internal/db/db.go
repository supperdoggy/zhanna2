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

type (
	DbStruct struct {
		dbSession                *mgo.Session
		logger                   *zap.Logger
		userFlowerDataCollection *mgo.Collection
		flowerCollection         *mgo.Collection
		mut                      sync.Mutex
		m                        []uint64
	}
	IDbStruct interface {
		AddFlower(f structs.Flower) (err error)
		RemoveFlower(id uint64) (err error)
		EditFlower(id uint64, f structs.Flower) (err error)
		GetFlower(id uint64) (result structs.Flower, err error)
		GetAllFlowers() (result []structs.Flower, err error)
		GetRandomFlower() (result structs.Flower, err error)
		GetUserCurrentFlower(owner int) (result structs.Flower, err error)
		CountFlowers(owner int) (total int, err error)
		GetUserFlowerByName(owner int, name string) (structs.Flower, error)
		GetAllUserFlowers(owner int) ([]structs.Flower, error)
		RemoveUserFlower(cryteria defaultCfg.Obj) error
		EditUserFlower(f structs.Flower) (err error)
		GetRandomID() uint64
		UserFlowerSlice(ids []int) (result []structs.Flower, err error)
		removeIDFromCache(val uint64) error
		GetUserFlowerDataCollection() *mgo.Collection
		GetLastUserFlower(owner int) (structs.Flower, error)
		CreateUserFlower(f structs.Flower) error
	}
)

func NewDB(logger *zap.Logger, url, dbName, userFlowerDataCollection, flowerCollection string) *DbStruct {
	s, err := mgo.Dial(url)
	if err != nil {
		logger.Fatal("error dialing with db", zap.Error(err))
	}

	DB := DbStruct{
		dbSession:                s,
		userFlowerDataCollection: s.DB(dbName).C(userFlowerDataCollection),
		flowerCollection:         s.DB(cfg.DBName).C(flowerCollection),
		logger:                   logger,
	}
	ai.Connect(DB.flowerCollection)
	ai.Connect(DB.userFlowerDataCollection)

	var allFlowerIDs []defaultCfg.Obj
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

func (d *DbStruct) GetUserFlowerDataCollection() *mgo.Collection {
	return d.userFlowerDataCollection
}

// AddFlower - adds new flower type
func (d *DbStruct) AddFlower(f structs.Flower) (err error) {
	f.ID = uint64(rand.Int63())
	err = d.flowerCollection.Insert(f)
	// check for if err is duplicate error
	if mgo.IsDup(err) {
		return d.AddFlower(f)
	} else if err != nil {
		return
	}
	d.mut.Lock()
	d.m = append(d.m, f.ID)
	d.mut.Unlock()
	return
}

func (d *DbStruct) RemoveFlower(id uint64) (err error) {
	err = d.flowerCollection.Remove(defaultCfg.Obj{"_id": id})
	if err != nil {
		return err
	}
	err = d.removeIDFromCache(id)
	if err != nil {
		d.logger.Error("error removing id from cache", zap.Error(err), zap.Any("id", id))
	}
	return err
}

func (d *DbStruct) EditFlower(id uint64, f structs.Flower) (err error) {
	return d.flowerCollection.Update(defaultCfg.Obj{"_id": id}, defaultCfg.Obj{"$set": f})
}

func (d *DbStruct) GetFlower(id uint64) (result structs.Flower, err error) {
	err = d.flowerCollection.Find(defaultCfg.Obj{"_id": id}).One(&result)
	return result, err
}

// getAllFlowers - returns all flower types
func (d *DbStruct) GetAllFlowers() (result []structs.Flower, err error) {
	err = d.flowerCollection.Find(nil).All(&result)
	return
}

func (d *DbStruct) GetRandomFlower() (result structs.Flower, err error) {
	err = d.flowerCollection.Find(defaultCfg.Obj{"_id": d.GetRandomID()}).One(&result)
	return
}

// returns growing user flower
func (d *DbStruct) GetUserCurrentFlower(owner int) (result structs.Flower, err error) {
	err = d.userFlowerDataCollection.Find(defaultCfg.Obj{"owner": owner, "hp": defaultCfg.Obj{"$ne": 100}, "dead": false}).One(&result)
	return
}

func (d *DbStruct) CountFlowers(owner int) (total int, err error) {
	total, err = d.userFlowerDataCollection.Find(defaultCfg.Obj{"owner": owner}).Count()
	return
}

func (d *DbStruct) GetUserFlowerByName(owner int, name string) (structs.Flower, error) {
	var f structs.Flower
	err := d.userFlowerDataCollection.Find(defaultCfg.Obj{"owner": owner, "name": name}).One(&f)
	return f, err
}

func (d *DbStruct) GetAllUserFlowers(owner int) ([]structs.Flower, error) {
	var result []structs.Flower
	err := d.userFlowerDataCollection.Find(defaultCfg.Obj{"owner": owner, "hp": 100, "dead": false}).Sort("lastTimeGrow").All(&result)
	return result, err
}

func (d *DbStruct) GetLastUserFlower(owner int) (structs.Flower, error) {
	// getting flowers
	flowers, err := d.GetAllUserFlowers(owner)
	if err != nil || len(flowers) == 0 { // if has no flower
		d.logger.Error("error GetAllUserFlowers", zap.Error(err))
		return structs.Flower{}, errors.New("user has no flowers")
	}
	return flowers[len(flowers)-1], nil
}

func (d *DbStruct) RemoveUserFlower(cryteria defaultCfg.Obj) error {
	return d.userFlowerDataCollection.Remove(cryteria)
}

// EditUserFlower - using for adding and updating new user flower
func (d *DbStruct) EditUserFlower(f structs.Flower) (err error) {
	_, err = d.userFlowerDataCollection.Upsert(defaultCfg.Obj{"_id": f.ID}, defaultCfg.Obj{"$set": f})
	if mgo.IsDup(err) {
		return d.EditUserFlower(f)
	}
	return
}

func (d *DbStruct) CreateUserFlower(f structs.Flower) error {
	// just made random int generator hope it will be fine....
	f.ID = uint64(rand.Int63())
	err := d.userFlowerDataCollection.Insert(f)
	if mgo.IsDup(err) {
		return d.CreateUserFlower(f)
	}
	return err
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
	err = d.userFlowerDataCollection.Find(defaultCfg.Obj{"$and": defaultCfg.Arr{defaultCfg.Obj{"$or": query}, defaultCfg.Obj{"dead": false}}}).Select(defaultCfg.Obj{"owner": 1, "hp": 1}).All(&result)
	return
}
