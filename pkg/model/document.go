//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package model

// Document represents the document file and its metadata.
type Document struct {
	ID          string `json:"id,omitempty"`
	Filename    string `json:"filename" validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
	FileBase64  string `json:"fileBase64,omitempty" validate:"required"`
	URL         string `json:"url,omitempty"`
	ShortLink   string `json:"shortLink,omitempty"`
}
