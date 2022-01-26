//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/robertwtucker/document-host/internal/document/mocks"
	"github.com/robertwtucker/document-host/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	input := &model.Document{
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		FileBase64:  "VGVzdEZpbGU=",
	}
	body, err := json.Marshal(input)
	assert.NoError(t, err)

	doc := &model.Document{
		ID:          "61f0023ee260d827b7156c55",
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		FileBase64:  "VGVzdEZpbGU=",
		URL:         "http://dev.local/v1/documents/61f0023ee260d827b7156c55",
		ShortLink:   "https://tiny.one/yckaxkhx",
	}

	e := echo.New()
	uc := new(mocks.UseCase)
	uc.On("Create", context.Background(), input).Return(doc, nil)
	RegisterHTTPHandlers(e, uc)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/documents", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c := e.NewContext(req, rec)
	c.SetPath("/v1/documents")
	h := NewHandler(uc)

	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, 201, rec.Code)
		assert.Equal(t, doc.URL, rec.Header().Get(echo.HeaderLocation))
	}
}

func TestGet(t *testing.T) {
	fileBytes := []byte("TestFile")
	file := &model.File{
		Filename: "test.pdf",
		Content:  fileBytes,
		Metadata: map[string]string{"contentType": "application/pdf"},
		Size:     42,
	}

	e := echo.New()
	uc := new(mocks.UseCase)
	uc.On("Get", context.Background(), "61f0023ee260d827b7156c55").Return(file, nil)
	RegisterHTTPHandlers(e, uc)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/documents/id", nil)

	c := e.NewContext(req, rec)
	c.SetPath("/v1/documents/:id")
	c.SetParamNames("id")
	c.SetParamValues("61f0023ee260d827b7156c55")
	h := NewHandler(uc)

	if assert.NoError(t, h.Get(c)) {
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "application/pdf", rec.Header().Get("Content-Type"))
		assert.Equal(t, string(fileBytes), rec.Body.String())
	}
}
