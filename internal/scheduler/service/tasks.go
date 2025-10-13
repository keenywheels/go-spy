package service

import (
	"context"

	"github.com/keenywheels/go-spy/pkg/scraper"
	"golang.org/x/sync/errgroup"
)

const workerCount = 5

// ScrapeTask is the task that will be executed by the scheduler
func (s *Service) ScrapeTask() {
	op := "Service.ScrapeTask"

	sitesCh := make(chan string)

	gr, ctx := errgroup.WithContext(s.ctx)

	// job producer
	gr.Go(func() error {
		for _, site := range s.sites {
			sitesCh <- site
		}

		close(sitesCh)

		return nil
	})

	// start workers
	for i := 0; i < workerCount; i++ {
		gr.Go(func() error {
			return s.scrapeWorker(ctx, i, sitesCh)
		})
	}

	if err := gr.Wait(); err != nil {
		s.logger.Errorf("[%s] scrape task failed: %v", op, err)
	}
}

// scrapeWorker is the worker that will perform the scraping
func (s *Service) scrapeWorker(ctx context.Context, workerNum int, sitesCh chan string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	for site := range sitesCh {
		s.logger.Infof("[WORKER %d] start scraping site: %s", workerNum, site)

		scraper, err := scraper.New(s.scraperCfg)
		if err != nil {
			s.logger.Errorf("failed to create scraper: %v", err)

			return err
		}

		scraper.Visit(site)
	}

	return nil
}
