package scheduler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/keenywheels/go-spy/internal/pkg/producer/kafka"
	"github.com/keenywheels/go-spy/internal/scheduler/repository/broker"
	"github.com/keenywheels/go-spy/internal/scheduler/service"
	"github.com/keenywheels/go-spy/pkg/logger"
	"github.com/keenywheels/go-spy/pkg/logger/zap"
	"golang.org/x/sync/errgroup"

	_ "net/http/pprof"
)

// App represent app environment
type App struct {
	opts *Options

	cfg    *Config
	logger logger.Logger
}

// New creates new application instance with options
func New() *App {
	opts := NewDefaultOpts()
	opts.LoadEnv()
	opts.LoadFlags()

	return &App{
		opts: opts,
	}
}

// Run starts application
func (app *App) Run() error {
	// read config
	cfg, err := LoadConfig(app.opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}

	app.cfg = cfg
	app.initLogger()
	defer func() {
		if err := app.logger.Close(); err != nil {
			log.Printf("failed to close logger: %v", err)
		}
	}()

	// create errgroup with signal context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// start system server if enabled
	if app.cfg.SchedulerCfg.SysSrvCfg.Enabled {
		app.logger.Info("starting system server")
		g.Go(func() error {
			return http.ListenAndServe(fmt.Sprintf(":%d", app.cfg.SchedulerCfg.SysSrvCfg.Port), nil)
		})
	}

	// create broker
	kafka, err := kafka.New(cfg.KafkaCfg.Brokers, kafka.Config{
		MaxRetry: cfg.KafkaCfg.MaxRetry,
	})
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}

	broker := broker.New(kafka, broker.Topics{
		ScraperData: cfg.KafkaCfg.Topics.ScraperData,
	})

	// create service layer
	srv, err := service.New(
		ctx,
		app.logger,
		&app.cfg.SchedulerCfg.ScraperCfg,
		app.cfg.SchedulerCfg.CronPattern,
		app.cfg.SchedulerCfg.WorkersCount,
		app.cfg.SchedulerCfg.Sites,
		broker,
	)
	if err != nil {
		return fmt.Errorf("failed to create service layer: %w", err)
	}

	g.Go(func() error {
		app.logger.Infof("starting scheduler with cron pattern: %s", app.cfg.SchedulerCfg.CronPattern)

		return srv.StartScheduler()
	})

	if err := g.Wait(); err != nil {
		app.logger.Error("app error: %v", err)

		return err
	}

	return nil
}

// initLogger create new Logger based on config
func (app *App) initLogger() {
	logCfg := app.cfg.SchedulerCfg.LoggerCfg
	opts := []zap.Option{}

	// setup opts
	if len(logCfg.LogLevel) != 0 {
		opts = append(opts, zap.LogLvl(logCfg.LogLevel))
	}

	if len(logCfg.Mode) != 0 {
		opts = append(opts, zap.Mode(logCfg.Mode))
	}

	if len(logCfg.Encoding) != 0 {
		opts = append(opts, zap.Encoding(logCfg.Encoding))
	}

	if len(logCfg.LogPath) != 0 {
		opts = append(opts, zap.LogPath(logCfg.LogPath))
	}

	if logCfg.MaxLogSize != 0 {
		opts = append(opts, zap.MaxLogSize(logCfg.MaxLogSize))
	}

	if logCfg.MaxLogBackups != 0 {
		opts = append(opts, zap.MaxLogBackups(logCfg.MaxLogBackups))
	}

	if logCfg.MaxLogAge != 0 {
		opts = append(opts, zap.MaxLogAge(logCfg.MaxLogAge))
	}

	// set logger
	app.logger = zap.New(opts...)
}
