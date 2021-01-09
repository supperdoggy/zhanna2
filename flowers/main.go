package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/addFlower", addNewFlower)
		apiv1.POST("/removeFlower", removeFlower)
	}
	// handlers
	fmt.Println("Handlers init start")

	if err := r.Run(":2345"); err != nil {
		fmt.Println("MAIN.GO -> RUN ERROR:", err.Error())
	}
}
