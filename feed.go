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

	"github.com/securityfirst/umbrella-api/models"
	"github.com/securityfirst/umbrella-api/utils"

	"github.com/PuerkitoBio/goquery"
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
		if inrange >= 0 && inrange < SourceCount {
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
			switch src := toRefresh[k]; src {
			case CDC:
				lenItems := len(items)
				fmt.Println("refresh cdc")
				// refresh cdc
				if len(items) == lenItems {
					diff = append(diff, src)
				}
				// go um.updateLastChecked(country.Iso2, CDC, time.Now().Unix())
			case ReliefWeb:
				lenItems := len(items)
				body, err := makeRequest(fmt.Sprintf("https://api.reliefweb.int/v1/countries/%v", country.ReliefWeb), "get", nil)
				var rwResp RWResponse
				err = json.Unmarshal(body, &rwResp)
				if err != nil {
					utils.CheckErr(err)
					fmt.Println(string(body[:]))
				} else {
					go um.updateLastChecked(country.Iso2, ReliefWeb, time.Now().Unix())
				}
				if len(rwResp.Data) < 1 {
					utils.CheckErr(errors.New("No data received"))
					continue
				}
				if rwResp.Data[0].Fields.DescriptionHTML != "" {
					doc, err := goquery.NewDocumentFromReader(strings.NewReader(rwResp.Data[0].Fields.DescriptionHTML))
					if err != nil {
						log.Fatal(err)
					}
					s := doc.Find("h3").First()
					s.Next().Children().Each(func(i int, t *goquery.Selection) {
						href, ok := t.Contents().Attr("href")
						if ok {
							item := models.FeedItem{
								Title:     t.Contents().Text(),
								URL:       href,
								Country:   country.Iso2,
								Source:    ReliefWeb,
								UpdatedAt: time.Now().Unix(),
							}
							segments := strings.Split(href, "/")
							if len(segments) > 0 && to.Int64(segments[len(segments)-1]) != 0 {
								nodeUrl := fmt.Sprintf("https://api.reliefweb.int/v1/reports/%v", segments[len(segments)-1])
								body, err := makeRequest(nodeUrl, "get", nil)
								var rwReport RWReport
								err = json.Unmarshal(body, &rwReport)
								if err != nil {
									utils.CheckErr(err)
									fmt.Println(string(body[:]))
								} else {
									if rwReport.Data[0].Fields.Headline.Summary != "" {
										item.Description = rwReport.Data[0].Fields.Headline.Summary
									} else {
										item.Description = rwReport.Data[0].Fields.BodyHTML
									}
									item.UpdatedAt = rwReport.Data[0].Fields.Date.Changed.Unix()
								}

							}
							items = append(items, item)
							go item.UpdateRelief(um.Db)
						}
					})
				}
				if len(items) == lenItems {
					diff = append(diff, src)
				}
			case GDASC:
				f := GdascFetcher{}
				srcItems, err := f.Fetch()
				if err != nil {
					utils.CheckErr(err)
					continue
				}
				var change bool
				for i, item := range srcItems {
					go srcItems[i].UpdateOthers(um.Db)
					if item.Country == country.Iso2 {
						items = append(items, item)
						if change {
							continue
						}
						change = true
						diff = append(diff, src)
					}
				}
			case CADATA:
				f := CadataFetcher{}
				srcItems, err := f.Fetch()
				if err != nil {
					utils.CheckErr(err)
					continue
				}
				var change bool
				for i, item := range srcItems {
					go srcItems[i].UpdateOthers(um.Db)
					if item.Country == country.Iso2 {
						items = append(items, item)
						if change {
							continue
						}
						change = true
						diff = append(diff, src)
					}
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

const (
	ReliefWeb = iota
	FCO
	UN
	CDC
	GDASC
	CADATA
	SourceCount
)
