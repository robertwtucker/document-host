//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package usecase

import (
	"context"
	"fmt"

	"github.com/robertwtucker/document-host/internal/config"
	"github.com/robertwtucker/document-host/internal/document"
	"github.com/robertwtucker/document-host/pkg/model"
	"github.com/robertwtucker/document-host/pkg/shortlink"
)

// DocumentUseCase is the concrete implementation the use cases for the document repository.
type DocumentUseCase struct {
	documentRepo document.Repository
	shortLinkSvc shortlink.Service
	config       *config.Configuration
}

// NewDocumentUseCase creates a new instance of the `DocumentUseCase`.
func NewDocumentUseCase(repo document.Repository, svc shortlink.Service, cfg *config.Configuration) *DocumentUseCase {
	return &DocumentUseCase{
		documentRepo: repo,
		shortLinkSvc: svc,
		config:       cfg,
	}
}

// Create implements the use case interface.
func (d DocumentUseCase) Create(ctx context.Context, doc *model.Document) (*model.Document, error) {
	doc, err := d.documentRepo.Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	doc.URL = fmt.Sprintf("%s/%s", d.config.App.URL, doc.ID)
	slRequest := &shortlink.ServiceRequest{URL: doc.URL}
	if slResponse := d.shortLinkSvc.Shorten(ctx, slRequest); slResponse != nil {
		doc.ShortLink = slResponse.ShortLink
	}

	return doc, nil
}

// Get implements the use case interface.
func (d DocumentUseCase) Get(ctx context.Context, id string) (*model.File, error) {
	return d.documentRepo.Get(ctx, id)
}
