package models

// ScraperEvent represents an event when the scraper gets data
type ScraperEvent struct {
	SiteName string `json:"site_name"`
	Msg      string `json:"msg"`
	Date     string `json:"date"`
	// TODO: подумать, что еще может понадобится; также согласовать формат с сервисом на ноде
}
