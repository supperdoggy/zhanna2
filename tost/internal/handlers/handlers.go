package handlers

import (
	tostdata "github.com/supperdoggy/superSecretDevelopement/structs/request/tost"
	"github.com/supperdoggy/superSecretDevelopement/tost/internal/tost"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Service tost.TostService
}

func (h *Handlers) GetRandomTost(c *gin.Context) {
	resp := h.Service.GetRandomTost()
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) GetTostById(c *gin.Context) {
	var req tostdata.GetTostByIdReq
	var resp tostdata.GetTostByIdResp
	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp = h.Service.GetTostById(req)
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) DeleteTost(c *gin.Context) {
	var req tostdata.DeleteTostReq
	var resp tostdata.DeleteTostResp
	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp = h.Service.DeleteTost(req)
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) AddTost(c *gin.Context) {
	var req tostdata.AddTostReq
	var resp tostdata.AddTostResp
	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if req.Text == "" {
		resp.Err = "text field is empty"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp = h.Service.AddTost(req)
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}
