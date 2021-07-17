package main

import (
	"fmt"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/tost"
	"github.com/supperdoggy/superSecretDevelopement/tost/internal/db"
	handlers1 "github.com/supperdoggy/superSecretDevelopement/tost/internal/handlers"
	"github.com/supperdoggy/superSecretDevelopement/tost/internal/tost"

	"github.com/gin-gonic/gin"
)

var handlers = handlers1.Handlers{Service: tost.TostService{DB: &db.DB}}

func main() {
	//TransferToMgo()
	r := gin.Default()
	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.GET(cfg.GetRandomTostURL, handlers.GetRandomTost) // checked, works fine
		apiv1.POST(cfg.GetTostByIdURL, handlers.GetTostById)    // checked, works fine
		apiv1.POST(cfg.DeleteTostURL, handlers.DeleteTost)
		apiv1.POST(cfg.AddTostURL, handlers.AddTost)
	}
	fmt.Println("Started anek server...")
	if err := r.Run(cfg.Port); err != nil {
		fmt.Println(err.Error())
	}
}
