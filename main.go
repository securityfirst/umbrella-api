package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

var (
	db *sql.DB
)

func main() {

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/segments", Auth(false), getSegments)
		v1.GET("/segments_raw", Auth(false), getSegmentsRaw)
		v1.GET("/segments_raw/:id/category", Auth(false), getSegmentsRawByCat)
		v1.GET("/segments/:id", Auth(false), getSegment)
		v1.GET("/segments/:id/category", Auth(false), getSegmentsByCat)
		v1.POST("/segments", Auth(true), addSegment)
		v1.PUT("/segments/:id/category", Auth(true), editSegmentByCat)
		v1.POST("/segments/:id/status", Auth(true), approveSegment)
		v1.DELETE("/segments/:id", Auth(true), deleteSegment)
		v1.GET("/check_items", Auth(false), getCheckItems)
		v1.GET("/check_items/:id", Auth(false), getCheckItem)
		v1.POST("/check_items", Auth(true), addCheckItem)
		v1.PUT("/check_items/:id", Auth(true), editCheckItem)
		v1.DELETE("/check_items/:id", Auth(true), deleteCheckItem)
		v1.POST("/check_items/:id/status", Auth(true), approveCheckItem)
		v1.GET("/account/login_check", Auth(true), loginCheck)
		v1.POST("/account/login", loginEndpoint)
		v1.GET("/categories", Auth(false), getCategories)
		v1.GET("/categories/:id/by_parent", Auth(false), getCategoryByParent)
		v1.GET("/categories/:id", Auth(false), getCategory)
		v1.POST("/categories", Auth(true), addCategory)
		v1.PUT("/categories/:id", Auth(true), editCategory)
		v1.POST("/categories/:id/status", Auth(true), approveCategory)
		v1.DELETE("/categories/:id", Auth(true), deleteCategory)
	}

	r.Static("/admin", "./admin")

	r.Run(":8080")

}

func newDb(connection string) *sql.DB {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}

	return db
}
