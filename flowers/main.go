package main

import (
	ai "github.com/night-codes/mgo-ai"
	"github.com/supperdoggy/superSecretDevelopement/flowers/internal/db"
	handlers2 "github.com/supperdoggy/superSecretDevelopement/flowers/internal/handlers"
	"github.com/supperdoggy/superSecretDevelopement/flowers/internal/service"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	logger, _ := zap.NewDevelopment()
	DB := db.NewDB(logger, "", cfg.DBName, cfg.UserFlowerDataCollection, cfg.FlowerCollection)
	Service := service.NewService(DB, logger)
	Handlers := handlers2.NewHandlers(Service, logger)

	ai.Connect(DB.GetUserFlowerDataCollection())

	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.POST(cfg.AddNewFlowerURL, Handlers.AddNewFlower)
		apiv1.GET(cfg.GetFlowerTypesURL, Handlers.GetFlowerTypes)
		apiv1.POST(cfg.RemoveFlowerURL, Handlers.RemoveFlower)
		apiv1.POST(cfg.GrowFlowerURL, Handlers.GrowFlower)
		apiv1.POST(cfg.GetUserFlowersURL, Handlers.GetUserFlowers)
		apiv1.POST(cfg.CanGrowFlowerURL, Handlers.CanGrowFlower)
		apiv1.POST(cfg.RemoveUserFlowerURL, Handlers.RemoveUserFlower)
		apiv1.POST(cfg.UserFlowerSliceURL, Handlers.UserFlowerSlice)
		apiv1.POST(cfg.GiveFlowerURL, Handlers.GiveFlower)
		apiv1.POST(cfg.AddUserFlowerURL, Handlers.AddUserFlower)
	}
	// Handlers
	logger.Info("Handlers init start")

	if err := r.Run(cfg.Port); err != nil {
		logger.Error("ERROR RUNNING SERVICE", zap.Error(err))
	}
}
