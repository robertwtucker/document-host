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
	Data struct {
		URL     string        `json:"url"`
		Domain  string        `json:"domain"`
		Alias   string        `json:"alias"`
		Tags    []interface{} `json:"tags,omitempty"`
		TinyURL string        `json:"tiny_url"`
	} `json:"data"`
	Code   int           `json:"code"`
	Errors []interface{} `json:"errors,omitempty"`
}

// Shorten implements the Short Link generation service interface
func (ts tinyURLService) Shorten(ctx context.Context, req *shortlink.ServiceRequest) *shortlink.ServiceResponse {
	postBody, _ := json.Marshal(map[string]string{
		"url":    req.URL,
		"domain": ts.Domain,
	})

	svcRequest, err := http.NewRequest(http.MethodPost, ts.ServiceURL, bytes.NewBuffer(postBody))
	if err != nil {
		return nil
	}
	svcRequest.Header.Set(headers.Accept, "application/json")
	svcRequest.Header.Set(headers.Authorization, "Bearer "+ts.APIKey)
	svcRequest.Header.Set(headers.ContentType, "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	svcResponse, err := client.Do(svcRequest)
	if err != nil {
		return nil
	}
	defer svcResponse.Body.Close()

	// Non-success responses won't parse. Log and retrun.
	if svcResponse.StatusCode != http.StatusOK {
		var respBody []byte
		respBody, _ = io.ReadAll(svcResponse.Body)
		log.Println("service returned non-ok status:", string(respBody))
		return nil
	}

	// Decode the service response
	tinyResponse := new(tinyURLServiceResponse)
	err = json.NewDecoder(svcResponse.Body).Decode(&tinyResponse)
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
