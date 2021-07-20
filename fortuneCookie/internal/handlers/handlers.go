package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/fortune"
	"net/http"
)

type obj map[string]interface{}

type Handlers struct {
	Service fortune.Service
}

func (h *Handlers) GetRandomFortuneCookieReq(c *gin.Context) {
	resp := h.Service.GetRandomFortuneCookie()
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
