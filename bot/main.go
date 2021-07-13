package main

import (
	"fmt"
	admin_handlers2 "github.com/supperdoggy/superSecretDevelopement/bot/internal/admin_handlers"
	handlers2 "github.com/supperdoggy/superSecretDevelopement/bot/internal/handlers"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

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
	handlers := handlers2.Handlers{Bot: bot}
	admin_handlers := admin_handlers2.AdminHandlers{Bot: bot}
	// handlers
	bot.Handle(Cfg.StartCommand, handlers.Start)
	bot.Handle(Cfg.FortuneCommand, handlers.FortuneCookie)
	bot.Handle(Cfg.AnekCommand, handlers.Anek)
	bot.Handle(Cfg.TostCommand, handlers.Tost)
	bot.Handle(Cfg.FlowerCommand, handlers.Flower)
	bot.Handle(Cfg.MyFlowersCommand, handlers.MyFlowers)
	bot.Handle(Cfg.GiveFlowerCommand, handlers.GiveOneFlower)
	bot.Handle(Cfg.FlowerTopCommand, handlers.Flowertop)
	bot.Handle(Cfg.DanetCommand, handlers.Danet)
	bot.Handle(Cfg.NHIECommand, handlers.Neverhaveiever)
	bot.Handle(telebot.OnText, handlers.OnTextHandler)

	// admin handlers
	bot.Handle(Cfg.AdminHelpCommand, admin_handlers.AdminHelp)
	bot.Handle(Cfg.AddFlowerCommand, admin_handlers.AddFlower)
	bot.Handle(Cfg.AdminCommand, admin_handlers.Admin)
	bot.Handle(Cfg.AllFlowersCommand, admin_handlers.AllFlowers)
	bot.Handle(Cfg.RemoveFlower, admin_handlers.RemoveFlower)

	log.Println("Bot is running...")
	bot.Start()
}
