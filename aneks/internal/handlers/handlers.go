package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/aneks/internal/aneks"
	aneksdata "github.com/supperdoggy/superSecretDevelopement/structs/request/aneks"
	"net/http"
)

type Handlers struct {
	Service *aneks.AneksService
}

func (h *Handlers) GetRandomAnekReq(c *gin.Context) {
	resp := h.Service.GetRandomAnek()
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *Handlers) GetAnekByID(c *gin.Context) {
	var req aneksdata.GetAnekByIdReq
	var resp aneksdata.GetAnekByIdResp
	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp = h.Service.GetAnekByID(req)
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) DeleteAnekByID(c *gin.Context) {
	var req aneksdata.DeleteAnekByIDReq
	var resp aneksdata.DeleteAnekByIDResp
	if err := c.Bind(&req); err != nil {
		fmt.Println(err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp = h.Service.DeleteAnekByID(req)
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) AddAnek(c *gin.Context) {
	var req aneksdata.AddAnekReq
	var resp aneksdata.AddAnekResp
	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp = h.Service.AddAnek(req)
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
