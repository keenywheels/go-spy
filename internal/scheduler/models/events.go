package models

import "time"

// ScraperEvent represents an event when the scraper gets data
type ScraperEvent struct {
	Site string    `json:"site"`
	Msg  string    `json:"msg"`
	Data time.Time `json:"data"`
	// TODO: подумать, что еще может понадобится; также согласовать формат с сервисом на ноде
}
