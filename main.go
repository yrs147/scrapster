package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// func getTitle(a *colly.HTMLElement) string {

// }

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.amazon.in"))

	var links []string

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
		r.Headers.Set("Accept-Language", "en-US, en;q=0.5")
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

	// Scraping Product title
	c.OnHTML("span#productTitle", func(e *colly.HTMLElement) {
		title = e.Text
		title = strings.TrimSpace(title)
	})

	//Scraping Product Price
	c.OnHTML("span.a-price.aok-align-center span.a-offscreen", func(e *colly.HTMLElement) {
		priceStr := e.Text
		priceStr = strings.TrimSpace(priceStr)
		priceStr = strings.ReplaceAll(priceStr, ",", "")
		priceStr = strings.TrimPrefix(priceStr, "₹")
		price, _ = strconv.ParseFloat(priceStr, 64)
	})

	c.OnHTML("span.a-size-base.a-color-base", func(e *colly.HTMLElement) {
		ratingStr := e.Text
		ratingStr = strings.TrimSpace(ratingStr)
		rating, _ = strconv.ParseFloat(ratingStr, 64)
	})

	c.OnHTML("span#acrCustomerReviewText", func(e *colly.HTMLElement) {
		reviewCountStr := e.Text
		reviewCountStr = strings.TrimSpace(reviewCountStr)
		reviewCountStr = strings.ReplaceAll(reviewCountStr, ",", "")
		reviewCountStr = strings.TrimSuffix(reviewCountStr, " ratings")
		reviewCount, _ = strconv.Atoi(reviewCountStr)
	})

	c.OnHTML("div#availability span", func(e *colly.HTMLElement) {
		availability = e.Text
		availability = strings.TrimSpace(availability)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError", err)
	})

	err := c.Visit("https://www.amazon.in/Apple-MacBook-Chip-13-inch-256GB/dp/B08N5T6CZ6/ref=sr_1_3?keywords=macbook&qid=1688532331&sr=8-3")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Title : ",title)
	fmt.Println("Price : ",price)
	fmt.Println("Rating : " ,rating)
	fmt.Println("No of Ratings : ",reviewCount)
	fmt.Println("Availiability : ",availability)

	// for _, link := range links {
	// 	fmt.Println(link)
	// }

}
