//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//
package api

import (
	"net/http"

	"github.com/robertwtucker/document-host/internal/api/healthcheck"
	"github.com/robertwtucker/document-host/internal/config"
	"github.com/robertwtucker/document-host/pkg/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func Routing(cfg *config.Configuration, logger log.Logger) http.Handler {
	logger.Debug("start: defining API routes")

	// Chi setup
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Register routes
	healthcheck.RegisterHandlers(r)

	logger.Debug("end: defining API routes")

	return r
}
