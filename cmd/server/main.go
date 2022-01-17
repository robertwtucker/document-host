//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//
package main

import (
	"context"
	"flag"
	logger "log"
	"os"

	"github.com/robertwtucker/document-host/internal/config"
	"github.com/robertwtucker/document-host/pkg/log"
)

// version is the application's current version
const version = "0.1.0"

var (
	// flagConfig specifies the location of the configuration file
	flagConfig = flag.String("config", "./config/local.yaml", "path to config file")
	// flagDebug indicates wheter or not to print debug information
	flagDebug = flag.Bool("debug", false, "output debug information")
)

// main entry point
func main() {
	if err := run(); err != nil {
		logger.Println("startup error:", err)
		os.Exit(1)
	}
}

// initialize sets up the application configuration and logging
func initialize() (*config.Configuration, log.Logger, error) {
	var logger log.Logger

	flag.Parse()
	if *flagDebug {
		logger = log.NewDebug().With(context.Background(), "version", version)
	} else {
		logger = log.New().With(context.Background(), "version", version)
	}

	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		return nil, nil, err
	}

	return cfg, logger, nil
}

// run creates app commpoennts and executes the application
func run() error {
	// Set up configuration and logging
	cfg, logger, err := initialize()
	if err != nil {
		logger.Errorf("failed to initialize app: %s", err)
		return err
	}

	// TODO: Set up DB connection

	// TODO: Set up HTTP listener & handlers
	logger.Infof("app will listen on port %s", cfg.Server.Port)

	return nil
}
