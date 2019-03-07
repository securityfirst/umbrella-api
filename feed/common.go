package feed

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/securityfirst/umbrella-api/models"
)

type Fetcher interface {
	Fetch() ([]models.FeedItem, error)
}

var client = http.Client{
	Timeout: time.Second * 10,
}

func makeRequest(uri string, method string, requestBody io.Reader) ([]byte, error) {
	req, err := http.NewRequest(strings.ToUpper("GET"), uri, requestBody)
	if err != nil {
		return nil, err
	}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// List of Sources
const (
	ReliefWeb = iota
	FCO
	UN
	CDC
	GDASC
	CADATA
	SourceCount
)
