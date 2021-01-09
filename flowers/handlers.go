package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func addNewFlower(c *gin.Context) {
	var req Flower
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> addNewFlower() -> binding error:", err.Error())
		c.JSON(400, obj{"err": "binding error"})
		return
	}
	if req.Name == "" || req.Icon == "" || req.Type == "" {
		c.JSON(400, obj{"err": "fill all fields"})
		return
	}

	req.ID = 2
	req.CreationTime = time.Now()
	if err := DB.addFlower(req); err != nil {
		fmt.Println("handlers.go -> addNewFlower() -> addFlower(req) error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	c.JSON(200, obj{"err": nil})
}

func removeFlower(c *gin.Context) {
	var req map[string]uint64

	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> removeFlower() -> bind error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	if _, ok := req["id"]; !ok {
		c.JSON(400, obj{"err": "no id field"})
		return
	}
	err := DB.removeFlower(req["id"])
	if err != nil {
		fmt.Println("handlers.go -> removeFlower() -> removeFlower() error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	c.JSON(200, obj{"err": nil})
}
