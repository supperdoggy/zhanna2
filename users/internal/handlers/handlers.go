package handlers

import (
	"fmt"
	"github.com/supperdoggy/parallel_running"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	usersdata "github.com/supperdoggy/superSecretDevelopement/structs/request/users"
	"github.com/supperdoggy/superSecretDevelopement/users/internal/service"
	"gopkg.in/night-codes/types.v1"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"log"
)

type Handlers struct {
	Service service.Service
	ParallelRunning *parallel_running.UserParallelRunning
}

func (h *Handlers) AddOrUpdateUser(c *gin.Context) {
	// todo think of something new
	var req structs.User
	var resp usersdata.AddOrUpdateUserResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> addOrUpdateUserReq() -> binding error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.AddOrUpdateUser(req)
	if err != nil {
		fmt.Println("handlers.go -> AddOrUpdateUser() ->", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) GetFortune(c *gin.Context) {
	var req usersdata.GetFortuneReq
	var resp usersdata.GetFortuneResp

	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getFortune() -> binding error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	// check parallel running
	h.ParallelRunning.Lock("fortune", types.String(req.ID))
	defer h.ParallelRunning.Unlock("fortune", types.String(req.ID))

	resp, err := h.Service.GetFortune(req)
	if err != nil {
		fmt.Println("handlers.go -> GetFortune() ->", err.Error())
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
		fmt.Println("binding error -> getRandomAnek():", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.GetRandomAnek(req)
	if err != nil {
		fmt.Println("handlers.go -> GetRandomAnek() ->", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) GetRandomTost(c *gin.Context) {
	var req usersdata.GetRandomTostReq
	var resp usersdata.GetRandomTostResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> getRandomTost() -> c.Bind() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.GetRandomTost(req)
	if err != nil {
		fmt.Println("handlers.go -> GetRandomTost() ->", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) AddFlower(c *gin.Context) {
	var req usersdata.AddFlowerReq
	var resp usersdata.AddFlowerResp
	if err := c.Bind(&req); err != nil {
		d, err := ioutil.ReadAll(c.Request.Body)
		fmt.Println("handlers.go -> addFlower() -> binding error:", err.Error(), string(d))
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.AddFlower(req)
	if err != nil {
		fmt.Println("handlers.go -> AddFlower() ->", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) Flower(c *gin.Context) {
	var req usersdata.FlowerReq
	var resp usersdata.FlowerResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("handlers.go -> flowerReq() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	h.ParallelRunning.Lock("flower", types.String(req.ID))

	defer h.ParallelRunning.Unlock("flower", types.String(req.ID))

	resp, err := h.Service.Flower(req)
	if err != nil {
		fmt.Println("handlers.go -> Flower() ->", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) DialogFlow(c *gin.Context) {
	var req usersdata.DialogFlowReq
	var resp usersdata.DialogFlowResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("dialogFlowReq() -> c.Bind() error", err.Error())
		resp.Err = "binding error"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.Service.DialogFlow(req)
	if err != nil {
		fmt.Println("handlers.go -> DialogFlow() ->", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) MyFlowers(c *gin.Context) {
	var req usersdata.MyFlowersReq
	var resp usersdata.MyFlowersResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("myflowers() -> c.Bind() error", err.Error())
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.MyFlowers(req)
	if err != nil {
		fmt.Println("handlers.go -> MyFlowers() ->", err.Error())
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
		fmt.Println("handlers.go -> give() -> binding error:", err.Error())
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.GiveFlower(req)
	if err != nil {
		fmt.Println("handlers.go -> GiveFlower() ->", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

// Flowertop - finds all users in chat and forms top users by total flowers
func (h *Handlers) Flowertop(c *gin.Context) {
	var req usersdata.FlowertopReq
	var resp usersdata.FlowertopResp
	if err := c.Bind(&req); err != nil {
		fmt.Println("flowertop() -> c.Bind() error", err.Error())
		resp.Err = "binding error"
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.Flowertop(req)
	if err != nil {
		fmt.Println("handlers.go -> Flowertop() ->", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}

func (h *Handlers) GetRandomNHIE(c *gin.Context) {
	var req usersdata.GetRandomNHIEreq
	var resp usersdata.GetRandomNHIEresp
	if err := c.Bind(&req); err != nil {
		log.Println("handlers.go -> getRandomNHIE() -> c.Bind() error:", err.Error())
		resp.Err = err.Error()
		c.JSON(400, resp)
		return
	}

	resp, err := h.Service.GetRandomNHIE(req)
	if err != nil {
		fmt.Println("handlers.go -> GetRandomNHIE() ->", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(200, resp)
}
