package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	l = log.New(os.Stdout, "HTTP server: ", log.LstdFlags|log.Lshortfile)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := &http.Server{
		Addr: ":10000",
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				if _, err := w.Write([]byte("Hello from AWS EC2!")); err != nil {
					l.Printf("response writing failed: %v", err)
				}
			},
		),
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Printf("starting failed: %v", err)
		}
	}()
	l.Println("started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		l.Println("shutting down gracefully")
	case <-ctx.Done():
		l.Println("context has terminated")
	}

	if err := s.Shutdown(ctx); err != nil {
		l.Printf("shutdown failed: %v", err)
	}
	l.Println("shut down")
}
