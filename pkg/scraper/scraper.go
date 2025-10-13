package scraper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// Scraper wrapper over gocolly package which provides scraper logic for html parse
type Scraper struct {
	c      *colly.Collector
	filter *regexp.Regexp
	tags   string

	output      []string
	outputEvery int

	headers map[string]string
}

// New creates new scraper instance with specified config
func New(cfg *Config) (*Scraper, error) {
	// create regexp
	re, err := regexp.Compile(cfg.FilterPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create regexp: %w", err)
	}

	// configure colly collector
	c := colly.NewCollector(
		colly.UserAgent(cfg.UserAgent),
		colly.MaxDepth(cfg.MaxDepth),
		colly.Async(cfg.IsAsync),
	)

	// if async then add limits
	if cfg.IsAsync {
		asyncDelay := cfg.AsyncDelay
		if asyncDelay == 0 {
			asyncDelay = defaultAsyncDelay
		}

		asyncRequestLimit := cfg.AsyncRequestLimit
		if asyncRequestLimit == 0 {
			asyncRequestLimit = defaultAsyncRequestLimit
		}

		c.Limit(&colly.LimitRule{
			Parallelism: asyncRequestLimit,
			RandomDelay: asyncDelay,
		})
	}

	return &Scraper{
		c:           c,
		filter:      re,
		tags:        strings.Join(cfg.TagsToParse, ", "),
		headers:     cfg.Headers,
		output:      make([]string, 0, cfg.OutputEvery),
		outputEvery: cfg.OutputEvery,
	}, nil
}

// NewDefault creates new scraper using default config
func NewDefault() (*Scraper, error) {
	// using default config
	cfg := DefaultConfig()

	// create regexp
	re, err := regexp.Compile(cfg.FilterPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create regexp: %w", err)
	}

	// configure colly collector
	c := colly.NewCollector(
		colly.UserAgent(cfg.UserAgent),
		colly.MaxDepth(cfg.MaxDepth),
	)

	// set limits
	c.Limit(&colly.LimitRule{
		Parallelism: defaultAsyncRequestLimit,
		RandomDelay: defaultAsyncDelay,
	})

	return &Scraper{
		c:           c,
		filter:      re,
		tags:        strings.Join(cfg.TagsToParse, ", "),
		headers:     cfg.Headers,
		output:      make([]string, 0, cfg.OutputEvery),
		outputEvery: cfg.OutputEvery,
	}, nil
}
