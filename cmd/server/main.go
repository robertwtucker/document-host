//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package main

import (
	"context"
	logpkg "log"
	"os"

	"github.com/robertwtucker/document-host/internal/api"
	"github.com/robertwtucker/document-host/internal/config"
	"github.com/robertwtucker/document-host/pkg/log"
	"github.com/spf13/viper"
)

// version is the application's current version
const version = "0.2.0"

// main entry point
func main() {
	if err := run(); err != nil {
		logpkg.Println("startup error:", err)
		os.Exit(1)
	}
}

// initialize sets up the application configuration and logging
func initialize() (*config.Configuration, log.Logger, error) {
	var logger log.Logger

	config.Init()

	if debug := viper.GetBool("log.debug"); debug {
		logger = log.NewDebug().With(context.Background(), "version", version)
	} else {
		logger = log.New().With(context.Background(), "version", version)
	}

	cfg, err := config.Load(logger)
	if err != nil {
		return nil, nil, err
	}

	return cfg, logger, nil
}

// run creates app commpoennts and executes the application
func run() error {
	// Set up configuration and logging
	cfg, logger, err := initialize()
	if err != nil {
		logpkg.Printf("failed to initialize app: %v \n", err)
		return err
	}
	logger.Debug("configuration:", cfg)

	// Initialize the API app
	app := api.NewApp(cfg, logger)

	// Run the app (server)
	app.Run()

	return nil
}
