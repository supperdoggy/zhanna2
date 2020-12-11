package main

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

// start() - handles /start command and sends text response
// todo below
func start(m *telebot.Message) {
	var response string
	// todo: create id checker and answer variations for different users
	response = "Привет, я пока что очень сырая, будь нежен со мной..."
	if _, err := bot.Send(m.Sender, response); err != nil {
		fmt.Println("handlers.go -> start() -> error:", err.Error(), ", user id:", m.Sender.ID)
		return
	}
}

// TODO:
func fortuneCookie(m *telebot.Message) {

}

// anek() - handles /anek command and sends anek text response
func anek(m *telebot.Message) {
	anekAnswer, err := MakeRandomAnekHttpReq()
	if err != nil{
		fmt.Println("handlers.go -> anek() -> make req error:", err.Error())
		return
	}
	if _, err := bot.Reply(m, anekAnswer.Text);err != nil{
		fmt.Println("handlers.go -> anek() -> reply error:", err.Error())
		return
	}

}
