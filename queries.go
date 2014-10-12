package main

import (
	"database/sql"

	"unicode/utf8"

	"github.com/coopernurse/gorp"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "apiuser:mEYP4JKdZeeZVbj5@tcp(localhost:1234)/umbrella?charset=utf8")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	return dbmap
}

func checkUser(c *gin.Context, dbmap *gorp.DbMap) (User, error) {

	var user User
	var err error
	token := c.Request.Header.Get("token")
	if utf8.RuneCountInString(token) > 0 {
		err = dbmap.SelectOne(&user, "select id, name, email, password, token from users where token=?", token)
	}
	return user, err
}
