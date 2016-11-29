package main

import (
	"fmt"
	"github.com/securityfirst/umbrella-api/models"

	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (um *Umbrella) loginEndpoint(c *gin.Context) {
	var json struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.Bind(&json); err == nil {
		var u models.User
		err := um.Db.SelectOne(&u, "select id, name, email, password, token, role from users where email=?", json.Email)
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
		err1 := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(json.Password))
		fmt.Println(err1)
		if err1 == nil {
			u.Token = randString(50)
			count, err := um.Db.Update(&u)
			um.checkErr(c, err)
			if err == nil && count == 1 {
				c.JSON(200, gin.H{"token": u.Token, "profile": u})
				return
			}
		}
	}
	c.JSON(401, gin.H{"error": "Email or password incorrect. Please try again"})
}

func (um *Umbrella) loginCheck(c *gin.Context) {
	user, err := um.checkUser(c)
	um.checkErr(c, err)
	loggedIn := user.Id != 0 && err == nil
	c.JSON(200, gin.H{"response": loggedIn})
}
