package adminHandlers

import (
	"encoding/json"
	flowersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/flowers"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	flowerscfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandlers struct {
	DB *db.DbStruct
}

func (ah *AdminHandlers) IsAdmin(c *gin.Context) {
	var req usersdata.IsAdminReq
	var resp usersdata.IsAdminResp
	err := c.Bind(&req)
	if err != nil || req.ID == 0 {
		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	u, err := ah.DB.GetUserByID(req.ID)
	if err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Result = u.Statuses.IsAdmin
	c.JSON(http.StatusOK, resp)
}

func (ah *AdminHandlers) ChangeAdmin(c *gin.Context) {
	var req usersdata.AdminReq
	var resp usersdata.AdminResp
	err := c.Bind(&req)
	if err != nil || req.ID == 0 {
		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	u, err := ah.DB.GetUserByID(req.ID)
	if err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	u.Statuses.IsAdmin = !u.Statuses.IsAdmin
	err = ah.DB.UpdateUser(u)
	if err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Admin = u.Statuses.IsAdmin
	resp.OK = true
	c.JSON(http.StatusOK, resp)
}

func (ah *AdminHandlers) GetAllFlowerTypes(c *gin.Context) {
	var resp usersdata.GetAllFlowerTypesResp
	data, err := communication.MakeReqToFlowers(flowerscfg.GetFlowerTypesURL, nil)
	if err != nil {
		log.Println("handlers_admin.go -> getAllFlowerTypes() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := json.Unmarshal(data, &resp); err != nil || resp.Err != ""{
		log.Printf("handlers_admin.go -> getAllFlowerTypes() -> unmarshal error:%v body: %v\n", err, string(data))
		resp.Err = "failed to make request to flowers"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ah *AdminHandlers) RemoveFlower(c *gin.Context) {
	var req flowersdata.RemoveFlowerReq
	var resp flowersdata.RemoveFlowerResp
	if err := c.Bind(&req); err != nil {
		log.Println("handlers_admin.go -> removeFlower() -> binding error:", err.Error())
		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	reqToFlowers := flowersdata.RemoveFlowerReq{ID: req.ID}
	respFromFlowers := flowersdata.RemoveFlowerResp{}
	data, err := communication.MakeReqToFlowers(flowerscfg.RemoveFlowerURL, reqToFlowers)
	if err != nil {
		log.Println("handlers_admin.go -> removeFlower() -> removeFlower req error:", err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if err := json.Unmarshal(data, &respFromFlowers); err != nil {
		log.Printf("handlers_admin.go -> removeFlower() -> unmarshal error: %v, body: %v\n", err.Error(), string(data))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if !respFromFlowers.OK {
		resp.Err = respFromFlowers.Err
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.OK = true
	c.JSON(http.StatusOK, resp)
}
