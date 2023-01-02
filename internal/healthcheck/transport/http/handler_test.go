//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/robertwtucker/document-host/internal/healthcheck/mocks"
	subject "github.com/robertwtucker/document-host/internal/healthcheck/transport/http"
	"github.com/stretchr/testify/assert"
)

func TestGetSuccess(t *testing.T) {
	e := echo.New()
	uc := new(mocks.UseCase)
	uc.On("Get", context.Background()).Return(nil)
	subject.RegisterHTTPHandlers(e, uc)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)

	c := e.NewContext(req, rec)
	c.SetPath("/health")
	h := subject.NewHandler(uc)

	if assert.NoError(t, h.Get(c)) {
		assert.Equal(t, 200, rec.Code)
	}
}
