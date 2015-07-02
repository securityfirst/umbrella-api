package main

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
)

func getUmbrella() (umbrella Umbrella) {
	umbrella = Umbrella{
		Db: initDb(),
	}
	if !isProduction {
		umbrella.Db.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	}
	return umbrella
}

func initDb() *gorp.DbMap {
	var connString string
	if isProduction {
		connString = os.Getenv("DB_PROD")
	} else {
		connString = os.Getenv("DB_DEV")
	}
	db, err := sql.Open("mysql", connString)
	if isProduction {
		db.SetMaxIdleConns(2000)
		db.SetMaxOpenConns(2000)
	}
	checkErr(err)
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Segment{}, "segments").SetKeys(true, "Id")
	dbmap.AddTableWithName(CheckItem{}, "check_items").SetKeys(true, "Id")
	dbmap.AddTableWithName(Category{}, "categories").SetKeys(true, "Id")
	dbmap.AddTableWithName(CategoryInsert{}, "categories").SetKeys(true, "Id")
	// dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	return dbmap
}

func (um *Umbrella) JSON(c *gin.Context, code int, obj interface{}) {
	if !c.Writer.Written() {
		if obj != nil {
			c.JSON(code, obj)
		} else {
			c.Writer.WriteHeader(code)
			c.Writer.Write([]byte(""))
			obj = gin.H{"": ""}
		}
	}
}

func (um *Umbrella) HTML(c *gin.Context, code int, name string, obj interface{}) {
	if !c.Writer.Written() {
		c.HTML(code, name, obj)
	}
}

func checkErr(err error) {
	if err != nil {
		info := color.New(color.FgGreen).SprintFunc()
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf(info(fmt.Sprintf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)))
	}
}

func redLog(toLog interface{}) {
	info := color.New(color.FgRed).SprintFunc()
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf(info(fmt.Sprintf("%s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, fmt.Sprint(toLog))))
}

func (um *Umbrella) checkErr(c *gin.Context, err error) {
	if err != nil && !c.Writer.Written() {
		checkErr(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		c.Abort()
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
func (um *Umbrella) Auth(strict bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err1 := um.checkUser(c)
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

func traceDb(dbmap *gorp.DbMap) {
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
}
