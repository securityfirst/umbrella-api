package feed

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gosexy/to"
	"github.com/securityfirst/umbrella-api/models"
)

type ReliefWebFetcher struct {
	Country *models.Country
}

func (r *ReliefWebFetcher) Fetch() ([]models.FeedItem, error) {
	body, err := makeRequest(fmt.Sprintf("https://api.reliefweb.int/v1/countries/%v", r.Country.ReliefWeb), http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	var resp reliefWebResp
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data) < 1 {
		return nil, errors.New("No data received")
	}
	if resp.Data[0].Fields.DescriptionHTML == "" {
		return nil, nil
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.Data[0].Fields.DescriptionHTML))
	if err != nil {
		return nil, err
	}
	var items []models.FeedItem
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

func (r *ReliefWebFetcher) parseItem(t *goquery.Selection) (*models.FeedItem, error) {
	href, ok := t.Contents().Attr("href")
	if !ok {
		return nil, errors.New("no href")
	}
	item := models.FeedItem{
		Title:     t.Contents().Text(),
		URL:       href,
		Country:   r.Country.Iso2,
		Source:    models.ReliefWeb,
		UpdatedAt: time.Now().Unix(),
	}
	segments := strings.Split(href, "/")
	if len(segments) == 0 || to.Int64(segments[len(segments)-1]) == 0 {
		return &item, nil
	}
	nodeURL := fmt.Sprintf("https://api.reliefweb.int/v1/reports/%v", segments[len(segments)-1])
	body, err := makeRequest(nodeURL, "get", nil)
	if err != nil {
		return nil, err
	}
	var rep reliefWebReportResp
	if err = json.Unmarshal(body, &rep); err != nil {
		return nil, err
	}
	if rep.Data[0].Fields.Headline.Summary != "" {
		item.Description = rep.Data[0].Fields.Headline.Summary
	} else {
		item.Description = rep.Data[0].Fields.BodyHTML
	}
	item.UpdatedAt = rep.Data[0].Fields.Date.Changed.Unix()
	return &item, nil
}

type reliefWebResp struct {
	Href  string `json:"href"`
	Time  int    `json:"time"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Collection struct {
			Href string `json:"href"`
		} `json:"collection"`
	} `json:"links"`
	TotalCount int `json:"totalCount"`
	Count      int `json:"count"`
	Data       []struct {
		Fields struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			Description     string `json:"description"`
			Status          string `json:"status"`
			Iso3            string `json:"iso3"`
			Featured        bool   `json:"featured"`
			VideoPlaylist   string `json:"video_playlist"`
			URL             string `json:"url"`
			URLAlias        string `json:"url_alias"`
			DescriptionHTML string `json:"description-html"`
			Current         bool   `json:"current"`
			Location        struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"location"`
		} `json:"fields"`
		ID string `json:"id"`
	} `json:"data"`
}

type reliefWebReportResp struct {
	Href  string `json:"href"`
	Time  int    `json:"time"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Collection struct {
			Href string `json:"href"`
		} `json:"collection"`
	} `json:"links"`
	TotalCount int `json:"totalCount"`
	Count      int `json:"count"`
	Data       []struct {
		Fields struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			Body     string `json:"body"`
			Headline struct {
				Title   string `json:"title"`
				Summary string `json:"summary"`
				Image   struct {
					ID        string `json:"id"`
					Width     string `json:"width"`
					Height    string `json:"height"`
					URL       string `json:"url"`
					Filename  string `json:"filename"`
					Mimetype  string `json:"mimetype"`
					Filesize  string `json:"filesize"`
					Copyright string `json:"copyright"`
					Caption   string `json:"caption"`
					URLLarge  string `json:"url-large"`
					URLSmall  string `json:"url-small"`
					URLThumb  string `json:"url-thumb"`
				} `json:"image"`
			} `json:"headline"`
			File []struct {
				ID          string `json:"id"`
				Description string `json:"description"`
				URL         string `json:"url"`
				Filename    string `json:"filename"`
				Mimetype    string `json:"mimetype"`
				Filesize    string `json:"filesize"`
				Preview     struct {
					URL      string `json:"url"`
					URLLarge string `json:"url-large"`
					URLSmall string `json:"url-small"`
					URLThumb string `json:"url-thumb"`
				} `json:"preview"`
			} `json:"file"`
			PrimaryCountry struct {
				Href     string `json:"href"`
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Iso3     string `json:"iso3"`
				Location struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"location"`
			} `json:"primary_country"`
			Country []struct {
				Href     string `json:"href"`
				ID       int    `json:"id"`
				Name     string `json:"name"`
				Iso3     string `json:"iso3"`
				Location struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"location"`
				Primary bool `json:"primary"`
			} `json:"country"`
			Source []struct {
				Href        string `json:"href"`
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Shortname   string `json:"shortname"`
				Longname    string `json:"longname,omitempty"`
				SpanishName string `json:"spanish_name,omitempty"`
				Homepage    string `json:"homepage"`
				Type        struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"type"`
			} `json:"source"`
			Language []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"language"`
			Theme []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"theme"`
			Format []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"format"`
			OchaProduct []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"ocha_product"`
			Disaster []struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Glide string `json:"glide"`
				Type  []struct {
					ID      int    `json:"id"`
					Name    string `json:"name"`
					Code    string `json:"code"`
					Primary bool   `json:"primary"`
				} `json:"type"`
			} `json:"disaster"`
			DisasterType []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"disaster_type"`
			URL      string `json:"url"`
			URLAlias string `json:"url_alias"`
			BodyHTML string `json:"body-html"`
			Date     struct {
				Original time.Time `json:"original"`
				Changed  time.Time `json:"changed"`
				Created  time.Time `json:"created"`
			} `json:"date"`
		} `json:"fields"`
		ID string `json:"id"`
	} `json:"data"`
}
