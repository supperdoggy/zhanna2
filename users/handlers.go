package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func addOrUpdateUserReq(c *gin.Context) {
	var newUser User
	if err := c.Bind(&newUser); err != nil {
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> binding error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}

	fmt.Println(newUser.Telebot.ID)
	exists, err := DB.userExists(newUser.Telebot.ID)
	if err != nil {
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> userExists() error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}

	newUser.LastOnlineTime = time.Now()
	newUser.LastOnline = newUser.LastOnlineTime.Unix()

	if !exists {
		// inserting
		err = DB.UsersCollection.Insert(newUser)
		if err != nil {
			fmt.Println("handlers.go -> addOrUpdateUserReq() -> insert error:", err.Error())
			c.JSON(400, obj{"err": err.Error()})
			return
		}
		return
	}

	var old User
	if err := DB.UsersCollection.Find(obj{"telebot.id": newUser.Telebot.ID}).One(&old); err != nil {
		fmt.Println(err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	// checks if chat already in user struct
	var in bool
	for _, v := range old.Chats {
		if v.Telebot.ID == newUser.Chats[0].Telebot.ID {
			in = true
		}
	}
	if !in {
		newUser.Chats = append(old.Chats, newUser.Chats...)
	}

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

	if err := DB.UsersCollection.Update(obj{"telebot.id": newUser.Telebot.ID}, newUser); err != nil {
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> update error:", err.Error())
		c.JSON(400, obj{"err": err})
		return
	}
	c.JSON(200, obj{"err": nil})
}

func getFortune(c *gin.Context) {
	var req struct {
		ID int `json:"id" bson:"id" form:"id"`
	}

	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> canGrowFlower() -> binding error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	u, err := DB.getUserFromDbById(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> canGrowFlower() -> cant find user:", err.Error())
		c.JSON(400, obj{"err": "cant find user"})
		return
	}
	fmt.Println((u.LastTimeGotFortuneCookie - 24*60*60) - time.Now().Unix())
	if (u.LastTimeGotFortuneCookie+24*60*60)-time.Now().Unix() > 0 {
		fmt.Println("Day didnt passed")
		c.JSON(400, obj{"err": "day didn`t pass"})
		return
	}

	resp, err := MakeHttpReq(fortuneCookieUrl+"/getRandomFortuneCookie", "GET", nil)
	if err != nil {
		fmt.Println("error making req:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	var result FortuneCookie
	if err := json.Unmarshal(resp, &result); err != nil {
		fmt.Println("fortune cookie unmarshal error")
		c.JSON(400, obj{"err": "unmarshal error"})
		return
	}
	if err := DB.updateLastTimeFortune(req.ID); err != nil {
		fmt.Println("error updating last time fortune:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}

	c.JSON(200, result)
	if ok := saveFortune(req.ID, result); !ok {
		fmt.Println("Failed to save fortune for user", req.ID)
	}
}

func getRandomAnek(c *gin.Context) {
	var req struct {
		ID int `json:"id" bson:"id"`
	}
	if err := c.Bind(&req); err != nil {
		fmt.Println("binding error -> getRandomAnek():", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	if req.ID == 0 {
		c.JSON(400, obj{"err": "id cannot be 0"})
		return
	}
	data, err := MakeReqToAnek("getRandomAnek", nil)
	if err != nil {
		fmt.Println("handlers.go -> getRandomAnek()-> req error", err.Error())
		c.JSON(400, obj{"err": "something went wrong, contact @supperdoggy"})
		return
	}
	var result Anek
	if err = json.Unmarshal(data, &result); err != nil {
		fmt.Println("handlers.go -> getRandomAnek() -> unmarshal error:", err.Error())
		c.JSON(400, obj{"err": "Something went wrong, contact @supperdoggy"})
		return
	}
	c.JSON(200, result)
	if ok := saveAnek(req.ID, result); !ok {
		fmt.Println("Not ok saving anek")
	}
}
