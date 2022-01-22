//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package mongo

import (
	"bytes"
	"context"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"

	"github.com/robertwtucker/document-host/pkg/model"
)

// DocumentRepository is the concrete implementation of the document repository
type DocumentRepository struct {
	db *mongo.Database
}

// NewDocumentRepository creates a new instance of the `DocumentRepository`
func NewDocumentRepository(db *mongo.Database) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create implements the use case interface
func (d DocumentRepository) Create(ctx context.Context, doc *model.Document) (*model.Document, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Decode and store the file
	bucket, _ := gridfs.NewBucket(d.db)
	var decoder = base64.NewDecoder(base64.StdEncoding, strings.NewReader(doc.FileBase64))
	opts := options.GridFSUpload().SetMetadata(bson.M{"contentType": doc.ContentType})
	fileID, err := bucket.UploadFromStream(doc.Filename, decoder, opts)
	if err != nil {
		return nil, err
	}

	// Update doc with ID and strip base64 element
	doc.ID = primitive.ObjectID.Hex(fileID)
	doc.FileBase64 = ""

	return doc, nil
}

// Get implements the use case interface
func (d DocumentRepository) Get(ctx context.Context, id string) (*model.File, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Get the file content
	fileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	bucket, _ := gridfs.NewBucket(d.db)
	var buffer bytes.Buffer
	_, err = bucket.DownloadToStream(fileID, &buffer)
	if err != nil {
		return nil, err
	}

	// Get the file meta
	var file = new(model.File)
	cursor, err := bucket.Find(bson.M{"_id": fileID})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			// Eat error and continue
		}
	}(cursor, ctx)
	// There can be only one...
	if cursor.Next(ctx) {
		err := cursor.Decode(&file)
		if err != nil {
			return nil, err
		}
		file.Content = buffer.Bytes()
	}

	return file, nil
}
