package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	cors "github.com/itsjamie/gin-cors"

	"github.com/securityfirst/tent"
	"github.com/securityfirst/tent/auth"
	"github.com/securityfirst/tent/repo"
	"github.com/securityfirst/umbrella-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/tmilewski/goenv"
)

var isProduction bool

var conf = auth.Config{
	OAuthHost: "http://127.0.0.1:8080",
	Login:     auth.HandleConf{"/auth/login", "/"},
	Logout:    auth.HandleConf{"/auth/logout", "/"},
	Callback:  auth.HandleConf{"/auth/callback", "/"},
}

func init() {
	err := goenv.Load()
	utils.CheckErr(err)
	// envArg := flag.String("env", "production", "Environment")
	envArg := flag.String("env", "development", "Environment")
	flag.Parse()
	if *envArg == "production" {
		isProduction = true
	}
	conf.Id = os.Getenv("GITHUB_ID")
	conf.Secret = os.Getenv("GITHUB_SECRET")
	if conf.Id == "" || conf.Secret == "" {
		fmt.Println("GITHUB_ID/GITHUB_SECRET not found. Please check your environment")
		os.Exit(1)
	}
	if host := os.Getenv("APP_HOST"); host != "" {
		conf.OAuthHost = host
	}
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
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
		admin.GET("/category/:cat_id", um.Category)
		admin.GET("/category/:cat_id/difficulty/:difficulty", um.Category)
		admin.POST("/login", um.LoginPost)
		admin.POST("/segment/edit/:id", um.WebAuth(), um.EditSegment)
		admin.POST("/segment/delete/:id", um.WebAuth(), um.DeleteSegment)
		admin.POST("/segment/add", um.WebAuth(), um.AddSegment)
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
		v1.GET("/categories/:id", um.Auth(false), um.getCategory)
		// v1.POST("/categories", um.Auth(true), um.addCategory)
		// v1.PUT("/categories/:id", um.Auth(true), um.editCategory)
		// v1.POST("/categories/:id/status", um.Auth(true), um.approveCategory)
		// v1.DELETE("/categories/:id", um.Auth(true), um.deleteCategory)
		v1.GET("/languages", um.Auth(false), um.getLanguages)
	}
	v2 := r.Group("/v2")
	{
		r, err := repo.New("securityfirst", "tent-content")
		if err != nil {
			log.Fatalf("Repo error: %s", err)
		}

		o := tent.New(r)
		o.Register(v2, conf)
	}
	r.Run(":8080")

}
