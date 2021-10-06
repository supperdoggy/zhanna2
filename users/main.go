package main

import (
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/users"
	adminHandlers "github.com/supperdoggy/superSecretDevelopement/users/internal/admin_handlers"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/db"
	handlers2 "github.com/supperdoggy/superSecretDevelopement/users/internal/handlers"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/service"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)


func main() {
	logger, _ := zap.NewDevelopment()
	DB := db.NewDB(cfg.DBName, cfg.UsersCollection, cfg.AdminCollection, cfg.MessagesCollection, logger)
	Service := service.NewService(DB, logger)
	Handlers := handlers2.NewHandlers(Service, logger)
	adminHandlers := adminHandlers.NewAdminHandlers(DB, logger)
	r := gin.Default()

	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.POST(cfg.AddOrUpdateUserURL, Handlers.AddOrUpdateUser)
		apiv1.POST(cfg.GetFortuneURL, Handlers.GetFortune)
		apiv1.POST(cfg.GetRandomAnekURL, Handlers.GetRandomAnek)
		apiv1.POST(cfg.GetRandomTostURL, Handlers.GetRandomTost)
		apiv1.POST(cfg.AddFlowerURL, Handlers.AddFlower)
		apiv1.POST(cfg.FlowerURL, Handlers.Flower)
		apiv1.POST(cfg.DialogFlowHandlerURL, Handlers.DialogFlow)
		apiv1.POST(cfg.MyFlowersURL, Handlers.MyFlowers)
		apiv1.POST(cfg.GiveFlowerURL, Handlers.GiveFlower)
		apiv1.POST(cfg.FlowertopURL, Handlers.Flowertop)
		apiv1.POST(cfg.GetRandomNHIEURL, Handlers.GetRandomNHIE)
		// todo add check if user is banned
	}

	// admin command handlers
	apiv1Admin := r.Group(defaultCfg.ApiV1Admin)
	{
		apiv1Admin.POST(cfg.IsAdminURL, adminHandlers.IsAdmin)
		apiv1Admin.POST(cfg.ChangeAdminURL, adminHandlers.ChangeAdmin)
		apiv1Admin.GET(cfg.GetAllFlowerTypesURL, adminHandlers.GetAllFlowerTypes)
		apiv1Admin.POST(cfg.RemoveFlowerURL, adminHandlers.RemoveFlower) // ??
	}

	if err := r.Run(cfg.Port); err != nil {
		logger.Error("error running service", zap.Error(err))
	}
}
