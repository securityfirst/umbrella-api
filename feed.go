package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func (um *Umbrella) getFeed(c *gin.Context) {
	since := to.Int64(c.Request.URL.Query().Get("since"))
	country := um.getCountry(c.Request.URL.Query().Get("country"))
	sources := strings.Split(c.Request.URL.Query().Get("sources"), ",")
	feedItems, err := um.getFeedItems(sources, country, since)
	um.checkErr(c, err)
	if err == nil {
		c.JSON(200, feedItems)
		return
	}
	c.JSON(400, gin.H{"error": err.Error()})
}
