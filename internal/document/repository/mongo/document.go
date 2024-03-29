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
	"strings"

	"github.com/robertwtucker/document-host/pkg/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DocumentRepository is the concrete implementation of the document repository.
type DocumentRepository struct {
	db *mongo.Database
}

// NewDocumentRepository creates a new instance of the `DocumentRepository`.
func NewDocumentRepository(db *mongo.Database) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create implements the use case interface.
func (d DocumentRepository) Create(_ context.Context, doc *model.Document) (*model.Document, error) {
	// Decode and store the file
	bucket, err := gridfs.NewBucket(d.db)
	if err != nil {
		log.Error("new bucket failed:", err)
		return nil, err
	}

	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(doc.FileBase64))
	opts := options.GridFSUpload().SetMetadata(bson.M{"contentType": doc.ContentType})
	fileID, err := bucket.UploadFromStream(doc.Filename, decoder, opts)
	if err != nil {
		log.Errorf("error uploading document to bucket: %v", err)
		return nil, err
	}

	// Update doc with ID and strip base64 element.
	doc.ID = primitive.ObjectID.Hex(fileID)
	doc.FileBase64 = ""

	return doc, nil
}

// Get implements the use case interface.
func (d DocumentRepository) Get(ctx context.Context, id string) (*model.File, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Get the file content.
	fileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("invalid id parameter '%s': %v", id, err)
		return nil, err
	}

	bucket, err := gridfs.NewBucket(d.db)
	if err != nil {
		log.Error("new bucket failed:", err)
		return nil, err
	}

	var buffer bytes.Buffer
	if _, err = bucket.DownloadToStream(fileID, &buffer); err != nil {
		log.Errorf("error streaming document from bucket: %v", err)
		return nil, err
	}

	// Get the file meta.
	cursor, err := bucket.Find(bson.M{"_id": fileID})
	if err != nil {
		log.Errorf("error finding document metadata: %v", err)
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	// There can be only one...
	var file = new(model.File)
	if cursor.Next(ctx) {
		if err = cursor.Decode(&file); err != nil {
			log.Errorf("error decoding document: %v", err)
			return nil, err
		}
		file.Content = buffer.Bytes()
	}

	return file, nil
}
