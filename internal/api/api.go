//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//
package api

import (
	"net/http"

	"github.com/robertwtucker/document-host/internal/config"
	"github.com/robertwtucker/document-host/pkg/log"

	"github.com/go-chi/chi/v5"
)

func Routing(cfg *config.Configuration, logger log.Logger) http.Handler {
	logger.Debug("start: defining API routes")

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

	logger.Debug("end: defining API routes")

	return r
}
