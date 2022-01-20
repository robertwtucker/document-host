//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//
package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	// dochttp "github.com/robertwtucker/document-host/internal/document/transport/http"
	health "github.com/robertwtucker/document-host/internal/healthcheck/transport/http"
	"github.com/robertwtucker/document-host/pkg/log"
)

type App struct {
	logger     log.Logger

	// documentUC document.UseCase
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
  if err := cv.validator.Struct(i); err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err.Error())
  }
  return nil
}

func NewApp(logger log.Logger) *App {
	// TODO:
	// db := initDB(cfg, logger)

	logger.Debug("start: wiring App components")


	logger.Debug("end: wiring App components")

	return &App{
		logger:     logger,

		// documentUC:
	}
}

func (a *App) Run() {

	a.logger.Debug("start: configuring server")
	timeout := viper.GetInt("server.timeout")

	// Echo setup
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(timeout) * time.Second,
	}))

	// HTTP endpoints
	health.RegisterHTTPHandlers(e)

	// API endpoints
	// r.Route("/v1", func(r chi.Router) {
	// 	dochttp.RegisterHTTPHandlers(r, a.documentUC)
	// })

	a.logger.Debug("end: configuring server")

	// Start server
	go func() {
		a.logger.Debug("starting server")
		err := e.Start(":" + viper.GetString("server.port"))
		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %+v", err)
		}
	}()

	// Channel for signal interrupts
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint
	// Interrupt received, attempt to shudtown gracefully
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
