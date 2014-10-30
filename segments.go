package main

import (
	"time"

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

func getSegment(c *gin.Context) {
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
		c.JSON(200, segment)
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}

func addSegment(c *gin.Context) {
	var json Segment
	dbmap := initDb()
	defer dbmap.Db.Close()
	if c.EnsureBody(&json) {
		user := c.MustGet("user").(User)
		segment := Segment{Title: json.Title, SubTitle: json.SubTitle, Body: json.Body, Category: json.Category, Status: "submitted", CreatedAt: time.Now().Unix(), Author: user.Id}
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
	var json Segment
	c.Bind(&json)
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 && (json.Title != "" || json.SubTitle != "" || json.Body != "" || json.Category != 0) {
		segment, err := getSegmentById(c, dbmap, segmentId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
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
			segment.Category = to.Int64(json.Category)
		}
		segment.Status = "submitted"
		segment.CreatedAt = time.Now().Unix()
		user := c.MustGet("user").(User)
		segment.Author = user.Id
		_, err = dbmap.Update(&segment)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, segment)
		return
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func approveSegment(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	var json JSONSegment
	c.Bind(&json)
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 {
		segment, err := getSegmentById(c, dbmap, segmentId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				c.JSON(404, gin.H{"error": "Requested resource could not be found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		segment.Status = json.Status
		user := c.MustGet("user").(User)
		if segment.Status == "published" {
			segment.ApprovedAt = time.Now().Unix()
			segment.ApprovedBy = user.Id
		}
		_, err = dbmap.Update(&segment)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, segment)
		return
	}
	c.JSON(400, gin.H{"error": "One or more parameters are missing"})
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
		c.Writer.WriteHeader(204)
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}
