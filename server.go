package joehttp

import (
	"context"
	"github.com/go-joe/joe"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type server struct {
	http   *http.Server
	logger *zap.Logger
	events joe.EventEmitter
}

func Server(path string) joe.Module {
	return func(conf *joe.Config) error {
		logger := conf.Logger("http")
		events := conf.EventEmitter()
		server, err := newServer(path, events, logger)
		if err != nil {
			return err
		}

		go server.Run()
		conf.RegisterHandler(func(joe.ShutdownEvent) {
			server.Shutdown()
		})

		return nil
	}
}

func newServer(addr string, events joe.EventEmitter, logger *zap.Logger) (*server, error) {
	srv := &server{
		logger: logger,
		events: events,
	}

	srv.http = &http.Server{
		Addr:              addr,
		Handler:           http.HandlerFunc(srv.HTTPHandler),
		ErrorLog:          zap.NewStdLog(logger),
		TLSConfig:         nil, // TODO
		ReadTimeout:       0,   // TODO
		ReadHeaderTimeout: 0,   // TODO
		WriteTimeout:      0,   // TODO
		IdleTimeout:       0,   // TODO
		MaxHeaderBytes:    0,   // TODO
	}

	return srv, nil
}

func (s *server) Run() {
	s.logger.Info("Starting HTTP server", zap.String("addr", s.http.Addr))
	err := s.http.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.Error("Failed to listen and serve requests", zap.Error(err)) // TODO: we want to see this error at startup!
	}
}

func (s *server) HTTPHandler(_ http.ResponseWriter, r *http.Request) {
	s.logger.Debug("Received HTTP request",
		zap.String("method", r.Method),
		zap.Stringer("url", r.URL),
		zap.String("remote_addr", r.RemoteAddr),
	)

	event := RequestEvent{
		Method:     r.Method,
		URL:        r.URL,
		RemoteAddr: r.RemoteAddr,
	}

	var err error
	if r.Body != nil {
		event.Body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			s.logger.Error("Failed to read request body")
		}
	}

	s.events.Emit(event)
}

func (s *server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.logger.Info("HTTP server is shutting down gracefully")
	err := s.http.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Failed to shutdown server", zap.Error(err))
	}
}
