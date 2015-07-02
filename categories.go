package main

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func (um *Umbrella) getCategories(c *gin.Context) {
	categories, err := um.getAllPublishedCategories(c)
	um.checkErr(c, err)
	c.JSON(200, categories)
}

func (um *Umbrella) getCategory(c *gin.Context) {
	categoryId := to.Int64(c.Params.ByName("id"))
	if categoryId != 0 {
		category, err := um.getCategoryById(c, categoryId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
		}
		c.JSON(200, category)
		return
	} else {
		c.JSON(404, gin.H{"error": "Not found"})
	}
}

func (um *Umbrella) getCategoryByParent(c *gin.Context) {
	categoryId := to.Int64(c.Params.ByName("id"))
	if categoryId != 0 {
		category, err := um.getAllPublishedCategoriesByParent(c, categoryId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
		}
		c.JSON(200, category)
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}

func (um *Umbrella) addCategory(c *gin.Context) {
	var json CategoryInsert
	c.Bind(&json)
	if true {
		user := c.MustGet("user").(User)
		category := CategoryInsert{Category: json.Category, Parent: json.Parent, Status: "submitted", CreatedAt: time.Now().Unix(), Author: user.Id}
		if user.Role == 1 {
			category.Status = "published"
			category.ApprovedBy = user.Id
			category.ApprovedAt = time.Now().Unix()
		}
		err := um.Db.Insert(&category)
		um.checkErr(c, err)
		c.JSON(200, category)
		return
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func (um *Umbrella) editCategory(c *gin.Context) {
	var json JSONCategory
	categoryId := to.Int64(c.Params.ByName("id"))
	if c.EnsureBody(&json) && categoryId != 0 {
		category, err := um.getCategoryById(c, categoryId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
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
			_, err := um.Db.Update(&category)
			um.checkErr(c, err)
			c.JSON(200, category)
			return
		}
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func (um *Umbrella) approveCategory(c *gin.Context) {
	var json JSONCategory
	c.Bind(&json)
	categoryId := to.Int64(c.Params.ByName("id"))
	if categoryId != 0 {
		category, err := um.getCategoryById(c, categoryId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not Found"})
				return
			}
			um.checkErr(c, err)
		}
		category.Status = json.Status
		if category.Status == "published" {
			category.ApprovedAt = time.Now().Unix()
			category.ApprovedBy = c.MustGet("user").(User).Id
		}
		_, err = um.Db.Update(&category)
		um.checkErr(c, err)
		c.JSON(200, category)
		return
	}
	c.JSON(400, gin.H{"error": "One or more parameters are missing"})
}

func (um *Umbrella) deleteCategory(c *gin.Context) {
	categoryId := to.Int64(c.Params.ByName("id"))
	if categoryId != 0 {
		category, err := um.getCategoryById(c, categoryId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
		}
		_, err = um.Db.Exec("delete from segments where category=?", category.Id)
		um.checkErr(c, err)
		_, err = um.Db.Delete(&category)
		um.checkErr(c, err)
		c.Writer.WriteHeader(204)
		return
	} else {
		c.JSON(404, gin.H{"error": "Not found"})
	}
}
