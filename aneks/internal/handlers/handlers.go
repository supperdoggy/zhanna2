package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/superSecretDevelopement/aneks/internal/db"
	aneksdata "github.com/supperdoggy/superSecretDevelopement/structs/request/aneks"
	"net/http"
)

type Handlers struct {
	DB db.DB
}

func (h *Handlers) GetRandomAnekReq(c *gin.Context) {
	var resp aneksdata.GetRandomAnekResp
	a, err := h.DB.GetRandomAnek()
	if err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Text = a.Text
	resp.ID = a.Id

	c.JSON(http.StatusOK, resp)
	return
}

func (h *Handlers) GetAnekByIDEndpiont(c *gin.Context) {
	var req aneksdata.GetAnekByIdReq
	var resp aneksdata.GetAnekByIdResp
	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	a := h.DB.GetAnekById(req.ID)
	resp.Anek = a
	c.JSON(http.StatusOK, a)
}

func (h *Handlers) DeleteAnekByIDEndpoint(c *gin.Context) {
	var req aneksdata.DeleteAnekByIDReq
	var resp aneksdata.DeleteAnekByIDResp
	if err := c.Bind(&req); err != nil {
		fmt.Println(err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := h.DB.DeleteAnek(req.ID)
	if err != nil {
		fmt.Println(err.Error())
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.OK = true

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) AddAnekEndpoint(c *gin.Context) {
	var req aneksdata.AddAnekReq
	var resp aneksdata.AddAnekResp
	if err := c.Bind(&req); err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if req.Text == "" {
		resp.Err = "text field cant be empty"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := h.DB.AddAnek(req.Text)
	if err != nil {
		resp.Err = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.OK = true
	c.JSON(http.StatusOK, resp)
}
