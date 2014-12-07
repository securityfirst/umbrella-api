package main

import (
	"fmt"

	"regexp"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gin-gonic/gin"
)

func loginEndpoint(c *gin.Context) {
	var json LoginJSON

	if c.EnsureBody(&json) {
		dbmap := initDb()
		defer dbmap.Db.Close()

		var u User
		err := dbmap.SelectOne(&u, "select id, name, email, password, token, role from users where email=?", json.Email)
		if err != nil {
			fmt.Println(err)
			match, _ := regexp.MatchString("connection refused", err.Error())
			if match {
				c.JSON(500, gin.H{"error": "Internal server error"})
			} else {
				c.JSON(401, gin.H{"error": "Email or password incorrect. Please try again"})
			}
			return
		}
		err1 := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(json.Password))
		fmt.Println(err1)
		if err1 == nil {
			u.Token = randString(50)
			count, err := dbmap.Update(&u)
			fmt.Println(err)
			if err == nil && count == 1 {
				c.JSON(200, gin.H{"token": u.Token, "profile": u})
				return
			}
		}
	}
	c.JSON(401, gin.H{"error": "Email or password incorrect. Please try again"})
}

func loginCheck(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()

	user, err1 := checkUser(c, dbmap)
	fmt.Println(err1)
	loggedIn := user.Id != 0 && err1 == nil
	c.JSON(200, gin.H{"response": loggedIn})
}
