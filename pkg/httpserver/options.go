package httpserver

import (
	"net"
	"time"
)

// Option -.
type Option func(*Server)

// Port -.
func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

// ReadTimeOut -.
func ReadTimeOut(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WriteTimeOut -.
func WriteTimeOut(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// ShutDownTimeOut -.
func ShutDownTimeOut(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}