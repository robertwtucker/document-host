//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package mongo

import (
	"context"

	"github.com/robertwtucker/document-host/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Document represents the repository version of the document model
type Document struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Filename   string             `bson:"filename"`
	MediaType  string             `bson:"mediaType"`
	FileBase64 string             `bson:"fileBase64"`
	URL        string             `bson:"url"`
	ShortLink  string             `bson:"shortLink"`
}

// DocumentRepository is the concrete implementation of the document repository
type DocumentRepository struct {
	db *mongo.Collection
}

// NewDocumentRepository creates a new instance of the `DocumentRepository`
func NewDocumentRepository(db *mongo.Database, collection string) *DocumentRepository {
	return &DocumentRepository{
		db: db.Collection(collection),
	}
}

// Create implements the use case interface
func (d DocumentRepository) Create(ctx context.Context, doc *model.Document) (*model.Document, error) {
	// TODO: implement me
	return nil, nil
}

// Get implements the use case interface
func (d DocumentRepository) Get(ctx context.Context, id string) (*model.Document, error) {
	// TODO: implement me
	return nil, nil
}
