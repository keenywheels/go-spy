package webapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	oas "github.com/keenywheels/go-spy/internal/ogen/api/v1"
	securityapi "github.com/keenywheels/go-spy/internal/webapp/delivery/http/security"
	api "github.com/keenywheels/go-spy/internal/webapp/delivery/http/v1"
	"github.com/keenywheels/go-spy/pkg/httpserver"
	"github.com/keenywheels/go-spy/pkg/httputils"
	"github.com/keenywheels/go-spy/pkg/logger"
	"github.com/keenywheels/go-spy/pkg/logger/zap"
	mw "github.com/keenywheels/go-spy/pkg/middleware"

	"golang.org/x/sync/errgroup"
)

// middleware allias for middleware funcs
type middleware func(http.Handler) http.Handler

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

// Run starts app
func (app *App) Run() error {
	// read config
	cfg, err := LoadConfig(app.opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}

	app.cfg = cfg
	app.initLogger()

	// create mux using ogen
	mux, err := app.initRouter()
	if err != nil {
		return fmt.Errorf("failed to create http ogen server: %v", err)
	}

	// create and run main http server
	apiSrv := app.createHttpServer(context.Background(), mux)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		app.logger.Infof("http api server is running on %s", apiSrv.GetAddr())
		return apiSrv.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		app.logger.Error("server error: %v", err)
		return err
	}

	return nil
}

// initLogger create new Logger based on config
func (app *App) initLogger() {
	logCfg := app.cfg.AppCfg.LoggerCfg
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

// initRouter creates router using ogen
func (app *App) initRouter() (http.Handler, error) {
	// prepare clients map for security handler
	clients := make(map[string]string, len(app.cfg.AppCfg.S2SCfg.Clients))
	for _, client := range app.cfg.AppCfg.S2SCfg.Clients {
		clients[client.Name] = client.Token
	}

	// create handler
	securityHandler := securityapi.New(clients)
	handler := api.New()

	// create custom handlers
	notFoundHandler := func(w http.ResponseWriter, r *http.Request) {
		httputils.NotFoundJSON(w)
	}

	errorHandler := func(_ context.Context, w http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, securityapi.ErrWrongToken) {
			httputils.ForbiddenJSON(w)
			return
		}

		httputils.BadRequestJSON(w)
	}

	// create ogen http server
	srv, err := oas.NewServer(handler, securityHandler,
		oas.WithNotFound(notFoundHandler),
		oas.WithErrorHandler(errorHandler),
	)
	if err != nil {
		return nil, err
	}

	// apply middlewares
	middlewares := app.prepareMiddlewares()

	var mux http.Handler = srv
	for _, m := range middlewares {
		mux = m(mux)
	}

	return mux, nil
}

// prepareMiddlewares generates all middleware to be used in app
func (app *App) prepareMiddlewares() []middleware {
	// prepare recover middleware
	recoverMw := func(next http.Handler) http.Handler {
		return mw.WithRecover(app.logger, next)
	}

	// prepare log mw
	logMw := func(next http.Handler) http.Handler {
		return mw.WithLogging(app.logger, next)
	}

	return []middleware{
		mw.WithContentTypeJSON,
		logMw,
		recoverMw,
	}
}

// createHttpServer creates new httpserver.Server instance using app config
func (app *App) createHttpServer(ctx context.Context, mux http.Handler) *httpserver.Server {
	opts := []httpserver.Option{}

	httpCfg := app.cfg.AppCfg.HttpCfg

	if httpCfg.ReadTimeout != 0 {
		opts = append(opts, httpserver.ReadTimeout(httpCfg.ReadTimeout))
	}

	if httpCfg.WriteTimeout != 0 {
		opts = append(opts, httpserver.WriteTimeout(httpCfg.WriteTimeout))
	}

	if httpCfg.ShutdownTimeout != 0 {
		opts = append(opts, httpserver.ShutdownTimeout(httpCfg.ShutdownTimeout))
	}

	if httpCfg.Port != "" || httpCfg.Host != "" {
		port := "8008"
		if httpCfg.Port != "" {
			port = httpCfg.Port
		}

		host := ""
		if httpCfg.Host != "" {
			host = httpCfg.Host
		}

		opts = append(opts, httpserver.Addr(host, port))
	}

	return httpserver.New(ctx, mux, opts...)
}
