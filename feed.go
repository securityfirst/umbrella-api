package main

import (
	"errors"
	"log"
	"time"

	"github.com/securityfirst/umbrella-api/feed"

	"sort"
	"strconv"
	"strings"

	"github.com/securityfirst/umbrella-api/models"
	"github.com/securityfirst/umbrella-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func (um *Umbrella) getFeed(c *gin.Context) {
	feedItems := []models.FeedItem{}
	since := to.Int64(c.Request.URL.Query().Get("since"))
	twoMonthsAgo := time.Now().AddDate(0, -2, 0).Unix()
	if since < twoMonthsAgo {
		since = twoMonthsAgo
	}
	country, err := um.getCountryInfo(c.Request.URL.Query().Get("country"))
	um.checkErr(c, err)
	if err != nil {
		log.Println("country", err)
		c.JSON(200, feedItems)
		return
	}
	sources := strings.Split(c.Request.URL.Query().Get("sources"), ",")
	feedItems, err = um.getFeedItems(sources, country, since)
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

func (um *Umbrella) getFeedItems(sources []string, country models.Country, since int64) ([]models.FeedItem, error) {
	items := []models.FeedItem{}
	var err error
	getLastChecked := um.getLastChecked(country.Iso2)
	var cleanSources, toRefresh, diff []int
	for i := range sources {
		inrange := to.Int64(strings.TrimSpace(sources[i]))
		if inrange >= 0 && inrange < feed.SourceCount {
			cleanSources = append(cleanSources, int(inrange))
			if getLastChecked[inrange] < time.Now().Add(-30*time.Minute).Unix() {
				toRefresh = append(toRefresh, int(inrange))
			}
		}
	}
	diff = difference(cleanSources, toRefresh)
	if len(cleanSources) == 0 {
		err = errors.New("No Valid sources selected")
	} else if len(toRefresh) > 0 {
		for _, src := range toRefresh {
			var fetcher interface {
				Fetch() ([]models.FeedItem, error)
			}
			var updater = (*models.FeedItem).UpdateOthers
			switch src {
			case feed.CDC:
				fetcher = &feed.CDCFetcher{}
			case feed.ReliefWeb:
				fetcher = &feed.ReliefWebFetcher{Country: &country}
				updater = (*models.FeedItem).UpdateRelief
			case feed.CADATA:
				fetcher = &feed.CadataFetcher{}
			case feed.GDASC:
				fetcher = &feed.GdascFetcher{}
			case feed.FCO:
				fetcher = &feed.FCOFetcher{}
			}
			if fetcher == nil {
				log.Printf("[%v] no match", src)
				continue
			}
			items, err := fetcher.Fetch()
			if err != nil {
				utils.CheckErr(err)
				continue
			}
			var change bool
			for i, item := range items {
				go updater(&items[i], um.Db)
				if item.Country != country.Iso2 {
					continue
				}
				items = append(items, item)
				if change {
					continue
				}
				change = true
				diff = append(diff, src)
			}
			go um.updateLastChecked(country.Iso2, src, time.Now().Unix())
		}
	}
	if len(diff) > 0 {
		var dbFeedSources []string
		for k := range diff {
			dbFeedSources = append(dbFeedSources, strconv.Itoa(diff[k]))
		}
		feedItems, err := um.getDbFeedItems(dbFeedSources, country.Iso2, 0)
		if err == nil {
			items = append(items, feedItems...)
		} else {
			utils.CheckErr(err)
		}
	}
	sort.Sort(SortFeedByDate(items))
	if since > 0 {
		for i := len(items) - 1; i >= 0; i-- {
			if items[i].UpdatedAt <= since {
				items = append(items[:i],
					items[i+1:]...)
			}
		}
	}
	return items, err
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
