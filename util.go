package main

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbols := big.NewInt(int64(len(alphanum)))
	states := big.NewInt(0)
	states.Exp(symbols, big.NewInt(int64(n)), nil)
	r, err := rand.Int(rand.Reader, states)
	if err != nil {
		panic(err)
	}
	var bytes = make([]byte, n)
	r2 := big.NewInt(0)
	symbol := big.NewInt(0)
	for i := range bytes {
		r2.DivMod(r, symbols, symbol)
		r, r2 = r2, r
		bytes[i] = alphanum[symbol.Int64()]
	}
	return string(bytes)
}

// Auth Middleware
func Auth(strict bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbmap := initDb()
		defer dbmap.Db.Close()

		user, err1 := checkUser(c, dbmap)
		if err1 != nil {
			user = User{Id: 0}
		}
		c.Set("user", user)
		if strict && user.Id == 0 {
			c.JSON(401, gin.H{"error": "Not Authorized"})
			c.Fail(401, errors.New("Not Authorized"))
		}
	}
}
