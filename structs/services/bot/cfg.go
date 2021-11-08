package Cfg

import (
	usercfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
)

const (
	TestbotId = 1442407913
	ProdBotID = 1058731629
	NeMoksID  = 424137509
	EdemID    = -1001419983908

	// todo something cleaver make
	UsersAdminURL = usercfg.UserURL + "/admin"

	// user

	StartCommand          = "/start"
	FortuneCommand        = "/fortune"
	AnekCommand           = "/anek"
	TostCommand           = "/tost"
	FlowerCommand         = "/flower"
	MyFlowersCommand      = "/myflowers"
	GiveLastFlowerCommand = "/givelastflower"
	GiveFlowerCommand     = "/give"
	FlowerTopCommand      = "/flowertop"
	DanetCommand          = "/danet"
	NHIECommand           = "/neverhaveiever"
	Den4ikGameCommand     = "/go"
	Den4ikGameReset       = "/resetden4ik"

	// admin

	AdminHelpCommand      = "/adminHelp"
	AddFlowerCommand      = "/addFlower"
	AdminCommand          = "/admin"
	AllFlowersCommand     = "/allFlowers"
	RemoveFlower          = "/removeFlower"
	AddUserFlowerMultiple = "/addUserFlowerMultiple"
	AddUserFlowerByID     = "/addUserFlowerByID"

	// db
	PicCollectionName = "Pic"
	DBName            = "Bot"
)
