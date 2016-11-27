package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gosexy/to"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (um *Umbrella) Index(c *gin.Context) {
	menuStruct, err := um.getAllPublishedCategories(c)
	checkErr(err)
	obj := gin.H{
		"title":  "Umbrella Dashboard",
		"menu":   menuStruct,
		"diff":   "",
		"cat_id": 0,
	}
	c.HTML(200, "index.tmpl", obj)
}

func (um *Umbrella) Login(c *gin.Context) {
	obj := gin.H{"title": "Login"}
	c.HTML(http.StatusOK, "login.tmpl", obj)
}

func (um *Umbrella) Category(c *gin.Context) {
	menuStruct, err := um.getAllPublishedCategories(c)
	checkErr(err)
	obj := gin.H{
		"title": "Category",
		"menu":  menuStruct,
	}
	if catId, ok := c.Params.Get("cat_id"); ok {
		obj["cat_id"] = to.Int64(catId)
	} else {
		obj["cat_id"] = int64(0)
	}
	if diff, ok := c.Params.Get("difficulty"); ok {
		obj["diff"] = diff
	} else {
		obj["diff"] = ""
	}
	c.HTML(http.StatusOK, "index.tmpl", obj)
}

func (um *Umbrella) LoginPost(c *gin.Context) {
	var login LoginForm
	var err error
	c.BindWith(&login, binding.Form)
	_, err = govalidator.ValidateStruct(login)
	if err != nil {
		obj := gin.H{"title": "Login", "error": err}
		c.HTML(http.StatusBadRequest, "login.tmpl", obj)
		return
	}
	var u User
	err = um.Db.SelectOne(&u, "select id, name, email, password, token, role from users where email=?", login.Email)
	if err != nil {
		fmt.Println(err)
		match, _ := regexp.MatchString("connection refused", err.Error())
		if match {
			um.checkErr(c, err)
		} else {
			c.JSON(401, gin.H{"error": "Email or password incorrect. Please try again"})
		}
		return
	}
	err1 := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(login.Password))
	fmt.Println(err1)
	if err1 == nil {
		u.Token = randString(50)
		_, err := um.Db.Update(&u)
		checkErr(err)
		c.Set("user", u)
		u.setCookie(c)
		c.Redirect(302, "/admin")
		return

	}
	obj := gin.H{"title": "Login", "error": err, "email": login.Email}
	c.HTML(http.StatusBadRequest, "login.tmpl", obj)
}

// LoginForm is a model that binds email login
type LoginForm struct {
	Email        string `form:"email" valid:"email,required"`
	Password     string `form:"password" valid:"valid_password,required"`
	PasswordHash string `form:"-"`
}

func (um *Umbrella) LogOut(c *gin.Context) {
	u := c.MustGet("user").(User)
	u.removeCookie(c)
	c.Redirect(302, "/admin/login")
	return
}
