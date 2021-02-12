package main

import (
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
