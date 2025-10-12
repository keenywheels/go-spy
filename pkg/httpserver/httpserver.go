package httpserver

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// defaultErrorLogger logrus as default error logger for http.Server
var defaultErrorLogger = log.New(logrus.New().Writer(), "[httpserver]", log.LstdFlags)

// default params for net/http server
const (
	defaultShutdownTimeout = 5 * time.Second
	defaultReadTimeout     = 3 * time.Second
	defaultWriteTimeout    = 3 * time.Second
	defaultAddr            = ":8080"
)

// Server wrapper over net/http server
type Server struct {
	srv             *http.Server
	shutdownTimeout time.Duration
}

// New create new Server instance
func New(ctx context.Context, h http.Handler, opts ...Option) *Server {
	s := &Server{
		srv: &http.Server{
			Handler:      h,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			Addr:         defaultAddr,
			ErrorLog:     defaultErrorLogger,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Run run server with gracefull shutdown
func (s *Server) Run(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	// run http server
	g.Go(func() error {
		return s.srv.ListenAndServe()
	})

	// gracefull shutdown
	g.Go(func() error {
		<-gCtx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// GetAddr return server addr
func (s *Server) GetAddr() string {
	return s.srv.Addr
}
