//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/robertwtucker/document-host/internal/healthcheck"
)

// Handler implements the use case for the healthcheck resource.
type Handler struct {
	useCase healthcheck.UseCase
}

// NewHandler creates a new `Handler` instance for the healthcheck use case.
func NewHandler(uc healthcheck.UseCase) *Handler {
	return &Handler{useCase: uc}
}

// Get implements the use case interface.
func (h *Handler) Get(c echo.Context) error {
	if err := h.useCase.Get(context.Background()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "OK")
}
