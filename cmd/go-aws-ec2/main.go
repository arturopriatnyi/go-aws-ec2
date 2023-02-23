package main

import (
	"context"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"go-aws-ec2/internal/http"
	"go-aws-ec2/pkg/counter"

	"github.com/prometheus/client_golang/prometheus"
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

	cms := counter.NewMemoryStore()
	cm := counter.NewManager(cms)

	http.MustRegisterMetrics(prometheus.DefaultRegisterer)

	s := &nethttp.Server{
		Addr:    ":10000",
		Handler: http.NewHandler(l, cm),
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
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
