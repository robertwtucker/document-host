//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package usecase

import (
	"context"
	"testing"

	"github.com/robertwtucker/document-host/internal/document/mocks"
	"github.com/robertwtucker/document-host/pkg/model"
	"github.com/robertwtucker/document-host/pkg/shortlink"
	slmocks "github.com/robertwtucker/document-host/pkg/shortlink/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	repoIn := &model.Document{
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		FileBase64:  "VGVzdEZpbGU=",
	}
	repoOut := &model.Document{
		ID:          "61f0023ee260d827b7156c55",
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		FileBase64:  "",
		URL:         "http://dev.local/v1/documents/61f0023ee260d827b7156c55",
	}
	svcIn := &shortlink.ServiceRequest{
		URL: "http://dev.local/v1/documents/61f0023ee260d827b7156c55\"",
	}
	svcOut := &shortlink.ServiceResponse{
		URL:       "http://dev.local/v1/documents/61f0023ee260d827b7156c55\"",
		ShortLink: "https://tiny.one/yckaxkhx",
	}

	repo := &mocks.Repository{}
	repo.On("Create", context.Background(), repoIn).Return(repoOut)

	svc := &slmocks.Service{}
	svc.On("Shorten", context.Background(), svcIn).Return(svcOut)

	uc := NewDocumentUseCase(repo, svc)
	doc, err := uc.Create(context.Background(), repoIn)
	if assert.NoError(t, err) {
		assert.Equal(t, doc.ShortLink, "https://tiny.one/yckaxkhx")
	}
}

func TestGet(t *testing.T) {
	id := "61f0023ee260d827b7156c55"
	fileBytes := []byte("TestFile")
	file := &model.File{
		Filename: "test.pdf",
		Content:  fileBytes,
		Metadata: map[string]string{"contentType": "application/pdf"},
		Size:     42,
	}

	repo := &mocks.Repository{}
	repo.On("Get", context.Background(), id).Return(file, nil)

	svc := &slmocks.Service{}
	svc.On("Shorten", context.Background(), nil).Return(nil)

	uc := NewDocumentUseCase(repo, svc)
	out, err := uc.Get(context.Background(), id)
	if assert.NoError(t, err) {
		assert.Equal(t, out.Filename, "test.pdf")
		assert.Equal(t, out.Content, fileBytes)
		assert.Equal(t, out.Metadata, map[string]string{"contentType": "application/pdf"})
		assert.Equal(t, out.Size, int64(42))
	}
}
