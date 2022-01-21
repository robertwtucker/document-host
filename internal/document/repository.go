//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package document

import (
	"context"

	"github.com/robertwtucker/document-host/pkg/model"
)

// Repository defines the operations supported by the document repository
type Repository interface {
	Create(ctx context.Context, doc *model.Document) (*model.Document, error)
	Get(ctx context.Context, id string) (*model.Document, error)
}
