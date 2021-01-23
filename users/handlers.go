package main

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/gin-gonic/gin"
)

// todo simplify
func addOrUpdateUserReq(c *gin.Context) {
	var newUser User
	if err := c.Bind(&newUser); err != nil {
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> binding error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}

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
	newUser.MessagesZhannaSent = append(old.MessagesZhannaSent, newUser.MessagesZhannaSent...)

	fieldsToSet := obj{
		"messagesUserSent":   newUser.MessagesUserSent,
		"messagesZhannaSent": newUser.MessagesZhannaSent,
		"lastOnlineTime":     newUser.LastOnlineTime,
		"lastOnline":         newUser.LastOnline,
		"chats":              newUser.Chats,
	}

	if err := DB.UsersCollection.Update(obj{"telebot.id": newUser.Telebot.ID}, obj{"$set": fieldsToSet}); err != nil {
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> update error:", err.Error())
		c.JSON(400, obj{"err": err})
		return
	}
	c.JSON(200, obj{"err": nil})
}

// todo simplify
func getFortune(c *gin.Context) {
	var req struct {
		ID int `json:"id" bson:"id" form:"id"`
	}

	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getFortune() -> binding error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}

	// checking if user exists if not then just create one
	exists, err := DB.userExists(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getFortune() -> userExists() error:", err.Error())
		c.JSON(400, obj{"err": "error making req"})
		return
	}
	if !exists {
		err := DB.UsersCollection.Insert(User{Telebot: telebot.User{ID: req.ID}})
		if err != nil {
			fmt.Println("handlers.go -> addOrUpdateUserReq() -> insert error:", err.Error())
			c.JSON(400, obj{"err": err.Error()})
			return
		}
	}

	u, err := DB.getUserFromDbById(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getFortune() -> cant find user:", err.Error())
		c.JSON(400, obj{"err": "cant find user"})
		return
	}
	// check if day passed to get new fortune
	if !CanGetFortune(u.LastTimeGotFortuneCookieTime) {
		fmt.Println("Day didnt pass")
		c.JSON(400, obj{"err": cantGetFortune})
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

// thats ok i guess
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
		fmt.Println("Not ok saving anek", req.ID)
	}
}

func getRandomTost(c *gin.Context) {
	var req struct {
		ID int `json:"id" bson:"_id"`
	}
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getRandomTost() -> c.Bind() error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}

	if req.ID == 0 {
		c.JSON(400, obj{"err": "binding error"})
		return
	}

	data, err := MakeReqToTost("getRandomTost", nil)
	if err != nil {
		fmt.Println("handlers.go -> getRandomTost() -> MakeReqToTost(\"getRandomTost\") error:", err.Error())
		return
	}
	var result Tost
	if err = json.Unmarshal(data, &result); err != nil {
		fmt.Println("handlers.go -> getRandomTost() -> json.Unmarshal error:", err.Error())
		c.JSON(400, obj{"err": "unmarshal error"})
		return
	}

	c.JSON(200, result)
	if ok := saveTost(req.ID, result); !ok {
		fmt.Println("not ok saving tost", req.ID)
	}
}

func addFlower(c *gin.Context) {
	var req struct {
		Icon string `json:"icon" bson:"icon"`
		Name string `json:"name" bson:"name"`
		Type string `json:"type" bson:"type"`
	}
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> addFlower() -> binding error:", err.Error())
		c.JSON(400, obj{"err": "binding error"})
		return
	}
	marshaled, err := json.Marshal(req)
	if err != nil {
		fmt.Println("handlers.go -> addFlower() -> marshal error:", err.Error())
		c.JSON(400, obj{"err": "marshal error"})
		return
	}
	data, err := MakeReqToFlowers("addFlower", marshaled)
	if err != nil {
		fmt.Println("handlers.go -> addFlower() -> MakeReqToFlowers error:", err.Error())
		c.JSON(400, obj{"err": "communication error"})
		return
	}

	var answer struct {
		Err string `json:"err"`
	}
	if err := json.Unmarshal(data, &answer); err != nil {
		fmt.Println("handlers.go -> addFlower() -> unmarshal error:", err.Error())
		c.JSON(400, obj{"err": "communication error"})
		return
	}
	c.JSON(200, obj{"err": nil})
}

func flowerReq(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> flowerReq() -> binding error:", err.Error())
		c.JSON(400, obj{"err": "binding error"})
		return
	}

	canGrow, err := canGrowFlower(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> canGrowFlower() error:", err.Error())
		c.JSON(400, obj{"err": "cant grow flower"})
		return
	}

	if !canGrow {
		c.JSON(400, obj{"err": "cant grow flower"})
		return
	}

	marshaledReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> marshal error:", err.Error())
		c.JSON(400, obj{"err": "marshal error"})
		return
	}
	data, err := MakeReqToFlowers("growFlower", marshaledReq)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> req error:", err.Error())
		c.JSON(400, obj{"err": "err req to flowers"})
		return
	}
	var answer Flower
	if err := json.Unmarshal(data, &answer); err != nil {
		fmt.Println("handlers.go -> flowerReq() -> unmarshal error:", err.Error())
		c.JSON(400, obj{"err": "communication error"})
		return
	}
	var resp struct {
		Flower
		Up   uint8 `json:"up"`
		Grew bool  `json:"grew"`
	}
	resp.Flower = answer
	resp.Up = answer.Grew
	resp.Grew = true
	c.JSON(200, resp)
}

func dialogFlowReq(c *gin.Context) {
	var req struct {
		Text string `json:"text"`
		ID   int    `json:"id"`
	}

	if err := c.Bind(&req); err != nil {
		fmt.Println("dialogFlowReq() -> c.Bind() error", err.Error())
		c.JSON(400, obj{"err": "binding error"})
		return
	}

	if req.Text == "" || req.ID == 0 {
		c.JSON(400, obj{"err": "fill all the fields"})
		return
	}

	answer, err := MakeReqToDialogFlow(req.Text)
	if err != nil {
		fmt.Println("dialogFlowReq() -> MakeReqToDialogFlow() error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	c.JSON(200, obj{"answer": answer})
}
