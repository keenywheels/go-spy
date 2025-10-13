package service

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
)

// initJobs initializes the scheduled job based on the cron pattern
func (s *Service) initJobs() error {
	// register job
	_, err := s.scheduler.NewJob(
		gocron.CronJob(s.cronPattern, false),
		gocron.NewTask(s.ScrapeTask),
	)
	if err != nil {
		return fmt.Errorf("failed to init job: %w", err)
	}

	return nil
}

// StartScheduler starts the job scheduler
func (s *Service) StartScheduler() error {
	s.scheduler.Start()

	<-s.ctx.Done()

	s.logger.Info("shutting down scheduler")

	return s.scheduler.Shutdown()
}
