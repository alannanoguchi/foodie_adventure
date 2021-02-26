package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Restaurant struct {
	Name string
	// Photo    string
	Cuisine  string
	Location string
	Summary  string
	Rating   string
}

type City struct {
	City        string
	Link        string
	Restaurants []*Restaurant //type list of restaurants
}

type State struct {
	State  string
	Cities map[string]*City
}

type Country struct {
	Country string
	States  map[string]*State
}

func scrapeCity(cityInfo *City) {
	//this function will scrape one specific city
	selector := "body > table:nth-child(7) > tbody > tr > td:nth-child(1) > table > tbody > tr > td > table >  tbody > tr > td > div:nth-child(1)" // grabs just the body

	fmt.Println(cityInfo.Link)

	y := colly.NewCollector(
		colly.AllowedDomains("zabihah.com", "www.zabihah.com"),
	)

	y.OnHTML(selector, func(p *colly.HTMLElement) {
		p.ForEach("tr > td:nth-child(3)", func(_ int, h *colly.HTMLElement) { // for loop, to get each individual restaurant
			link := h.ChildAttr("a", "href")
			fmt.Printf("Link found: -> %s\n", link)

		})
		fmt.Println("print")
	})

	y.Visit(cityInfo.Link)

}

func main() {
	var cities []City
	city_selector := "body > table:nth-child(7) > tbody > tr > td:nth-child(1) > table:nth-child(11) > tbody > tr"
	// link_selector := "#header > table > tbody > tr > td:nth-child(3) > div.titleBS > a"

	c := colly.NewCollector(
		colly.AllowedDomains("zabihah.com", "www.zabihah.com"),
	)

	c.OnHTML(city_selector, func(e *colly.HTMLElement) {
		// cityLink := e.ChildText("td:nth-child(1) > a > b")

		tmpCity := City{}
		tmpCity.City = e.ChildText("td:nth-child(1) > a > b")
		tmpCity.Link = "https://www.zabihah.com" + e.ChildAttr("tr > td:nth-child(1) > a", "href")
		// tmpCity.Restaurants = e.ChildAttr("#alertBox2")
		if tmpCity.City == "" {
			tmpCity.City = e.ChildText("td:nth-child(1) > div > a > b")
			tmpCity.Link = "https://www.zabihah.com" + e.ChildAttr("td:nth-child(1) > div > a", "href")
		}

		cities = append(cities, tmpCity)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	startUrl := fmt.Sprintf("https://www.zabihah.com/reg/United-States/California/C3Jynwv1mE")
	c.Visit(startUrl)
	for _, city := range cities {
		scrapeCity(&city)
	}
}
