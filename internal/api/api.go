//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//
package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	// dochttp "github.com/robertwtucker/document-host/internal/document/transport/http"
	health "github.com/robertwtucker/document-host/internal/healthcheck/transport/http"
	"github.com/robertwtucker/document-host/pkg/log"
	"github.com/robertwtucker/document-host/pkg/server"
)

type App struct {
	logger     log.Logger
	validate   *validator.Validate

	// documentUC document.UseCase
}

var validate *validator.Validate

// func NewApp(cfg *config.Configuration, logger log.Logger) http.Handler {
func NewApp(logger log.Logger) *App {
	validate = validator.New()
	// TODO:
	// db := initDB(cfg, logger)

	logger.Debug("start: wiring App components")


	logger.Debug("end: wiring App components")

	return &App{
		logger:     logger,
		validate:   validate,

		// documentUC:
	}
}

func (a *App) Run() {

	a.logger.Debug("start: configuring server")

	// Chi setup
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// HTTP endpoints
	health.RegisterHTTPHandlers(r)

	// API endpoints
	// r.Route("/v1", func(r chi.Router) {
	// 	dochttp.RegisterHTTPHandlers(r, a.documentUC)
	// })

	// Set up HTTP listener config
	serverConfig := &server.Config{
		Addr:                   viper.GetString("Server.Addr"),
		ReadTimeoutSeconds:     viper.GetInt("Server.ReadTimeoutSeconds"),
		ShutdownTimeoutSeconds: viper.GetInt("Server.ShutdownTimeoutSeconds"),
		WriteTimeoutSeconds:    viper.GetInt("Server.WriteTimeoutSeconds"),
	}

	a.logger.Debug("end: configuring server")

	// Start server
	if err := server.Start(*serverConfig, r, a.logger); err != nil {
		a.logger.Errorf("server error: %s", err)
	}
}
