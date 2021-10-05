package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/fortune"
	"go.uber.org/zap"
	"net/http"
)

type obj map[string]interface{}

type Handlers struct {
	Service fortune.Service
	Logger *zap.Logger
}

func (h *Handlers) GetRandomFortuneCookieReq(c *gin.Context) {
	resp := h.Service.GetRandomFortuneCookie()
	if resp.Err != "" {
		h.Logger.Error("error GetRandomFortuneCookie", zap.Any("error", resp.Err))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
