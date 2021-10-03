package main

import (
	"fmt"
	"github.com/supperdoggy/parallel_running"
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
	logger := zap.L()
	parallelRunning := parallel_running.InitUserParallelRunning(*logger)
	handlers := handlers2.Handlers{Service: service.Service{DB: db.DB}, ParallelRunning: parallelRunning}
	admin_handlers := adminHandlers.AdminHandlers{DB: &db.DB}

	r := gin.Default()

	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.POST(cfg.AddOrUpdateUserURL, handlers.AddOrUpdateUser)
		apiv1.POST(cfg.GetFortuneURL, handlers.GetFortune)
		apiv1.POST(cfg.GetRandomAnekURL, handlers.GetRandomAnek)
		apiv1.POST(cfg.GetRandomTostURL, handlers.GetRandomTost)
		apiv1.POST(cfg.AddFlowerURL, handlers.AddFlower)
		apiv1.POST(cfg.FlowerURL, handlers.Flower)
		apiv1.POST(cfg.DialogFlowHandlerURL, handlers.DialogFlow)
		apiv1.POST(cfg.MyFlowersURL, handlers.MyFlowers)
		apiv1.POST(cfg.GiveFlowerURL, handlers.GiveFlower)
		apiv1.POST(cfg.FlowertopURL, handlers.Flowertop)
		apiv1.POST(cfg.GetRandomNHIEURL, handlers.GetRandomNHIE)
		// todo add check if user is banned
	}

	// admin command handlers
	apiv1_admin := r.Group(defaultCfg.ApiV1Admin)
	{
		apiv1_admin.POST(cfg.IsAdminURL, admin_handlers.IsAdmin)
		apiv1_admin.POST(cfg.ChangeAdminURL, admin_handlers.ChangeAdmin)
		apiv1_admin.GET(cfg.GetAllFlowerTypesURL, admin_handlers.GetAllFlowerTypes)
		apiv1_admin.POST(cfg.RemoveFlowerURL, admin_handlers.RemoveFlower) // ??
	}

	if err := r.Run(cfg.Port); err != nil {
		fmt.Println(err.Error())
	}
}
