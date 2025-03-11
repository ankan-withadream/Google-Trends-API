package models

type TrendingItem struct {
	Keyword              string `json:"keyword"`
	SearchVolume         string `json:"search_volume"`
	SearchVolumeIncrease string `json:"search_volume_increase"`
	StartedFromNow       string `json:"started_from_now"`
	ActiveLasted         string `json:"active_lasted"`
}
