package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

// MakeRandomAnekHttpReq - sends http req to anek server and unmarshals it to RandomAnekAnswer struct
func MakeRandomAnekHttpReq(id int) (response RandomAnekAnswer, err error) {
	req := struct {
		ID int `json:"id" bson:"id"`
	}{ID: id}
	data, err := json.Marshal(req)
	if err != nil {
		return response, err
	}
	resp, err := MakeHttpReq(userUrl+"/getRandomAnek", "POST", data)
	if err != nil {
		fmt.Println("handlers.go -> MakeRandomAnekHttpReq() -> MakeHttpReq ->", err.Error())
		return
	}

	if err = json.Unmarshal(resp, &response); err != nil {
		fmt.Println("communication.go -> MakeRandomAnekHttpReq() -> error ->", err.Error())
		return
	}

	return
}

func UpdateUser(usermsg, botmsg *telebot.Message) {
	var user User = User{
		Telebot: *usermsg.Sender,
		Chats: []Chat{{
			Telebot:    *usermsg.Chat,
			LastOnline: time.Now().Unix(),
		}},
		MessagesUserSent:   []telebot.Message{*usermsg},
		MessagesZhannaSent: []telebot.Message{*botmsg},
	}
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("communication -> UpdateUser() -> marshal error:", err.Error())
		return
	}
	respStruct := obj{}
	resp, err := MakeHttpReq(userUrl+"/addOrUpdateUser", "POST", data)
	if err != nil {
		fmt.Println("communication -> UpdateUser() -> req error:", err.Error())
		return
	}
	err = json.Unmarshal(resp, &respStruct)
	if err != nil || respStruct != nil {
		fmt.Println("communication -> UpdateUser() -> unmarshal error:", err, respStruct)
		return
	}
}

// MakeHttpReq - func for sending http req with given path, method(get or post!) and data
func MakeHttpReq(path, method string, data []byte) (answer []byte, err error) {
	var resp *http.Response
	switch method {
	case "GET":
		resp, err = http.Get(path)
	case "POST":
		resp, err = http.Post(path, "application/json", bytes.NewReader(data))
	default:
		err = fmt.Errorf("method not supported, use get or post methods")
	}
	if err != nil {
		return nil, err
	}

	answer, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return
}
