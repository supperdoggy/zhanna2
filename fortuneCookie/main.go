package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/getRandomFortuneCookie", getRandomFortuneCookieReq)
	}

	if err := r.Run(":1010"); err != nil {
		fmt.Println("error running server!")
	}
}
