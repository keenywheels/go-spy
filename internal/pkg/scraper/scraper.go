package scraper

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

// outputCallback is a callback function that processes output
type outputCallback func(string)

// Scraper wrapper over gocolly package which provides scraper logic for html parse
type Scraper struct {
	c *colly.Collector
	q *queue.Queue

	filter *regexp.Regexp
	tags   string

	siteName   string
	siteDomain string
	visited    map[string]struct{}

	output      []string
	outputEvery int
	cb          outputCallback
	mu          sync.Mutex

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

	// set queue if enabled
	var q *queue.Queue

	if cfg.Queue.Enabled {
		q, _ = queue.New(
			cfg.Queue.ThreadNumber,
			&queue.InMemoryQueueStorage{MaxSize: cfg.Queue.MaxSize},
		)
	}

	return &Scraper{
		c:           c,
		q:           q,
		filter:      re,
		tags:        strings.Join(cfg.TagsToParse, ", "),
		headers:     cfg.Headers,
		output:      make([]string, 0, cfg.OutputEvery),
		outputEvery: cfg.OutputEvery,
		cb:          defaultOutputCallback,
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
		cb:          defaultOutputCallback,
	}, nil
}

// defaultOutputCallback is the default output callback function
func defaultOutputCallback(msg string) {
	fmt.Printf("RESULT: %s\n", msg)
}
