package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/fortuneCookie/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	fortunedata "github.com/supperdoggy/superSecretDevelopement/structs/request/fortune"
	"net/http"
)

type obj map[string]interface{}

type Handlers struct {
	DB db.DbStruct
}

func (h *Handlers) GetRandomFortuneCookieReq(c *gin.Context) {
	var resp fortunedata.GetRandomFortuneCookieResp
	var cookie structs.Cookie

	cookie, err := h.DB.GetRandomFortune()
	if err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Text = cookie.Text
	resp.ID = cookie.ID
	c.JSON(http.StatusOK, resp)
}
