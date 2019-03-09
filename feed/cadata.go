package feed

import (
	"encoding/xml"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/securityfirst/umbrella-api/country"
	"github.com/securityfirst/umbrella-api/models"
)

type CadataFetcher struct{}

func (g *CadataFetcher) Fetch() ([]models.FeedItem, error) {
	const baseURL = "https://cadatacatalog.state.gov/storage/f/2013-11-24T21%3A00%3A"
	var feeds []models.FeedItem
	for _, src := range []string{"30.424Z/tas.xml", "58.223Z/tws.xml"} {
		body, err := makeRequest(baseURL+src, http.MethodGet, nil)
		if err != nil {
			return nil, err
		}
		var v cadataResp
		if err := xml.Unmarshal(body, &v); err != nil {
			return nil, err
		}
		for i := range v.Title {
			code := strings.TrimSpace(v.Identifier[i])
			if code == "" {
				continue
			}
			c, err := country.ByFips(code)
			if err != nil {
				log.Printf("Cannot find country: %s", err)
				continue
			}
			t, _ := time.Parse(time.RFC1123, v.PubDate[i])
			feeds = append(feeds, models.FeedItem{
				Title:       v.Title[i],
				Description: v.Description[i],
				URL:         v.Link[i],
				Country:     strings.ToLower(c.Codes.Alpha2),
				Source:      CADATA,
				UpdatedAt:   t.Unix(),
			})
		}
	}
	return feeds, nil
}

type cadataResp struct {
	Link        []string `xml:"channel>item>link"`
	Description []string `xml:"channel>item>description"`
	Title       []string `xml:"channel>item>title"`
	PubDate     []string `xml:"channel>item>pubDate"`
	Identifier  []string `xml:"channel>item>identifier"`
}
