package main

import (
	"time"

	"gopkg.in/mgo.v2"
)

type DbStruct struct {
	DbSession       *mgo.Session
	UsersCollection *mgo.Collection
	AdminCollection *mgo.Collection
}

type obj map[string]interface{}

func connectToDB() *mgo.Session {
	d, err := mgo.Dial("")
	if err != nil {
		panic(err.Error())
	}
	return d
}

func connectToUsersCollection() *mgo.Collection {
	return DB.DbSession.DB(mainDbName).C("users")
}

func connectToAdminCollection() *mgo.Collection {
	return DB.DbSession.DB(mainDbName).C("admin")
}

func (d *DbStruct) getUserFromDbById(id int) (result User, err error) {
	err = DB.UsersCollection.Find(obj{"telebot.id": id, "statuses.isBanned": false}).One(&result)
	return
}

func (d *DbStruct) userExists(id int) (bool, error) {
	var u User
	if err := d.UsersCollection.Find(obj{"telebot.id": id}).One(&u); err != nil {
		if err.Error() == "not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *DbStruct) updateLastTimeFortune(id int) error {
	return d.UsersCollection.Update(obj{"telebot.id": id}, obj{"$set": obj{"lastTimeGotFortuneCookie": time.Now().Unix(), "lastTimeGotFortuneCookieTime": time.Now()}})
}

func (d *DbStruct) updateUser(u User) error {
	return d.UsersCollection.Update(obj{"telebot.id": u.Telebot.ID, "statuses.isBanned": false}, u)
}

// getChatUsersIDs - returns all users ids which are in given chat
func (d *DbStruct) getChatUsersIDs(chatid int) (ids []int, err error) {
	users := []User{}
	err = d.UsersCollection.Find(obj{"chats.telebot.id": chatid}).Select(obj{"telebot.id": 1}).All(&users)
	if err != nil {
		return
	}
	for _, v := range users {
		ids = append(ids, v.Telebot.ID)
	}
	return
}

// getChatUsers - returns all users which are in given chat
func (d *DbStruct) getChatUsers(chatid int) (users []User, err error) {
	err = d.UsersCollection.Find(obj{"chats.telebot.id": chatid}).All(&users)
	return
}
