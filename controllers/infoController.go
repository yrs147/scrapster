package controllers

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)


func GetTitle(e *colly.HTMLElement) string {
	title := e.Text
	title = strings.TrimSpace(title)

	return title
}

func GetPrice(e *colly.HTMLElement) float64 {

	priceStr := e.Text
	priceStr = strings.TrimSpace(priceStr)
	priceStr = strings.ReplaceAll(priceStr, ",", "")
	priceStr = strings.TrimPrefix(priceStr, "₹")
	price, _ := strconv.ParseFloat(priceStr, 64)

	return price
}

func GetRating(e *colly.HTMLElement) float64 {
	ratingStr := e.Text
	ratingStr = strings.TrimSpace(ratingStr)
	rating, _ := strconv.ParseFloat(ratingStr, 64)

	return rating
}

func GetReviewCount(e *colly.HTMLElement) int {
	reviewCountStr := e.Text
	reviewCountStr = strings.TrimSpace(reviewCountStr)
	reviewCountStr = strings.ReplaceAll(reviewCountStr, ",", "")
	reviewCountStr = strings.TrimSuffix(reviewCountStr, " ratings")
	reviewCount, _ := strconv.Atoi(reviewCountStr)

	return reviewCount
}

func GetAvailability(e *colly.HTMLElement) string {
	availability := e.Text
	availability = strings.TrimSpace(availability)

	return availability
}

// func Createlink(url string) string{
// 	baselink := "https://www.amazon.in/"
// }
