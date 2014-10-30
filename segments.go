package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func getSegments(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	segments, err := getAllPublishedSegments(c, dbmap)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, segments)
}

func addSegment(c *gin.Context) {
	var json Segment
	dbmap := initDb()
	defer dbmap.Db.Close()
	fmt.Println(json)
	if c.EnsureBody(&json) {
		segment := Segment{Title: json.Title, SubTitle: json.SubTitle, Body: json.Body, Category: json.Category}
		err := dbmap.Insert(&segment)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, segment)
		return
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func editSegment(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	var json EditSegment
	if c.EnsureBody(&json) {
		segment, err := getSegmentById(c, dbmap, json.Id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if json.Title != "" || json.SubTitle != "" || json.Body != "" || json.Category != 0 {
			if json.Title != "" {
				segment.Title = json.Title
			}
			if json.SubTitle != "" {
				segment.SubTitle = json.SubTitle
			}
			if json.Body != "" {
				segment.Body = json.Body
			}
			if json.Category != 0 {
				segment.Category = json.Category
			}
			count, err := dbmap.Update(&segment)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			} else if count < 1 {
				c.JSON(500, gin.H{"error": "Could not update the resource"})
			}
			c.JSON(200, segment)
			return
		}
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func getSegment(c *gin.Context) {
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 {
		dbmap := initDb()
		defer dbmap.Db.Close()
		segment, err := getSegmentById(c, dbmap, segmentId)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, segment)
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}

func deleteSegment(c *gin.Context) {
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 {
		dbmap := initDb()
		defer dbmap.Db.Close()
		segment, err := getSegmentById(c, dbmap, segmentId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		_, err = dbmap.Delete(&segment)
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
