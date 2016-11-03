package main

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"
)

type CadataFetcher struct{}

func (g *CadataFetcher) Fetch() ([]FeedItem, error) {
	var feeds []FeedItem
	for _, src := range []string{
		"https://cadatacatalog.state.gov/storage/f/2013-11-24T21%3A00%3A30.424Z/tas.xml",
		"https://cadatacatalog.state.gov/storage/f/2013-11-24T21%3A00%3A58.223Z/tws.xml",
	} {
		resp, err := http.Get(src)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var v CadataResp
		if err := xml.NewDecoder(resp.Body).Decode(&v); err != nil {
			return nil, err
		}
		for i := range v.Title {
			t, _ := time.Parse(time.RFC1123, v.PubDate[i])
			feeds = append(feeds, FeedItem{
				Title:       v.Title[i],
				Description: v.Description[i],
				URL:         v.Link[i],
				Country:     strings.TrimSpace(v.Identifier[i]),
				Source:      GDASC,
				UpdatedAt:   t.Unix(),
			})
		}
	}
	return feeds, nil
}

type CadataResp struct {
	Link        []string `xml:"channel>item>link"`
	Description []string `xml:"channel>item>description"`
	Title       []string `xml:"channel>item>title"`
	PubDate     []string `xml:"channel>item>pubDate"`
	Identifier  []string `xml:"channel>item>identifier"`
}
