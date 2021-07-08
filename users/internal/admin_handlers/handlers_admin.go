package adminHandlers

import (
	"encoding/json"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/db"
	"log"

	"github.com/gin-gonic/gin"
)

type AdminHandlers struct {
	DB *db.DbStruct
}

func (ah *AdminHandlers) IsAdminReq(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	err := c.Bind(&req)
	if err != nil || req.ID == 0 {
		c.JSON(400, obj{"err": "no id field"})
		return
	}

	u, err := ah.DB.GetUserByID(req.ID)
	if err != nil {
		c.JSON(400, obj{"result": false, "err": err.Error()})
		return
	}

	c.JSON(200, obj{"result": u.Statuses.IsAdmin})
}

func (ah *AdminHandlers) AdminReq(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}
	err := c.Bind(&req)
	if err != nil || req.ID == 0 {
		c.JSON(400, obj{"err": "no id field"})
		return
	}

	u, err := ah.DB.GetUserByID(req.ID)
	if err != nil {
		c.JSON(400, obj{"err": err.Error()})
		return
	}

	u.Statuses.IsAdmin = !u.Statuses.IsAdmin
	err = ah.DB.UpdateUser(u)
	if err != nil {
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	c.JSON(200, obj{"err": "", "admin": u.Statuses.IsAdmin})
}

func (ah *AdminHandlers) GetAllFlowerTypes(c *gin.Context) {
	data, err := communication.MakeReqToFlowers("getFlowerTypes", nil)
	if err != nil {
		log.Println("handlers_admin.go -> getAllFlowerTypes() error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	var resp struct {
		Result []structs.Flower `json:"result"`
		Err    string   `json:"err"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		log.Printf("handlers_admin.go -> getAllFlowerTypes() -> unmarshal error:%v body: %v\n", err.Error(), string(data))
		c.JSON(400, obj{"err": "unmarhsal error"})
		return
	}
	c.JSON(200, resp)
}

func (ah *AdminHandlers) RemoveFlower(c *gin.Context) {
	var req struct {
		ID uint64 `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		log.Println("handlers_admin.go -> removeFlower() -> binding error:", err.Error())
		c.JSON(400, obj{"err": "no id field"})
		return
	}

	data, err := communication.MakeReqToFlowers("removeFlower", req)
	if err != nil {
		log.Println("handlers_admin.go -> removeFlower() -> removeFlower req error:", err.Error())
		c.JSON(400, obj{"err": err.Error()})
		return
	}
	var resp struct {
		Err string `json:"err"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Printf("handlers_admin.go -> removeFlower() -> unmarshal error: %v, body: %v\n", err.Error(), string(data))
		c.JSON(400, obj{"err": err.Error()})
	}

	c.JSON(200, obj{"err": resp.Err})
}
