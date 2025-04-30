package services

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"

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

	// Wait for the page to load and the data to be available
	// ...existing code...
	page.MustWaitStable()
	fmt.Println("page load complete")

	const maxRetries = 3
	const retryDelay = 2 // seconds

	for i := 0; ; i++ { // Loop indefinitely until break or error
		var nextPageButton *rod.Element
		var err error

		// --- Find Next Page Button with Retry ---
		for attempt := 0; attempt < maxRetries; attempt++ {
			nextPageButton, err = page.ElementX(`//*[@id="trend-table"]/div[2]/div/div[2]/span[3]/button`)
			if err == nil {
				break // Found the button
			}
			fmt.Printf("Attempt %d: Error finding next page button: %v. Retrying in %d seconds...\n", attempt+1, err, retryDelay)
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}
		if err != nil {
			fmt.Printf("Error finding next page button after %d attempts: %v. Assuming no more pages or page structure changed.\n", maxRetries, err)
			// Decide if this is a critical error or just means no more pages.
			// If it might just be the end, we can break the loop.
			break
		}

		// --- Check if Button is Disabled ---
		isDisabledAttr, err := nextPageButton.Attribute("disabled")
		if err != nil {
			// Error getting attribute, might be transient or element changed
			fmt.Printf("Error getting 'disabled' attribute: %v. Assuming button is not disabled and trying to proceed.\n", err)
			// Optionally add retry logic here as well
		}
		// Check if the attribute exists and is not nil (meaning it's disabled)
		if isDisabledAttr != nil {
			fmt.Println("Next page button is disabled. Reached the last page.")
			break // Exit the loop, no more pages
		}

		fmt.Println("Next page button not disabled")

		// --- Wait for Load/Stable State (Optional but recommended) ---
		// Consider adding a wait here if content loading is slow after interaction
		// page.MustWaitLoad() or page.MustWaitStable() with error handling

		// --- Get Table Body HTML with Retry ---
		var tableBody *rod.Element
		for attempt := 0; attempt < maxRetries; attempt++ {
			tableBody, err = page.ElementX(`/html/body/c-wiz/div/div[5]/div[1]/c-wiz/div/div[2]/div[1]/div[1]/div[1]/table/tbody[2]`)
			if err == nil {
				break // Found the table body
			}
			fmt.Printf("Attempt %d: Error finding table body: %v. Retrying in %d seconds...\n", attempt+1, err, retryDelay)
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}
		if err != nil {
			fmt.Printf("Error finding table body after %d attempts: %v. Skipping this page.\n", maxRetries, err)
			// Decide how to handle: skip page, return error?
			// For now, let's try clicking next anyway, maybe the table appears later.
		} else {
			// Extract HTML only if tableBody was found
			html, err := tableBody.HTML()
			if err != nil {
				fmt.Printf("Error getting table body HTML: %v. Skipping this page's data.\n", err)
			} else {
				rawData += html
			}
		}

		// --- Click Next Page Button with Retry ---
		for attempt := 0; attempt < maxRetries; attempt++ {
			err = nextPageButton.Click(proto.InputMouseButtonLeft, 1) // Using non-Must variant
			if err == nil {
				fmt.Println("Next page button clicked")
				break // Click successful
			}
			fmt.Printf("Attempt %d: Error clicking next page button: %v. Retrying in %d seconds...\n", attempt+1, err, retryDelay)
			// Wait before retrying click
			time.Sleep(time.Duration(retryDelay) * time.Second)
			// Re-find the button in case the page re-rendered after failed click attempt
			nextPageButton, err = page.ElementX(`//*[@id="trend-table"]/div[2]/div/div[2]/span[3]/button`)
			if err != nil {
				fmt.Printf("Error re-finding next page button after failed click: %v\n", err)
				break // Break inner retry loop if button can't be re-found
			}
		}
		if err != nil {
			fmt.Printf("Error clicking next page button after %d attempts: %v. Stopping pagination for this URL.\n", maxRetries, err)
			break // Exit the main pagination loop if click fails repeatedly
		}

		// --- Wait for Page Load/Stable after Click ---
		// Use WaitLoad or WaitStable with proper error handling and timeouts
		err = page.WaitLoad() // Using non-Must variant
		if err != nil {
			fmt.Printf("Error waiting for page load after click: %v. Proceeding cautiously.\n", err)
			// Decide if this is critical. Maybe add a short sleep as fallback.
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}

		// The loop continues to check the button state again
	}

	// fmt.Println("All rawData: ", rawData)
	data := sanitizeHTML(rawData)
	Adddata(data, countryName, regionName)
}
