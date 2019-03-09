package feed

import (
	"testing"

	"github.com/securityfirst/umbrella-api/models"
)

func baseTest(t *testing.T, f Fetcher) {
	items, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%T: %d feeds", f, len(items))
}

func TestRefiweb(t *testing.T) {
	baseTest(t, &ReliefWebFetcher{&models.Country{ReliefWeb: 241, Iso2: "UA"}})
}

func TestGdasc(t *testing.T) {
	baseTest(t, &GdascFetcher{})
}

func TestCadata(t *testing.T) {
	baseTest(t, &CadataFetcher{})
}

func TestCDC(t *testing.T) {
	baseTest(t, &CDCFetcher{})
}

func TestFCO(t *testing.T) {
	baseTest(t, &FCOFetcher{})
}
