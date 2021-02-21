package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type obj map[string]interface{}

// redo!!!!
func getRandomFortuneCookieReq(c *gin.Context) {
	var result Cookie
	rand.Seed(time.Now().Unix())
	size, err := DB.CookieCollection.Count()
	if err != nil {
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	randomId := rand.Intn(size - 1)
	if err := DB.CookieCollection.Find(obj{"_id": randomId}).One(&result); err != nil {
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	c.JSON(200, result)
}
