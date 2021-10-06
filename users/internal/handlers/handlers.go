package handlers

import (
	"github.com/supperdoggy/superSecretDevelopement/structs"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/service"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Service service.IService
	Logger  *zap.Logger
}

func NewHandlers(s service.IService, logger *zap.Logger) *Handlers {
	return &Handlers{
		Service: s,
		Logger:  logger,
	}
}

func (h *Handlers) AddOrUpdateUser(c *gin.Context) {
	// todo tink of something new
	var req structs.User
	var resp usersdata.AddOrUpdateUserResp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.AddOrUpdateUser(req)
	if err != nil {
		h.Logger.Error("add or update error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

// todo simplify
func (h *Handlers) GetFortune(c *gin.Context) {
	var req usersdata.GetFortuneReq
	var resp usersdata.GetFortuneResp

	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.GetFortune(req)
	if err != nil {
		h.Logger.Error("GetFortune error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

// thats ok i guess
func (h *Handlers) GetRandomAnek(c *gin.Context) {
	var req usersdata.GetRandomAnekReq
	var resp usersdata.GetRandomAnekResp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.GetRandomAnek(req)
	if err != nil {
		h.Logger.Error("GetRandomAnek error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) GetRandomTost(c *gin.Context) {
	var req usersdata.GetRandomTostReq
	var resp usersdata.GetRandomTostResp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.GetRandomTost(req)
	if err != nil {
		h.Logger.Error("GetRandomTost error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) AddFlower(c *gin.Context) {
	var req usersdata.AddFlowerReq
	var resp usersdata.AddFlowerResp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.AddFlower(req)
	if err != nil {
		h.Logger.Error("AddFlower error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) Flower(c *gin.Context) {
	var req usersdata.FlowerReq
	var resp usersdata.FlowerResp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.Flower(req)
	if err != nil {
		h.Logger.Error("Flower error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) DialogFlow(c *gin.Context) {
	var req usersdata.DialogFlowReq
	var resp usersdata.DialogFlowResp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.Service.DialogFlow(req)
	if err != nil {
		h.Logger.Error("DialogFlow error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) MyFlowers(c *gin.Context) {
	var req usersdata.MyFlowersReq
	var resp usersdata.MyFlowersResp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.MyFlowers(req)
	if err != nil {
		h.Logger.Error("MyFlowers error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

// path for giving flower
func (h *Handlers) GiveFlower(c *gin.Context) {
	var req usersdata.GiveFlowerReq
	var resp usersdata.GiveFlowerResp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.GiveFlower(req)
	if err != nil {
		h.Logger.Error("GiveFlower error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

// Flowertop - finds all users in chat and forms top users by total flowers
// TODO Simplify
// TODO dude rewrite this for the love of god
func (h *Handlers) Flowertop(c *gin.Context) {
	var req usersdata.FlowertopReq
	var resp usersdata.FlowertopResp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.Flowertop(req)
	if err != nil {
		h.Logger.Error("Flowertop error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) GetRandomNHIE(c *gin.Context) {
	var req usersdata.GetRandomNHIEreq
	var resp usersdata.GetRandomNHIEresp
	if err := c.Bind(&req); err != nil {
		d, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("error binding request body", zap.Error(err), zap.String("body", string(d)))

		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.GetRandomNHIE(req)
	if err != nil {
		h.Logger.Error("GetRandomNHIE error", zap.Error(err), zap.Any("request", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}
