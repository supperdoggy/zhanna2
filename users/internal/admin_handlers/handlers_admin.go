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
	db     db.IDbStruct
	logger *zap.Logger
}

func NewAdminHandlers(d db.IDbStruct, logger *zap.Logger) *AdminHandlers {
	return &AdminHandlers{
		db:     d,
		logger: logger,
	}
}

func (ah *AdminHandlers) IsAdmin(c *gin.Context) {
	var req usersdata.IsAdminReq
	var resp usersdata.IsAdminResp
	err := c.Bind(&req)
	if err != nil || req.ID == 0 {
		data, _ := ioutil.ReadAll(c.Request.Body)
		ah.logger.Error("binding error is admin", zap.Error(err), zap.Any("data", string(data)))
		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	ah.logger.Info("ADMIN req: isAdmin", zap.Any("req", req))

	u, err := ah.db.GetUserByID(req.ID)
	if err != nil {
		ah.logger.Error("error getting user by id admin", zap.Error(err), zap.Any("req", req))
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
		ah.logger.Error("binding error is admin", zap.Error(err), zap.Any("data", string(data)))

		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	ah.logger.Info("ADMIN req: ChangeAdmin", zap.Any("req", req))

	u, err := ah.db.GetUserByID(req.ID)
	if err != nil {
		ah.logger.Error("error getting user by id", zap.Error(err), zap.Any("req", req))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	u.Statuses.IsAdmin = !u.Statuses.IsAdmin
	err = ah.db.UpdateUser(u)
	if err != nil {
		ah.logger.Error("error updating user", zap.Error(err), zap.Any("user", u))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Admin = u.Statuses.IsAdmin
	resp.OK = true
	ah.logger.Info("ADMIN req: ChangeAdmin", zap.Any("req", req), zap.Any("resp", resp))

	c.JSON(http.StatusOK, resp)
}

// todo should be flowersdata.GetAllFlowerTypesResp and only after that usersdata.GetAllFlowerTypesResp
func (ah *AdminHandlers) GetAllFlowerTypes(c *gin.Context) {
	var resp usersdata.GetAllFlowerTypesResp
	err := communication.MakeReqToFlowers(flowerscfg.GetFlowerTypesURL, nil, &resp)
	if err != nil {
		ah.logger.Error("error making req to flowers getAllFlowerTypes", zap.Error(err))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if resp.Err != "" {
		ah.logger.Error("unmarshal error", zap.Error(err), zap.Any("resp", resp))
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
		ah.logger.Error("binding error is admin", zap.Error(err), zap.Any("data", string(data)))

		resp.Err = "no id field"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := communication.MakeReqToFlowers(flowerscfg.RemoveFlowerURL, req, &resp)
	if err != nil {
		ah.logger.Error("error making req to flowers RemoveFlowerURL", zap.Error(err), zap.Any("req", req))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !resp.OK {
		ah.logger.Error("got error from flowers RemoveFlowerURL",
			zap.Error(err),
			zap.Any("req", req),
		)

		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.OK = true
	c.JSON(http.StatusOK, resp)
}

func (ah *AdminHandlers) AddUserFlower(c *gin.Context) {
	var req flowersdata.AddUserFlowerReq
	var resp flowersdata.AddUserFlowerResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		ah.logger.Error("error binding req", zap.Any("json", string(data)))
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := communication.MakeReqToFlowers(flowerscfg.AddUserFlowerURL, req, &resp)
	if err != nil {
		ah.logger.Error("error making request to users",
			zap.Error(err),
			zap.String("url", flowerscfg.AddUserFlowerURL),
			zap.Any("req", req))
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if resp.Error != "" {
		ah.logger.Error("got error from flowers",
			zap.String("error", resp.Error),
			zap.String("url", flowerscfg.AddUserFlowerURL),
			zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}
