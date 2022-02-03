//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// HealthCheckDatabaseHelper is the concrete implementation of the MongoDB helper
type HealthCheckDatabaseHelper struct {
	db *mongo.Database
}

// NewHealthCheckDatabaseHelper creates a new instance of the `HealthCheckDatabaseHelper`
func NewHealthCheckDatabaseHelper(db *mongo.Database) *HealthCheckDatabaseHelper {
	return &HealthCheckDatabaseHelper{db: db}
}

// CheckDB implements the use case interface
func (h HealthCheckDatabaseHelper) CheckDB(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return h.db.Client().Ping(ctx, nil)
}
