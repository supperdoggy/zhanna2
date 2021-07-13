package Cfg

import (
	usercfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
)

const (
	TestbotId = 1442407913
	ProdBotID = 1058731629
	NeMoksID  = 424137509
	EdemID    = -1001419983908

	UsersURL = usercfg.UserURL
	// todo something cleaver make
	UsersAdminURL = usercfg.UserURL + "/admin"

	// user

	StartCommand = "/start"
	FortuneCommand = "/fortune"
	AnekCommand = "/anek"
	TostCommand = "/tost"
	FlowerCommand = "/flower"
	MyFlowersCommand = "/myflowers"
	GiveFlowerCommand = "/giveoneflower"
	FlowerTopCommand = "/flowertop"
	DanetCommand = "/danet"
	NHIECommand = "/neverhaveiever"


	// admin

	AdminHelpCommand = "/adminHelp"
	AddFlowerCommand = "/addFlower"
	AdminCommand = "/admin"
	AllFlowersCommand = "/allFlowers"
	RemoveFlower = "/removeFlower"
)
