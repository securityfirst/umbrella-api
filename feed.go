package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gosexy/to"
)

func (um *Umbrella) getFeed(c *gin.Context) {
	since := to.Int64(c.Request.URL.Query().Get("since"))
	country, err := um.getCountryInfo(c.Request.URL.Query().Get("country"))
	um.checkErr(c, err)
	sources := strings.Split(c.Request.URL.Query().Get("sources"), ",")
	feedItems, err := um.getFeedItems(sources, country, since)
	um.checkErr(c, err)
	if err == nil {
		c.JSON(200, feedItems)
		return
	}
	c.JSON(400, gin.H{"error": err.Error()})
}

func (um *Umbrella) getFeedItems(sources []string, country Country, since int64) (items []FeedItem, err error) {
	getLastChecked := um.getLastChecked(country.Iso2)
	var cleanSources, toRefresh, diff []int
	for i := range sources {
		inrange := to.Int64(strings.TrimSpace(sources[i]))
		if inrange >= 0 && inrange <= 3 {
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
		for k := range toRefresh {
			switch toRefresh[k] {
			case CDC:
				fmt.Println("refresh cdc")
				// refresh cdc
				go um.updateLastChecked(country.Iso2, CDC, time.Now().Unix())
			case ReliefWeb:
				body, err := makeRequest(fmt.Sprintf("http://api.rwlabs.org/v0/country/%v", country.ReliefWeb), "get", nil)
				var rwResp RWResponse
				err = json.Unmarshal(body, &rwResp)
				if err != nil {
					checkErr(err)
					fmt.Println(string(body[:]))
				} else {
					go um.updateLastChecked(country.Iso2, ReliefWeb, time.Now().Unix())
				}
				if rwResp.Data.Item.DescriptionHTML != "" {
					doc, err := goquery.NewDocumentFromReader(strings.NewReader(rwResp.Data.Item.DescriptionHTML))
					if err != nil {
						log.Fatal(err)
					}
					s := doc.Find("h3").First()
					s.Next().Children().Each(func(i int, t *goquery.Selection) {
						href, ok := t.Contents().Attr("href")
						if ok {
							item := FeedItem{
								Title:     t.Contents().Text(),
								URL:       href,
								Country:   country.Iso2,
								Source:    ReliefWeb,
								UpdatedAt: time.Now().Unix(),
							}
							segments := strings.Split(href, "/")
							if len(segments) > 0 && to.Int64(segments[len(segments)-1]) != 0 {
								nodeUrl := fmt.Sprintf("http://api.rwlabs.org/v0/report/%v", segments[len(segments)-1])
								body, err := makeRequest(nodeUrl, "get", nil)
								var rwReport RWReport
								err = json.Unmarshal(body, &rwReport)
								if err != nil {
									checkErr(err)
									fmt.Println(string(body[:]))
								} else {
									if rwReport.Data.Item.Headline.Summary != "" {
										item.Description = rwReport.Data.Item.Headline.Summary
									} else {
										item.Description = rwReport.Data.Item.BodyHTML
									}
									item.UpdatedAt = rwReport.Data.Item.Date.Changed / 1000
								}

							}
							items = append(items, item)
							go item.updateRelief(um)
						}
					})
				}
			}
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
			checkErr(err)
		}
	}
	sort.Sort(SortFeedByDate(items))
	return items, err
}

type SortFeedByDate []FeedItem

func (slice SortFeedByDate) Len() int {
	return len(slice)
}

func (slice SortFeedByDate) Less(i, j int) bool {
	return slice[i].UpdatedAt > slice[j].UpdatedAt
}

func (slice SortFeedByDate) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

const (
	ReliefWeb = iota
	FCO
	UN
	CDC
)
