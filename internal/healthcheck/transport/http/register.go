// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterHTTPHandlers maps the handler for the health check endpoint
func RegisterHTTPHandlers(e *echo.Echo) {
	e.GET("/health", healthCheckHandler)
}

func healthCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
