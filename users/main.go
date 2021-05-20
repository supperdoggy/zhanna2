package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// DB global object
var DB = &DbStruct{}

func main() {
	DB.initDB()
	fmt.Printf("%+v\n", DB)

	r := gin.Default()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/addOrUpdateUser", addOrUpdateUserReq)
		apiv1.POST("/getFortune", getFortune)
		apiv1.POST("/getRandomAnek", getRandomAnek)
		apiv1.POST("/getRandomTost", getRandomTost)
		apiv1.POST("/addFlower", addFlower)
		apiv1.POST("/flower", flowerReq)
		apiv1.POST("/getAnswer", dialogFlowReq)
		apiv1.POST("/myflowers", myflowers)
		apiv1.POST("/give", give)
		apiv1.POST("/flowertop", flowertop)
		apiv1.POST("/getRandomNHIE", getRandomNHIE)
		// add check if user is banned
	}

	// admin command handlers
	apiv1_admin := r.Group("/api/v1/admin")
	{
		apiv1_admin.POST("/isAdmin", isAdminReq)
		apiv1_admin.POST("/admin", adminReq)
		apiv1_admin.GET("/getAllFlowerTypes", getAllFlowerTypes)
		apiv1_admin.POST("/removeFlower", removeFlower)
	}

	if err := r.Run(":1488"); err != nil {
		fmt.Println(err.Error())
	}
}
