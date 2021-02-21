package main

import (
	"github.com/gin-gonic/gin"
)

type obj map[string]interface{}

func getRandomNeverHaveIEver(c *gin.Context) {
	resp := getRandomNHIE()
	c.JSON(200, obj{"result": resp})
}
