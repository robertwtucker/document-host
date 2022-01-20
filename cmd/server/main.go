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
	"github.com/robertwtucker/document-host/pkg/server"
	"github.com/spf13/viper"
)

// version is the application's current version
const version = "0.1.0"

// main entry point
func main() {
	if err := run(); err != nil {
		logpkg.Println("startup error:", err)
		os.Exit(1)
	}
}

// initialize sets up the application configuration and logging
func initialize() (log.Logger, error) {
	var logger log.Logger

	config.Init()

	if debug := viper.GetBool("log.debug"); debug == true {
		logger = log.NewDebug().With(context.Background(), "version", version)
	} else {
		logger = log.New().With(context.Background(), "version", version)
	}

	if err := config.Load(logger); err != nil {
		return nil, err
	}

	return logger, nil
}

// run creates app commpoennts and executes the application
func run() error {
	// Set up configuration and logging
	logger, err := initialize()
	if err != nil {
		logpkg.Printf("failed to initialize app: %v \n", err)
		return err
	}

	// TODO: Set up DB connection

	// Set up API handlers
	r := api.Routing(logger)

	// Set up HTTP listener config
	serverConfig := &server.Config{
		Addr:                   viper.GetString("Server.Addr"),
		ReadTimeoutSeconds:     viper.GetInt("Server.ReadTimeoutSeconds"),
		ShutdownTimeoutSeconds: viper.GetInt("Server.ShutdownTimeoutSeconds"),
		WriteTimeoutSeconds:    viper.GetInt("Server.WriteTimeoutSeconds"),
	}

	// Start server
	if err := server.Start(*serverConfig, r, logger); err != nil {
		logger.Errorf("server error: %s", err)
	}

	return nil
}
