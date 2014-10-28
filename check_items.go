package main

import "github.com/gin-gonic/gin"

type CheckItem struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	Value    int    `json:"value"`
	Parent   int    `json:"parent"`
	Category int    `json:"category"`
}

type Segment struct {
	Title    string `json:"title"`
	SubTitle string `json:"subtitle"`
	Body     string `json:"body"`
	Category int    `json:"category"`
}

func getCheckItems(c *gin.Context) {
	var checkItems []CheckItem
	checkItems = append(checkItems, CheckItem{Title: "First", Text: "", Value: 0, Parent: 0, Category: 1})
	checkItems = append(checkItems, CheckItem{Title: "Second", Text: "", Value: 0, Parent: 0, Category: 1})
	checkItems = append(checkItems, CheckItem{Title: "Third", Text: "", Value: 0, Parent: 0, Category: 1})
	checkItems = append(checkItems, CheckItem{Title: "Fourth", Text: "", Value: 0, Parent: 0, Category: 1})

	c.JSON(200, checkItems)
}

func getSegments(c *gin.Context) {
	var segments []Segment
	segments = append(segments, Segment{Title: "Title1", SubTitle: "Subtitle1", Body: "Body1", Category: 1})
	segments = append(segments, Segment{Title: "Title2", SubTitle: "Subtitle2", Body: "Body2", Category: 1})
	segments = append(segments, Segment{Title: "Title3", SubTitle: "Subtitle3", Body: "Body3", Category: 1})
	segments = append(segments, Segment{Title: "Title4", SubTitle: "Subtitle4", Body: "Body4", Category: 1})

	c.JSON(200, segments)
}
