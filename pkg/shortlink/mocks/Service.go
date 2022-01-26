// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	shortlink "github.com/robertwtucker/document-host/pkg/shortlink"
	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Shorten provides a mock function with given fields: ctx, req
func (_m *Service) Shorten(ctx context.Context, req *shortlink.ServiceRequest) *shortlink.ServiceResponse {
	ret := _m.Called(ctx, req)

	var r0 *shortlink.ServiceResponse
	if rf, ok := ret.Get(0).(func(context.Context, *shortlink.ServiceRequest) *shortlink.ServiceResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*shortlink.ServiceResponse)
		}
	}

	return r0
}
