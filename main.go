package main

import (
	"flag"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/tmilewski/goenv"
)

var isProduction bool

func main() {
	err := goenv.Load()
	checkErr(err)
	// envArg := flag.String("env", "production", "Environment")
	envArg := flag.String("env", "development", "Environment")
	flag.Parse()
	if *envArg == "production" {
		isProduction = true
	}
	um := getUmbrella()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*") // Make templates available
	r.Static("/assets", "./assets")
	r.Use(cors.Middleware(cors.Config{
		Origins:         "http://localhost:8000, http://127.0.0.1:8000",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "Access-Control-Allow-Origin",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	admin := r.Group("/admin")
	{
		admin.GET("/", um.WebAuth(), um.Index)
		admin.GET("/login", um.Login)
		admin.POST("/login", um.LoginPost)
		admin.GET("/logout", um.WebAuth(), um.LogOut)
	}
	v1 := r.Group("/v1")
	{
		v1.GET("/account/login_check", um.Auth(true), um.loginCheck)
		v1.POST("/account/login", um.loginEndpoint)

		v1.GET("/feed", um.getFeed)

		v1.GET("/segments", um.getSegments)
		// v1.GET("/segments/:id", um.Auth(false), um.getSegment)
		// v1.GET("/segments/:id/category", um.Auth(false), um.getSegmentsByCat)
		// v1.POST("/segments", um.Auth(true), um.addSegment)
		// v1.PUT("/segments/:id/category", um.Auth(true), um.editSegmentByCat)
		// v1.POST("/segments/:id/status", um.Auth(true), um.approveSegment)
		// v1.DELETE("/segments/:id", um.Auth(true), um.deleteSegment)

		v1.GET("/check_items", um.Auth(false), um.getCheckItems)
		// v1.GET("/check_items/:id", um.Auth(false), um.getCheckItem)
		// v1.GET("/check_items/:id/category", um.Auth(false), um.getCheckItemsByCat)
		// v1.POST("/check_items", um.Auth(true), um.addCheckItem)
		// v1.DELETE("/check_items/:id", um.Auth(true), um.deleteCheckItem)
		// v1.POST("/check_items/:id/status", um.Auth(true), um.approveCheckItem)

		v1.GET("/categories", um.Auth(false), um.getCategories)
		// v1.GET("/categories/:id/by_parent", um.Auth(false), um.getCategoryByParent)
		// v1.GET("/categories/:id", um.Auth(false), um.getCategory)
		// v1.POST("/categories", um.Auth(true), um.addCategory)
		// v1.PUT("/categories/:id", um.Auth(true), um.editCategory)
		// v1.POST("/categories/:id/status", um.Auth(true), um.approveCategory)
		// v1.DELETE("/categories/:id", um.Auth(true), um.deleteCategory)
		v1.GET("/languages", um.Auth(false), um.getLanguages)
	}
	r.Run(":8080")

}
