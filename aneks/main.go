package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db2 "github.com/supperdoggy/superSecretDevelopement/aneks/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/aneks/internal/handlers"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/aneks"
)

// TODO: created admin check for DeleteAnekByID and addAnekEndpoint
// TODO: maybe create logs and user request id`s

func main() {
	var h = handlers.Handlers{
		DB: &db2.DB,
	}
	r := gin.Default()
	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.GET(cfg.GetRandomAnekURL, h.GetRandomAnekReq) // checked, works fine
		apiv1.POST(cfg.GetAnekByIdURL, h.GetAnekByID)       // checked, works fine
		apiv1.POST(cfg.DeleteAnekURL, h.DeleteAnekByID)     // checked, works fine
		apiv1.POST(cfg.AddAnekURL, h.AddAnek)               // checked, works fine
	}
	fmt.Println("Started anek server...")
	if err := r.Run(cfg.Port); err != nil {
		fmt.Println(err.Error())
	}
}
