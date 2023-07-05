package main

import (
	"fmt"
	"log"
	"math/rand"
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

func randUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

// func getTitle(a *colly.HTMLElement) string {

// }

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

	var (
		title        string
		price        float64
		rating       float64
		reviewCount  int
		availability string
	)

	c.OnHTML("span#productTitle", func(e *colly.HTMLElement) {
		title = controller.GetTitle(e)
	})

	c.OnHTML("span.a-price.aok-align-center span.a-offscreen", func(e *colly.HTMLElement) {
		price = controller.GetPrice(e)	
	})

	c.OnHTML("span.a-size-base.a-color-base", func(e *colly.HTMLElement) {
		rating = controller.GetRating(e)
	})

	c.OnHTML("span#acrCustomerReviewText", func(e *colly.HTMLElement) {
		reviewCount = controller.GetReviewCount(e)	
	})

	c.OnHTML("div#availability span", func(e *colly.HTMLElement) {
		availability = controller.GetAvailability(e)	
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError", err)
	})

	err := c.Visit("https://www.amazon.in/Apple-Cellular-Rugged-Titanium-Orange/dp/B0BDKGNMJQ/ref=sr_1_1_sspa?crid=3T3NEK3BBR4P2&keywords=apple+watch&qid=1688543984&sprefix=apple+watch%2Caps%2C218&sr=8-1-spons&sp_csd=d2lkZ2V0TmFtZT1zcF9hdGY&psc=1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Title:", title)
	fmt.Println("Price:", price)
	fmt.Println("Rating:", rating)
	fmt.Println("Review Count:", reviewCount)
	fmt.Println("Availability:", availability)

	// for _, link := range links {
	// 	fmt.Println(link)
	// }

}
