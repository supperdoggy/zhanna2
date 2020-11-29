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
	}
	a := getRandomAnek()
	fmt.Println(a.Id, a.Text)
}