package main

import "github.com/gin-gonic/gin"

func main(){
	r := gin.Default()

	v1 := r.Group("/api/v1")
	v1.POST("/sendText", )
}
