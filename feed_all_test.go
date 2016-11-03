package main

import "testing"

func TestGdasc(t *testing.T) {
	var f = GdascFetcher{}
	items, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range items {
		t.Logf("%+v", item)
	}
}

func TestCadata(t *testing.T) {
	var f = CadataFetcher{}
	items, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range items {
		t.Logf("%q", item.Country)
	}
}
