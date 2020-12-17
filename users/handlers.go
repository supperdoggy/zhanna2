package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func addOrUpdateUserReq(c *gin.Context) {
	var newUser User
	if err := c.Bind(&newUser); err != nil{
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> binding error:", err.Error())
		c.JSON(400, obj{"err":err.Error()})
		return
	}

	fmt.Println(newUser.Telebot.ID)
	exists, err := DB.userExists(newUser.Telebot.ID)
	if err != nil{
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> userExists() error:", err.Error())
		c.JSON(400, obj{"err":err.Error()})
		return
	}

	newUser.LastOnlineTime = time.Now()
	newUser.LastOnline = newUser.LastOnlineTime.Unix()

	if !exists{
		// inserting
		err = DB.UsersCollection.Insert(newUser)
		if err != nil{
			fmt.Println("handlers.go -> addOrUpdateUserReq() -> insert error:", err.Error())
			c.JSON(400, obj{"err":err.Error()})
			return
		}
		return
	}

	var old User
	if err := DB.UsersCollection.Find(obj{"telebot.id":newUser.Telebot.ID}).One(&old); err != nil{
		fmt.Println(err.Error())
		c.JSON(400, obj{"err":err.Error()})
		return
	}

	newUser.Chats = append(old.Chats, newUser.Chats...)
	newUser.MessagesUserSent = append(old.MessagesUserSent, newUser.MessagesUserSent...)
	newUser.Aneks = append(old.Aneks, newUser.Aneks...)
	newUser.FortuneCookies = append(old.FortuneCookies, newUser.FortuneCookies...)
	newUser.Balance = old.Balance
	newUser.Statuses = old.Statuses
	newUser.LastTimeGotAnek = old.LastTimeGotAnek
	newUser.LastTimeGotAnekTime = old.LastTimeGotAnekTime
	newUser.MessagesZhannaSent = append(old.MessagesZhannaSent, newUser.MessagesZhannaSent...)
	newUser.LastTimeGotFortuneCookie = old.LastTimeGotFortuneCookie
	newUser.LastTimeGotFortuneCookieTime = old.LastTimeGotFortuneCookieTime

	if err := DB.UsersCollection.Update(obj{"telebot.id":newUser.Telebot.ID}, newUser); err != nil{
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> update error:", err.Error())
		c.JSON(400, obj{"err":err})
		return
	}
	c.JSON(200, obj{"err":nil})
}
