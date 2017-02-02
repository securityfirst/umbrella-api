package main

import (
	"database/sql"

	"github.com/gosexy/to"

	"github.com/gin-gonic/gin"
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
