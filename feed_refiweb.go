package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gosexy/to"
)

type RefiWebFetcher struct {
	Country *Country
}

func (r *RefiWebFetcher) Fetch() ([]FeedItem, error) {
	body, err := makeRequest(fmt.Sprintf("http://api.rwlabs.org/v0/country/%v", r.Country.ReliefWeb), "get", nil)
	if err != nil {
		return nil, err
	}
	var rwResp RWResponse
	if err = json.Unmarshal(body, &rwResp); err != nil {
		return nil, err
	}
	if rwResp.Data.Item.DescriptionHTML == "" {
		return nil, nil
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rwResp.Data.Item.DescriptionHTML))
	if err != nil {
		return nil, err
	}
	var items []FeedItem
	doc.Find("h3").First().Next().Children().Each(func(i int, t *goquery.Selection) {
		item, err := r.parseItem(t)
		if err != nil {
			log.Println("", err)
			return
		}
		items = append(items, *item)
	})
	return items, nil
}

func (r *RefiWebFetcher) parseItem(t *goquery.Selection) (*FeedItem, error) {
	href, ok := t.Contents().Attr("href")
	if !ok {
		return nil, errors.New("no href")
	}
	item := FeedItem{
		Title:     t.Contents().Text(),
		URL:       href,
		Country:   r.Country.Iso2,
		Source:    ReliefWeb,
		UpdatedAt: time.Now().Unix(),
	}
	segments := strings.Split(href, "/")
	if len(segments) == 0 || to.Int64(segments[len(segments)-1]) == 0 {
		return &item, nil
	}
	nodeUrl := fmt.Sprintf("http://api.rwlabs.org/v0/report/%v", segments[len(segments)-1])
	body, err := makeRequest(nodeUrl, "get", nil)
	if err != nil {
		return nil, err
	}
	var rwReport RWReport
	if err = json.Unmarshal(body, &rwReport); err != nil {
		return nil, err
	}
	if rwReport.Data.Item.Headline.Summary != "" {
		item.Description = rwReport.Data.Item.Headline.Summary
	} else {
		item.Description = rwReport.Data.Item.BodyHTML
	}
	item.UpdatedAt = rwReport.Data.Item.Date.Changed / 1000
	return &item, nil
}
