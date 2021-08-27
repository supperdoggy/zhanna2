package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/den4ik/internal/db"
	handlers2 "github.com/supperdoggy/superSecretDevelopement/den4ik/internal/handlers"
	"github.com/supperdoggy/superSecretDevelopement/den4ik/internal/service"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	den4ikcfg "github.com/supperdoggy/superSecretDevelopement/structs/services/den4ik"
)

func main() {
	handlers := handlers2.Handlers{
		Service: service.Service{
			DB: &db.DB,
		},
	}
	r := gin.Default()

	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.POST(den4ikcfg.GetCardURL, handlers.GetCard)
	}

	if err := r.Run(den4ikcfg.Port); err != nil {
		fmt.Println("main.go -> run error:", err.Error())
	}
}
