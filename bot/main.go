package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

type obj map[string]interface{}

var (
	bot *telebot.Bot
	err error
)

func init() {
	timeout := time.Second
	bot, err = telebot.NewBot(telebot.Settings{
		// todo maybe pass token as env variable?
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: timeout},
	})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Bot created!")
	fmt.Println("Timeout:", timeout)
}

func main() {
	// handlers
	bot.Handle("/start", start)
	bot.Handle("/fortune", fortuneCookie)
	bot.Handle("/anek", anek)
	bot.Handle("/tost", tost)
	bot.Handle("/flower", flower)
	// just text handler
	bot.Handle(telebot.OnText, onTextHandler)
	bot.Handle("/myflowers", myflowers)
	bot.Handle("/giveoneflower", giveOneFlower)
	bot.Handle("/testMessage", testMessage)
	bot.Handle("/flowertop", flowertop)
	bot.Handle("/danet", danet)
	bot.Handle("/neverhaveiever", neverhaveiever)

	// admin handlers
	bot.Handle("/adminHelp", adminHelp)
	bot.Handle("/addFlower", addFlower)
	bot.Handle("/admin", admin)
	bot.Handle("/allFlowers", allFlowers)
	bot.Handle("/removeFlower", removeFlower)
	bot.Handle("/danet", danet)

	log.Println("Bot is running...")
	bot.Start()
}
