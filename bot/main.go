package main

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

func main(){
	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{Timeout: time.Second},
	})
	if err != nil{
		panic(err.Error())
	}
	fmt.Println("Bot created!")

	bot.Handle("/start", func(m *telebot.Message) {
		if _, err := bot.Send(m.Sender, "Hello world!");err!=nil{
			fmt.Println("failed sending msg to user", m.Chat.ID)
		}
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message){

	})

	fmt.Println("Bot running...")
	bot.Start()
}