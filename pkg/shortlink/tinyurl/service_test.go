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
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/robertwtucker/document-host/pkg/shortlink"
	"github.com/robertwtucker/document-host/pkg/shortlink/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShorten(t *testing.T) {
	// Mock HTTP Call
	mockBodyIn := []byte(`{"url":"http://dev.local/v1/documents/61f0023ee260d827b7156c55"}`)
	httpRequest, err := http.NewRequest(http.MethodPost, tinyServiceURL, bytes.NewBuffer(mockBodyIn))
	assert.NoError(t, err, "error creating request")
	httpRequest.Header.Set(headers.Accept, "application/json")
	httpRequest.Header.Set(headers.Authorization, "Bearer token")
	httpRequest.Header.Set(headers.ContentType, "application/json")

	mockBodyOut := `
	{
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
	httpResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(mockBodyOut))),
	}

	client := new(mocks.Client)
	client.On("Do", httpRequest).Return(httpResponse)

	// Mock Service Request
	expectedInput := &shortlink.ServiceRequest{URL: "http://dev.local/v1/documents/61f0023ee260d827b7156c55"}

	var data tinyURLData
	jsonData := `{
		"url": "http://dev.local/v1/documents/61f0023ee260d827b7156c55",
		"domain": "tiny.one",
		"alias": "yckaxkhx",
		"tags": [],
		"tiny_url": "https://tiny.one/yckaxkhx"
	}`
	assert.NoError(t, json.Unmarshal([]byte(jsonData), &data))
	response := &tinyURLServiceResponse{
		Data:   data,
		Code:   0,
		Errors: nil,
	}

	svc := new(mocks.Service)
	svc.On("Shorten", context.Background(), expectedInput).Return(response)

	// Test
	//tiny := NewTinyURLService("token", "tiny.one")
	//svcRequest := &shortlink.ServiceRequest{URL: "http://dev.local/v1/documents/61f0023ee260d827b7156c55"}
	//svcResponse := tiny.Shorten(context.Background(), svcRequest)
	//
	//assert.Equal(t, svcResponse.ShortLink, "https://tiny.one/yckaxkhx")

}
