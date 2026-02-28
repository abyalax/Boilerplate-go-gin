package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/abyalax/Boilerplate-go-gin/internal/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env") // load .env

	// Get configuration from environment
	dbURL := getEnv("DATABASE_URL", "")
	port := getEnvInt("PORT", 4000)

	// Initialize application
	app, err := bootstrap.NewApp(dbURL, port)
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

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.TrimSpace(value)
	}
	return defaultValue
}

// getEnvInt gets an environment variable as integer with a default value
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Sscanf(strings.TrimSpace(value), "%d", &defaultValue)
	}
	return defaultValue
}
