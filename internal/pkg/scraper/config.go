package scraper

import (
	"time"
)

const (
	defaultOutputEvery       = 1000
	defaultMaxDepth          = 10
	defaultFilterPattern     = "^[A-Za-zА-Яа-яЁё]+$"
	defaultAsyncDelay        = 5 * time.Second
	defaultAsyncRequestLimit = 5
)

var (
	defaultTags      = []string{"div", "span", "p", "a", "h1", "h2", "h3", "h4", "h5", "h6"}
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	defaultHeaders   = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Encoding": "gzip, deflate, br, zstd",
		"Connection":      "keep-alive",
	}
)

// Config contains scraper setting
type Config struct {
	// OutputEvery specifies how often to output results
	OutputEvery int `mapstructure:"output_every"`
	// MaxDepth specifies the maximum depth to crawl
	MaxDepth int `mapstructure:"max_depth"`
	// FilterPattern specifies the regex pattern to filter words
	FilterPattern string `mapstructure:"filter_pattern"`
	// TagsToParse specifies the HTML tags to parse
	TagsToParse []string `mapstructure:"tags_to_parse"`

	// IsAsync shows is scraper should work in async mode
	IsAsync bool `mapstructure:"is_async"`
	// AsyncDelay shows delay between async requests
	AsyncDelay time.Duration `mapstructure:"async_delay"`
	// AsyncRequestLimit specifies the maximum number of concurrent requests
	AsyncRequestLimit int `mapstructure:"async_request_limit"`

	// Headers specifies the headers to include in requests
	Headers map[string]string `mapstructure:"headers"`
	// AllowedDomains allowed domains for scraping
	AllowedDomains []string `mapstructure:"allowed_domains"`
	// CacheDir specifies directory for cache
	CacheDir string `mapstructure:"cache_dir"`
	// CacheExpiration specifies the duration for cache expiration
	CacheExpiration time.Duration `mapstructure:"cache_expiration"`
	// ProxyURLs defines list of proxy urls
	ProxyURLs []string `mapstructure:"proxy_urls"`
	// UserAgent specifies the user agent for requests
	UserAgent string `mapstructure:"user_agent"`
}

// DefaultConfig returns new config with default values
func DefaultConfig() *Config {
	cfg := Config{
		IsAsync:           true,
		AsyncDelay:        defaultAsyncDelay,
		AsyncRequestLimit: defaultAsyncRequestLimit,
		OutputEvery:       defaultOutputEvery,
		MaxDepth:          defaultMaxDepth,
		FilterPattern:     defaultFilterPattern,
		TagsToParse:       defaultTags,
		Headers:           defaultHeaders,
		UserAgent:         defaultUserAgent,
	}

	return &cfg
}
