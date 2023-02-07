package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap logger is not created: %v", err)
	}
	undo := zap.ReplaceGlobals(l)
	defer undo()

	s := &http.Server{
		Addr: ":10000",
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				if _, err := w.Write([]byte("Hello from AWS EC2!")); err != nil {
					l.Warn("response writing failed", zap.Error(err))
				}
			},
		),
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("HTTP server didn't start", zap.Error(err))
		}
	}()
	l.Info("HTTP server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		l.Info("shutting down gracefully")
	case <-ctx.Done():
		l.Info("context has terminated")
	}

	if err := s.Shutdown(ctx); err != nil {
		l.Fatal("HTTP server shutdown failed", zap.Error(err))
	}
	l.Info("HTTP server shut down")
}
