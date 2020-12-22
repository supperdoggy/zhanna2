package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// db itself
var (
	DB = DbStruct{
		DbSession: connectToDB(),
	}
)

func main() {
	DB.AdminCollection = connectToAdminCollection()
	DB.UsersCollection = connectToUsersCollection()
	r := gin.Default()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/addOrUpdateUser", addOrUpdateUserReq)
		apiv1.POST("/getFortune", getFortune)
	}

	if err := r.Run(":1488"); err != nil {
		fmt.Println(err.Error())
	}
}
