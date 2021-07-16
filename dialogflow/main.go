package main

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/dialogflow/internal/dialogflow"
	handlers2 "github.com/supperdoggy/superSecretDevelopement/dialogflow/internal/handlers"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/dialogflow"
	"log"
)

func main() {
	handlers := handlers2.Handlers{D: dialogflow.DF}
	r := gin.Default()
	
	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.POST(cfg.GetAnswerURL, handlers.GetAnswer)
	}

	if err := r.Run(cfg.Port); err != nil {
		log.Println("Error runnning gin service!!!!!!!!")
	}
}