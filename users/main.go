package main

import "github.com/gin-gonic/gin"

// db itself
var (
	DB = DbStruct{
		DbSession: connectToDB(),
	}
)

func main(){
	r := gin.Default()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/addOrUpdateUser", addOrUpdateUserReq)
	}

}