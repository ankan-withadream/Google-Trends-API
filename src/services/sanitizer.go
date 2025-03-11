package services

import (
	"encoding/json"
	"fmt"
	"google-trends-api/src/api/models"

	// "google-trends-api/src/services"
	"strings"
)

func SanitizeHTML() {
	htmlContent := RawHTML
	table_rows := strings.SplitN(htmlContent, "</tr>", -1)
	fmt.Println("Rows: ", len(table_rows))
	items := []models.TrendingItem{}
	// var relatedTerms []string

	for i := 0; i < len(table_rows)-1; i++ {
		row := table_rows[i]
		var currentItem models.TrendingItem
		table_datas := strings.SplitN(row, "</td>", 5)

		// div1 holidng keyword and search voulme
		div1 := table_datas[1]
		div1 = strings.Replace(div1, "</div", "", -1)
		div1Slices := strings.SplitN(div1, ">", -1)
		// div2 holding search volume increase
		div2 := table_datas[2]
		div2 = strings.Replace(div2, "</div", "", -1)
		div2Slices := strings.SplitN(div2, ">", -1)
		// div3 holding started from now and active lasted
		div3 := table_datas[3]
		div3 = strings.Replace(div3, "</div", "", -1)
		div3Slices := strings.SplitN(div3, ">", -1)

		currentItem.Keyword = div1Slices[2]
		currentItem.SearchVolume = div1Slices[6]
		currentItem.SearchVolumeIncrease = div2Slices[8]
		currentItem.StartedFromNow = div3Slices[2]
		currentItem.ActiveLasted = div1Slices[14]

		items = append(items, currentItem)
	}
	fmt.Println("Items: ", items)

	jsonBytes, err := json.Marshal(items)
	if err != nil {
		// return "", err
		fmt.Println("Error: ", err)
	}

	SanitizedData = string(jsonBytes)
}
