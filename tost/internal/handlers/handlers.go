package handlers

import (
	tostdata "github.com/supperdoggy/superSecretDevelopement/structs/request/tost"
	"github.com/supperdoggy/superSecretDevelopement/tost/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	DB *db.DbStruct
}

func (h *Handlers) GetRandomTost(c *gin.Context) {
	var resp tostdata.GetRandomTostResp
	a, err := h.DB.GetRandomTost()
	if err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Text = a.Text
	resp.ID = a.ID
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

	a := h.DB.GetTostById(req.ID)
	if a.ID == 0 && a.Text == "" {
		resp.Err = "Not Found"
		c.JSON(http.StatusNotFound, resp)
		return
	}
	resp.Text = a.Text
	resp.ID = a.ID
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

	err := h.DB.DeleteTost(req.ID)
	if err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.OK = true
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

	err := h.DB.AddTost(req.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp)
		return
	}

}
