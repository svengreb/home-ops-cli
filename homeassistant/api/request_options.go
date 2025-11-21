// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

package api

import (
	"github.com/go-resty/resty/v2"
)

// RequestOption is an option for a Home Assistant API HTTP client request.
type RequestOption func(*RequestOptions)

// RequestOptions are options for a Home Assistant API HTTP client request.
type RequestOptions struct {
	// retryConditionFuncs are functions to run to determine if the request should be retried when any function returns true and the error is
	// nil.
	retryConditionFuncs []resty.RetryConditionFunc
}

// newDefaultRequestOptions creates new default options for a Home Assistant API HTTP client request.
func newDefaultRequestOptions(opts ...RequestOption) *RequestOptions {
	opt := &RequestOptions{
		retryConditionFuncs: []resty.RetryConditionFunc{},
	}
	for _, o := range opts {
		o(opt)
	}

	return opt
}

// WithRetryConditionFuncs sets functions to run to determine if the request should be retried when any function returns true and the error
// is nil.
func WithRetryConditionFuncs(funcs ...resty.RetryConditionFunc) RequestOption {
	return func(o *RequestOptions) {
		o.retryConditionFuncs = funcs
	}
}
