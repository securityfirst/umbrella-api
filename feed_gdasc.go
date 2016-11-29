package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/securityfirst/umbrella-api/country"
	"github.com/securityfirst/umbrella-api/models"
)

type GdascFetcher struct{}

func (g *GdascFetcher) Fetch() ([]models.FeedItem, error) {
	resp, err := http.Get("http://www.gdacs.org/xml/rss.xml")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var v GdascResp
	if err := xml.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	var feeds = make([]models.FeedItem, len(v.Title))
	for i := range v.Title {
		if v.Country[i] == "" {
			continue
		}
		t, _ := time.Parse(time.RFC1123, v.PubDate[i])
		for _, name := range strings.Split(v.Country[i], ", ") {
			c, err := country.ByName(name)
			if err != nil {
				log.Printf("Country %q: %s", name, err)
				continue
			}
			feeds = append(feeds, models.FeedItem{
				Title:       v.Title[i],
				Description: v.Description[i],
				URL:         v.Link[i],
				Country:     strings.ToLower(c.Codes.Alpha2),
				Source:      GDASC,
				UpdatedAt:   t.Unix(),
			})
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
