//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package tinyurl

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-http-utils/headers"
	"github.com/robertwtucker/document-host/pkg/shortlink"
)

// TinyServiceURL is the API endpoint.
const TinyServiceURL = "https://api.tinyurl.com/create"

// TinyURLService is the short link generation service implementation for TinyURL.
type TinyURLService struct {
	APIKey     string
	Domain     string
	ServiceURL string
}

// NewTinyURLService returns a new instance of the TinyURL short link service.
func NewTinyURLService(apiKey string, domain string) *TinyURLService {
	return &TinyURLService{
		APIKey:     apiKey,
		Domain:     domain,
		ServiceURL: TinyServiceURL,
	}
}

// tinyURLServiceResponse represents the response payload expected from the TinyURL service.
type tinyURLServiceResponse struct {
	Data   tinyURLData   `json:"data"`
	Code   int           `json:"code"`
	Errors []interface{} `json:"errors,omitempty"`
}

// tinyURLData represents the TinyURL service's data payload.
type tinyURLData struct {
	URL     string        `json:"url"`
	Domain  string        `json:"domain"`
	Alias   string        `json:"alias"`
	Tags    []interface{} `json:"tags,omitempty"`
	TinyURL string        `json:"tiny_url"`
}

// Shorten implements the Short Link generation service interface.
func (ts TinyURLService) Shorten(_ context.Context, req *shortlink.ServiceRequest) *shortlink.ServiceResponse {
	postBody, _ := json.Marshal(map[string]string{
		"url":    req.URL,
		"domain": ts.Domain},
	)
	request, err := http.NewRequest(http.MethodPost, ts.ServiceURL, bytes.NewBuffer(postBody))
	if err != nil {
		return nil
	}
	request.Header.Set(headers.Accept, "application/json")
	request.Header.Set(headers.Authorization, "Bearer "+ts.APIKey)
	request.Header.Set(headers.ContentType, "application/json")

	client := &http.Client{Timeout: time.Second * 5}
	response, err := client.Do(request)
	if err != nil {
		return nil
	}
	defer func() { _ = response.Body.Close }()

	// Non-success responses won't parse. Log and return.
	if response.StatusCode != http.StatusOK {
		var respBody []byte
		respBody, _ = io.ReadAll(response.Body)
		log.Println("service returned non-ok status:", string(respBody))
		return nil
	}

	// Decode the service response
	tinyResponse := new(tinyURLServiceResponse)
	if err = json.NewDecoder(response.Body).Decode(&tinyResponse); err != nil {
		return nil
	}

	return toShortLink(tinyResponse)
}

// toShortLink is a utility function that converts the TinyURL response to the standard form.
func toShortLink(tiny *tinyURLServiceResponse) *shortlink.ServiceResponse {
	return &shortlink.ServiceResponse{
		URL:       tiny.Data.URL,
		ShortLink: tiny.Data.TinyURL,
	}
}
