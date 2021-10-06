package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/fortune"
	"go.uber.org/zap"
	"net/http"
)

type Handlers struct {
	service fortune.IService
	logger  *zap.Logger
}

func NewHandlers(s fortune.IService, l *zap.Logger) *Handlers {
	return &Handlers{
		service: s,
		logger:  l,
	}
}

func (h *Handlers) GetRandomFortuneCookieReq(c *gin.Context) {
	resp := h.service.GetRandomFortuneCookie()
	if resp.Err != "" {
		h.logger.Error("error GetRandomFortuneCookie", zap.Any("error", resp.Err))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
