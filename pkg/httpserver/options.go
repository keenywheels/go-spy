package httpserver

import (
	"log"
	"net"
	"time"
)

type Option func(*Server)

// ShutdownTimeout set shutdownTimeout
func ShutdownTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = t
	}
}

// ReadTimeout set ReadTimeout in net/http server
func ReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.srv.ReadTimeout = t
	}
}

// WriteTimeout set WriteTimeout in net/http server
func WriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.srv.WriteTimeout = t
	}
}

// Addr set addr in "[host]:[port]" format
func Addr(host, port string) Option {
	return func(s *Server) {
		s.srv.Addr = net.JoinHostPort(host, port)
	}
}

// ErrorLog set ErrorLog in net/http server
func ErrorLog(l *log.Logger) Option {
	return func(s *Server) {
		s.srv.ErrorLog = l
	}
}
