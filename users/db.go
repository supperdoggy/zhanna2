package main

import (
	"fmt"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"gopkg.in/mgo.v2"
)

// DbStruct - the main db struct
type DbStruct struct {
	DbSession         *mgo.Session
	UsersCollection   *mgo.Collection
	AdminCollection   *mgo.Collection
	MessageCollection *mgo.Collection
}

type obj map[string]interface{}

// GetDB - returns db object
func (ds *DbStruct) initDB() {
	d, err := mgo.Dial("")
	if err != nil || d == nil {
		panic(err.Error())
	}
	*ds = DbStruct{
		DbSession:         d,
		UsersCollection:   d.DB(mainDbName).C("users"),
		AdminCollection:   d.DB(mainDbName).C("admin"),
		MessageCollection: d.DB(mainDbName).C("messages"),
	}
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
	fmt.Println("all good")
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

func (d *DbStruct) writeMessage(userMsg, botMsg telebot.Message) error {
	if userMsg.Sender.ID == 0 {
		return fmt.Errorf("user id is 0")
	}
	var msg Message = Message{UserID: userMsg.Sender.ID, Message: userMsg, Reply: botMsg, Time: time.Now()}
	return d.MessageCollection.Insert(msg)
}

// getUserMsgCount - returns number of msgs user wrote to zhanna :p
func (d *DbStruct) getUserMsgCount(id int) (int, error) {
	if count, err := d.MessageCollection.Find(obj{"userID": id}).Count(); err != nil {
		return 0, err
	} else {
		return count, nil
	}
}
