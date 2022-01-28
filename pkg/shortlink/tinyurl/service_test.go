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
	"net/http"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/jarcoal/httpmock"
	"github.com/robertwtucker/document-host/pkg/shortlink"
	"github.com/stretchr/testify/assert"
)

type shortLinkRequestBody struct {
	URL    string `json:"url"`
	Domain string `json:"domain"`
}

func TestShorten(t *testing.T) {
	const apiKey = "${API_KEY}"
	const domain = "tiny.one"
	const url = "http://dev.local/v1/documents/61f0023ee260d827b7156c55"

	// Mock HTTP Endpoint
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, tinyServiceURL,
		func(req *http.Request) (*http.Response, error) {

			assert.Equal(t, req.Header.Get(headers.Accept), "application/json")
			assert.Equal(t, req.Header.Get(headers.Authorization), "Bearer "+apiKey)
			assert.Equal(t, req.Header.Get(headers.ContentType), "application/json")

			var reqBody = &shortLinkRequestBody{}
			err := json.NewDecoder(req.Body).Decode(&reqBody)
			defer req.Body.Close()

			assert.NoError(t, err, "error decoding request")
			assert.Equal(t, reqBody.URL, url)
			assert.Equal(t, reqBody.Domain, domain)

			body := `{
				"data": {
					"url": "http://dev.local/v1/documents/61f0023ee260d827b7156c55",
					"domain": "tiny.one",
					"alias": "yckaxkhx",
					"tags": [],
					"tiny_url": "https://tiny.one/yckaxkhx"
				},
				"code": 0,
				"errors": []
			}`

			return httpmock.NewStringResponse(200, body), nil
		},
	)

	svc := NewTinyURLService(apiKey, domain)
	svcRequest := &shortlink.ServiceRequest{URL: url}
	svcResponse := svc.Shorten(context.Background(), svcRequest)

	assert.NotNil(t, svcResponse, "service response should not be nil")
	assert.Equal(t, svcResponse.ShortLink, "https://tiny.one/yckaxkhx")
}
