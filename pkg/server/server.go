package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerConfig struct {
	Addr            string
	ShutdownTimeout time.Duration

	OnShutdown func()
}

type HTTPServer interface {
	Run(ctx context.Context) error
}

type server struct {
	srv *http.Server
	ln  net.Listener
	cfg ServerConfig
}

func NewHTTPServer(cfg ServerConfig, api http.Handler) (HTTPServer, error) {
	// start server
	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	mux := newRouter(api)
	srv := &http.Server{Handler: mux}

	// Ensure resources are released when the server shuts down
	srv.RegisterOnShutdown(func() {
		if cfg.OnShutdown != nil {
			cfg.OnShutdown()
		}
		log.Println("shutting down server")
	})

	return &server{
		cfg: cfg,
		srv: srv,
		ln:  ln,
	}, nil
}

func (s *server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", s.cfg.Addr)
		errCh <- s.srv.Serve(s.ln)
	}()

	// Make a signal context to also react to OS signals
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		// proceed to shutdown
	case <-sigCtx.Done():
		// proceed to shutdown
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}

	sdCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	_ = s.srv.Shutdown(sdCtx)

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	case <-time.After(s.cfg.ShutdownTimeout + time.Second):
		return context.DeadlineExceeded
	}
}
