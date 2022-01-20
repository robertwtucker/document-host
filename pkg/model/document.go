//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//
package model

// Document represents the file and its metadata
type Document struct {
	ID         string `json:"id"`
	Filename   string `json:"filename"`
	MediaType  string `json:"mediaType"`
	FileBase64 string `json:"fileBase64"`
	URL        string `json:"url"`
	ShortLink  string `json:"shortLink"`
}
