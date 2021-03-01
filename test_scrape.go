package main

import (
	"fmt"
	"testing"
)

type MockCity struct {
}

func TestScrapeCity(t *testing.T) {
	fmt.Println("Test City Scraper")
	expected := "https://www.zabihah.com/sub/United-States/California/Los-Angeles/High-Desert/KpKdyDDQ6G"
	result := scrapeCity(MockCity.link)
	if expected != result {
		t.Error("Test Failed")
	}
}
