package service

import (
	"context"
	"fmt"
	"time"

	"github.com/keenywheels/go-spy/internal/pkg/scraper"
	"github.com/keenywheels/go-spy/internal/scheduler/models"
	"golang.org/x/sync/errgroup"
)

// siteMsg represents message from scraped site
type siteMsg struct {
	SiteName string
	SiteUrl  string
}

// ScrapeTask is the task that will be executed by the scheduler
func (s *Service) ScrapeTask() {
	op := "Service.ScrapeTask"

	sitesCh := make(chan siteMsg)

	gr, ctx := errgroup.WithContext(s.ctx)

	// job producer
	gr.Go(func() error {
		for name, url := range s.sites {
			sitesCh <- siteMsg{
				SiteName: name,
				SiteUrl:  url,
			}
		}

		close(sitesCh)

		return nil
	})

	scrapeStart := time.Now().Format("02-01-2006")

	// start workers
	for i := 0; i < s.workersCount; i++ {
		gr.Go(func() error {
			s.logger.Infof("[%s] starting scrape worker %d", op, i)
			return s.scrapeWorker(ctx, i, scrapeStart, sitesCh)
		})
	}

	if err := gr.Wait(); err != nil {
		s.logger.Errorf("[%s] scrape task failed: %v", op, err)
	}
}

// scrapeWorker is the worker that will perform the scraping
func (s *Service) scrapeWorker(
	ctx context.Context,
	workerNum int,
	start string,
	sitesCh chan siteMsg,
) error {
	op := fmt.Sprintf("WORKER %d", workerNum)

	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("[%s] received done signal", op)
			return ctx.Err()
		case site, ok := <-sitesCh:
			if !ok {
				s.logger.Infof("[%s] no available sites to scrape, stopping...", op)
				return nil
			}

			s.logger.Infof("[%s] start scraping site: %v", op, site)

			scraper, err := scraper.New(s.scraperCfg)
			if err != nil {
				s.logger.Errorf("failed to create scraper: %v", err)
				continue
			}

			cb := func(msg string) {
				s.logger.Infof("[%s] sending data to kafka", op)

				if err := s.broker.SendScraperData(models.ScraperEvent{
					SiteName: site.SiteName,
					Msg:      msg,
					Date:     start,
				}); err != nil {
					s.logger.Errorf("[%s] failed to send data to kafka: %v", op, err)
				}
			}

			scraper.SetOutputCallback(cb)
			scraper.Init()

			if err := scraper.VisitWithSiteName(site.SiteUrl, site.SiteName); err != nil {
				s.logger.Errorf("[%s] failed to visit site %s: %v", op, site, err)
				scraper.Flush()

				continue
			}

			scraper.Flush()
		}
	}
}
