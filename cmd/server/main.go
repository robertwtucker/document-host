package main

import (
	"context"
	logger "log"
	"os"

	"github.com/robertwtucker/document-host/pkg/log"
)

// The current version of the application
const version = "0.1.0"

// The main program entry point
func main() {
	if err := run(); err != nil {
		logger.Println("startup error:", err)
		os.Exit(1)
	}
}

// Initializes the configuration (to-do) and logging
func initialize() (log.Logger, error) {
	logger := log.New().With(context.Background(), "version", version)

	return logger, nil
}

// Creates commpoennts and runs the application
func run() error {
	logger, err := initialize()
	if err != nil {
		logger.Errorf("failed to initialize app: %s", err)
		return err
	}

	logger.Info("Here's a standard logging message ")
	logger.Info("Here's an extended example ", "key:[", "value", "]")
	return nil
}
