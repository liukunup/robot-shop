package http

import (
	"backend/pkg/log"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	httpSrv  *http.Server
	host     string
	port     int
	certFile string
	keyFile  string
	logger   *log.Logger
}
type Option func(s *Server)

func NewServer(engine *gin.Engine, logger *log.Logger, opts ...Option) *Server {
	s := &Server{
		Engine: engine,
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}
func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}
func WithCertFiles(certFile string, keyFile string) Option {
	return func(s *Server) {
		s.certFile = certFile
		s.keyFile = keyFile
	}
}

func (s *Server) Start(ctx context.Context) error {
	if s.certFile != "" && s.keyFile != "" { // HTTPS
		if _, err := tls.LoadX509KeyPair(s.certFile, s.keyFile); err != nil {
			s.logger.Sugar().Fatalf("certificate or key file error: %v", err)
		}
		s.httpSrv = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
			Handler: s,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}
		if err := s.httpSrv.ListenAndServeTLS(s.certFile, s.keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Sugar().Fatalf("listen: %s\n", err)
		}
	} else { // HTTP
		s.httpSrv = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
			Handler: s,
		}
		if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Sugar().Fatalf("listen: %s\n", err)
		}
	}

	return nil
}
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Sugar().Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		s.logger.Sugar().Fatal("Server forced to shutdown: ", err)
	}

	s.logger.Sugar().Info("Server exiting")
	return nil
}
