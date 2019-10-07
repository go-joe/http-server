package joehttp

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-joe/joe"
	"go.uber.org/zap"
)

// RequestEvent corresponds to an HTTP request that was received by the server.
type RequestEvent struct {
	Header     http.Header
	Method     string
	URL        *url.URL
	RemoteAddr string
	Body       []byte
}

type server struct {
	http   *http.Server
	logger *zap.Logger
	conf   config
	events joe.EventEmitter
}

// Server returns a joe Module that runs an HTTP server to receive HTTP requests
// and emit them as events. This Module is mainly meant to be used to integrate
// a Bot with other systems that send events via HTTP (e.g. pull requests on GitHub).
func Server(path string, opts ...Option) joe.Module {
	return joe.ModuleFunc(func(joeConf *joe.Config) error {
		conf, err := newConf(path, joeConf, opts)
		if err != nil {
			return err
		}

		events := joeConf.EventEmitter()
		server := newServer(conf, events)

		server.logger.Info("Starting HTTP server", zap.String("addr", server.http.Addr))
		started := make(chan bool)
		go func() {
			started <- true
			server.Run()
		}()

		<-started

		joeConf.RegisterHandler(func(joe.ShutdownEvent) {
			server.Shutdown()
		})

		return nil
	})
}

func newServer(conf config, events joe.EventEmitter) *server {
	srv := &server{
		logger: conf.logger,
		events: events,
		conf:   conf,
	}

	srv.http = &http.Server{
		Addr:         conf.listenAddr,
		Handler:      http.HandlerFunc(srv.HTTPHandler),
		ErrorLog:     zap.NewStdLog(conf.logger),
		TLSConfig:    conf.tlsConf,
		ReadTimeout:  conf.readTimeout,
		WriteTimeout: conf.writeTimeout,
	}

	return srv
}

// Run starts the HTTP server to handle any incoming requests on the listen
// address that was configured.
func (s *server) Run() {
	var err error
	if s.conf.certFile == "" {
		err = s.http.ListenAndServe()
	} else {
		err = s.http.ListenAndServeTLS(s.conf.certFile, s.conf.keyFile)
	}

	if err != nil && err != http.ErrServerClosed {
		s.logger.Error("Failed to listen and serve requests", zap.Error(err)) // TODO: we want to see this error at startup!
	}
}

// HTTPHandler receives any incoming requests and emits them as events to the
// bots Brain.
func (s *server) HTTPHandler(_ http.ResponseWriter, r *http.Request) {
	clientIP := s.clientAddress(r)
	s.logger.Debug("Received HTTP request",
		zap.String("method", r.Method),
		zap.Stringer("url", r.URL),
		zap.String("remote_addr", clientIP),
	)

	event := RequestEvent{
		Header:     r.Header,
		Method:     r.Method,
		URL:        r.URL,
		RemoteAddr: clientIP,
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

// Shutdown gracefully shuts down the HTTP server without interrupting any
// active connections.
func (s *server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.logger.Info("HTTP server is shutting down gracefully")
	err := s.http.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Failed to shutdown server", zap.Error(err))
	}
}

func (s *server) clientAddress(req *http.Request) string {
	rip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		s.logger.Error("Error parsing RemoteAddr", zap.String("RemoteAddr", req.RemoteAddr))
		return req.RemoteAddr
	}
	ip := req.Header.Get(s.conf.trustedHeader)
	if ip == "" {
		ip = rip
	}
	return net.JoinHostPort(ip, port)
}
