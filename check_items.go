package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func getCheckItems(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	checkItems, err := getAllPublishedCheckItems(c, dbmap)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, checkItems)
}

func getCheckItem(c *gin.Context) {
	checkItemId := to.Int64(c.Params.ByName("id"))
	if checkItemId != 0 {
		dbmap := initDb()
		defer dbmap.Db.Close()
		checkItem, err := getCheckItemById(c, dbmap, checkItemId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, checkItem)
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}

func addCheckItem(c *gin.Context) {
	var json CheckItem
	dbmap := initDb()
	defer dbmap.Db.Close()
	c.Bind(&json)
	fmt.Println(json)
	if true {
		user := c.MustGet("user").(User)
		checkItem := CheckItem{Title: json.Title, Text: json.Text, Value: json.Value, Parent: json.Parent, Category: json.Category, Status: "submitted", CreatedAt: time.Now().Unix(), Author: user.Id}
		err := dbmap.Insert(&checkItem)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, checkItem)
		return
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func editCheckItem(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	var json JSONCheckItem
	checkItemId := to.Int64(c.Params.ByName("id"))
	if c.EnsureBody(&json) && checkItemId != 0 {
		checkItem, err := getCheckItemById(c, dbmap, checkItemId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if json.Title != "" || json.Text != "" || json.Category != 0 {
			if json.Title != "" {
				checkItem.Title = json.Title
			}
			if json.Text != "" {
				checkItem.Text = json.Text
			}
			checkItem.Value = to.Int64(json.Value)
			checkItem.Parent = to.Int64(json.Parent)
			if json.Category != 0 {
				checkItem.Category = to.Int64(json.Category)
			}
			checkItem.Status = "submitted"
			checkItem.CreatedAt = time.Now().Unix()
			user := c.MustGet("user").(User)
			checkItem.Author = user.Id
			count, err := dbmap.Update(&checkItem)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			} else if count < 1 {
				c.JSON(500, gin.H{"error": "Could not update the resource"})
			}
			c.JSON(200, checkItem)
			return
		}
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func approveCheckItem(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	var json JSONCheckItem
	c.Bind(&json)
	checkItemId := to.Int64(c.Params.ByName("id"))
	if checkItemId != 0 {
		checkItem, err := getCheckItemById(c, dbmap, checkItemId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		checkItem.Status = json.Status
		user := c.MustGet("user").(User)
		if checkItem.Status == "published" {
			checkItem.ApprovedAt = time.Now().Unix()
			checkItem.ApprovedBy = user.Id
		}
		_, err = dbmap.Update(&checkItem)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, checkItem)
		return
	}
	c.JSON(400, gin.H{"error": "One or more parameters are missing"})
}

func deleteCheckItem(c *gin.Context) {
	checkItemId := to.Int64(c.Params.ByName("id"))
	if checkItemId != 0 {
		dbmap := initDb()
		defer dbmap.Db.Close()
		checkItem, err := getCheckItemById(c, dbmap, checkItemId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		_, err = dbmap.Delete(&checkItem)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Writer.WriteHeader(204)
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}
