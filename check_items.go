package main

import "github.com/gin-gonic/gin"

func (um *Umbrella) getCheckItems(c *gin.Context) {
	checkItems, err := um.getAllPublishedCheckItems(c)
	um.checkErr(c, err)
	c.JSON(200, checkItems)
}
