package main

import (
	"log"
	"time"

	"github.com/securityfirst/umbrella-api/feed"

	"strings"

	"github.com/securityfirst/umbrella-api/models"
	"github.com/securityfirst/umbrella-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func (um *Umbrella) UpdateFeeds() error {
	countries, err := um.getCountries()
	if err != nil {
		return err
	}
	var feeds = make([]feed.Fetcher, 0, len(countries)+4)
	for i := range countries {
		feeds = append(feeds, &feed.ReliefWebFetcher{
			Country: &countries[i],
		})
	}
	feeds = append(feeds, &feed.GdascFetcher{},
		&feed.FCOFetcher{}, &feed.CadataFetcher{}, &feed.CDCFetcher{})
	for _, f := range feeds {
		results, err := f.Fetch()
		if err != nil {
			utils.CheckErr(err)
			continue
		}
		for _, item := range results {
			item.Update(um.Db)
		}
	}
	return nil
}

func (um *Umbrella) getFeed(c *gin.Context) {
	since := to.Int64(c.Request.URL.Query().Get("since"))
	// twoMonthsAgo := time.Now().AddDate(0, -2, 0).Unix()
	// if since < twoMonthsAgo {
	// 	since = twoMonthsAgo
	// }
	country, err := um.getCountryInfo(c.Request.URL.Query().Get("country"))
	um.checkErr(c, err)
	if err != nil {
		log.Println("country", err)
		c.JSON(200, []interface{}{})
		return
	}
	sources := strings.Split(c.Request.URL.Query().Get("sources"), ",")
	feedItems, err := um.getDbFeedItems(sources, country.Iso2, since)
	um.checkErr(c, err)
	if err != nil {
		log.Println("sources", err)
		c.JSON(200, feedItems)
		return
	}
	feedLog := models.FeedRequestLog{
		Country:   country.Iso2,
		Sources:   c.Request.URL.Query().Get("sources"),
		CheckedAt: time.Now().Unix(),
	}
	if err != nil {
		utils.CheckErr(err)
		feedLog.Status = 500
		go utils.CheckErr(um.Db.Insert(&feedLog))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	feedLog.Status = 200
	go utils.CheckErr(um.Db.Insert(&feedLog))
	c.JSON(200, feedItems)
}

type SortFeedByDate []models.FeedItem

func (slice SortFeedByDate) Len() int {
	return len(slice)
}

func (slice SortFeedByDate) Less(i, j int) bool {
	return slice[i].UpdatedAt > slice[j].UpdatedAt
}

func (slice SortFeedByDate) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
