package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Restaurant struct {
	Name     string
	Photo    string
	Cuisine  string
	Location string
	// Summary  string
	Rating string
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

// func scrapeRestaurant(restaurantInfo *Restaurant) {
// 	// link_selector := "#header > table > tbody > tr > td:nth-child(3) > div.titleBS > a"
// 	// link_selector := "#header > table > tbody > tr"

// 	fmt.Println("Restaurant: ", restaurantInfo.Name)
// 	fmt.Println("Photo: ", restaurantInfo.Photo)
// 	fmt.Println("Cuisine: ", restaurantInfo.Cuisine)
// 	fmt.Println("Location: ", restaurantInfo.Location)
// 	fmt.Println("Rating: ", restaurantInfo.Rating)

// 	x := colly.NewCollector(
// 		colly.AllowedDomains("zabihah.com", "www.zabihah.com"),
// 	)

// 	// x.OnHTML(link_selector, func(p *colly.HTMLElement) {
// 	// 	// for loop, to get each individual restaurant
// 	// 	p.ForEach("tr > td:nth-child(3)", func(_ int, h *colly.HTMLElement) {
// 	// 		link := h.ChildAttr("a", "href")
// 	// 		fmt.Printf("Restaurant link found: -> %s\n", link)

// 	// 	})
// 	// })

// 	x.Visit(restaurantInfo.Name)
// }

func scrapeCity(cityInfo *City) {

	// channel <- City{City: city, Link: link, Restaurants: []*Restaurant }

	var restaurants []Restaurant

	//this function will scrape one specific city
	selector := "body > table:nth-child(7) > tbody > tr > td:nth-child(1) > table > tbody > tr > td > table >  tbody > tr > td > div:nth-child(1)"

	fmt.Printf("City link found: -> %s\n", cityInfo.Link)

	y := colly.NewCollector(
		colly.AllowedDomains("zabihah.com", "www.zabihah.com"),
	)

	y.OnHTML(selector, func(p *colly.HTMLElement) {

		tmpRestaurant := Restaurant{}
		tmpRestaurant.Name = p.ChildText("#header > table > tbody > tr > td:nth-child(3) > div.titleBS > a")
		tmpRestaurant.Photo = p.ChildAttr("td:nth-child(1) > a > img", "src")
		tmpRestaurant.Cuisine = p.ChildText("#alertBox2")
		tmpRestaurant.Location = p.ChildText("#header > table > tbody > tr > td:nth-child(3) > div.tinyLink")
		tmpRestaurant.Rating = p.ChildText("#badge_score")

		restaurants = append(restaurants, tmpRestaurant)

		// for loop, to get each individual restaurant
		p.ForEach("tr > td:nth-child(3)", func(_ int, h *colly.HTMLElement) {
			link := h.ChildAttr("a", "href")
			fmt.Printf("Restaurant link found: -> %s\n", link)

		})

		js, err := json.MarshalIndent(tmpRestaurant, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))
	})

	// y.OnRequest(func(h *colly.Request) {
	// 	fmt.Println("Visiting City Link: ", h.URL.String())
	// })

	y.Visit(cityInfo.Link)
}

func scrapeWorker(cities chan *City) {
	for city := range cities {
		scrapeCity(city)
	}
}

func main() {

	cityChan := make(chan *City)

	// Created 5 workers
	for i := 0; i < 5; i++ {
		go scrapeWorker(cityChan)
	}
	// go scrapeWorker(cityChan)
	// channel <- City{City: city, Link: link, Restaurants: []*Restaurant }

	var cities []*City // a list of references
	city_selector := "body > table:nth-child(7) > tbody > tr > td:nth-child(1) > table:nth-child(11) > tbody > tr"

	time.Sleep(5 * time.Second)

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

		cityChan <- &tmpCity
		cities = append(cities, &tmpCity)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	startUrl := fmt.Sprintf("https://www.zabihah.com/reg/United-States/California/C3Jynwv1mE")
	c.Visit(startUrl)

	close(cityChan)

}
