package main

import (
	"fmt"

	ai "github.com/night-codes/mgo-ai"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	ai.Connect(DB.FlowerCollection)
	ai.Connect(DB.UserFlowerDataCollection)
	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/addFlower", addNewFlower)
		apiv1.POST("/removeFlower", removeFlower)
		apiv1.POST("/growFlower", growFlowerReq)
		apiv1.POST("/getUserFlowers", getUserFlowers)
		apiv1.POST("/canGrow", canGrowFlower)
		apiv1.POST("/removeUserFlower", removeUserFlower)
	}
	// handlers
	fmt.Println("Handlers init start")

	if err := r.Run(":2345"); err != nil {
		fmt.Println("MAIN.GO -> RUN ERROR:", err.Error())
	}
}
