package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var db = DbStruct{
	DbSession: connectToDB(),
}

func main() {
	db.TostCollection = db.connectToTostCollection()
	//TransferToMgo()
	r := gin.Default()
	apiv1 := r.Group("api/v1")
	{
		apiv1.GET("/getRandomTost", getRandomTostReq) // checked, works fine
		apiv1.POST("/getTostById", getTostByIdReq)    // checked, works fine
		apiv1.POST("/deleteTost", deleteTostReq)      // checked, works fine
		apiv1.POST("/addTost", addTostReq)            // checked, works fine
	}
	fmt.Println("Started anek server...")
	if err := r.Run(":9393"); err != nil {
		fmt.Println(err.Error())
	}
}
