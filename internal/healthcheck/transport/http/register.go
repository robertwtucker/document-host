// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package http

import (
	"github.com/labstack/echo/v4"
	"github.com/robertwtucker/document-host/internal/healthcheck"
)

// RegisterHTTPHandlers maps the handler for the health check endpoint
func RegisterHTTPHandlers(e *echo.Echo, uc healthcheck.UseCase) {
	h := NewHandler(uc)
	e.GET("/health", h.Get)
}
