package db

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	"go.uber.org/zap"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"gopkg.in/mgo.v2"
)

// DbStruct - the main aneks struct
type DbStruct struct {
	dbSession         *mgo.Session
	usersCollection   *mgo.Collection
	adminCollection   *mgo.Collection
	messageCollection *mgo.Collection
	logger            *zap.Logger
}

type IDbStruct interface {
	GetUserByID(id int) (result structs.User, err error)
	UserExists(id int) (bool, error)
	UpdateLastTimeFortune(id int) error
	UpdateUser(u structs.User) error
	UpdateUserWithFields(search, set defaultCfg.Obj) error
	GetChatUsersIDs(chatid int) (ids []int, err error)
	GetChatUsers(chatid int) (users []structs.User, err error)
	WriteMessage(userMsg, botMsg telebot.Message) error
	GetUserMsgCount(userID int) (int, error)
	SaveAnek(userID int, a structs.Anek) error
	SaveFortune(userID int, a structs.Cookie) error
	SaveTost(userID int, a structs.Tost) bool
	InserUser(u structs.User) error
}

func NewDB(url, dbName, usersCollection, adminCollection, messagesCollection string, logger *zap.Logger) *DbStruct {
	d, err := mgo.Dial(url)
	if err != nil || d == nil {
		logger.Fatal("error connecting to db", zap.Error(err))
	}
	return &DbStruct{
		dbSession:         d,
		usersCollection:   d.DB(dbName).C(usersCollection),
		adminCollection:   d.DB(dbName).C(adminCollection),
		messageCollection: d.DB(dbName).C(messagesCollection),
		logger:            logger,
	}
}

func (d *DbStruct) GetUserByID(id int) (result structs.User, err error) {
	err = d.usersCollection.Find(defaultCfg.Obj{"telebot.id": id, "statuses.isBanned": false}).One(&result)
	return
}

func (d *DbStruct) UserExists(id int) (bool, error) {
	var u structs.User
	if err := d.usersCollection.Find(defaultCfg.Obj{"telebot.id": id}).One(&u); err != nil {
		if err.Error() == "not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *DbStruct) UpdateLastTimeFortune(id int) error {
	return d.usersCollection.Update(defaultCfg.Obj{"telebot.id": id}, defaultCfg.Obj{"$set": defaultCfg.Obj{"lastTimeGotFortuneCookie": time.Now().Unix(), "lastTimeGotFortuneCookieTime": time.Now()}})
}

func (d *DbStruct) UpdateUser(u structs.User) error {
	return d.usersCollection.Update(defaultCfg.Obj{"telebot.id": u.Telebot.ID, "statuses.isBanned": false}, u)
}

func (d *DbStruct) UpdateUserWithFields(search, set defaultCfg.Obj) error {
	return d.usersCollection.Update(search, set)
}

// getChatUsersIDs - returns all users ids which are in given chat
func (d *DbStruct) GetChatUsersIDs(chatid int) (ids []int, err error) {
	var users []structs.User
	err = d.usersCollection.Find(defaultCfg.Obj{"chats.telebot.id": chatid}).Select(defaultCfg.Obj{"telebot.id": 1}).All(&users)
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
	err = d.usersCollection.Find(defaultCfg.Obj{"chats.telebot.id": chatid}).All(&users)
	return
}

func (d *DbStruct) WriteMessage(userMsg, botMsg telebot.Message) error {
	if userMsg.Sender.ID == 0 {
		return fmt.Errorf("user id is 0")
	}
	var msg structs.Message = structs.Message{UserID: userMsg.Sender.ID, Message: userMsg, Reply: botMsg, Time: time.Now()}
	return d.messageCollection.Insert(msg)
}

// getUserMsgCount - returns number of msgs user wrote to zhanna :p
func (d *DbStruct) GetUserMsgCount(userID int) (int, error) {
	if count, err := d.messageCollection.Find(defaultCfg.Obj{"userID": userID}).Count(); err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

// appends anek to anek slice and saves user
func (d *DbStruct) SaveAnek(userID int, a structs.Anek) error {
	u, err := d.GetUserByID(userID)
	if err != nil {
		d.logger.Error("Failed to get user by id", zap.Error(err), zap.Any("user_id", userID), zap.Any("anek", a))
		return err
	}
	u.Aneks = append(u.Aneks, a)
	u.LastTimeGotAnek = time.Now().Unix()
	u.LastTimeGotAnekTime = time.Now()
	err = d.UpdateUser(u)
	if err != nil {
		d.logger.Error("Failed to save anek", zap.Error(err), zap.Any("user", u))
		return err
	}
	return err
}

func (d *DbStruct) SaveFortune(userID int, a structs.Cookie) error {
	u, err := d.GetUserByID(userID)
	if err != nil {
		d.logger.Error("failed to get user in saveFortune", zap.Error(err), zap.Any("user_id", userID), zap.Any("cookie", a))
		return err
	}
	u.FortuneCookies = append(u.FortuneCookies, a)
	err = d.UpdateUser(u)
	if err != nil {
		d.logger.Error("error updating user saving fortune", zap.Error(err), zap.Any("user", u))
		return err
	}
	return err
}

func (d *DbStruct) SaveTost(userID int, a structs.Tost) bool {
	u, err := d.GetUserByID(userID)
	if err != nil {
		d.logger.Error("failed to get user in saveTost", zap.Error(err), zap.Any("user_id", userID), zap.Any("tost", a))
		return false
	}
	u.Tosts = append(u.Tosts, a)
	err = d.UpdateUser(u)
	if err != nil {
		d.logger.Error("error updating user saving tost", zap.Error(err), zap.Any("user", u))
		return false
	}
	return true
}

func (d *DbStruct) InserUser(u structs.User) error {
	return d.usersCollection.Insert(u)
}
