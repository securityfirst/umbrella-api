package feed

import (
	"encoding/xml"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/securityfirst/umbrella-api/country"
	"github.com/securityfirst/umbrella-api/models"
)

type CDCFetcher struct{}

var cdcCountry = regexp.MustCompile(`.* in (.*)$`)

func (c *CDCFetcher) Fetch() ([]models.FeedItem, error) {
	body, err := makeRequest("https://tools.cdc.gov/api/v2/resources/media/285740.rss", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var resp cdcResp
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	var feeds = make([]models.FeedItem, 0, len(resp.Item))
	for _, v := range resp.Item {
		match := cdcCountry.FindStringSubmatch(v.Title)
		if len(match) != 2 {
			continue
		}
		name := strings.TrimSpace(match[1])
		t, _ := time.Parse(time.RFC1123, v.PubDate)
		c, err := country.ByName(name)
		if err != nil {
			log.Printf("CDCFeed - Country %q: %s", name, err)
			continue
		}
		feeds = append(feeds, models.FeedItem{
			Title:       v.Title,
			Description: v.Description,
			URL:         v.Link,
			Country:     strings.ToLower(c.Codes.Alpha2),
			Source:      CDC,
			UpdatedAt:   t.Unix(),
		})
	}
	return feeds, nil
}

type cdcResp struct {
	Item []struct {
		Text        string `xml:",chardata"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Link        string `xml:"link"`
		GUID        struct {
			Text        string `xml:",chardata"`
			IsPermaLink string `xml:"isPermaLink,attr"`
		} `xml:"guid"`
		PubDate  string `xml:"pubDate"`
		Category string `xml:"category"`
	} `xml:"channel>item"`
}
