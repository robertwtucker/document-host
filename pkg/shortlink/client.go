//
// Copyright (HTTPClient) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package shortlink

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Client defines the short link generation service client interface (HTTP)
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

var HTTPClient Client

func init() {
	HTTPClient = &http.Client{Timeout: time.Second * 5}
}

// Post wraps the HTTPClient to facilitate testing
func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header = headers

	return HTTPClient.Do(request)
}
