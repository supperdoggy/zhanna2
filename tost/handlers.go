package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/night-codes/types.v1"
)

func getRandomTostReq(c *gin.Context) {
	a, err := getRandomTost()
	if err != nil {
		resp := obj{"err": err.Error()}
		c.JSON(200, resp)
		return
	}

	c.JSON(200, obj{"text": a.Text, "id": a.ID})
	return
}

func getTostByIdReq(c *gin.Context) {
	var req struct {
		ID int `json:"id" bson:"_id"`
	}
	if err := c.Bind(&req); err != nil {
		c.JSON(200, obj{"err": err.Error()})
		return
	}

	a := getTostById(req.ID)
	c.JSON(200, a)
}

func deleteTostReq(c *gin.Context) {
	var req obj
	if err := c.Bind(&req); err != nil {
		fmt.Println(err.Error())
		c.String(400, err.Error())
		return
	}
	v, ok := req["id"]
	if !ok {
		fmt.Println("no id field")
		c.String(400, "no id field")
		return
	}
	id, ok := v.(float64)
	if !ok {
		fmt.Println("wrong id type (need int)")
		c.String(400, "wrong id type (need int)")
		return
	}
	err := deleteTost(types.Int(id))
	if err != nil {
		fmt.Println(err.Error())
		c.String(400, err.Error())
		return
	}
}

func addTostReq(c *gin.Context) {
	req := map[string]string{}

	if err := c.Bind(&req); err != nil {
		c.String(400, err.Error())
		return
	}

	if v, ok := req["text"]; !ok || v == "" {
		c.String(400, "not text field")
		return
	}

	err := addTost(req["text"])
	if err != nil {
		c.String(400, err.Error())
		return
	}

}
