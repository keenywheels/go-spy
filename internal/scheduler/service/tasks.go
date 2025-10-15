package service

import (
	"context"
	"time"

	"github.com/keenywheels/go-spy/internal/pkg/scraper"
	"github.com/keenywheels/go-spy/internal/scheduler/models"
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

	scrapeStart := time.Now()

	// start workers
	for i := 0; i < workerCount; i++ {
		gr.Go(func() error {
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
	start time.Time,
	sitesCh chan string,
) error {
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

			cb := func(msg string) {
				s.logger.Infof("[WORKER %d] sending data to kafka\n", workerNum)

				if err := s.broker.SendScraperData(models.ScraperEvent{
					Site: site,
					Msg:  msg,
					Data: start,
				}); err != nil {
					s.logger.Errorf("[WORKER %d] failed to send data to kafka: %v", workerNum, err)
				}
			}

			scraper.SetOutputCallback(cb)
			scraper.Init()

			if err := scraper.Visit(site); err != nil {
				s.logger.Errorf("[WORKER %d] failed to visit site %s: %v", workerNum, site, err)
				continue
			}
		}
	}
}
