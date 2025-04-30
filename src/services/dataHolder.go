package services

import (
	"google-trends-api/src/api/models"
)

var Data map[string][]models.TrendingItem

func Adddata(item []models.TrendingItem, countryName string, regionName string) {
	if Data == nil {
		Data = make(map[string][]models.TrendingItem)
	}
	if regionName != "" {
		Data[countryName+"-"+regionName] = item
	} else {
		Data[countryName] = item
	}
}
