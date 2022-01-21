//
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

// Handler implements the use case for the document resource
type Handler struct {
	useCase document.UseCase
}

// NewHandler creates a new `Handler` instance for the document use case
func NewHandler(uc document.UseCase) *Handler {
	return &Handler{useCase: uc}
}

// Create implements the use case interface
func (h *Handler) Create(c echo.Context) error {
	// TODO: implement me
	return c.JSON(201, nil)
}

// Get implements the use case interface
func (h *Handler) Get(c echo.Context) error {
	// TODO: implement me
	return c.JSON(200, nil)
}
