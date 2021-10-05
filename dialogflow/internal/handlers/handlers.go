package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/dialogflow/internal/dialogflow"
	dialogflowdata "github.com/supperdoggy/superSecretDevelopement/structs/request/dialogflow"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type Handlers struct {
	D dialogflow.Dialogflow
	Logger *zap.Logger
}

func (h *Handlers) GetAnswer(c *gin.Context) {
	var req dialogflowdata.GetAnswerReq
	var resp dialogflowdata.GetAnswerResp
	if err := c.Bind(&req); err != nil {
		data, _ := ioutil.ReadAll(c.Request.Body)
		h.Logger.Error("binding error", zap.Error(err), zap.String("req", string(data)))
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp = h.D.DetectIntentText(req)
	if resp.Err != "" {
		h.Logger.Error("error DetectIntentText", zap.String("error", resp.Err), zap.Any("req", req))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
