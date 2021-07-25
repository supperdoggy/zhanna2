package service

import (
	"encoding/json"
	"errors"
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
	"github.com/supperdoggy/superSecretDevelopement/users/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"sync"
	"time"
)

type obj map[string]interface{}

type Service struct {
	DB db.DbStruct
}

// todo simplify
func (s *Service) AddOrUpdateUser(req structs.User) (resp usersdata.AddOrUpdateUserResp, err error) {
	if req.Telebot.ID == 0 || req.Telebot.FirstName == ""{
		resp.Err = "fill all fields"
		return resp, errors.New(resp.Err)
	}
	old, err := s.DB.GetUserByID(req.Telebot.ID)
	// if we get mongo error
	if err != nil && err != mgo.ErrNotFound {
		resp.Err = err.Error()
		return
		// if we dont s.ve this user in db
	} else if err == mgo.ErrNotFound {
		// inserting
		err = s.DB.UsersCollection.Insert(req)
		if err != nil {
			resp.Err = err.Error()
			return
		}
		resp.OK = true
		return
	}

	req.LastOnlineTime = time.Now()
	req.LastOnline = req.LastOnlineTime.Unix()

	// checks if chat already in user struct
	var in bool
	for _, v := range old.Chats {
		if v.Telebot.ID == req.Chats[0].Telebot.ID {
			in = true
		}
	}
	if !in {
		req.Chats = append(old.Chats, req.Chats...)
	}

	// add message
	if len(req.MessagesUserSent) != 0 {
		err = s.DB.WriteMessage(req.MessagesUserSent[0], req.MessagesZhannaSent[0])
		if err != nil {
			log.Println("handlers.go -> addOrUpdateUserReq() -> DB.writeMessage() error:", err.Error())
		}
	}

	fieldsToSet := obj{
		"lastOnlineTime":     req.LastOnlineTime,
		"lastOnline":         req.LastOnline,
		"chats":              req.Chats,
		"telebot.username":   req.Telebot.Username,
		"telebot.first_name": req.Telebot.FirstName,
		"telebot.last_name":  req.Telebot.LastName,
	}

	if err := s.DB.UsersCollection.Update(obj{"telebot.id": req.Telebot.ID}, obj{"$set": fieldsToSet}); err != nil {
		resp.Err = err.Error()
		return resp, err
	}
	resp.OK = true
	return
}

func (s *Service) GetFortune(req usersdata.GetFortuneReq) (resp usersdata.GetFortuneResp, err error) {
	// checking if user exists if not then just create one
	exists, err := s.DB.UserExists(req.ID)
	if err != nil {
		resp.Err = "error getting user"
		return
	}
	if !exists {
		err = s.DB.UsersCollection.Insert(structs.User{Telebot: telebot.User{ID: req.ID}})
		if err != nil {
			fmt.Println("handlers.go -> addOrUpdateUserReq() -> insert error:", err.Error())
			resp.Err = err.Error()
			return
		}
	}

	u, err := s.DB.GetUserByID(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getFortune() -> cant find user:", err.Error())
		resp.Err = "cant find user"
		return
	}
	// check if day passed to get new fortune
	if !s.CanGetFortune(u.LastTimeGotFortuneCookieTime) {
		resp.Err = "Попробуй завтра!"
		// getting last fortune
		resp.Fortune = u.FortuneCookies[len(u.FortuneCookies)-1]
		return
	}
	var respFromFortune fortunedata.GetRandomFortuneCookieResp

	data, err := communication.MakeHttpReq(cfg.FortuneCookieURL+fortunecfg.GetRandomFortuneCookieURL, "GET", nil)
	if err != nil {
		fmt.Println("error making req:", err.Error())
		resp.Err = err.Error()
		return
	}
	if err = json.Unmarshal(data, &respFromFortune); err != nil {
		fmt.Println("fortune cookie unmarshal error")
		resp.Err = "unmarshal error"
		return
	}
	if respFromFortune.Err != "" {
		resp.Err = respFromFortune.Err
		return resp, errors.New(resp.Err)
	}
	if err = s.DB.UpdateLastTimeFortune(req.ID); err != nil {
		fmt.Println("error updating last time fortune:", err.Error())
		resp.Err = err.Error()
		return
	}
	resp.Fortune = structs.Cookie{
		ID:   respFromFortune.ID,
		Text: respFromFortune.Text,
	}

	// saving fotune
	if ok := s.DB.SaveFortune(req.ID, resp.Fortune); !ok {
		fmt.Println("Failed to save fortune for user", req.ID)
	}
	return
}

func (s *Service) GetRandomAnek(req usersdata.GetRandomAnekReq) (resp usersdata.GetRandomAnekResp, err error) {
	if req.ID == 0 {
		resp.Err = "id cannot be 0"
		return resp, errors.New(resp.Err)
	}

	var respFromAneks aneksdata.GetRandomAnekResp
	err = communication.MakeReqToAnek(anekscfg.GetRandomAnekURL, nil, &respFromAneks)
	if err != nil {
		resp.Err = "something went wrong, contact @supperdoggy"
		return
	}
	if respFromAneks.Err != "" {
		resp.Err = respFromAneks.Err
		return resp, errors.New(resp.Err)
	}
	resp.Id = respFromAneks.ID
	resp.Text = respFromAneks.Text
	// saving anek
	if ok := s.DB.SaveAnek(req.ID, resp.Anek); !ok {
		fmt.Println("Not ok saving anek", req.ID)
	}
	return
}

func (s *Service) GetRandomTost(req usersdata.GetRandomTostReq) (resp usersdata.GetRandomTostResp, err error) {
	if req.ID == 0 {
		resp.Err = "binding error"
		return resp, errors.New(resp.Err)
	}

	var respFromTost tostdata.GetRandomTostResp
	err = communication.MakeReqToTost(tostcfg.GetRandomTostURL, nil, &respFromTost)
	if err != nil {
		resp.Err = err.Error()
		return
	}

	if respFromTost.Err != "" {
		resp.Err = respFromTost.Err
		return
	}
	resp.ID = respFromTost.ID
	resp.Text = respFromTost.Text

	// saving tost
	if ok := s.DB.SaveTost(req.ID, resp.Tost); !ok {
		fmt.Println("not ok saving tost", req.ID)
	}

	return
}

func (s *Service) AddFlower(req usersdata.AddFlowerReq) (resp usersdata.AddFlowerResp, err error) {
	var reqToFlowers flowersdata.AddNewFlowerReq
	var respFromFlowers flowersdata.AddNewFlowerResp
	reqToFlowers.Name = req.Name
	reqToFlowers.Icon = req.Icon
	reqToFlowers.Type = req.Type
	err = communication.MakeReqToFlowers(flowercfg.AddNewFlowerURL, reqToFlowers, &respFromFlowers)
	if err != nil {
		resp.Err = "communication error"
		return
	}

	if !respFromFlowers.OK {
		resp.Err = respFromFlowers.Err
		return resp, errors.New(resp.Err)
	}
	resp.OK = true
	return
}

func (s *Service) Flower(req usersdata.FlowerReq) (resp usersdata.FlowerResp, err error) {
	canGrow, err := s.CanGrowFlower(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> canGrowFlower() error:", err.Error())
		resp.Err = "cant grow flower"
		return
	}

	if !canGrow {
		resp.Err = "cant grow flower"
		return resp, errors.New(resp.Err)
	}

	req.MsgCount, err = s.DB.GetUserMsgCount(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> getUserMsgCount error:", err.Error())
	}
	var reqToFlower flowersdata.GrowFlowerReq
	var respFromFlower flowersdata.GrowFlowerResp
	reqToFlower.ID = req.ID
	reqToFlower.NonDying = req.NonDying
	reqToFlower.MsgCount = req.MsgCount
	err = communication.MakeReqToFlowers(flowercfg.GrowFlowerURL, reqToFlower, &respFromFlower)
	if err != nil {
		fmt.Println("handlers.go -> flowerReq() -> req error:", err.Error())
		resp.Err = "err req to flowers"
		return
	}

	resp.Flower = respFromFlower.Flower
	resp.Up = respFromFlower.Flower.Grew
	// grew successful
	resp.Grew = true
	resp.Extra = respFromFlower.Extra
	return
}

func (s *Service) DialogFlow(req usersdata.DialogFlowReq) (resp usersdata.DialogFlowResp, err error) {
	if req.Text == "" || req.ID == "" {
		resp.Err = "fill all the fields"
		return resp, errors.New(resp.Err)
	}

	resp = communication.MakeReqToDialogFlow(req)
	if resp.Err != "" {
		return resp, errors.New(resp.Err)
	}
	return
}

func (s *Service) MyFlowers(req usersdata.MyFlowersReq) (resp usersdata.MyFlowersResp, err error) {
	if req.ID == 0 {
		fmt.Println("myflowers() -> id is 0")
		resp.Err = "no id field"
		return resp, errors.New(resp.Err)
	}

	var reqToFlower flowersdata.GetUserFlowersReq
	var respFromFlower flowersdata.GetUserFlowersResp
	reqToFlower.ID = req.ID
	err = communication.MakeReqToFlowers(flowercfg.GetUserFlowersURL, reqToFlower, &respFromFlower)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	if resp.Err != "" {
		fmt.Println("myflowers() -> response error:", resp.Err)
		resp.Err = respFromFlower.Err
		return resp, errors.New(resp.Err)
	}
	resp.Flowers = respFromFlower.Flowers
	resp.Last = respFromFlower.Last
	resp.Total = respFromFlower.Total
	return
}

func (s *Service) GiveFlower(req usersdata.GiveFlowerReq) (resp usersdata.GiveFlowerResp, err error) {
	if req.Owner == 0 || req.Reciever == 0 || !req.Last && req.ID == 0 {
		resp.Err = "fill all fields"
		return resp, errors.New(resp.Err)
	}

	var reqToFlowers flowersdata.GiveFlowerReq
	var respFromFlowers flowersdata.GiveFlowerResp
	reqToFlowers.ID = req.ID
	reqToFlowers.Owner = req.Owner
	reqToFlowers.Reciever = req.Reciever
	reqToFlowers.Last = true
	err = communication.MakeReqToFlowers(flowercfg.GiveFlowerURL, reqToFlowers, &respFromFlowers)
	if err != nil {
		fmt.Println("handlers.go -> give() -> MakeReqToFlowers error:", err.Error())
		resp.Err = "err making req"
		return
	}
	if respFromFlowers.Err != "" {
		fmt.Println("handlers.go -> give() -> Unmarshal error:", respFromFlowers.Err)
		resp.Err = respFromFlowers.Err
		return resp, errors.New(resp.Err)
	}
	resp.OK = true
	return
}

func (s *Service) Flowertop(req usersdata.FlowertopReq) (resp usersdata.FlowertopResp, err error) {
	if req.ChatId == 0 {
		fmt.Println("flowertop() -> ChatId is 0")
		resp.Err = "no id field"
		return resp, errors.New(resp.Err)
	}
	// getting chat users
	users, err := s.DB.GetChatUsers(req.ChatId)
	fmt.Println(len(users))
	if err != nil {
		fmt.Println("flowertop() -> getChatUsers() error:", err.Error(), req.ChatId)
		resp.Err = "error getting users from chat"
		return
	}
	if len(users) == 0 {
		resp.Err = "no users in chat"
		return resp, errors.New(resp.Err)
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
	err = communication.MakeReqToFlowers(flowercfg.UserFlowerSliceURL, reqToFlowers, &respFromFlowers)
	if err != nil {
		fmt.Println("flowertop() -> MakeReqToFlowers(\"userFlowerSlice\") error:", err.Error())
		resp.Err = "error making req"
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
	return
}

func (s *Service) GetRandomNHIE(req usersdata.GetRandomNHIEreq) (resp usersdata.GetRandomNHIEresp, err error) {
	var respFromNHIE NHIEdata.GetRandomNHIEResponse
	data, err := communication.MakeHttpReq(cfg.NHIE_URL+nhiecfg.GetRandomNeverHaveIEverURL, "GET", nil)
	if err != nil {
		log.Println("handlers.go -> getRandomNHIE() -> c.Bind() error:", err.Error())
		resp.Err = err.Error()
		return
	}

	err = json.Unmarshal(data, &respFromNHIE)
	if err != nil {
		log.Printf("handlers.go -> getRandomNHIE() -> Unmarshal error:%v, body:%v\n", err.Error(), string(data))
		resp.Err = "unmarshal error"
		return
	}
	resp.Result.Text = respFromNHIE.Text
	resp.Result.ID = respFromNHIE.ID
	return
}

func (s *Service) CanGetFortune(date time.Time) bool {
	now := time.Now()
	return date.Day() != now.Day() || date.Month() != now.Month() || date.Year() != now.Year()
}

func (s *Service) CanGrowFlower(id int) (bool, error) {
	var reqToFlowers flowersdata.CanGrowFlowerReq
	var respFromFlowers flowersdata.CanGrowFlowerResp
	reqToFlowers.ID = id
	err := communication.MakeReqToFlowers(flowercfg.CanGrowFlowerURL, reqToFlowers, &respFromFlowers)
	if err != nil {
		fmt.Println("canGrowFlower() -> MakeReqToFlower(canGrowFlower) error:", err.Error())
		return false, err
	}

	if respFromFlowers.Err != "" {
		fmt.Println("canGrowFlower() -> got error from flower:", respFromFlowers.Err)
		return false, fmt.Errorf(respFromFlowers.Err)
	}
	return respFromFlowers.Answer, nil

}