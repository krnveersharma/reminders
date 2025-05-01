package dto

type MonthCount struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type DayCount struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}
