//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package model

// File represents the document as retrieved from the repository.
type File struct {
	Filename string            `bson:"filename" json:"filename"`
	Content  []byte            `bson:"-" json:"-"`
	Metadata map[string]string `bson:"metadata,omitempty" json:"metadata,omitempty"`
	Size     int64             `bson:"length,omitempty" json:"size,omitempty"`
}
