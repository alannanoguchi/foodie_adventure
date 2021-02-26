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

func scrapeRestaurant(restaurantInfo *Restaurant) {
	// link_selector := "#header > table > tbody > tr > td:nth-child(3) > div.titleBS > a"
	link_selector := "#header > table > tbody > tr"

	fmt.Println(restaurantInfo.Name)

	x := colly.NewCollector(
		colly.AllowedDomains("zabihah.com", "www.zabihah.com"),
	)

	x.OnHTML(link_selector, func(p *colly.HTMLElement) {
		// for loop, to get each individual restaurant
		p.ForEach("tr > td:nth-child(3)", func(_ int, h *colly.HTMLElement) {
			link := h.ChildAttr("a", "href")
			fmt.Printf("Restaurant link found: -> %s\n", link)

		})
	})

	x.Visit(restaurantInfo.Name)
}

func scrapeCity(cityInfo *City) {

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
		tmpRestaurant.Cuisine = p.ChildText("#alertBox2")
		tmpRestaurant.Location = p.ChildText("#header > table > tbody > tr > td:nth-child(3) > div.tinyLink")
		tmpRestaurant.Rating = p.ChildText("#badge_score")

		restaurants = append(restaurants, tmpRestaurant)

		// for loop, to get each individual restaurant
		p.ForEach("tr > td:nth-child(3)", func(_ int, h *colly.HTMLElement) {
			link := h.ChildAttr("a", "href")
			fmt.Printf("Restaurant link found: -> %s\n", link)

		})
	})

	y.Visit(cityInfo.Link)
	for _, restaurant := range restaurants {
		scrapeRestaurant(&restaurant)
	}

}

func main() {
	var cities []City
	city_selector := "body > table:nth-child(7) > tbody > tr > td:nth-child(1) > table:nth-child(11) > tbody > tr"

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