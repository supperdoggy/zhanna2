package db

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"gopkg.in/mgo.v2"
)

// DbStruct - the main aneks struct
type DbStruct struct {
	DbSession         *mgo.Session
	UsersCollection   *mgo.Collection
	AdminCollection   *mgo.Collection
	MessageCollection *mgo.Collection
}

type obj map[string]interface{}

var DB DbStruct

func init() {
	d, err := mgo.Dial("")
	if err != nil || d == nil {
		panic(err)
	}
	DB = DbStruct{
		DbSession:         d,
		UsersCollection:   d.DB(cfg.DBName).C("users"),
		AdminCollection:   d.DB(cfg.DBName).C("admin"),
		MessageCollection: d.DB(cfg.DBName).C("messages"),
	}
}

func (d *DbStruct) GetUserByID(id int) (result structs.User, err error) {
	err = d.UsersCollection.Find(obj{"telebot.id": id, "statuses.isBanned": false}).One(&result)
	return
}

func (d *DbStruct) UserExists(id int) (bool, error) {
	var u structs.User
	if err := d.UsersCollection.Find(obj{"telebot.id": id}).One(&u); err != nil {
		if err.Error() == "not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *DbStruct) UpdateLastTimeFortune(id int) error {
	return d.UsersCollection.Update(obj{"telebot.id": id}, obj{"$set": obj{"lastTimeGotFortuneCookie": time.Now().Unix(), "lastTimeGotFortuneCookieTime": time.Now()}})
}

func (d *DbStruct) UpdateUser(u structs.User) error {
	return d.UsersCollection.Update(obj{"telebot.id": u.Telebot.ID, "statuses.isBanned": false}, u)
}

// getChatUsersIDs - returns all users ids which are in given chat
func (d *DbStruct) GetChatUsersIDs(chatid int) (ids []int, err error) {
	var users []structs.User
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
func (d *DbStruct) GetChatUsers(chatid int) (users []structs.User, err error) {
	err = d.UsersCollection.Find(obj{"chats.telebot.id": chatid}).All(&users)
	return
}

func (d *DbStruct) WriteMessage(userMsg, botMsg telebot.Message) error {
	if userMsg.Sender.ID == 0 {
		return fmt.Errorf("user id is 0")
	}
	var msg structs.Message = structs.Message{UserID: userMsg.Sender.ID, Message: userMsg, Reply: botMsg, Time: time.Now()}
	return d.MessageCollection.Insert(msg)
}

// getUserMsgCount - returns number of msgs user wrote to zhanna :p
func (d *DbStruct) GetUserMsgCount(userID int) (int, error) {
	if count, err := d.MessageCollection.Find(obj{"userID": userID}).Count(); err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

// appends anek to anek slice and saves user
func (d *DbStruct) SaveAnek(userID int, a structs.Anek) bool {
	u, err := d.GetUserByID(userID)
	if err != nil {
		fmt.Println("Failed to get user", err.Error())
		return false
	}
	u.Aneks = append(u.Aneks, a)
	u.LastTimeGotAnek = time.Now().Unix()
	u.LastTimeGotAnekTime = time.Now()
	err = d.UpdateUser(u)
	if err != nil {
		fmt.Println("Failed to save anek to user")
		return false
	}
	return true
}

func (d *DbStruct) SaveFortune(userID int, a structs.Cookie) bool {
	u, err := d.GetUserByID(userID)
	if err != nil {
		fmt.Println("failed to get user in saveFortune", err.Error())
		return false
	}
	u.FortuneCookies = append(u.FortuneCookies, a)
	err = d.UpdateUser(u)
	if err != nil {
		fmt.Println("error updating user saving fortune")
		return false
	}
	return true
}

func (d *DbStruct) SaveTost(userID int, a structs.Tost) bool {
	u, err := d.GetUserByID(userID)
	if err != nil {
		fmt.Println("failed to get user in saveTost", err.Error())
		return false
	}
	u.Tosts = append(u.Tosts, a)
	err = DB.UpdateUser(u)
	if err != nil {
		fmt.Println("error updating user saving tost", err.Error())
		return false
	}
	return true
}
