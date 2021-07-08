package main

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/neverHaveIEver/internal/db"
	handlers2 "github.com/supperdoggy/superSecretDevelopement/neverHaveIEver/internal/handlers"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/NHIE"

	"github.com/gin-gonic/gin"
)

func main() {
	handlers := handlers2.Handlers{DB: &db.DB}
	r := gin.Default()

	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.GET(cfg.GetRandomNeverHaveIEverURL, handlers.GetRandomNeverHaveIEver)
	}

	if err := r.Run(cfg.Port); err != nil {
		fmt.Println("error running server!")
	}
}
