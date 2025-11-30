package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/keenywheels/go-spy/pkg/logger"
)

// SetOutputCallback sets output callback function
func (s *Scraper) SetOutputCallback(cb outputCallback) {
	s.cb = cb
}

// Init initializes scraper
func (s *Scraper) Init(l logger.Logger) {
	// set headers
	s.c.OnRequest(func(r *colly.Request) {
		for k, v := range s.headers {
			r.Headers.Set(k, v)
		}
	})

	// parse tags for text
	s.c.OnHTML(s.tags, func(e *colly.HTMLElement) {
		txt := s.getDirectText(e)
		s.saveWords(s.filterText(txt))
	})

	// parse links
	s.c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))

		// filter links
		if s.filterLink(link) {
			s.mu.Lock()
			s.visited[link] = struct{}{}
			s.mu.Unlock()

			// add to queue if exists
			if s.q != nil {
				req, err := getCollyRequest(link, e.Request.Depth+1)
				if err != nil {
					l.Errorf("[scraper]: failed to create colly request: %v", err)
					return
				}

				s.q.AddRequest(req)
				return
			}

			// no queue -> call visit directly
			s.c.Visit(link)
		}
	})

	if s.isLogErrors {
		// logging errors
		s.c.OnError(func(r *colly.Response, err error) {
			l.Errorf("[Scraper]: got scraper error: URL=%s, Error=%v, Status=%d\n",
				r.Request.URL, err, r.StatusCode)
		})

		// logging not ok responses
		s.c.OnResponse(func(r *colly.Response) {
			if r.StatusCode < 200 && r.StatusCode >= 300 {
				l.Warnf("[Scraper]: got not ok status: URL=%s, Body=%s, Status=%d\n",
					r.Request.URL, string(r.Body[:100]), r.StatusCode)
			}
		})

	}
}

// Visit start scraping from specified url
func (s *Scraper) Visit(url string) error {
	s.prepareScraper("", url)

	// using queue if exists
	if s.q != nil {
		req, err := getCollyRequest(url, 0)
		if err != nil {
			return fmt.Errorf("failed to create colly request: %w", err)
		}

		s.q.AddRequest(req)
		s.q.Run(s.c)
	}

	return s.c.Visit(url)
}

// VisitWithSiteName start scraping from specified url
func (s *Scraper) VisitWithSiteName(url string, siteName string) error {
	s.prepareScraper(siteName, url)

	// using queue if exists
	if s.q != nil {
		req, err := getCollyRequest(url, 0)
		if err != nil {
			return fmt.Errorf("failed to create colly request: %w", err)
		}

		s.q.AddRequest(req)
		s.q.Run(s.c)
	}

	return s.c.Visit(url)
}

// Flush flushes remaining output
func (s *Scraper) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.output) > 0 {
		msg := strings.Join(s.output, " ")
		s.cb(msg)
		s.output = s.output[:0]
	}
}

// getDirectText get only direct text in element
func (s *Scraper) getDirectText(e *colly.HTMLElement) string {
	// check if leaf elem
	isLeaf := true
	e.DOM.Children().Each(func(i int, s *goquery.Selection) {
		isLeaf = false
	})

	// leaf -> return text
	if isLeaf {
		return e.Text
	}

	// not leaf -> get only direct text in element
	var txt strings.Builder

	e.DOM.Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "#text" {
			txt.WriteString(s.Text())
		}
	})

	return txt.String()
}

// filterText filters text and returns only valid words
func (s *Scraper) filterText(txt string) []string {
	parsedWords := make([]string, 0, 100)

	words := strings.FieldsSeq(txt)
	for w := range words {
		w = strings.ToLower(w)
		if !s.filter.MatchString(w) {
			continue
		}

		parsedWords = append(parsedWords, w)
	}

	return parsedWords
}

// saveWords saves words to output
func (s *Scraper) saveWords(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.output = append(s.output, words...)

	if len(s.output) >= s.outputEvery {
		msg := strings.Join(s.output, " ")
		s.cb(msg)
		s.output = s.output[:0]
	}
}

// prepareScraper prepares the scraper
func (s *Scraper) prepareScraper(name, siteURL string) {
	s.visited = make(map[string]struct{})

	var domain string
	parsedURL, err := url.Parse(siteURL)
	if err == nil {
		domain = parsedURL.Hostname()
	}

	s.siteName = name
	s.siteDomain = domain
}

// filterLink checks if link belongs to the same domain
func (s *Scraper) filterLink(link string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.visited[link]; ok {
		return false
	}

	if strings.Contains(link, s.siteDomain) || strings.Contains(link, s.siteName) {
		return true
	}

	return false
}

// getCollyRequest creates colly request from url
func getCollyRequest(u string, d int) (*colly.Request, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url=%s: %w", u, err)
	}

	return &colly.Request{
		URL:    parsedURL,
		Depth:  d,
		Method: "GET",
	}, nil
}
