package adminHandlers

import (
	flowersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/flowers"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	flowerscfg "github.com/supperdoggy/superSecretDevelopement/structs/services/flowers"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/db"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandlers struct {
	DB     *db.DbStruct
	Logger *zap.Logger
}

func (ah *AdminHandlers) IsAdmin(c *gin.Context) {
	var req usersdata.IsAdminReq
	var resp usersdata.IsAdminResp
	err := c.Bind(&req)
	if err != nil || req.ID == 0 {
		data, _ := ioutil.ReadAll(c.Request.Body)
		ah.Logger.Error("binding error is admin", zap.Error(err), zap.Any("data", string(data)))
		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	ah.Logger.Info("ADMIN req: isAdmin", zap.Any("req", req))

	u, err := ah.DB.GetUserByID(req.ID)
	if err != nil {
		ah.Logger.Error("error getting user by id admin", zap.Error(err), zap.Any("req", req))
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
		data, _ := ioutil.ReadAll(c.Request.Body)
		ah.Logger.Error("binding error is admin", zap.Error(err), zap.Any("data", string(data)))

		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	ah.Logger.Info("ADMIN req: ChangeAdmin", zap.Any("req", req))

	u, err := ah.DB.GetUserByID(req.ID)
	if err != nil {
		ah.Logger.Error("error getting user by id", zap.Error(err), zap.Any("req", req))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	u.Statuses.IsAdmin = !u.Statuses.IsAdmin
	err = ah.DB.UpdateUser(u)
	if err != nil {
		ah.Logger.Error("error updating user", zap.Error(err), zap.Any("user", u))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Admin = u.Statuses.IsAdmin
	resp.OK = true
	ah.Logger.Info("ADMIN req: ChangeAdmin", zap.Any("req", req), zap.Any("resp", resp))

	c.JSON(http.StatusOK, resp)
}

// todo should be flowersdata.GetAllFlowerTypesResp and only after that usersdata.GetAllFlowerTypesResp
func (ah *AdminHandlers) GetAllFlowerTypes(c *gin.Context) {
	var resp usersdata.GetAllFlowerTypesResp
	err := communication.MakeReqToFlowers(flowerscfg.GetFlowerTypesURL, nil, &resp)
	if err != nil {
		ah.Logger.Error("error making req to flowers getAllFlowerTypes", zap.Error(err))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if resp.Err != "" {
		ah.Logger.Error("unmarshal error", zap.Error(err), zap.Any("resp", resp))
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
		data, _ := ioutil.ReadAll(c.Request.Body)
		ah.Logger.Error("binding error is admin", zap.Error(err), zap.Any("data", string(data)))

		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	reqToFlowers := flowersdata.RemoveFlowerReq{ID: req.ID}
	respFromFlowers := flowersdata.RemoveFlowerResp{}
	err := communication.MakeReqToFlowers(flowerscfg.RemoveFlowerURL, reqToFlowers, &respFromFlowers)
	if err != nil {
		ah.Logger.Error("error making req to flowers RemoveFlowerURL", zap.Error(err), zap.Any("req", req))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !respFromFlowers.OK {
		ah.Logger.Error("got error from flowers RemoveFlowerURL",
			zap.Error(err),
			zap.Any("req", req),
			zap.Any("resp", respFromFlowers),
		)

		resp.Err = respFromFlowers.Err
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.OK = true
	c.JSON(http.StatusOK, resp)
}
