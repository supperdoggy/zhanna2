package main

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

func initBot() *telebot.Bot {
	timeout := time.Second
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: timeout},
	})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Bot created!")
	fmt.Println("Timeout:", timeout)
	return bot
}
