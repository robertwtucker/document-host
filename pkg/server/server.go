//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//
package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/robertwtucker/document-host/pkg/log"
)

// Default configuration settings
var (
	defaultAddr                   = ":8080"
	defaultReadTimeoutSeconds     = 10
	defaultShutdownTimeoutSeconds = 10
	defaultWriteTimeoutSeconds    = 5
)

// Config represents server-specific configuration items
type Config struct {
	Addr                       string
	ReadTimeoutSeconds         int
	ShutdownTimeoutSeconds int
	WriteTimeoutSeconds        int
}

// Start instantiates, configures and runs the server
func Start(cfg Config, r http.Handler, logger log.Logger) error {

	// Validate Config settings and update w/defaults, if needed
	if len(cfg.Addr) == 0 {
		cfg.Addr = defaultAddr
	}
	if cfg.ReadTimeoutSeconds == 0 {
		cfg.ReadTimeoutSeconds = defaultReadTimeoutSeconds
	}
	if cfg.ShutdownTimeoutSeconds == 0 {
		cfg.ShutdownTimeoutSeconds = defaultShutdownTimeoutSeconds
	}
	if cfg.WriteTimeoutSeconds == 0 {
		cfg.WriteTimeoutSeconds = defaultWriteTimeoutSeconds
	}

	// Initialize the server
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}

	// Start the listener using a goroutine
	go func() {
		logger.Infof("starting server, listening at %s", cfg.Addr)
		server.ListenAndServe()
	}()

	// Channel for server errors
	serverError := make(chan error, 1)

	// Channel for signal interrupts
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	select {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	case err := <-serverError:
		return errors.Wrap(err, "server error: %s")

	// Block main and wait for shutdown.
	case sig := <-sigint:
		logger.Infof("start shutdown: %v", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.ShutdownTimeoutSeconds)*time.Second)
		defer cancel()

		// Ask listener to shutdown and shed load
		err := server.Shutdown(ctx)
		if err != nil {
			logger.Infof(
				"graceful shutdown did not complete in %v sec : %v",
				cfg.ShutdownTimeoutSeconds,
				err)
			err = server.Close()
			return err
		}

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("system integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
