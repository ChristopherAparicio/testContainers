// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/christapa/tinyurl/internal/tinyurl/domain"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// URL is an autogenerated mock type for the URL type
type URL struct {
	mock.Mock
}

// CreateShortenUrl provides a mock function with given fields: ctx, url, expiration
func (_m *URL) CreateShortenUrl(ctx context.Context, url string, expiration time.Time) (domain.Url, error) {
	ret := _m.Called(ctx, url, expiration)

	if len(ret) == 0 {
		panic("no return value specified for CreateShortenUrl")
	}

	var r0 domain.Url
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) (domain.Url, error)); ok {
		return rf(ctx, url, expiration)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) domain.Url); ok {
		r0 = rf(ctx, url, expiration)
	} else {
		r0 = ret.Get(0).(domain.Url)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time) error); ok {
		r1 = rf(ctx, url, expiration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
 
// GetOriginalUrl provides a mock function with given fields: ctx, shortUrl
func (_m *URL) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	ret := _m.Called(ctx, shortUrl)

	if len(ret) == 0 {
		panic("no return value specified for GetOriginalUrl")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, shortUrl)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, shortUrl)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, shortUrl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetURLMetadata provides a mock function with given fields: ctx, url
func (_m *URL) GetURLMetadata(ctx context.Context, url string) (domain.Url, error) {
	ret := _m.Called(ctx, url)

	if len(ret) == 0 {
		panic("no return value specified for GetURLMetadata")
	}

	var r0 domain.Url
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.Url, error)); ok {
		return rf(ctx, url)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Url); ok {
		r0 = rf(ctx, url)
	} else {
		r0 = ret.Get(0).(domain.Url)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewURL creates a new instance of URL. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURL(t interface {
	mock.TestingT
	Cleanup(func())
}) *URL {
	mock := &URL{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
