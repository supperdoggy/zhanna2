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

func MakeRandomTostHttpReq(id int) (response Tost, err error) {
	req := struct {
		ID int `json:"id" bson:"id"`
	}{ID: id}
	data, err := json.Marshal(req)
	if err != nil {
		return response, err
	}
	resp, err := MakeHttpReq(userUrl+"/getRandomTost", "POST", data)
	if err != nil {
		fmt.Println("comunication.go -> MakeRandomTostHttpReq() -> MakeHttpReq ->", err.Error())
		return
	}

	if err = json.Unmarshal(resp, &response); err != nil {
		fmt.Println("communication.go -> MakeRandomTostHttpReq() -> error ->", err.Error())
		return
	}

	return
}

// MakeUserHttpReq - method handler for users req
// TODO: refactor it!!!
func MakeUserHttpReq(method string, req interface{}) (answer []byte, err error) {
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	path := fmt.Sprintf("%s/%s", userUrl, method)
	switch method {
	case "addFlower":
		answer, err = MakeHttpReq(path, "POST", data)
	case "getAnswer":
		answer, err = MakeHttpReq(path, "POST", data)
	case "myflowers":
		answer, err = MakeHttpReq(path, "POST", data)
	case "give":
		answer, err = MakeHttpReq(path, "POST", data)
	case "flowertop":
		answer, err = MakeHttpReq(path, "POST", data)
	default:
		err = fmt.Errorf("no such method")
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
	if err != nil {
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

// grow flower
func MakeFlowerReq(id int) (msg string, err error) {
	var data = struct {
		ID int `json:"id"`
	}{ID: id}

	marshaled, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("communication.go -> flowerReq() -> json.Marshal() error: %v user %v\n", err.Error(), data.ID)
		return "communication error", err
	}
	resp, err := MakeHttpReq(userUrl+"/flower", "POST", marshaled)
	if err != nil {
		fmt.Println("communication.go -> flowerReq() -> json.MakeHttpReq() error", err.Error())
		return "communication error", err
	}
	var answer struct {
		Flower
		Up   uint8  `json:"up"`
		Grew bool   `json:"grew"`
		Err  string `json:"err"`
	}
	if err := json.Unmarshal(resp, &answer); err != nil {
		fmt.Printf("communication.go -> flowerReq() -> json.Unmarshal() error: %v body %v\n", err.Error(), string(resp))
		return "communication error", err
	}
	if answer.Err == "cant grow flower" {
		return "Ты уже сегодня поливал цветочки!\nПопробуй позже", nil
	}

	if answer.Err != "" {
		fmt.Println("communication.go -> flowerReq() -> answer.Err != '', err:", answer.Err)
		return "communication error", err
	}
	if answer.HP == 100 {
		return fmt.Sprintf("Поздравляю! Твой %v вырос!", answer.Icon), err
	}
	if answer.Grew {
		return fmt.Sprintf("Твой цветок вырос на %v единиц, теперь его размер %v единиц!", answer.Up, answer.HP), err
	}
	return "its not time, try again later...", err
}
