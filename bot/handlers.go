package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/tucnak/telebot.v2"
)

// start() - handles /start command and sends text response
// todo below
func start(m *telebot.Message) {
	var response string
	// todo: create id checker and answer variations for different users
	response = "Привет, я пока что очень сырая, будь нежен со мной..."
	botmsg, err := bot.Send(m.Sender, response)
	if err != nil {
		fmt.Println("handlers.go -> start() -> error:", err.Error(), ", user id:", m.Sender.ID)
		return
	}
	go UpdateUser(m, botmsg)
}

// TODO:
func fortuneCookie(m *telebot.Message) {
	var resp struct {
		FortuneCookie
		Err string `json:"err"`
	}
	data := obj{"id": m.Sender.ID}
	readyData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error unmarshalling")
		return
	}
	r, err := MakeHttpReq(userUrl+"/getFortune", "POST", readyData)
	if err != nil {
		fmt.Println("error making request")
		return
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		fmt.Println("error unmarshalling")
		return
	}
	msg := resp.Text
	if resp.Err != "" {
		msg = resp.Err
	}

	botmsg, err := bot.Reply(m, msg)
	if err != nil {
		fmt.Println("error sending answer, FortuneCookie:", err.Error())
		return
	}
	go UpdateUser(m, botmsg)
}

// anek() - handles /anek command and sends anek text response
func anek(m *telebot.Message) {
	anekAnswer, err := MakeRandomAnekHttpReq(m.Sender.ID)
	if err != nil {
		fmt.Println("handlers.go -> anek() -> make req error:", err.Error())
		return
	}
	botmsg, err := bot.Reply(m, anekAnswer.Text)
	if err != nil {
		fmt.Println("handlers.go -> anek() -> reply error:", err.Error())
		return
	}
	go UpdateUser(m, botmsg)
}
