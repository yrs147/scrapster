package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gocolly/colly"
	controller "github.com/yrs147/scrapster/controllers"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func randUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.amazon.in"))

	var links []string

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

	err := c.Visit("https://www.amazon.in/s?k=macbook")
	if err != nil {
		log.Fatal(err)
	}

	count := 1

	for _, link := range links {
		a := "https://www.amazon.in" + link

		// Create a new collector instance for each product
		productCollector := colly.NewCollector()

		productCollector.OnRequest(func(r *colly.Request) {
			r.Headers.Set("User-Agent", randUserAgent())
		})

		fmt.Println("------------Printing Details of Product", count, "-----------")
		productCollector.OnHTML("span#productTitle", func(e *colly.HTMLElement) {
			title := controller.GetTitle(e)
			fmt.Println("Title:", title)
		})

		productCollector.OnHTML("span.a-price.aok-align-center span.a-offscreen", func(e *colly.HTMLElement) {
			price := controller.GetPrice(e)
			fmt.Println("Price:", price)
		})

		productCollector.OnHTML("span.a-size-base.a-color-base", func(e *colly.HTMLElement) {
			rating := controller.GetRating(e)
			if rating != 0 {
				fmt.Println("Rating:", rating)
			}
		})
		

		productCollector.OnHTML("span#acrCustomerReviewText", func(e *colly.HTMLElement) {
			reviewCount := controller.GetReviewCount(e)
			fmt.Println("Review Count:", reviewCount)
		})

		productCollector.OnHTML("div#availability span", func(e *colly.HTMLElement) {
			availability := controller.GetAvailability(e)
			fmt.Println("Availability:", availability)
		})

		err := productCollector.Visit(a)
		if err != nil {
			log.Fatal(err)
		}

		count++
	}
}
