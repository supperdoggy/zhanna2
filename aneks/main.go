package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type obj map[string]interface{}

func main() {
	r := gin.Default()
	apiv1 := r.Group("api/v1")
	{
		apiv1.GET("/getRandomAnek", getRandomAnekReq)
		apiv1.POST("/getAnekById", getAnekByIdReq)
	}
	fmt.Println("Started anek server...")
	if err := r.Run(":9090");err!=nil{
		fmt.Println(err.Error())
	}
}