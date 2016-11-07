package main

import "testing"

type Fetcher interface {
	Fetch() ([]FeedItem, error)
}

func baseTest(t *testing.T, f Fetcher) {
	items, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%T: %d feeds", f, len(items))
}

func TestRefiweb(t *testing.T) {
	baseTest(t, &RefiWebFetcher{&Country{ReliefWeb: 241, Iso2: "UA"}})
}

func TestGdasc(t *testing.T) {
	baseTest(t, &GdascFetcher{})
}

func TestCadata(t *testing.T) {
	baseTest(t, &CadataFetcher{})
}
