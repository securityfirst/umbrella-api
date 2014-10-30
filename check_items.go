package main

import (
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
	if c.EnsureBody(&json) {
		checkItem := CheckItem{Title: json.Title, Text: json.Text, Value: json.Value, Parent: json.Parent, Category: json.Category}
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
	var json EditCheckItem
	checkItemId := to.Int64(c.Params.ByName("id"))
	if c.EnsureBody(&json) && checkItemId != 0 {
		checkItem, err := getCheckItemById(c, dbmap, checkItemId)
		if err != nil {
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
		c.JSON(204, gin.H{})
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}
