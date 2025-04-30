package services

import (
	"fmt"

	"github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/input"
	// "github.com/go-rod/rod/lib/proto"

	"google-trends-api/internal/config"
)

func Autoscrap() {
	fmt.Println("Starting AutoScrap...")
	for {
		fmt.Println("AutoScrap loop...")
		// Iterate over countries
		for countryName, locationInfo := range config.CountryAbbreviations {
			fmt.Printf("Processing Country: %s (%s)\n", countryName, locationInfo.Code)

			// TODO: Add logic to scrape data for the country using locationInfo.Code
			countryURL := fmt.Sprintf("%s?geo=%s", config.GOOGLE_TRENDS_BASE_URL, locationInfo.Code)
			// Scraping for the country level
			scrapeData(countryURL, countryName, "")

			// Check if the country has regions defined
			if locationInfo.Regions != nil {
				// Iterate over regions within the country
				for regionName, regionCode := range locationInfo.Regions {
					fmt.Printf("  Processing Region: %s (%s)\n", regionName, regionCode)
					stateURL := fmt.Sprintf("%s?geo=%s-%s", config.GOOGLE_TRENDS_BASE_URL, locationInfo.Code, regionCode)
					scrapeData(stateURL, countryName, regionName)
				}
			} else {
				fmt.Printf("  No specific regions defined for %s.\n", countryName)
			}
		}
	}
}

func scrapeData(url string, countryName string, regionName string) {
	var rawData string
	fmt.Println("Scraping data from URL:", url)
	// Create a new browser instance in headless mode
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage(url)

	// Wait for the page to load
	page.MustWaitStable()
	fmt.Println("page load complete")

	nextPageButton := page.MustElementX(`//*[@id="trend-table"]/div[2]/div/div[2]/span[3]/button`)
	ifdisabled := nextPageButton.MustAttribute("disabled")
	for ifdisabled == nil {
		fmt.Println("next page button not disabled")
		page.MustWaitLoad()
		// Copy the table body
		tableBody := page.MustElementX(`/html/body/c-wiz/div/div[5]/div[1]/c-wiz/div/div[2]/div[1]/div[1]/div[1]/table/tbody[2]`)

		// RawHTML = tableBody.MustHTML()
		rawData += tableBody.MustHTML()

		nextPageButton.MustClick()
		fmt.Println("next page button clicked")
		page.MustWaitLoad()
		ifdisabled = nextPageButton.MustAttribute("disabled")
	}
	fmt.Println("All rawData: ", rawData)
	data := sanitizeHTML(rawData)
	Adddata(data, countryName, regionName)
}
