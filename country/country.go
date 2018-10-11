package country

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pariz/gountries"
)

var fipsIndex map[string]string

var exceptions = map[string]string{
	"Russian Federation": "Russia",
}

var query = gountries.New()

func ByName(name string) (*gountries.Country, error) {
	if v, ok := exceptions[name]; ok {
		name = v
	}
	c, err := query.FindCountryByName(name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func ByFips(name string) (*gountries.Country, error) {
	v, ok := fipsIndex[name]
	if !ok {
		return nil, fmt.Errorf("No FIBS entry for %q", name)
	}
	return ByName(v)
}

func init() {
	path := os.Getenv("FIPS")
	if path == "" {
		path = "fips.csv"
	}
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open %s: %s", path, err)
	}
	defer f.Close()
	fipsIndex = make(map[string]string, 265)
	r := csv.NewReader(f)
	r.Comma = '	'
	for count := 0; ; count++ {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Cannot red %s: %s", path, err)
		}
		fipsIndex[record[0]] = record[1]
	}
}
