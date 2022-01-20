// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func RegisterHTTPHandlers(r chi.Router) {
	r.Get("/health", HealthCheckHandler)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}
