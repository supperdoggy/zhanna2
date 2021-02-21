package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/getRandomNeverHaveIEver", getRandomNeverHaveIEver)
	}

	if err := r.Run(":1122"); err != nil {
		fmt.Println("error running server!")
	}
}
