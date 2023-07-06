package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	controller "github.com/yrs147/scrapster/controllers"

	"github.com/gocolly/colly"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

type Product struct {
	Title        string
	Price        float64
	Rating       float64
	ReviewCount  int
	Availability string
}

func randUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.amazon.in"))

	var links []string
	var wg sync.WaitGroup

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", randUserAgent())
	})

	c.OnHTML("a.a-link-normal.s-underline-text.s-underline-link-text.s-link-style.a-text-normal", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		links = append(links, link)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError", err)
	})

	err := c.Visit("https://www.amazon.in/s?k=iphone")
	if err != nil {
		log.Fatal(err)
	}

	products := []Product{} // Initialize an empty slice of products

	for _, link := range links {
		a := "https://www.amazon.in" + link

		// Create a new collector instance for each product
		productCollector := colly.NewCollector()

		productCollector.OnRequest(func(r *colly.Request) {
			r.Headers.Set("User-Agent", randUserAgent())
		})

		product := Product{} // Create a new product instance

		wg.Add(1)
		go func(productCollector *colly.Collector, product *Product) {
			defer wg.Done()
			err := productCollector.Visit(a)
			if err != nil {
				log.Fatal(err)
			}
			products = append(products, *product) // Append the product to the products slice
		}(productCollector, &product)

		productCollector.OnHTML("span#productTitle", func(e *colly.HTMLElement) {
			title := controller.GetTitle(e)
			product.Title = title
		})

		productCollector.OnHTML("span.a-price-whole", func(e *colly.HTMLElement) {
			price := controller.GetPrice(e)
			product.Price = price
		})

		productCollector.OnHTML("span.a-size-base.a-color-base", func(e *colly.HTMLElement) {
			rating := controller.GetRating(e)
			if rating != 0 {
				product.Rating = rating
			}
		})

		productCollector.OnHTML("span#acrCustomerReviewText", func(e *colly.HTMLElement) {
			reviewCount := controller.GetReviewCount(e)
			product.ReviewCount = reviewCount
		})

		productCollector.OnHTML("div#availability span", func(e *colly.HTMLElement) {
			availability := controller.GetAvailability(e)
			product.Availability = availability
		})

	}

	wg.Wait() // Wait for all goroutines to finish

	// Print the collected products
	// for i, product := range products {
	// 	fmt.Println("------------Printing Details of Product", i+1, "-----------")
	// 	fmt.Println("Title:", product.Title)
	// 	fmt.Println("Price:", product.Price)
	// 	fmt.Println("Rating:", product.Rating)
	// 	fmt.Println("Review Count:", product.ReviewCount)
	// 	fmt.Println("Availability:", product.Availability)
	// 	fmt.Println("----------------------------------------------")
	// 	time.Sleep(time.Second *10)
	// }

	fmt.Println(products[0])
}
