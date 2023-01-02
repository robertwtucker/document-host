//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package usecase_test

import (
	"context"
	"testing"

	"github.com/robertwtucker/document-host/internal/healthcheck/mocks"
	subject "github.com/robertwtucker/document-host/internal/healthcheck/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	dbh := new(mocks.DatabaseHelper)
	dbh.On("CheckDB", mock.Anything).Return(nil)

	uc := subject.NewHealthCheckUseCase(dbh)
	assert.NoError(t, uc.Get(context.Background()))
}
