package services

import (
	"fmt"
	"google-trends-api/src/api/models"
)

var Data map[string][]models.TrendingItem

func Adddata(item []models.TrendingItem, countryName string, regionName string) {
	if Data == nil {
		fmt.Println("Initializing Data map")
		Data = make(map[string][]models.TrendingItem)
	}
	if regionName != "" {
		fmt.Printf("Adding data for %s-%s\n", countryName, regionName)
		Data[countryName+"-"+regionName] = item
	} else {
		fmt.Printf("Adding data for %s\n", countryName)
		Data[countryName] = item
	}
}
