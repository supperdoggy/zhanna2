package main

import (
	"fmt"
	"math/rand"
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

	id, _ := DB.FlowerCollection.Count()
	req.ID = uint64(id) + 3
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

func growFlowerReq(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> growFlowerReq) -> binding error:", err.Error())
		c.JSON(400, obj{"err": "binding error"})
		return
	}
	fmt.Println("bind ok")
	flower, err := DB.getUserFlower(req.ID)
	fmt.Println(flower, err)
	if err != nil && err.Error() != "not found" {
		fmt.Println("handlers.go -> growFlowerReq) -> getUserFlower() error:", err.Error())
		c.JSON(400, obj{"err": "error getting flower"})
		return
	}
	// not found flower, creating new
	if err != nil && err.Error() == "not found" {
		flower, err = DB.getRandomFlower()
		if err != nil {
			fmt.Println("handlers.go -> growFlowerReq) -> getRandomFlower() error:", err.Error())
			c.JSON(400, obj{"err": err.Error()})
			return
		}
		id, err := DB.UserFlowerDataCollection.Count()
		if err != nil {
			fmt.Println("handlers.go -> growFlowerReq) -> Count() error:", err.Error())
			c.JSON(400, obj{"err": err.Error()})
			return
		}
		flower.ID = uint64(id + 1)
		flower.Owner = req.ID
	}
	flower.HP += uint8(rand.Intn(31))
	if flower.HP > 100 {
		flower.HP = 100
	}
	flower.LastTimeGrow = time.Now()
	if _, err := DB.UserFlowerDataCollection.Upsert(obj{"_id": flower.ID}, flower); err != nil {
		fmt.Println("handlers.go -> growFlowerReq) -> Upsert() error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	c.JSON(200, flower)

}

func getUserFlowers(c *gin.Context) {
	var req struct {
		ID int `json:"id" bson:"owner"`
	}
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getUserFlowers() -> binding error:", err.Error())
		c.JSON(400, obj{"err": "binding error"})
		return
	}
	flowers, err := DB.getAllUserFlowers(req.ID)
	if err != nil {
		fmt.Println("handlers.go -> getUserFlowers() -> getAllUserFlowers() error:", err.Error())
		c.JSON(400, obj{"err": "error getting flowers"})
	}
	c.JSON(200, obj{"flowers": flowers})
}

func canGrowFlower(c *gin.Context) {
	var req struct {
		ID int `json:"id" bson:"own"`
	}
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> canGrowFlower() -> binding error:", err.Error())
		c.JSON(400, obj{"answer": false, "err": "binding error"})
		return
	}
	flower, err := DB.getUserFlower(req.ID)
	if err != nil {
		// if we cant find flower in the collection we return true
		if err.Error() == "not found" {
			c.JSON(200, obj{"answer": true})
			return
		}
		// if we cant find due to mongo error then return error
		fmt.Println("handlers.go -> canGrowFlower() -> getUserFlower() error:", err.Error())
		c.JSON(400, obj{"answer": false, "err": "get flower error"})
		return
	}
	canGrow := int(time.Now().Sub(flower.LastTimeGrow).Hours())/growTimeout >= 1
	c.JSON(200, obj{"answer": canGrow})
}
