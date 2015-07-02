package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"code.google.com/p/go.net/html"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func (um *Umbrella) getSegmentsRaw(c *gin.Context) {
	segmentList, err := um.getAllPublishedSegments(c)
	checkErr(err)
	um.checkErr(c, err)
	um.JSON(c, 200, segmentList)
}

func (um *Umbrella) getSegments(c *gin.Context) {
	segmentList, err := um.getAllPublishedSegments(c)
	checkErr(err)
	um.checkErr(c, err)
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

func (um *Umbrella) getSegmentsByCat(c *gin.Context) {
	category := to.Int64(c.Params.ByName("id"))
	if category != 0 {
		segmentList, err := um.getAllPublishedSegmentsByCat(c, category)
		um.checkErr(c, err)
		c.JSON(200, parseForSegments(segmentList))
		return
	}
	c.JSON(404, gin.H{"error": "Not found"})
}

func (um *Umbrella) getSegmentsRawByCat(c *gin.Context) {
	category := to.Int64(c.Params.ByName("id"))
	if category != 0 {
		segmentList, err := um.getAllPublishedSegmentsByCat(c, category)
		um.checkErr(c, err)
		c.JSON(200, segmentList)
		return
	}
	c.JSON(404, gin.H{"error": "Not found"})
}

func (um *Umbrella) getSegment(c *gin.Context) {
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 {
		segment, err := um.getSegmentById(c, segmentId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
		}
		c.JSON(200, segment)
		return
	} else {
		c.JSON(404, gin.H{"error": "Not found"})
	}
}

func (um *Umbrella) addSegment(c *gin.Context) {
	var json Segment
	if c.EnsureBody(&json) {
		user := c.MustGet("user").(User)
		segment := Segment{Title: json.Title, SubTitle: json.SubTitle, Body: json.Body, Category: json.Category, Status: "submitted", CreatedAt: time.Now().Unix(), Author: user.Id}
		if user.Role == 1 {
			segment.Status = "published"
			segment.ApprovedBy = user.Id
			segment.ApprovedAt = time.Now().Unix()
		}
		err := um.Db.Insert(&segment)
		um.checkErr(c, err)
		c.JSON(200, segment)
		return
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func (um *Umbrella) editSegment(c *gin.Context) {
	var json Segment
	c.Bind(&json)
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 && (json.Title != "" || json.SubTitle != "" || json.Body != "" || json.Category != 0) {
		segment, err := um.getSegmentById(c, segmentId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
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

		_, err = um.Db.Update(&segment)
		um.checkErr(c, err)
		c.JSON(200, segment)
		return
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func (um *Umbrella) editSegmentByCat(c *gin.Context) {
	var json Segment
	c.Bind(&json)
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 && (json.Title != "" || json.SubTitle != "" || json.Body != "" || json.Category != 0) {
		fmt.Println(segmentId)
		segment, err := um.getSegmentByCatId(c, segmentId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
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
		segment.CreatedAt = time.Now().Unix()
		user := c.MustGet("user").(User)
		segment.Author = user.Id
		segment.Id = 0
		err = um.Db.Insert(&segment)
		um.checkErr(c, err)
		c.JSON(200, segment)
		return
	}
	c.JSON(400, gin.H{"error": "One or several fields missing. Please check and try again"})
}

func (um *Umbrella) approveSegment(c *gin.Context) {
	var json JSONSegment
	c.Bind(&json)
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 {
		segment, err := um.getSegmentById(c, segmentId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
		}
		um.checkErr(c, err)
		segment.Status = json.Status
		user := c.MustGet("user").(User)
		if segment.Status == "published" {
			segment.ApprovedAt = time.Now().Unix()
			segment.ApprovedBy = user.Id
		}
		_, err = um.Db.Update(&segment)
		um.checkErr(c, err)
		c.JSON(200, segment)
		return
	}
	c.JSON(400, gin.H{"error": "One or more parameters are missing"})
}

func (um *Umbrella) deleteSegment(c *gin.Context) {
	segmentId := to.Int64(c.Params.ByName("id"))
	if segmentId != 0 {
		segment, err := um.getSegmentById(c, segmentId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Not found"})
				return
			}
			um.checkErr(c, err)
		}
		_, err = um.Db.Delete(&segment)
		um.checkErr(c, err)
		c.Writer.WriteHeader(204)
		return
	} else {
		c.JSON(404, gin.H{"error": "Requested resource could not be found"})
	}
}
