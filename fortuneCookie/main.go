package main

import (
	"fmt"
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/fortune"
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/handlers"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/fortune"

	"github.com/gin-gonic/gin"
)

func main() {
	h := handlers.Handlers{Service: fortune.Service{DB: db.DB}}
	r := gin.Default()

	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.GET(cfg.GetRandomFortuneCookieURL, h.GetRandomFortuneCookieReq)
	}

	if err := r.Run(cfg.Port); err != nil {
		fmt.Println("error running server!")
	}
}
