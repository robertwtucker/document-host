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
	"github.com/robertwtucker/document-host/internal/document"
	"github.com/robertwtucker/document-host/pkg/model"
)

// Handler implements the use case for the document resource.
type Handler struct {
	useCase document.UseCase
}

// NewHandler creates a new `Handler` instance for the document use case.
func NewHandler(uc document.UseCase) *Handler {
	return &Handler{useCase: uc}
}

// Create implements the use case interface.
func (h *Handler) Create(c echo.Context) error {
	input := new(model.Document)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	result, err := h.useCase.Create(context.Background(), input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set(echo.HeaderLocation, result.URL)
	return c.JSON(http.StatusCreated, result)
}

// Get implements the use case interface.
func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	file, err := h.useCase.Get(context.Background(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	contentType := file.Metadata["contentType"]
	if contentType == "" {
		contentType = echo.MIMEOctetStream
	}

	return c.Blob(http.StatusOK, contentType, file.Content)
}
