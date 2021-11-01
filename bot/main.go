package main

import (
	admin_handlers2 "github.com/supperdoggy/superSecretDevelopement/bot/internal/admin_handlers"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/config"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/db"
	handlers2 "github.com/supperdoggy/superSecretDevelopement/bot/internal/handlers"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/service"
	Cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/bot"
	"go.uber.org/zap"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

var (
	bot *telebot.Bot
	err error
)

func initUserHandlers(bot *telebot.Bot, handlers *handlers2.Handlers) {
	// handlers
	bot.Handle(Cfg.StartCommand, handlers.Start)
	bot.Handle(Cfg.FortuneCommand, handlers.FortuneCookie)
	bot.Handle(Cfg.AnekCommand, handlers.Anek)
	bot.Handle(Cfg.TostCommand, handlers.Tost)
	bot.Handle(Cfg.FlowerCommand, handlers.Flower)
	bot.Handle(Cfg.MyFlowersCommand, handlers.MyFlowers)
	bot.Handle(Cfg.GiveLastFlowerCommand, handlers.GiveLastFlower)
	bot.Handle(Cfg.GiveFlowerCommand, handlers.GiveFlower)
	bot.Handle(Cfg.FlowerTopCommand, handlers.Flowertop)
	bot.Handle(Cfg.DanetCommand, handlers.Danet)
	bot.Handle(Cfg.NHIECommand, handlers.Neverhaveiever)
	bot.Handle(Cfg.Den4ikGameCommand, handlers.Den4ikGame)
	bot.Handle(Cfg.Den4ikGameReset, handlers.ResetDen4ik)
	// menu with user flowers
	bot.Handle(telebot.OnQuery, handlers.InlineHandler)
	// handlers text messages
	bot.Handle(telebot.OnText, handlers.OnTextHandler)
}

func initAdminHandlers(bot *telebot.Bot, adminHandlers *admin_handlers2.AdminHandlers) {
	// admin handlers
	bot.Handle(Cfg.AdminHelpCommand, adminHandlers.AdminHelp)
	bot.Handle(Cfg.AddFlowerCommand, adminHandlers.AddFlower)
	bot.Handle(Cfg.AdminCommand, adminHandlers.Admin)
	bot.Handle(Cfg.AllFlowersCommand, adminHandlers.AllFlowers)
	bot.Handle(Cfg.RemoveFlower, adminHandlers.RemoveFlower)
}

func main() {
	logger, _ := zap.NewDevelopment()
	conf := config.GetConfig(logger)
	timeout := time.Second
	bot, err = telebot.NewBot(telebot.Settings{
		Token:  conf.Token,
		Poller: &telebot.LongPoller{Timeout: timeout},
	})
	if err != nil {
		logger.Fatal("error creating bot", zap.Error(err))
	}

	logger.Info("bot created!",
		zap.Any("timeout", timeout),
		zap.Bool("error_notify", conf.ErrorAdminNotification),
		zap.Bool("is_prod", conf.IsProd))

	DB := db.NewDbStruct(logger, "", Cfg.DBName, Cfg.PicCollectionName)
	Service := *service.NewService(logger, DB)
	handlers := handlers2.NewHandlers(bot, Service, logger)
	adminHandlers := admin_handlers2.NewAdminHandlers(bot, logger)

	initUserHandlers(bot, handlers)
	initAdminHandlers(bot, adminHandlers)

	bot.Start()
}
