package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Restaurant struct {
	Name     string `json:"name"`
	Link     string `json:"link:`
	Photo    string `json:"photo"`
	Cuisine  string `json:"cuisine"`
	Location string `json:"location"`
	Rating   string `json:"rating"`
}

type City struct {
	City        string        `json:"city"`
	Link        string        `json:"link"`
	Restaurants []*Restaurant `json:"restaurant"` //type list of restaurants
}

type State struct {
	State  string           `json:"state"`
	Cities map[string]*City `json:"cities"`
}

type Country struct {
	Country string            `json:"country"`
	States  map[string]*State `json:"state"`
}

func scrapeCity(cityInfo *City) {

	selector := "body > table:nth-child(7) > tbody > tr > td:nth-child(1) > table > tbody > tr > td > table >  tbody > tr > td > div:nth-child(1)"

	fmt.Printf("City link found: -> %s\n", cityInfo.Link)

	y := colly.NewCollector(
		colly.AllowedDomains("zabihah.com", "www.zabihah.com"),
	)

	y.OnHTML(selector, func(p *colly.HTMLElement) {

		tmpRestaurant := Restaurant{}
		tmpRestaurant.Name = p.ChildText("#header > table > tbody > tr > td:nth-child(3) > div.titleBS > a")
		tmpRestaurant.Link = p.ChildAttr("a", "href")
		tmpRestaurant.Photo = p.ChildAttr("td:nth-child(1) > a > img", "src")
		tmpRestaurant.Cuisine = p.ChildText("#alertBox2")
		tmpRestaurant.Location = p.ChildText("#header > table > tbody > tr > td:nth-child(3) > div.tinyLink")
		tmpRestaurant.Rating = p.ChildText("#badge_score")

		cityInfo.Restaurants = append(cityInfo.Restaurants, &tmpRestaurant)

		// for loop, to get each individual restaurant
		p.ForEach("tr > td:nth-child(3)", func(_ int, h *colly.HTMLElement) {
			link := h.ChildAttr("a", "href")
			fmt.Printf("Restaurant link found: -> %s\n", link)

		})

	})

	y.Visit(cityInfo.Link)

}

// A function that creates workers to scrape city data
func scrapeWorker(cities chan *City) {
	for city := range cities {
		scrapeCity(city)
	}
}

func createJSONFile(stateInfo *State) {
	file, err := json.MarshalIndent(stateInfo, "", "    ")
	if err != nil {
		log.Println("JSON file not created")
		return
	}

	_ = ioutil.WriteFile("scraped_data.json", file, 0644)
}

func main() {

	california := State{
		State:  "California",
		Cities: make(map[string]*City),
	}

	cityChan := make(chan *City)

	// Created 5 workers
	for i := 0; i < 5; i++ {
		go scrapeWorker(cityChan)
	}

	city_selector := "body > table:nth-child(7) > tbody > tr > td:nth-child(1) > table:nth-child(11) > tbody > tr"

	time.Sleep(5 * time.Second)

	c := colly.NewCollector(
		colly.AllowedDomains("zabihah.com", "www.zabihah.com"),
	)

	c.OnHTML(city_selector, func(e *colly.HTMLElement) {

		tmpCity := City{}
		tmpCity.City = e.ChildText("td:nth-child(1) > a > b")
		tmpCity.Link = "https://www.zabihah.com" + e.ChildAttr("tr > td:nth-child(1) > a", "href")
		if tmpCity.City == "" {
			tmpCity.City = e.ChildText("td:nth-child(1) > div > a > b")
			tmpCity.Link = "https://www.zabihah.com" + e.ChildAttr("td:nth-child(1) > div > a", "href")
		}

		california.Cities[tmpCity.City] = &tmpCity
		cityChan <- &tmpCity

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	startUrl := fmt.Sprintf("https://www.zabihah.com/reg/United-States/California/C3Jynwv1mE")
	c.Visit(startUrl)

	close(cityChan)

	createJSONFile(&california)

}
