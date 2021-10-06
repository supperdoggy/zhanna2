package main

import (
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/fortune"
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/handlers"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/fortune"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func main() {
	logger, _ := zap.NewDevelopment()
	DB := db.NewDB(logger, "", cfg.DBName, cfg.FortuneCollection)
	Service := fortune.NewService(logger, DB)
	h := handlers.NewHandlers(Service, logger)
	r := gin.Default()

	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.GET(cfg.GetRandomFortuneCookieURL, h.GetRandomFortuneCookieReq)
	}

	if err := r.Run(cfg.Port); err != nil {
		logger.Error("error running server", zap.Error(err))
	}
}
