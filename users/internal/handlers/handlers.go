package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	aneksdata "github.com/supperdoggy/superSecretDevelopement/structs/request/aneks"
	flowersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/flowers"
	fortunedata "github.com/supperdoggy/superSecretDevelopement/structs/request/fortune"
	NHIEdata "github.com/supperdoggy/superSecretDevelopement/structs/request/nhie"
	tostdata "github.com/supperdoggy/superSecretDevelopement/structs/request/tost"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	nhiecfg "github.com/supperdoggy/superSecretDevelopement/structs/services/NHIE"
	anekscfg "github.com/supperdoggy/superSecretDevelopement/structs/services/aneks"
	flowercfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	fortunecfg "github.com/supperdoggy/superSecretDevelopement/structs/services/fortune"
	tostcfg "github.com/supperdoggy/superSecretDevelopement/structs/services/tost"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"

	"github.com/supperdoggy/superSecretDevelopement/users/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/db"
	"log"
	"os/exec"
	"sync"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	DB *db.DbStruct
}

type obj map[string]interface{}

// todo simplify
func (h *Handlers) AddOrUpdateUser(c *gin.Context) {
	// todo tink of something new
	var userReq structs.User
	var resp usersdata.AddOrUpdateUserResp
	if err := c.Bind(&userReq); err != nil {
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> binding error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}
	old, err := h.DB.GetUserByID(userReq.Telebot.ID)
	// if we get mongo error
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> userExists() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
		// if we dont have this user in db
	} else if err == mgo.ErrNotFound {
		// inserting
		err = h.DB.UsersCollection.Insert(userReq)
		if err != nil {
			fmt.Println("handlers.go -> addOrUpdateUserReq() -> insert error:", err.Error())
			resp.Err = err.Error()
			c.JSON(400, resp)
			return
		}
		resp.OK = true
		c.JSON(http.StatusOK, resp)
		return
	}

	userReq.LastOnlineTime = time.Now()
	userReq.LastOnline = userReq.LastOnlineTime.Unix()

	// checks if chat already in user struct
	var in bool
	for _, v := range old.Chats {
		if v.Telebot.ID == userReq.Chats[0].Telebot.ID {
			in = true
		}
	}
	if !in {
		userReq.Chats = append(old.Chats, userReq.Chats...)
	}

	// add message
	if len(userReq.MessagesUserSent) != 0 {
		err = h.DB.WriteMessage(userReq.MessagesUserSent[0], userReq.MessagesZhannaSent[0])
		if err != nil {
			log.Println("handlers.go -> addOrUpdateUserReq() -> DB.writeMessage() error:", err.Error())
		}
	}

	fieldsToSet := obj{
		"lastOnlineTime":     userReq.LastOnlineTime,
		"lastOnline":         userReq.LastOnline,
		"chats":              userReq.Chats,
		"telebot.username":   userReq.Telebot.Username,
		"telebot.first_name": userReq.Telebot.FirstName,
		"telebot.last_name":  userReq.Telebot.LastName,
	}

	if err := h.DB.UsersCollection.Update(obj{"telebot.id": userReq.Telebot.ID}, obj{"$set": fieldsToSet}); err != nil {
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> update error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}
	resp.OK = true
	c.JSON(200, resp)
}

// todo simplify
func (h *Handlers) GetFortune(c *gin.Context) {
	var req usersdata.GetFortuneReq
	var resp usersdata.GetFortuneResp

	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getFortune() -> binding error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	// checking if user exists if not then just create one
	exists, err := h.DB.UserExists(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getFortune() -> userExists() error:", err.Error())
		resp.Err = "error getting user"
		c.JSON(400, resp)
		return
	}
	if !exists {
		err := h.DB.UsersCollection.Insert(structs.User{Telebot: telebot.User{ID: req.ID}})
		if err != nil {
			fmt.Println("handlers.go -> addOrUpdateUserReq() -> insert error:", err.Error())
			resp.Err = err.Error()
			c.JSON(400, resp)
			return
		}
	}

	u, err := h.DB.GetUserByID(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getFortune() -> cant find user:", err.Error())
		resp.Err = "cant find user"
		c.JSON(400, resp)
		return
	}
	// check if day passed to get new fortune
	if !CanGetFortune(u.LastTimeGotFortuneCookieTime) {
		resp.Err = "Попробуй завтра!"
		// getting last fortune
		resp.Fortune = u.FortuneCookies[len(u.FortuneCookies)-1]
		c.JSON(400, resp)
		return
	}
	var respFromFortune fortunedata.GetRandomFortuneCookieResp

	data, err := communication.MakeHttpReq(cfg.FortuneCookieURL+fortunecfg.GetRandomFortuneCookieURL, "GET", nil)
	if err != nil {
		fmt.Println("error making req:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}
	if err := json.Unmarshal(data, &respFromFortune); err != nil {
		fmt.Println("fortune cookie unmarshal error")
		resp.Err = "unmarshal error"
		c.JSON(400, resp)
		return
	}
	if respFromFortune.Err != "" {
		resp.Err = respFromFortune.Err
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if err := h.DB.UpdateLastTimeFortune(req.ID); err != nil {
		fmt.Println("error updating last time fortune:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}
	resp.Fortune = structs.Cookie{
		ID:   respFromFortune.ID,
		Text: respFromFortune.Text,
	}

	c.JSON(200, resp)
	if ok := h.DB.SaveFortune(req.ID, resp.Fortune); !ok {
		fmt.Println("Failed to save fortune for user", req.ID)
	}
}

// thats ok i guess

func (h *Handlers) GetRandomAnek(c *gin.Context) {
	var req usersdata.GetRandomAnekReq
	var resp usersdata.GetRandomAnekResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("binding error -> getRandomAnek():", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}
	if req.ID == 0 {
		resp.Err = "id cannot be 0"
		c.JSON(400, resp)
		return
	}
	data, err := communication.MakeReqToAnek(anekscfg.GetRandomAnekURL, nil)
	if err != nil {
		fmt.Println("handlers.go -> getRandomAnek()-> req error", err.Error())
		resp.Err = "something went wrong, contact @supperdoggy"
		c.JSON(400, resp)
		return
	}
	var respFromAneks aneksdata.GetRandomAnekResp
	if err = json.Unmarshal(data, &respFromAneks); err != nil {
		fmt.Println("handlers.go -> getRandomAnek() -> unmarshal error:", err.Error())
		resp.Err = "something went wrong, contact @supperdoggy"
		c.JSON(400, resp)
		return
	}
	if respFromAneks.Err != "" {
		resp.Err = respFromAneks.Err
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Id = respFromAneks.ID
	resp.Text = respFromAneks.Text
	c.JSON(200, resp)
	if ok := h.DB.SaveAnek(req.ID, resp.Anek); !ok {
		fmt.Println("Not ok saving anek", req.ID)
	}
}

func (h *Handlers) GetRandomTost(c *gin.Context) {
	var req usersdata.GetRandomTostReq
	var resp usersdata.GetRandomTostResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getRandomTost() -> c.Bind() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	if req.ID == 0 {
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	data, err := communication.MakeReqToTost(tostcfg.GetRandomTostURL, nil)
	if err != nil {
		fmt.Println("handlers.go -> getRandomTost() -> MakeReqToTost(\"getRandomTost\") error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	var respFromTost tostdata.GetRandomTostResp
	if err = json.Unmarshal(data, &respFromTost); err != nil {
		fmt.Println("handlers.go -> getRandomTost() -> json.Unmarshal error:", err.Error())
		resp.Err = "unmarshal error"
		c.JSON(400, resp)
		return
	}
	if respFromTost.Err != "" {
		resp.Err = respFromTost.Err
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.ID = respFromTost.ID
	resp.Text = respFromTost.Text

	c.JSON(200, resp)
	if ok := h.DB.SaveTost(req.ID, resp.Tost); !ok {
		fmt.Println("not ok saving tost", req.ID)
	}
}

func (h *Handlers) AddFlower(c *gin.Context) {
	var req usersdata.AddFlowerReq
	var resp usersdata.AddFlowerResp
	if err := c.Bind(&req); err != nil {
		d, err := ioutil.ReadAll(c.Request.Body)
		fmt.Println("handlers.go -> addFlower() -> binding error:", err.Error(), string(d))
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	var reqToFlowers flowersdata.AddNewFlowerReq
	var respFromFlowers flowersdata.AddNewFlowerResp
	reqToFlowers.Name = req.Name
	reqToFlowers.Icon = req.Icon
	reqToFlowers.Type = req.Type
	data, err := communication.MakeReqToFlowers(flowercfg.AddNewFlowerURL, reqToFlowers)
	if err != nil {
		fmt.Println("handlers.go -> addFlower() -> MakeReqToFlowers error:", err.Error())
		resp.Err = "communication error"
		c.JSON(400, resp)
		return
	}

	if err := json.Unmarshal(data, &respFromFlowers); err != nil {
		fmt.Println("handlers.go -> addFlower() -> unmarshal error:", err.Error())
		resp.Err = "communication error"
		c.JSON(400, resp)
		return
	}
	if !respFromFlowers.OK {
		resp.Err = respFromFlowers.Err
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.OK = true
	c.JSON(200, resp)
}

func (h *Handlers) Flower(c *gin.Context) {
	var req usersdata.FlowerReq
	var resp usersdata.FlowerResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> flowerReq() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	canGrow, err := canGrowFlower(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> canGrowFlower() error:", err.Error())
		resp.Err = "cant grow flower"
		c.JSON(400, resp)
		return
	}

	if !canGrow {
		resp.Err = "cant grow flower"
		c.JSON(400, resp)
		return
	}

	req.MsgCount, err = h.DB.GetUserMsgCount(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> getUserMsgCount error:", err.Error())
	}
	var reqToFlower flowersdata.GrowFlowerReq
	var respFromFlower flowersdata.GrowFlowerResp
	reqToFlower.ID = req.ID
	reqToFlower.NonDying = req.NonDying
	reqToFlower.MsgCount = req.MsgCount
	data, err := communication.MakeReqToFlowers(flowercfg.GrowFlowerURL, req)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> req error:", err.Error())
		resp.Err = "err req to flowers"
		c.JSON(400, resp)
		return
	}
	if err := json.Unmarshal(data, &respFromFlower); err != nil {
		fmt.Println("handlers.go -> flowerReq() -> unmarshal error:", err.Error())
		resp.Err = "communication error"
		c.JSON(400, resp)
		return
	}
	resp.Flower = respFromFlower.Flower
	resp.Up = respFromFlower.Flower.Grew
	// grew successful
	resp.Grew = true
	resp.Extra = respFromFlower.Extra
	c.JSON(200, resp)
}

func (h *Handlers) DialogFlow(c *gin.Context) {
	var req usersdata.DialogFlowReq
	var resp usersdata.DialogFlowResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("dialogFlowReq() -> c.Bind() error", err.Error())
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	if req.Text == "" || req.ID == 0 {
		resp.Err = "fill all the fields"
		c.JSON(400, resp)
		return
	}

	answer, err := communication.MakeReqToDialogFlow(req.Text)
	if err != nil {
		// if we dont get proper answer from python service we just restart it)
		// shitcode but it works for now
		// p.s sorry programmer reading this
		fmt.Println("dialogFlowReq() -> MakeReqToDialogFlow() error:", err.Error())
		fmt.Println("dialogFlowReq() -> starting python service again....")
		// restarts service)))
		go exec.Command("python3", "/root/dialogflow/main.py").Run()
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}
	resp.Answer = answer
	c.JSON(200, resp)
}

func (h *Handlers) MyFlowers(c *gin.Context) {
	var req usersdata.MyFlowersReq
	var resp usersdata.MyFlowersResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("myflowers() -> c.Bind() error", err.Error())
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}
	if req.ID == 0 {
		fmt.Println("myflowers() -> id is 0")
		resp.Err = "no id field"
		c.JSON(400, resp)
		return
	}

	var reqToFlower flowersdata.GetUserFlowersReq
	var respFromFlower flowersdata.GetUserFlowersResp
	reqToFlower.ID = req.ID
	answer, err := communication.MakeReqToFlowers(flowercfg.GetUserFlowersURL, reqToFlower)
	if err != nil {
		fmt.Println("myflowers() -> MakeHttpReq(getUserFlowers) error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}
	if err := json.Unmarshal(answer, &respFromFlower); err != nil {
		fmt.Println("myflowers() -> unmarshal error:", err.Error(), string(answer))
		resp.Err = "unmarshal error"
		c.JSON(400, resp)
		return
	}
	if resp.Err != "" {
		fmt.Println("myflowers() -> response error:", resp.Err)
		resp.Err = respFromFlower.Err
		c.JSON(400, resp)
		return
	}
	resp.Flowers = respFromFlower.Flowers
	resp.Last = respFromFlower.Last
	resp.Total = respFromFlower.Total
	c.JSON(200, resp)
}

// path for giving flower
func (h *Handlers) GiveFlower(c *gin.Context) {
	var req usersdata.GiveFlowerReq
	var resp usersdata.GiveFlowerResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> give() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}
	if req.Owner == 0 || req.Reciever == 0 || !req.Last && req.ID == 0 {
		resp.Err = "fill all fields"
		c.JSON(400, resp)
		return
	}

	var reqToFlowers flowersdata.GiveFlowerReq
	var respFromFlowers flowersdata.GiveFlowerResp
	reqToFlowers.ID = req.ID
	reqToFlowers.Owner = req.Owner
	reqToFlowers.Reciever = req.Reciever
	answer, err := communication.MakeReqToFlowers(flowercfg.GiveFlowerURL, req)
	if err != nil {
		fmt.Println("handlers.go -> give() -> MakeReqToFlowers error:", err.Error())
		resp.Err = "err making req"
		c.JSON(400, resp)
		return
	}
	if err := json.Unmarshal(answer, &resp); err != nil || respFromFlowers.Err != "" {
		fmt.Println("handlers.go -> give() -> Unmarshal error:", err, string(answer))
		resp.Err = "flower error"
		c.JSON(400, resp)
		return
	}
	resp.OK = true
	c.JSON(200, resp)
}

// Flowertop - finds all users in chat and forms top users by total flowers
// TODO Simplify
// TODO dude rewrite this for the love of god
func (h *Handlers) Flowertop(c *gin.Context) {
	var req usersdata.FlowertopReq
	var resp usersdata.FlowertopResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("flowertop() -> c.Bind() error", err.Error())
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	if req.ChatId == 0 {
		fmt.Println("flowertop() -> ChatId is 0")
		resp.Err = "no id field"
		c.JSON(400, resp)
		return
	}
	// getting chat users
	users, err := h.DB.GetChatUsers(req.ChatId)
	fmt.Println(len(users))
	if err != nil {
		fmt.Println("flowertop() -> getChatUsers() error:", err.Error(), req.ChatId)
		resp.Err = "error getting users from chat"
		c.JSON(400, resp)
		return
	}
	if len(users) == 0 {
		resp.Err = "no users in chat"
		c.JSON(400, resp)
		return
	}

	// creating map of users and slice of ids
	m := struct { // map
		m   map[int]structs.User
		mut sync.Mutex
	}{m: map[int]structs.User{}, mut: sync.Mutex{}}
	ids := []int{} // ids

	m.mut.Lock()
	for _, v := range users {
		m.m[v.Telebot.ID] = v
		ids = append(ids, v.Telebot.ID)
	}
	m.mut.Unlock()

	var reqToFlowers flowersdata.UserFlowerSliceReq
	var respFromFlowers flowersdata.UserFlowerSliceResp
	reqToFlowers.ID = ids
	answer, err := communication.MakeReqToFlowers(flowercfg.UserFlowerSliceURL, reqToFlowers)
	if err != nil {
		fmt.Println("flowertop() -> MakeReqToFlowers(\"userFlowerSlice\") error:", err.Error())
		resp.Err = "error making req"
		c.JSON(400, req)
		return
	}
	if err := json.Unmarshal(answer, &respFromFlowers); err != nil {
		fmt.Println("flowertop() -> unmarshal error:", err.Error(), string(answer))
		resp.Err = "unmarshal error"
		c.JSON(400, resp)
		return
	}

	// so fucking bad
	// i really dont have any idea how this works i am not joking
	m.mut.Lock()
	for i := range respFromFlowers.Result {
		if user, ok := m.m[respFromFlowers.Result[i].Key]; ok {
			data := struct {
				Username string `json:"username"`
				Total    int    `json:"total"`
			}{Username: user.Telebot.Username, Total: respFromFlowers.Result[i].Value}
			if data.Username == "" {
				data.Username = fmt.Sprintf("%v %v", user.Telebot.FirstName, user.Telebot.LastName)
			}

			resp.Result = append(resp.Result, data)
		}
	}
	m.mut.Unlock()
	c.JSON(200, resp)

}

func (h *Handlers) GetRandomNHIE(c *gin.Context) {
	var req usersdata.GetRandomNHIEreq
	var resp usersdata.GetRandomNHIEresp
	if err := c.Bind(&req); err != nil {
		log.Println("handlers.go -> getRandomNHIE() -> c.Bind() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	var respFromNHIE NHIEdata.GetRandomNHIEResponse
	data, err := communication.MakeHttpReq(cfg.NHIE_URL+nhiecfg.GetRandomNeverHaveIEverURL, "GET", nil)
	if err != nil {
		log.Println("handlers.go -> getRandomNHIE() -> c.Bind() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	err = json.Unmarshal(data, &respFromNHIE)
	if err != nil {
		log.Printf("handlers.go -> getRandomNHIE() -> Unmarshal error:%v, body:%v\n", err.Error(), string(data))
		resp.Err = "unmarshal error"
		c.JSON(400, resp)
		return
	}
	resp.Result.Text = respFromNHIE.Text
	resp.Result.ID = respFromNHIE.ID

	c.JSON(200, resp)
}

func CanGetFortune(date time.Time) bool {
	now := time.Now()
	return date.Day() != now.Day() || date.Month() != now.Month() || date.Year() != now.Year()
}

func canGrowFlower(id int) (bool, error) {
	var reqToFlowers flowersdata.CanGrowFlowerReq
	var respFromFlowers flowersdata.CanGrowFlowerResp
	reqToFlowers.ID = id
	answer, err := communication.MakeReqToFlowers(flowercfg.CanGrowFlowerURL, reqToFlowers)
	if err != nil {
		fmt.Println("canGrowFlower() -> MakeReqToFlower(canGrowFlower) error:", err.Error())
		return false, err
	}

	if err := json.Unmarshal(answer, &respFromFlowers); err != nil {
		fmt.Println("canGrowFlower() -> Unmarshal error:", err.Error(), string(answer))
		return false, err
	}

	if respFromFlowers.Err != "" {
		fmt.Println("canGrowFlower() -> got error from flower:", respFromFlowers.Err)
		return false, fmt.Errorf(respFromFlowers.Err)
	}
	return respFromFlowers.Answer, nil

}
