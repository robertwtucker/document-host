// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package http

import (
	"github.com/labstack/echo/v4"
	"github.com/robertwtucker/document-host/internal/document"
)

// RegisterHTTPHandlers maps the handlers for the document resource use cases
func RegisterHTTPHandlers(e *echo.Echo, uc document.UseCase) {
	h := NewHandler(uc)
	r := e.Group("/v1")
	r.POST("/documents", h.Create)
	r.GET("/documents/:id", h.Get)
}
