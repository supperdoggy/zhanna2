package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/flowers/internal/service"
	flowersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/flowers"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type obj map[string]interface{}

type Handlers struct {
	service service.IService
	logger  *zap.Logger
}

func NewHandlers(s *service.Service, l *zap.Logger) *Handlers {
	return &Handlers{
		service: s,
		logger:  l,
	}
}

// adds new flower type
func (h Handlers) AddNewFlower(c *gin.Context) {
	var req flowersdata.AddNewFlowerReq
	var resp flowersdata.AddNewFlowerResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.AddNewFlower(req)
	if err != nil {
		h.logger.Error("error AddNewFlower", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// removes flower type
func (h Handlers) RemoveFlower(c *gin.Context) {
	var req flowersdata.RemoveFlowerReq
	var resp flowersdata.RemoveFlowerResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.RemoveFlower(req)
	if err != nil {
		h.logger.Error("error RemoveFlower", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// grows user flower
func (h Handlers) GrowFlower(c *gin.Context) {
	var req flowersdata.GrowFlowerReq
	var resp flowersdata.GrowFlowerResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.GrowFlower(req)
	if err != nil {
		h.logger.Error("error GrowFlower", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)

}

// returns map of user flowers and quantity of different type
func (h Handlers) GetUserFlowers(c *gin.Context) {
	var req flowersdata.GetUserFlowersReq
	var resp flowersdata.GetUserFlowersResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.GetUserFlowers(req)
	if err != nil {
		h.logger.Error("error GetUserFlowers", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// returns bool value if user can grow flower
func (h Handlers) CanGrowFlower(c *gin.Context) {
	var req flowersdata.CanGrowFlowerReq
	var resp flowersdata.CanGrowFlowerResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.CanGrowFlower(req)
	if err != nil {
		h.logger.Error("error CanGrowFlower", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// removeUserFlower - removes current user flower
func (h Handlers) RemoveUserFlower(c *gin.Context) {
	var req flowersdata.RemoveUserFlowerReq
	var resp flowersdata.RemoveUserFlowerResp

	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.RemoveUserFlower(req)
	if err != nil {
		h.logger.Error("error RemoveUserFlower", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// returns int quantity of user grown flowers
func (h Handlers) GetUserFlowerTotal(c *gin.Context) {
	var req flowersdata.GetUserFlowerTotalReq
	var resp flowersdata.GetUserFlowerTotalResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.GetUserFlowerTotal(req)
	if err != nil {
		h.logger.Error("error GetUserFlowerTotal", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handlers) GetLastFlower(c *gin.Context) {
	var req flowersdata.GetLastFlowerReq
	var resp flowersdata.GetLastFlowerResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.GetLastFlower(req)
	if err != nil {
		h.logger.Error("error GetLastFlower", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// userFlowerSlice - returns slice of users flowers
func (h Handlers) UserFlowerSlice(c *gin.Context) {
	var req flowersdata.UserFlowerSliceReq
	var resp flowersdata.UserFlowerSliceResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.UserFlowerSlice(req)
	if err != nil {
		h.logger.Error("error UserFlowerSlice", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// gives flower to other user
func (h Handlers) GiveFlower(c *gin.Context) {
	var req flowersdata.GiveFlowerReq
	var resp flowersdata.GiveFlowerResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.logger.Error("error binding req", zap.Error(err), zap.Any("body", string(data)))

		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.service.GiveFlower(req)
	if err != nil {
		h.logger.Error("error GiveFlower", zap.Error(err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// getFlowerTypes - for admin, returns all flower types
func (h Handlers) GetFlowerTypes(c *gin.Context) {
	resp, err := h.service.GetFlowerTypes()
	if err != nil {
		h.logger.Error("error GetFlowerTypes", zap.Error(err))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
