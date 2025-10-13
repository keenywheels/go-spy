package service

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/keenywheels/go-spy/pkg/logger"
	"github.com/keenywheels/go-spy/pkg/scraper"
)

// Service represent service layer of the application
type Service struct {
	cronPattern string
	scheduler   gocron.Scheduler

	sites []string

	ctx        context.Context
	logger     logger.Logger
	scraperCfg *scraper.Config
}

// New creates new service instance
func New(
	ctx context.Context,
	logger logger.Logger,
	scraperCfg *scraper.Config,
	cronPattern string,
	sites []string,
) (*Service, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	srv := Service{
		cronPattern: cronPattern,
		sites:       sites,
		scheduler:   scheduler,
		ctx:         ctx,
		logger:      logger,
		scraperCfg:  scraperCfg,
	}

	if err := srv.initJobs(); err != nil {
		return nil, fmt.Errorf("failed to init job: %w", err)
	}

	return &srv, nil
}
