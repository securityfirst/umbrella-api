package models

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id       int64  `db:"id" json:"-"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `json:"-"`
	Token    string `json:"-"`
	Role     int    `db:"role" json:"-"`
}

func (u *User) SetCookie(c *gin.Context) {
	expiration := time.Now().Add(time.Hour * 24 * 365)
	cookie := http.Cookie{Name: "token", Value: u.Token, Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
}

func (u *User) RemoveCookie(c *gin.Context) {
	expiration := time.Now().Add(time.Hour * -1)
	cookie := http.Cookie{Name: "token", Value: "", Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
}
