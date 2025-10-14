package service

import (
	"context"

	"github.com/keenywheels/go-spy/internal/pkg/scraper"
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
	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("[WORKER %d] received done signal", workerNum)
			return ctx.Err()
		case site, ok := <-sitesCh:
			if !ok {
				s.logger.Infof("[WORKER %d] no available sites to scrape, stopping...", workerNum)
				return nil
			}

			s.logger.Infof("[WORKER %d] start scraping site: %s", workerNum, site)

			scraper, err := scraper.New(s.scraperCfg)
			if err != nil {
				s.logger.Errorf("failed to create scraper: %v", err)
				continue
			}

			// TODO: заменить на cb с кафкой
			cb := func(msg string) {
				s.logger.Infof("RESULT: %s\n", msg)
			}

			scraper.SetOutputCallback(cb)
			// TODO: заменить на cb с кафкой

			scraper.Init()

			if err := scraper.Visit(site); err != nil {
				s.logger.Errorf("[WORKER %d] failed to visit site %s: %v", workerNum, site, err)
				continue
			}
		}
	}
}
