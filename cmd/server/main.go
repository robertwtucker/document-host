//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package main

import (
	"context"
	stdlog "log"
	"os"

	"github.com/robertwtucker/document-host/internal/api"
	"github.com/robertwtucker/document-host/internal/config"
	"github.com/robertwtucker/document-host/pkg/log"
	"github.com/spf13/viper"
)

// main entry point
func main() {
	if err := run(); err != nil {
		stdlog.Println("startup error:", err)
		os.Exit(1)
	}
}

// initialize sets up the application configuration and logging
func initialize() (*config.Configuration, log.Logger, error) {
	var logger log.Logger

	version := config.AppVersion()
	fmtVersion := version.Version + "-" + version.Revision
	stdlog.Printf("starting %s:%s\n", config.AppName, fmtVersion)

	config.Init()

	if debug := viper.GetBool("log.debug"); debug {
		logger = log.NewDebug().With(context.Background(), "version", fmtVersion)
	} else {
		logger = log.New().With(context.Background(), "version", fmtVersion)
	}

	cfg, err := config.Load(logger)
	if err != nil {
		return nil, nil, err
	}
	cfg.App.Version = fmtVersion

	return cfg, logger, nil
}

// run creates app commpoennts and executes the application
func run() error {
	// Set up configuration and logging
	cfg, logger, err := initialize()
	if err != nil {
		stdlog.Printf("failed to initialize app: %v \n", err)
		return err
	}
	logger.Debug("configuration:", cfg)

	// Initialize the API app
	app, err := api.NewApp(cfg, logger)
	if err != nil {
		logger.Errorf("error initializing app: %v", err)
		return err
	}

	// Run the app (server)
	app.Run()

	return nil
}
