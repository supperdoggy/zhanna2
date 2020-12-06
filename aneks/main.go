package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// TODO: created admin check for deleteAnekReq and addAnekReq

type obj map[string]interface{}

func main() {
	r := gin.Default()
	apiv1 := r.Group("api/v1")
	{
		apiv1.GET("/getRandomAnek", getRandomAnekReq) // checked, works fine
		apiv1.POST("/getAnekById", getAnekByIdReq) // checked, works fine
		apiv1.POST("/deleteAnek", deleteAnekReq) // checked, works fine
		apiv1.POST("/addAnek", addAnekReq) // checked, works fine
	}
	fmt.Println("Started anek server...")
	if err := r.Run(":9090");err!=nil{
		fmt.Println(err.Error())
	}
}