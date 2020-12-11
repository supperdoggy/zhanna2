package main

import (
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

type Statuses struct {
	IsBanned bool `json:"isBanned" bson:"isBanned"`
	LastBan uint64 `json:"lastBan" bson:"lastBan"` // writes last time user got banned
	LastOnlineTime time.Time `json:"lastOnlineTime" bson:"lastOnlineTime"`

	IsAdmin bool `json:"isAdmin" bson:"isAdmin"`
	LastAdminPromotion uint64 `json:"lastAdminPromotion" bson:"lastAdminPromotion"` // writes last time user was promoted to admin
	LastAdminPromotionTime time.Time `json:"lastAdminPromotionTime" bson:"lastAdminPromotionTime"`

	IsPrime bool `json:"isPrime" bson:"isPrime"`
	// writes last time user changed prime status
	LastPrimeStatusChange uint64 `json:"lastPrimeStatusChange" bson:"lastPrimeStatusChange"`
	LastPrimeStatusChangeTime time.Time `json:"lastPrimeStatusChangeTime" bson:"lastPrimeStatusChangeTime"`

	IsDeleted bool `json:"isDeleted" bson:"isDeleted"`
}

type User struct {
	// telebot api user structure
	telebot.User
	Statuses
	LastOnline uint64 `json:"lastOnline" bson:"lastOnline"` // writes last time user wrote to bot
	LastOnlineTime time.Time `json:"lastOnlineTime" bson:"lastOnlineTime"`
	Chats []telebot.Chat `json:"chats" bson:"chats"` // chats where user wrote messages
	// all messages user wrote will be stored here
	MessagesUserSent []telebot.Message `json:"messagesUserSent" bson:"messagesUserSent"`
	// all messages Zhanna sent to user
	MessagesZhannaSent []telebot.Message `json:"messagesZhannaSent" bson:"messagesZhannaSent"`
	// all aneks user got
	Aneks []Anek `json:"aneks" bson:"aneks"`
	LastTimeGotAnek uint64 `json:"lastTimeGotAnek" bson:"lastTimeGotAnek"`
	LastTimeGotAnekTime time.Time `json:"lastTimeGotAnekTime" bson:"lastTimeGotAnekTime"`
	// all fortune cookies user got
	FortuneCookies []FortuneCookie `json:"fortuneCookies" bson:"fortuneCookies"`
	LastTimeGotFortuneCookie uint64 `json:"lastTimeGotFortuneCookie" bson:"lastTimeGotFortuneCookie"`
	LastTimeGotFortuneCookieTime time.Time `json:"lastTimeGotFortuneCookieTime" bson:"lastTimeGotFortuneCookieTime"`
	// todo flowers struct
}

type Chat struct {
	telebot.Chat
	Users []User `json:"users" bson:"users"`
	LastOnline uint64 `json:"lastOnline" bson:"lastOnline"`
}
