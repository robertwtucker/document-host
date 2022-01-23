//
// Copyright (c) 2022 Quadient Group AG
//
// This file is subject to the terms and conditions defined in the
// 'LICENSE' file found in the root of this source code package.
//

package shortlink

import "context"

// Service defines the short link generation service interface
type Service interface {
	Shorten(ctx context.Context, req *ServiceRequest) *ServiceResponse
}

// ServiceRequest represents the service short link generation inputs
type ServiceRequest struct {
	URL string `json:"url,omitempty"`
}

// ServiceResponse represents the short link generation service outputs
type ServiceResponse struct {
	URL       string `json:"url,omitempty"`
	ShortLink string `json:"shortLink,omitempty"`
}
