//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package usecase

import (
	"context"

	"github.com/robertwtucker/document-host/internal/document"
	"github.com/robertwtucker/document-host/pkg/model"
)

// DocumentUseCase is the concrete implementation the use cases for the document repository
type DocumentUseCase struct {
	documentRepo document.Repository
}

// NewDocumentUseCase creates a new instance of the `DocumentUseCase`
func NewDocumentUseCase(documentRepo document.Repository) *DocumentUseCase {
	return &DocumentUseCase{
		documentRepo: documentRepo,
	}
}

// Create implements the use case interface
func (d DocumentUseCase) Create(ctx context.Context, doc *model.Document) (*model.Document, error) {
	aDoc, err := d.documentRepo.Create(ctx, doc)
	if err != nil {
		return nil, err
	}
	// TODO: Implement short link service call
	/*
		aDoc.URL = strings.sprintf(%s/%s, viper.GetString("app.url"), doc.ID)
		shortLink, err := url.shorten(doc.URL)
		if err != nil {
			return nil, err
		}
	*/
	return aDoc, nil
}

// Get implements the use case interface
func (d DocumentUseCase) Get(ctx context.Context, id string) (*model.File, error) {
	return d.documentRepo.Get(ctx, id)
}
