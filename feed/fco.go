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

type FCOFetcher struct{}

func (g *FCOFetcher) Fetch() ([]models.FeedItem, error) {
	body, err := makeRequest("https://www.gov.uk/government/organisations/latest.atom?organisations%5B%5D=foreign-commonwealth-office", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var v fcoResp
	if err := xml.Unmarshal(body, &v); err != nil {
		return nil, err
	}
	var feeds = make([]models.FeedItem, 0, len(v.Entry))
	for _, e := range v.Entry {
		if !strings.HasSuffix(e.Title, " travel advice") {
			continue
		}
		name := strings.TrimSpace(strings.TrimSuffix(e.Title, " travel advice"))
		c, err := country.ByName(name)
		if err != nil {
			log.Printf("FCO - Country %q: %s", name, err)
			continue
		}
		t, err := time.Parse(time.RFC3339, e.Updated)
		if err != nil {
			log.Printf("FCO - Time %q: %s", e.Updated, err)
			continue
		}
		feeds = append(feeds, models.FeedItem{
			Title:       e.Title,
			Description: e.Summary.Text,
			URL:         e.Link.Href,
			Country:     strings.ToLower(c.Codes.Alpha2),
			Source:      models.FCO,
			UpdatedAt:   t.Unix(),
		})
	}
	return feeds, nil
}

type fcoResp struct {
	Entry []struct {
		Text    string `xml:",chardata"`
		ID      string `xml:"id"`
		Updated string `xml:"updated"`
		Link    struct {
			Text string `xml:",chardata"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Title   string `xml:"title"`
		Summary struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"summary"`
	} `xml:"entry"`
}
