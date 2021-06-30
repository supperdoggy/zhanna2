package main

import (
	"fmt"
	defaultCfg "github.com/supperdoggy/superSecretDevelopement/structs/request/default"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/tost"
	"github.com/supperdoggy/superSecretDevelopement/tost/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/tost/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	h := handlers.Handlers{DB: &db.DB}
	//TransferToMgo()
	r := gin.Default()
	apiv1 := r.Group(defaultCfg.ApiV1)
	{
		apiv1.GET(cfg.GetRandomTostURL, h.GetRandomTost) // checked, works fine
		apiv1.POST(cfg.GetTostByIdURL, h.GetTostById)    // checked, works fine
		apiv1.POST(cfg.DeleteTostURL, h.DeleteTost)
		apiv1.POST(cfg.AddTostURL, h.AddTost)
	}
	fmt.Println("Started anek server...")
	if err := r.Run(cfg.Port); err != nil {
		fmt.Println(err.Error())
	}
}
