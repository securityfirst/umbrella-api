package main

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/pariz/gountries"
)

var query = gountries.New()

type GdascFetcher struct{}

func (g *GdascFetcher) Fetch() ([]FeedItem, error) {
	resp, err := http.Get("http://www.gdacs.org/xml/rss.xml")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var v GdascResp
	if err := xml.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	var feeds = make([]FeedItem, len(v.Title))
	for i := range v.Title {
		t, _ := time.Parse(time.RFC1123, v.PubDate[i])
		c, _ := query.FindCountryByName(v.Country[i])
		feeds[i] = FeedItem{
			Title:       v.Title[i],
			Description: v.Description[i],
			URL:         v.Link[i],
			Country:     c.Codes.Alpha2,
			Source:      GDASC,
			UpdatedAt:   t.Unix(),
		}
	}
	return feeds, nil
}

type GdascResp struct {
	Title       []string `xml:"channel>item>title"`
	Description []string `xml:"channel>item>description"`
	Country     []string `xml:"channel>item>country"`
	Link        []string `xml:"channel>item>link"`
	PubDate     []string `xml:"channel>item>pubDate"`
}
