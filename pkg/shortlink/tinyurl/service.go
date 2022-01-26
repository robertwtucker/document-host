//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package tinyurl

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-http-utils/headers"
	"github.com/robertwtucker/document-host/pkg/shortlink"
)

// tinyServiceURL is the API endpoint
const tinyServiceURL = "https://api.tinyurl.com/create"

// tinyURLService is the short link generation service implementation for TinyURL
type tinyURLService struct {
	APIKey     string
	Domain     string
	ServiceURL string
}

// NewTinyURLService returns a new instance of the TinyURL short link service
func NewTinyURLService(apiKey string, domain string) *tinyURLService {
	return &tinyURLService{
		APIKey:     apiKey,
		Domain:     domain,
		ServiceURL: tinyServiceURL,
	}
}

// tinyURLServiceResponse represents the response payload expected from the TinyURL service
type tinyURLServiceResponse struct {
	Data   tinyURLData   `json:"data"`
	Code   int           `json:"code"`
	Errors []interface{} `json:"errors,omitempty"`
}

// tinyURLData represents the TinyURL service's data payload
type tinyURLData struct {
	URL     string        `json:"url"`
	Domain  string        `json:"domain"`
	Alias   string        `json:"alias"`
	Tags    []interface{} `json:"tags,omitempty"`
	TinyURL string        `json:"tiny_url"`
}

// Shorten implements the Short Link generation service interface
func (ts tinyURLService) Shorten(ctx context.Context, req *shortlink.ServiceRequest) *shortlink.ServiceResponse {
	postBody, _ := json.Marshal(map[string]string{
		"url":    req.URL,
		"domain": ts.Domain,
	})

	hdr := http.Header{
		headers.Accept:        []string{"application/json"},
		headers.Authorization: []string{"Bearer " + ts.APIKey},
		headers.ContentType:   []string{"application/json"},
	}

	response, err := shortlink.Post(ts.ServiceURL, postBody, hdr)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	// Non-success responses won't parse. Log and return.
	if response.StatusCode != http.StatusOK {
		var respBody []byte
		respBody, _ = io.ReadAll(response.Body)
		log.Println("service returned non-ok status:", string(respBody))
		return nil
	}

	// Decode the service response
	tinyResponse := new(tinyURLServiceResponse)
	err = json.NewDecoder(response.Body).Decode(&tinyResponse)
	if err != nil {
		return nil
	}

	return toShortLink(tinyResponse)
}

// toShortLink is a utility function that converts the TinyURL response to the standard form
func toShortLink(tiny *tinyURLServiceResponse) *shortlink.ServiceResponse {
	return &shortlink.ServiceResponse{
		URL:       tiny.Data.URL,
		ShortLink: tiny.Data.TinyURL,
	}
}
