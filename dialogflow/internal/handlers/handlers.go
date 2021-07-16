package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/dialogflow/internal/dialogflow"
	dialogflowdata "github.com/supperdoggy/superSecretDevelopement/structs/request/dialogflow"
	"net/http"
)

type Handlers struct {
	D dialogflow.Dialogflow
}

func (h *Handlers) GetAnswer(c *gin.Context) {
	var req dialogflowdata.GetAnswerReq
	var resp dialogflowdata.GetAnswerResp
	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp = h.D.DetectIntentText(req)
	if resp.Err != "" {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
