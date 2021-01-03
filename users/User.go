package main

import (
	"fmt"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

type Statuses struct {
	IsBanned       bool      `json:"isBanned" bson:"isBanned"`
	LastBan        int64     `json:"lastBan" bson:"lastBan"` // writes last time user got banned
	LastOnlineTime time.Time `json:"lastOnlineTime" bson:"lastOnlineTime"`

	IsAdmin                bool      `json:"isAdmin" bson:"isAdmin"`
	LastAdminPromotion     int64     `json:"lastAdminPromotion" bson:"lastAdminPromotion"` // writes last time user was promoted to admin
	LastAdminPromotionTime time.Time `json:"lastAdminPromotionTime" bson:"lastAdminPromotionTime"`

	IsPrime bool `json:"isPrime" bson:"isPrime"`
	// writes last time user changed prime status
	LastPrimeStatusChange     int64     `json:"lastPrimeStatusChange" bson:"lastPrimeStatusChange"`
	LastPrimeStatusChangeTime time.Time `json:"lastPrimeStatusChangeTime" bson:"lastPrimeStatusChangeTime"`

	IsDeleted bool `json:"isDeleted" bson:"isDeleted"`
}

type User struct {
	// telebot api user structure
	Telebot        telebot.User `json:"telebot" bson:"telebot"`
	Statuses       Statuses     `json:"statuses" bson:"statuses"`
	LastOnline     int64        `json:"lastOnline" bson:"lastOnline"` // writes last time user wrote to bot
	LastOnlineTime time.Time    `json:"lastOnlineTime" bson:"lastOnlineTime"`
	Chats          []Chat       `json:"chats" bson:"chats"` // chats where user wrote messages
	// all messages user wrote will be stored here
	MessagesUserSent []telebot.Message `json:"messagesUserSent" bson:"messagesUserSent"`
	// all messages Zhanna sent to user
	MessagesZhannaSent []telebot.Message `json:"messagesZhannaSent" bson:"messagesZhannaSent"`
	// all aneks user got
	Aneks               []Anek    `json:"aneks" bson:"aneks"`
	LastTimeGotAnek     int64     `json:"lastTimeGotAnek" bson:"lastTimeGotAnek"`
	LastTimeGotAnekTime time.Time `json:"lastTimeGotAnekTime" bson:"lastTimeGotAnekTime"`
	// all fortune cookies user got
	FortuneCookies               []FortuneCookie `json:"fortuneCookies" bson:"fortuneCookies"`
	LastTimeGotFortuneCookie     int64           `json:"lastTimeGotFortuneCookie" bson:"lastTimeGotFortuneCookie"`
	LastTimeGotFortuneCookieTime time.Time       `json:"lastTimeGotFortuneCookieTime" bson:"lastTimeGotFortuneCookieTime"`
	// tosts
	Tosts               []Tost    `json:"tosts" bson:"tosts"`
	LastTimeGotTost     int64     `json:"lastTimeGotTost" bson:"lastTimeGotTost"`
	LastTimeGotTostTime time.Time `json:"lastTimeGotTostTime" bson:"lastTimeGotTostTime"`
	// todo flowers struct
	Balance uint64 `json:"balance" bson:"balance"`
}

type Chat struct {
	Telebot    telebot.Chat `json:"telebot" bson:"telebot"`
	Users      []User       `json:"users" bson:"users"`
	LastOnline int64        `json:"lastOnline" bson:"lastOnline"`
	Deleted    bool         `json:"deleted" bson:"deleted"`
}

// appends anek to anek slice and saves user
func saveAnek(id int, a Anek) bool {
	u, err := DB.getUserFromDbById(id)
	if err != nil {
		fmt.Println("Failed to get user", err.Error())
		return false
	}
	u.Aneks = append(u.Aneks, a)
	u.LastTimeGotAnek = time.Now().Unix()
	u.LastTimeGotAnekTime = time.Now()
	err = DB.updateUser(u)
	if err != nil {
		fmt.Println("Failed to save anek to user")
		return false
	}
	return true
}

func saveFortune(id int, a FortuneCookie) bool {
	u, err := DB.getUserFromDbById(id)
	if err != nil {
		fmt.Println("failed to get user in saveFortune", err.Error())
		return false
	}
	u.FortuneCookies = append(u.FortuneCookies, a)
	err = DB.updateUser(u)
	if err != nil {
		fmt.Println("error updating user saving fortune")
		return false
	}
	return true
}

func saveTost(id int, a Tost) bool {
	u, err := DB.getUserFromDbById(id)
	if err != nil {
		fmt.Println("failed to get user in saveTost", err.Error())
		return false
	}
	u.Tosts = append(u.Tosts, a)
	err = DB.updateUser(u)
	if err != nil {
		fmt.Println("error updating user saving tost", err.Error())
		return false
	}
	return true
}
