package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/den4ik/internal/service"
	den4ikdata "github.com/supperdoggy/superSecretDevelopement/structs/request/den4ik"
	"net/http"
)

type Handlers struct {
	Service service.Service
}

func (h Handlers) GetCard(c *gin.Context) {
	var req den4ikdata.GetCardReq
	var resp den4ikdata.GetCardResp

	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.Service.GetCard(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handlers) ResetSession(c *gin.Context) {
	var req den4ikdata.ResetSessionReq
	var resp den4ikdata.ResetSessionResp

	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.Service.ResetSession(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
