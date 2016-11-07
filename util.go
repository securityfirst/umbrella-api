package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

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
	dbmap.AddTableWithName(FeedLastChecked{}, "feed_last_checked").SetKeys(true, "Id")
	dbmap.AddTableWithName(FeedItem{}, "feed_items").SetKeys(true, "Id")
	dbmap.AddTableWithName(FeedRequestLog{}, "feed_log").SetKeys(true, "Id")
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
		if err != nil {
			info := color.New(color.FgGreen).SprintFunc()
			pc, fn, line, _ := runtime.Caller(1)
			log.Printf(info(fmt.Sprintf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)))
		}
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
			c.AbortWithError(401, errors.New("Not Authorized"))
		}
	}
}

func traceDb(dbmap *gorp.DbMap) {
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
}

func colorLog(toLog interface{}, col ...color.Attribute) {
	info := color.New(col...).SprintFunc()
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf(info(fmt.Sprintf("%s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, fmt.Sprint(toLog))))
}

func makeRequest(uri string, method string, requestBody io.Reader) (response []byte, err error) {
	req, err := http.NewRequest(strings.ToUpper("GET"), uri, requestBody)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		checkErr(err)
	}
	return response, err
}

func difference(slice1 []int, slice2 []int) []int {
	var diff []int
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

func (um *Umbrella) WebAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("token")
		colorLog(fmt.Sprintf("%+v, %+v", cookie, err), color.FgCyan)
		if err != nil || cookie.Value == "" {
			c.Redirect(302, "/admin/login")
			c.Abort()
		} else {
			if user, err := um.checkWebUser(cookie.Value); err == nil {
				c.Set("user", user)
			} else {
				c.Redirect(302, "/admin/login")
				c.Abort()
			}
		}
	}
}

func (u *User) setCookie(c *gin.Context) {
	expiration := time.Now().Add(time.Hour * 24 * 365)
	cookie := http.Cookie{Name: "token", Value: u.Token, Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
}

func (u *User) removeCookie(c *gin.Context) {
	expiration := time.Now().Add(time.Hour * -1)
	cookie := http.Cookie{Name: "token", Value: "", Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
}

func buildFips() map[string]string {
	const fileName = "fips.csv"

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Cannot open %s: %s", fileName, err)
	}
	var m = make(map[string]string, 265)
	r := csv.NewReader(f)
	r.Comma = '	'
	for count := 0; ; count++ {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Cannot red %s: %s", fileName, err)
		}
		m[record[0]] = record[1]
	}
	return m
}
