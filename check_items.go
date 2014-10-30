package main

import "github.com/gin-gonic/gin"

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
