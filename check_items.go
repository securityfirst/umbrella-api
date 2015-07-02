package main

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func (um *Umbrella) getCheckItems(c *gin.Context) {
	checkItems, err := um.getAllPublishedCheckItems(c)
	um.checkErr(c, err)
	c.JSON(200, checkItems)
}

func (um *Umbrella) getCheckItemsByCat(c *gin.Context) {
	categoryId := to.Int64(c.Params.ByName("id"))
	if categoryId != 0 {
		checkItems, err := um.getAllPublishedCheckItemsByCat(c, categoryId)
		um.checkErr(c, err)
		c.JSON(200, checkItems)
		return
	}
	c.JSON(404, gin.H{"error": "Not found"})
}

func (um *Umbrella) getCheckItem(c *gin.Context) {
	checkItemId := to.Int64(c.Params.ByName("id"))
	if checkItemId != 0 {
		checkItem, err := um.getCheckItemById(c, checkItemId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
		}
		c.JSON(200, checkItem)
		return
	}
	c.JSON(404, gin.H{"error": "Requested resource could not be found"})
}

func (um *Umbrella) addCheckItem(c *gin.Context) {
	var json CheckItem
	c.Bind(&json)
	user := c.MustGet("user").(User)
	checkItem := CheckItem{Title: json.Title, Text: json.Text, Value: json.Value, Parent: json.Parent, Category: json.Category, Status: "submitted", CreatedAt: time.Now().Unix(), Author: user.Id}
	err := um.Db.Insert(&checkItem)
	um.checkErr(c, err)
	c.JSON(200, checkItem)
}

func (um *Umbrella) editCheckItem(c *gin.Context) {
	var json JSONCheckItem
	checkItemId := to.Int64(c.Params.ByName("id"))
	if c.EnsureBody(&json) && checkItemId != 0 {
		checkItem, err := um.getCheckItemById(c, checkItemId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
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
			_, err := um.Db.Update(&checkItem)
			um.checkErr(c, err)
			c.JSON(200, checkItem)
			return
		}
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func (um *Umbrella) approveCheckItem(c *gin.Context) {
	var json JSONCheckItem
	c.Bind(&json)
	checkItemId := to.Int64(c.Params.ByName("id"))
	if checkItemId != 0 {
		checkItem, err := um.getCheckItemById(c, checkItemId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
		}
		checkItem.Status = json.Status
		user := c.MustGet("user").(User)
		if checkItem.Status == "published" {
			checkItem.ApprovedAt = time.Now().Unix()
			checkItem.ApprovedBy = user.Id
		}
		_, err = um.Db.Update(&checkItem)
		um.checkErr(c, err)
		c.JSON(200, checkItem)
		return
	}
	c.JSON(400, gin.H{"error": "One or more parameters are missing"})
}

func (um *Umbrella) deleteCheckItem(c *gin.Context) {
	checkItemId := to.Int64(c.Params.ByName("id"))
	if checkItemId != 0 {
		checkItem, err := um.getCheckItemById(c, checkItemId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
		}
		_, err = um.Db.Delete(&checkItem)
		um.checkErr(c, err)
		c.Writer.WriteHeader(204)
		return
	}
	c.JSON(404, gin.H{"error": "Requested resource could not be found"})
}
