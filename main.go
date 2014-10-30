package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/segments", Auth(false), getSegments)
		v1.GET("/segment/:id", Auth(false), getSegment)
		v1.POST("/segments", Auth(true), addSegment)
		v1.PUT("/segments", Auth(true), editSegment)
		v1.DELETE("/segments/:id", Auth(true), deleteSegment)
		v1.GET("/check_items", Auth(false), getCheckItems)
		v1.GET("/account/login_check", Auth(false), loginCheck)
		v1.POST("/account/login", loginEndpoint)
	}
	r.Run(":8080")

}
