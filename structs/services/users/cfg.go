package cfg

import (
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	nhiecfg "github.com/supperdoggy/superSecretDevelopement/structs/services/NHIE"
	anekscfg "github.com/supperdoggy/superSecretDevelopement/structs/services/aneks"
	flowercfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	cookieCfg "github.com/supperdoggy/superSecretDevelopement/structs/services/fortune"
	tostcfg "github.com/supperdoggy/superSecretDevelopement/structs/services/tost"
)

const (
	DBName = "Zhanna2"
	Port   = ":1488"

	// move it to default cfg
	FortuneCookieURL = "http://localhost" + cookieCfg.Port + "/" + defaultCfg.ApiV1
	FlowersURL        = "http://localhost" + flowercfg.Port + "/" + defaultCfg.ApiV1
	AnekURL          = "http://localhost" + anekscfg.Port + "/" + defaultCfg.ApiV1
	NHIE_URL         = "http://localhost" + nhiecfg.Port + "/" + defaultCfg.ApiV1
	TostURL          = "http://localhost" + tostcfg.Port + "/" + defaultCfg.ApiV1
	UserURL          = "http://localhost" + Port + "/" + defaultCfg.ApiV1
	DialogFlowURL  = "http://localhost:5000/" + defaultCfg.ApiV1

	// user handlers

	AddOrUpdateUserURL = "/addOrUpdateUser"
	GetFortuneURL = "/getFortune"
	GetRandomAnekURL = "/getRandomAnek"
	GetRandomTostURL = "/getRandomTost"
	AddFlowerURL = "/addFlower"
	FlowerURL = "/flower"
	DialogFlowHandlerURL = "/getAnswer"
	MyFlowersURL = "/myflowers"
	GiveFlowerURL = "/give"
	FlowertopURL = "/flowertop"
	GetRandomNHIEURL = "/getRandomNHIE"


	// admin handlers

	IsAdminURL = "/isAdmin"
	ChangeAdminURL = "/admin"
	GetAllFlowerTypesURL = "/getAllFlowerTypes"
	RemoveFlowerURL = "/removeFlower"
)
