package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abyalax/Boilerplate-go-gin/src/bootstrap"
	"github.com/abyalax/Boilerplate-go-gin/src/config/env"
)

func main() {
	cfg, err := env.Load()

	// Initialize application
	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize app: %v\n", err)
		os.Exit(1)
	}

	// Start application in a goroutine
	go func() {
		if err := app.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to stop app: %v\n", err)
		os.Exit(1)
	}
}
