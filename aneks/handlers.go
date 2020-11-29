package main

import "github.com/gin-gonic/gin"

func getRandomAnekReq(c *gin.Context){
	a, err := getRandomAnek()
	if err != nil{
		resp := obj{"err":err.Error()}
		c.JSON(200, resp)
		return
	}

	c.JSON(200, obj{"resp":a})
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

	id, ok := v.(int)
	if !ok{
		c.JSON(200, obj{"err":"wrong id type (need int)"})
		return
	}

	a := getAnekById(id)
	c.JSON(200, a)
}
