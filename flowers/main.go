package main

import (
	"github.com/supperdoggy/superSecretDevelopement/flowers/internal/db"
	handlers2 "github.com/supperdoggy/superSecretDevelopement/flowers/internal/handlers"
	"github.com/supperdoggy/superSecretDevelopement/flowers/internal/service"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func main() {
	logger, _ := zap.NewDevelopment()
	r := gin.Default()
	handlers := handlers2.Handlers{
		Service: service.Service{
			DB: db.DB,
			Logger: logger,
		},
		Logger: logger,
	}

	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.POST(cfg.AddNewFlowerURL, handlers.AddNewFlower)
		apiv1.GET(cfg.GetFlowerTypesURL, handlers.GetFlowerTypes)
		apiv1.POST(cfg.RemoveFlowerURL, handlers.RemoveFlower)
		apiv1.POST(cfg.GrowFlowerURL, handlers.GrowFlower)
		apiv1.POST(cfg.GetUserFlowersURL, handlers.GetUserFlowers)
		apiv1.POST(cfg.CanGrowFlowerURL, handlers.CanGrowFlower)
		apiv1.POST(cfg.RemoveUserFlowerURL, handlers.RemoveUserFlower)
		apiv1.POST(cfg.UserFlowerSliceURL, handlers.UserFlowerSlice)
		apiv1.POST(cfg.GiveFlowerURL, handlers.GiveFlower)
	}
	// handlers
	logger.Info("Handlers init start")

	if err := r.Run(cfg.Port); err != nil {
		logger.Error("ERROR RUNNING SERVICE", zap.Error(err))
	}
}
