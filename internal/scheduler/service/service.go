package service

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/keenywheels/go-spy/internal/pkg/scraper"
	"github.com/keenywheels/go-spy/internal/scheduler/models"
	"github.com/keenywheels/go-spy/pkg/logger"
)

// IBroker represents broker interface
type IBroker interface {
	SendScraperData(event models.ScraperEvent) error
}

// Service represent service layer of the application
type Service struct {
	cronPattern string
	scheduler   gocron.Scheduler

	sites map[string]string

	ctx        context.Context
	logger     logger.Logger
	scraperCfg *scraper.Config

	broker IBroker
}

// New creates new service instance
func New(
	ctx context.Context,
	logger logger.Logger,
	scraperCfg *scraper.Config,
	cronPattern string,
	sites map[string]string,
	broker IBroker,
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
		broker:      broker,
	}

	if err := srv.initJobs(); err != nil {
		return nil, fmt.Errorf("failed to init job: %w", err)
	}

	return &srv, nil
}
