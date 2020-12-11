package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/night-codes/types.v1"
)

func getRandomAnekReq(c *gin.Context){
	a, err := getRandomAnek()
	if err != nil{
		resp := obj{"err":err.Error()}
		c.JSON(200, resp)
		return
	}

	c.JSON(200, obj{"text":a.Text, "id":a.Id})
	return
}

func getAnekByIdReq(c *gin.Context){
	var req obj
	if err := c.Bind(&req); err != nil{
		c.JSON(200, obj{"err":err.Error()})
		return
	}
	v, ok := req["id"]
	if !ok{
		c.JSON(200, obj{"err":"no id field"})
		return
	}

	id, ok := v.(float64)
	if !ok{
		c.JSON(200, obj{"err":"wrong id type (need int)"})
		return
	}

	a := getAnekById(types.Int(id))
	c.JSON(200, a)
}

func deleteAnekReq(c *gin.Context){
	var req obj
	if err := c.Bind(&req);err != nil{
		fmt.Println(err.Error())
		c.String(400, err.Error())
		return
	}
	v, ok := req["id"]
	if !ok{
		fmt.Println("no id field")
		c.String(400, "no id field")
		return
	}
	id, ok := v.(float64)
	if !ok{
		fmt.Println("wrong id type (need int)")
		c.String(400, "wrong id type (need int)")
		return
	}
	err := deleteAnek(types.Int(id))
	if err != nil{
		fmt.Println(err.Error())
		c.String(400, err.Error())
		return
	}
}

func addAnekReq(c *gin.Context){
	req := map[string]string{}

	if err := c.Bind(&req);err != nil{
		c.String(400, err.Error())
		return
	}

	if v, ok := req["text"];!ok || v ==""{
		c.String(400, "not text field")
		return
	}

	err := addAnek(req["text"])
	if err != nil{
		c.String(400, err.Error())
		return
	}

}
