package models


type Trends struct {
	ID          uint   `gorm:"primaryKey"`
	Keyword     string `json:"keyword"`
	SearchVolume int    `json:"search_volume"`
	SearchVolumeIncrease int `json:"search_volume_increase"`
	StartedFromNow int `json:"started_from_now"`
	ActiveLasted int `json:"active_lasted"`
	SimilarKeywords []string `json:"similar_keywords"`
}