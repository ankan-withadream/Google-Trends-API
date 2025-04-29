package services

import (
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/input"
	// "github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/robfig/cron/v3"

	"google-trends-api/internal/config"
)

func AutoScrap() {

	c := cron.New()
	c.AddFunc("@every 10m", func() { ExtractGoogleTrends() })
	c.Start()
	// for {
	// fmt.Println("Auto Scraping Google Trends")
	// ExtractGoogleTrends()
	// time.Sleep(1 * time.Minute)
	// }
}

func ExtractGoogleTrends() {
	fmt.Println("Extracting Google Trends")
	// Create a new browser instance in headless mode
	browser := rod.New().MustConnect()
	fmt.Println("browser")
	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage(config.GOOGLE_TRENDS_BASE_URL)

	// Wait for the page to load
	page.MustWaitStable()

	fmt.Println("page stable")

	// expand page length
	// expandDropdownPoint := page.MustElementX(`/html/body/c-wiz/div/div[5]/div[1]/c-wiz/div/div[2]/div[1]/div[1]/div[2]/div/div[1]/div[2]/div/div[1]/span[3]`).MustShape().OnePointInside()
	// page.Mouse.MustMoveTo(expandDropdownPoint.X, expandDropdownPoint.Y).MustClick(proto.InputMouseButtonRight)

	// fmt.Println("dropdown clicked")
	// sleep for 1 seconds
	// time.Sleep(1 * time.Second)

	// select page length option
	// xelement := page.MustElementX(`/html/body/c-wiz/div/div[5]/div[1]/c-wiz/div/div[2]/div[1]/div[1]/div[2]/div/div[1]/div[2]/div/div[2]/div/ul/li[3]/span[4]/span`)
	// fmt.Println("element focused", xelement)
	// xelement.Focus()
	// fmt.Println("element focused", xelement)
	// err2 := page.Keyboard.Press(input.Enter)
	// if err2 != nil {
	// 	fmt.Println("Error: ", err2)
	// }

	// pageLengthOptionPoint := page.MustElementX(`/html/body/c-wiz/div/div[5]/div[1]/c-wiz/div/div[2]/div[1]/div[1]/div[2]/div/div[1]/div[2]/div/div[2]/div/ul/li[3]/span[4]/span`).MustShape().OnePointInside()
	// page.Mouse.MustMoveTo(pageLengthOptionPoint.X, pageLengthOptionPoint.Y).MustClick(proto.InputMouseButtonRight)
	// page.MustElementX(`/html/body/c-wiz/div/div[5]/div[1]/c-wiz/div/div[2]/div[1]/div[1]/div[2]/div/div[1]/div[2]/div/div[2]/div/ul/li[3]/span[4]/span`).MustClick()

	// expand page length
	// expandPageLengthOption, err := page.ElementX(`/html/body/c-wiz/div/div[5]/div[1]/c-wiz/div/div[2]/div[1]/div[1]/div[2]/div/div[1]/div[2]/div/div[2]/div/ul/li[3]`)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// fmt.Println("option: ", pageLengthOptionPoint)
	// expandPageLengthOption.MustClick(proto.InputMouseButtonRight)
	// fmt.Println("option clicked")

	// page.MustElementX(`/html/body/c-wiz/div/div[5]/div[1]/c-wiz/div/div[2]/div[1]/div[1]/div[2]/div/div[1]/div[2]/div/div[2]/div/ul/li[3]`).MustClick()
	// fmt.Println("page expanded")
	nextPageButton := page.MustElementX(`//*[@id="trend-table"]/div[2]/div/div[2]/span[3]/button`)
	// nextPageButtonDisabled := page.MustElementX(`//*[@id="trend-table"]/div[2]/div/div[2]/span[3]/button[@disabled]`)
	var allData string
	ifdisabled := nextPageButton.MustAttribute("disabled")
	for ifdisabled == nil {
		fmt.Println("next page button not disabled")
		// Copy the table body
		tableBody := page.MustElementX(`/html/body/c-wiz/div/div[5]/div[1]/c-wiz/div/div[2]/div[1]/div[1]/div[1]/table/tbody[2]`)

		// RawHTML = tableBody.MustHTML()
		allData += tableBody.MustHTML()

		nextPageButton.MustClick()
		fmt.Println("next page button clicked")
		page.MustWaitLoad()
		ifdisabled = nextPageButton.MustAttribute("disabled")
	}
	fmt.Println("All data: ", allData)
	sanitizeHTML(allData)
}

func ExportGoogleTrends() {
	fmt.Println("Exporting Google Trends")
	// Create a new browser instance in headless mode
	browser := rod.New().MustConnect()
	fmt.Println("browser")
	defer browser.MustClose()

	// Create a new page
	page := browser.MustPage("https://trends.google.com/trending?geo=IN&status=active")

	// Wait for the page to load
	page.MustWaitStable()

	// Wait for and click the export button using a specific XPath
	exportButton := page.MustElementX(`//button[.//span[contains(text(), "Export")]]`)
	fmt.Println(exportButton)
	exportButton.MustClick()

	// Wait a moment for the dropdown to appear
	time.Sleep(1 * time.Second)

	// Click the Download CSV option
	csvOption := page.MustElementX(`//li[@role="menuitem" and @aria-label="Download CSV"]`)
	fmt.Println(csvOption)
	csvOption.MustClick()

	// Wait a bit longer to ensure the download starts
	// time.Sleep(5 * time.Second)
	// Wait for the download to start
	wait := browser.MustWaitDownload()

	// Save the downloaded file
	err := utils.OutputFile("downloaded_file1.csv", wait())
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded file saved")
	// Example: Read the contents of the downloaded file
	file, err := os.Open("downloaded_file1.csv")
	if err != nil {
		panic(err)
	}
	fmt.Println("file: ", file)
	defer file.Close()

	// Read the file contents
	buf := make([]byte, 1024)
	_, err = file.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))
}
