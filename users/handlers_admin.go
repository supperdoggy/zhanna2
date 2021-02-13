package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

func isAdminReq(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	err := c.Bind(&req)
	if err != nil || req.ID == 0 {
		c.JSON(400, obj{"err": "no id field"})
		return
	}

	u, err := DB.getUserFromDbById(req.ID)
	if err != nil {
		c.JSON(400, obj{"result": false, "err": err.Error()})
		return
	}

	c.JSON(200, obj{"result": u.Statuses.IsAdmin})
}

func adminReq(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	err := c.Bind(&req)
	if err != nil || req.ID == 0 {
		c.JSON(400, obj{"err": "no id field"})
		return
	}

	u, err := DB.getUserFromDbById(req.ID)
	if err != nil {
		c.JSON(400, obj{"err": err.Error()})
		return
	}

	u.Statuses.IsAdmin = !u.Statuses.IsAdmin
	err = DB.updateUser(u)
	if err != nil {
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	c.JSON(200, obj{"err": "", "admin": u.Statuses.IsAdmin})
}

func getAllFlowerTypes(c *gin.Context) {
	data, err := MakeReqToFlowers("getFlowerTypes", nil)
	if err != nil {
		log.Println("handlers_admin.go -> getAllFlowerTypes() error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	var resp struct {
		Result []Flower `json:"result"`
		Err    string   `json:"err"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		log.Printf("handlers_admin.go -> getAllFlowerTypes() -> unmarshal error:%v body: %v\n", err.Error(), string(data))
		c.JSON(400, obj{"err": "unmarhsal error"})
		return
	}
	c.JSON(200, resp)
}
