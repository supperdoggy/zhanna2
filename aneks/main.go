package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/aneks/internal/aneks"
	db2 "github.com/supperdoggy/superSecretDevelopement/aneks/internal/db"
	handlers1 "github.com/supperdoggy/superSecretDevelopement/aneks/internal/handlers"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/aneks"
)

// TODO: created admin check for DeleteAnekByID and addAnekEndpoint
// TODO: maybe create logs and user request id`s

var handlers = handlers1.Handlers{
	Service: &aneks.AneksService{DB: &db2.DB},
}

func main() {
	r := gin.Default()
	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.GET(cfg.GetRandomAnekURL, handlers.GetRandomAnekReq) // checked, works fine
		apiv1.POST(cfg.GetAnekByIdURL, handlers.GetAnekByID)       // checked, works fine
		apiv1.POST(cfg.DeleteAnekURL, handlers.DeleteAnekByID)     // checked, works fine
		apiv1.POST(cfg.AddAnekURL, handlers.AddAnek)               // checked, works fine
	}
	fmt.Println("Started anek server...")
	if err := r.Run(cfg.Port); err != nil {
		fmt.Println(err.Error())
	}
}
