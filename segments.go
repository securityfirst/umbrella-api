package main

import (
	"fmt"
	"strings"
	"time"

	"code.google.com/p/go.net/html"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func getSegmentsRaw(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	segmentList, err := getAllPublishedSegments(c, dbmap)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, segmentList)
}

func getSegments(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	segmentList, err := getAllPublishedSegments(c, dbmap)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, parseForSegments(segmentList))
}

func parseForSegments(segmentList []Segment) []Segment {
	var segments []Segment
	for i := 0; i < len(segmentList); i++ {
		doc, err := html.Parse(strings.NewReader(segmentList[i].Body))
		if err != nil {
			fmt.Println(err.Error())
		}
		relSegment := false
		var f func(*html.Node)
		f = func(n *html.Node) {
			if relSegment {
				newSeg := Segment{Id: int64(len(segments)), Title: segmentList[i].Title, SubTitle: segmentList[i].SubTitle, Body: n.Data, Category: segmentList[i].Category}
				segments = append(segments, newSeg)
				relSegment = false
			}
			if n.Type == html.ElementNode && n.Data == "p" {
				relSegment = true
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
	}
	return segments
}

func getSegmentsByCat(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	category := to.Int64(c.Params.ByName("id"))
	segmentList, err := getAllPublishedSegmentsByCat(c, dbmap, category)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, parseForSegments(segmentList))
}

func getSegmentsRawByCat(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	category := to.Int64(c.Params.ByName("id"))
	segmentList, err := getAllPublishedSegmentsByCat(c, dbmap, category)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, segmentList)
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
	fmt.Println(json)
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

func editSegmentByCat(c *gin.Context) {
	dbmap := initDb()
	defer dbmap.Db.Close()
	var json Segment
	c.Bind(&json)
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 && (json.Title != "" || json.SubTitle != "" || json.Body != "" || json.Category != 0) {
		fmt.Println(segmentId)
		segment, err := getSegmentByCatId(c, dbmap, segmentId)
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
		// segment.Status = "submitted"
		segment.CreatedAt = time.Now().Unix()
		user := c.MustGet("user").(User)
		segment.Author = user.Id
		segment.Id = 0
		err = dbmap.Insert(&segment)
		// _, err = dbmap.Update(&segment)
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
