package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// db itself
var (
	DB = DbStruct{
		DbSession: connectToDB(),
	}
)

func main() {
	DB.AdminCollection = connectToAdminCollection()
	DB.UsersCollection = connectToUsersCollection()
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
	}

	if err := r.Run(":1488"); err != nil {
		fmt.Println(err.Error())
	}
}
