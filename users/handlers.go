package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"sync"
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
	fmt.Printf("%+v\n", newUser.Telebot.ID)

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

	// add message
	err = DB.writeMessage(newUser.MessagesUserSent[0], newUser.MessagesZhannaSent[0])
	if err != nil {
		log.Println("handlers.go -> addOrUpdateUserReq() -> DB.writeMessage() error:", err.Error())
	}

	fieldsToSet := obj{
		"lastOnlineTime":     newUser.LastOnlineTime,
		"lastOnline":         newUser.LastOnline,
		"chats":              newUser.Chats,
		"telebot.username":   newUser.Telebot.Username,
		"telebot.first_name": newUser.Telebot.FirstName,
		"telebot.last_name":  newUser.Telebot.LastName,
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
		c.JSON(400, obj{"err": cantGetFortune, "fortune": u.FortuneCookies[len(u.FortuneCookies)-1]})
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

	c.JSON(200, obj{"fortune": result})
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

	data, err := MakeReqToFlowers("addFlower", req)
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
		ID       int  `json:"id"`
		NonDying bool `json:"nonDying"`
		MsgCount int  `json:"msg_count"`
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

	req.MsgCount, err = DB.getUserMsgCount(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> getUserMsgCount error:", err.Error())
	}
	data, err := MakeReqToFlowers("growFlower", req)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> req error:", err.Error())
		c.JSON(400, obj{"err": "err req to flowers"})
		return
	}
	var answer struct{
		Flower Flower `json:"flower"`
		Extra int `json:"extra"`
	}
	if err := json.Unmarshal(data, &answer); err != nil {
		fmt.Println("handlers.go -> flowerReq() -> unmarshal error:", err.Error())
		c.JSON(400, obj{"err": "communication error"})
		return
	}
	var resp struct {
		Flower
		Up   uint8 `json:"up"`
		Grew bool  `json:"grew"`
		Extra int `json:"extra"`
	}
	resp.Flower = answer.Flower
	resp.Up = answer.Flower.Grew
	// grew successful
	resp.Grew = true
	resp.Extra = answer.Extra
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
		fmt.Println("dialogFlowReq() -> starting python service again....")
		// restarts service)))
		go exec.Command("python3", "/root/dialogflow/main.py").Run()
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	c.JSON(200, obj{"answer": answer})
}

func myflowers(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}

	if err := c.Bind(&req); err != nil {
		fmt.Println("myflowers() -> c.Bind() error", err.Error())
		c.JSON(400, obj{"err": "binding error"})
		return
	}
	if req.ID == 0 {
		fmt.Println("myflowers() -> id is 0")
		c.JSON(400, obj{"err": "no id field"})
		return
	}

	answer, err := MakeReqToFlowers("getUserFlowers", req)
	if err != nil {
		fmt.Println("myflowers() -> MakeHttpReq(getUserFlowers) error:", err.Error())
		c.JSON(400, obj{"err": "req error"})
		return
	}
	var resp struct {
		Flowers map[string]int `json:"flowers"`
		Last    uint8          `json:"last"`
		Total   int            `json:"total"`
		Err     string         `json:"err"`
	}
	if err := json.Unmarshal(answer, &resp); err != nil {
		fmt.Println("myflowers() -> unmarshal error:", err.Error(), string(answer))
		c.JSON(400, obj{"err": "unmarshal error"})
		return
	}
	if resp.Err != "" {
		fmt.Println("myflowers() -> response error:", resp.Err)
		c.JSON(400, obj{"err": resp.Err})
		return
	}
	c.JSON(200, resp)
}

// path for giving flower
func give(c *gin.Context) {
	var req struct {
		Owner    int    `json:"owner"`
		Count    int    `json:"count"`
		Reciever int    `json:"reciever"`
		Last     bool   `json:"last"`
		ID       uint64 `json:"id"`
	}

	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> give() -> binding error:", err.Error())
		c.JSON(400, obj{"err": "binding error"})
		return
	}
	if req.Owner == 0 || req.Reciever == 0 || !req.Last && req.ID == 0 {
		c.JSON(400, obj{"err": "fill all the fields"})
		return
	}

	answer, err := MakeReqToFlowers("giveFlower", req)
	if err != nil {
		fmt.Println("handlers.go -> give() -> MakeReqToFlowers error:", err.Error())
		c.JSON(400, obj{"err": "err making req"})
		return
	}

	var resp struct {
		Err string `json:"err"`
	}
	if err := json.Unmarshal(answer, &resp); err != nil || resp.Err != "" {
		fmt.Println("handlers.go -> give() -> Unmarshal error:", err.Error(), string(answer))
		c.JSON(400, obj{"err": "flower error"})
		return
	}

	c.JSON(200, obj{"err": ""})
}

// flowertop - finds all users in chat and forms top users by total flowers
// todo: fimplify
func flowertop(c *gin.Context) {
	var req struct {
		ChatId int `json:"chatid"`
	}

	if err := c.Bind(&req); err != nil {
		fmt.Println("flowertop() -> c.Bind() error", err.Error())
		c.JSON(400, obj{"err": "binding error"})
		return
	}

	if req.ChatId == 0 {
		fmt.Println("flowertop() -> ChatId is 0")
		c.JSON(400, obj{"err": "no id field"})
		return
	}
	// getting chat users
	users, err := DB.getChatUsers(req.ChatId)
	fmt.Println(len(users))
	if err != nil {
		fmt.Println("flowertop() -> getChatUsers() error:", err.Error(), req.ChatId)
		c.JSON(400, obj{"err": "error getting users from chat"})
		return
	}
	if len(users) == 0 {
		c.JSON(400, obj{"err": "no users in chat"})
		return
	}

	// creating map of users and slice of ids
	m := struct { // map
		m   map[int]User
		mut sync.Mutex
	}{m: map[int]User{}, mut: sync.Mutex{}}
	ids := []int{} // ids

	m.mut.Lock()
	for _, v := range users {
		m.m[v.Telebot.ID] = v
		ids = append(ids, v.Telebot.ID)
	}
	m.mut.Unlock()

	answer, err := MakeReqToFlowers("userFlowerSlice", obj{"id": ids})
	if err != nil {
		fmt.Println("flowertop() -> MakeReqToFlowers(\"userFlowerSlice\") error:", err.Error())
		c.JSON(400, obj{"err": "err making req"})
		return
	}
	var resp struct {
		Result []struct {
			ID       int    `json:"id"`
			Total    int    `json:"total"`
			Username string `json:"username"`
		} `json:"result"`
	}
	if err := json.Unmarshal(answer, &resp); err != nil {
		fmt.Println("flowertop() -> unmarshal error:", err.Error(), string(answer))
		c.JSON(400, obj{"err": "err making req"})
		return
	}

	// so fucking bad
	result := []struct {
		Username string `json:"username"`
		Total    int    `json:"total"`
	}{}
	m.mut.Lock()
	for i := range resp.Result {
		if user, ok := m.m[resp.Result[i].ID]; ok {
			data := struct {
				Username string `json:"username"`
				Total    int    `json:"total"`
			}{Username: user.Telebot.Username, Total: resp.Result[i].Total}
			if data.Username == "" {
				data.Username = fmt.Sprintf("%v %v", user.Telebot.FirstName, user.Telebot.LastName)
			}
			log.Println(resp.Result[i].ID)

			result = append(result, data)
		}
	}
	m.mut.Unlock()
	c.JSON(200, obj{"result": result})

}

func getRandomNHIE(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		log.Println("handlers.go -> getRandomNHIE() -> c.Bind() error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}

	// u, err := DB.getUserFromDbById(req.ID)
	// if err != nil {
	// 	log.Println("handlers.go -> getRandomNHIE() -> c.Bind() error:", err.Error())
	// }

	data, err := MakeHttpReq(NHIEUrl+"/getRandomNeverHaveIEver", "GET", nil)
	if err != nil {
		log.Println("handlers.go -> getRandomNHIE() -> c.Bind() error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}

	var resp struct {
		Err    string `json:"err"`
		Result NHIE   `json:"result"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.Printf("handlers.go -> getRandomNHIE() -> Unmarshal error:%v, body:%v\n", err.Error(), string(data))
		c.JSON(400, obj{"err": "error unmarshaling"})
		return
	}

	c.JSON(200, resp)

}
