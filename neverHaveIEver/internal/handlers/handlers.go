package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/neverHaveIEver/internal/db"
	"net/http"
)

type Handlers struct {
	DB *db.DBStruct
}

func (h *Handlers) GetRandomNeverHaveIEver(c *gin.Context) {
	resp := h.DB.GetRandomNHIE()
	c.JSON(http.StatusOK, resp)
}
