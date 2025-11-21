// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

// Package api provides functionality to interact with the Home Assistant API.
package api

import (
	"net/url"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
)

const (
	// DefaultHomeAssistantAPIAddress is the default address for the Home Assistant API that consists of the scheme, hostname and port.
	DefaultHomeAssistantAPIAddress = "http://localhost:8123"

	// apiPathSegmentBase is the name of the API base path segment.
	apiPathSegmentBase = "api"
)

// APIError represents an Home Assistant API error.
type APIError struct {
	err    error
	msg    string
	status int
}

// Client is an HTTP client to interact with the Home Assistant API.
type Client struct {
	logger *log.Logger
	rc     *resty.Client
}

func (e *APIError) Error() string {
	if e.err != nil {
		if e.msg != "" {
			return e.msg + "" + e.err.Error()
		}
		return e.err.Error()
	}
	return e.msg
}

// Status returns the API error status.
func (e *APIError) Status() int {
	return e.status
}

// Unwrap unwraps the API error.
func (e *APIError) Unwrap() error {
	return e.err
}

func (c *Client) prepareRequestWithOptions(opts ...RequestOption) *resty.Request {
	opt := newDefaultRequestOptions(opts...)
	req := c.rc.R()

	switch {
	case len(opt.retryConditionFuncs) > 0:
		for _, f := range opt.retryConditionFuncs {
			req = req.AddRetryCondition(f)
		}
	}

	return req
}

// NewClient returns a new HTTP client to interact with the Home Assistant API.
func NewClient(address, apiToken string, logger *log.Logger) (*Client, error) {
	apiURL, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.JoinPath(apiURL.String(), apiPathSegmentBase)
	if err != nil {
		return nil, err
	}

	rc := resty.New().SetAuthToken(apiToken).SetBaseURL(baseURL)

	return &Client{
		logger: logger,
		rc:     rc,
	}, nil
}
