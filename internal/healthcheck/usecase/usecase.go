//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package usecase

import (
	"context"
	"time"

	"github.com/robertwtucker/document-host/internal/healthcheck"
)

// HealthCheckUseCase is the concrete implementation the use cases for the document repository
type HealthCheckUseCase struct {
	helper healthcheck.DatabaseHelper
}

// NewHealthCheckUseCase creates a new instance of the `HealthCheckUseCase`
func NewHealthCheckUseCase(dbh healthcheck.DatabaseHelper) *HealthCheckUseCase {
	return &HealthCheckUseCase{helper: dbh}
}

// Get implements the use case interface
func (h HealthCheckUseCase) Get(ctx context.Context) error {
	// Set a short timeout (default readiness timeout is usually only a second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*750)
	defer cancel()

	return h.helper.CheckDB(ctx)
}
