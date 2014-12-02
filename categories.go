package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func getCategories(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	categories, err := getAllPublishedCategories(c, dbmap)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, categories)
}

func getCategory(c *gin.Context) {
	categoryId := to.Int64(c.Params.ByName("id"))
	if categoryId != 0 {
		dbmap := initDb()
		defer dbmap.Db.Close()
		category, err := getCategoryById(c, dbmap, categoryId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, category)
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}

func getCategoryByParent(c *gin.Context) {
	categoryId := to.Int64(c.Params.ByName("id"))
	if categoryId != 0 {
		dbmap := initDb()
		defer dbmap.Db.Close()
		category, err := getAllPublishedCategoriesByParent(c, dbmap, categoryId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, category)
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}

func addCategory(c *gin.Context) {
	var json Category
	dbmap := initDb()
	defer dbmap.Db.Close()
	c.Bind(&json)
	fmt.Println(json)
	if true {
		user := c.MustGet("user").(User)
		category := Category{Category: json.Category, Parent: json.Parent, Status: "submitted", CreatedAt: time.Now().Unix(), Author: user.Id}
		if user.Role == 1 {
			category.Status = "published"
		}
		err := dbmap.Insert(&category)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, category)
		return
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func editCategory(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	var json JSONCategory
	categoryId := to.Int64(c.Params.ByName("id"))
	if c.EnsureBody(&json) && categoryId != 0 {
		category, err := getCategoryById(c, dbmap, categoryId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if json.Category != "" {
			user := c.MustGet("user").(User)
			category.Category = json.Category
			category.Parent = to.Int64(json.Parent)
			if user.Role != 1 {
				category.Status = "submitted"
			}
			category.CreatedAt = time.Now().Unix()
			category.Author = user.Id
			count, err := dbmap.Update(&category)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			} else if count < 1 {
				c.JSON(500, gin.H{"error": "Could not update the resource"})
			}
			c.JSON(200, category)
			return
		}
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func approveCategory(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	var json JSONCategory
	c.Bind(&json)
	categoryId := to.Int64(c.Params.ByName("id"))
	if categoryId != 0 {
		category, err := getCategoryById(c, dbmap, categoryId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		category.Status = json.Status
		user := c.MustGet("user").(User)
		if category.Status == "published" {
			category.ApprovedAt = time.Now().Unix()
			category.ApprovedBy = user.Id
		}
		_, err = dbmap.Update(&category)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, category)
		return
	}
	c.JSON(400, gin.H{"error": "One or more parameters are missing"})
}

func deleteCategory(c *gin.Context) {
	categoryId := to.Int64(c.Params.ByName("id"))
	if categoryId != 0 {
		dbmap := initDb()
		defer dbmap.Db.Close()
		category, err := getCategoryById(c, dbmap, categoryId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		_, err = dbmap.Delete(&category)
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
