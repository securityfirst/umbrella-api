package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/securityfirst/umbrella-api/utils"
	"gopkg.in/securityfirst/tent.v1"
	"gopkg.in/securityfirst/tent.v1/auth"
	"gopkg.in/securityfirst/tent.v1/repo"

	tent2 "gopkg.in/securityfirst/tent.v2"
	auth2 "gopkg.in/securityfirst/tent.v2/auth"
	repo2 "gopkg.in/securityfirst/tent.v2/repo"

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
	conf.RandomString = os.Getenv("APP_SECRET")
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	um := getUmbrella()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080", "https://api.secfirst.org"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/v1")
	{
		v1.GET("/feed", um.getFeed)
		v1.GET("/segments", um.getSegments)
		v1.GET("/check_items", um.Auth(false), um.getCheckItems)
		v1.GET("/categories", um.Auth(false), um.getCategories)
		v1.GET("/categories/:id", um.Auth(false), um.getCategory)
		v1.GET("/languages", um.Auth(false), um.getLanguages)
	}
	v2 := r.Group("/v2")
	{
		r, err := repo.New("securityfirst", "tent-content", "master")
		if err != nil {
			log.Fatalf("Repo error: %s", err)
		}
		o := tent.New(r)
		// No authentication
		v2.GET("/feed", um.getFeed)
		// Autentication
		o.Register(v2, conf)
	}
	v3 := r.Group("/v3")
	{
		r, err := repo2.New("securityfirst", "tent-content", "difficulty")
		if err != nil {
			log.Fatalf("Repo error: %s", err)
		}
		o := tent2.New(r)
		// No authentication
		v3.GET("/feed", um.getFeed)
		// Autentication

		o.Register(v3, auth2.Config{
			ID:        conf.Id,
			Secret:    conf.Secret,
			OAuthHost: conf.OAuthHost,
			Host:      conf.Host,
			State:     conf.RandomString,
			Login:     auth2.HandleConf{conf.Login.Endpoint, conf.Login.Redirect},
			Logout:    auth2.HandleConf{conf.Logout.Endpoint, conf.Logout.Redirect},
			Callback:  auth2.HandleConf{conf.Callback.Endpoint, conf.Callback.Redirect},
		})
	}
	r.Run(":8080")

}
