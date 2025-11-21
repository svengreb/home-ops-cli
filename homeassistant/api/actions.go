// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

package api

import (
	"fmt"
	"maps"
	"net/url"
	"slices"
)

const (
	// apiPathActions is the name of the API path segment for actions.
	// Note that actions were previously named "Services" which are now part of an action, but the API path still reflects the previous naming
	// scheme due to backward-compatibility reasons.
	apiPathSegmentActions = "services"
)

// Action represents a Home Assistant "Action", previously named "Services".
type Action struct {
	// Domain is the domain of the action.
	Domain string `json:"domain"`

	// Fields is a map to hold generic fields that are specific to each action.
	// It can be mapped to a dedicated struct using the "github.com/go-viper/mapstructure" package.
	Fields map[string]any `json:"fields"`

	// Services is a map to hold generic services that are specific to each action.
	// It can be mapped to a dedicated struct using the "github.com/go-viper/mapstructure" package.
	Services map[string]any `json:"services"`
}

// ActionError represents an error for Home Assistant actions.
type ActionError struct {
	err error
	msg string
}

// HasService checks whether the action has the named service.
func (a *Action) HasService(name string) bool {
	if len(a.Services) == 0 {
		return false
	}

	for svc := range maps.Keys(a.Services) {
		if svc == name {
			return true
		}
	}

	return false
}

func (e ActionError) Error() string {
	if e.err != nil {
		if e.msg != "" {
			return e.msg + "" + e.err.Error()
		}
		return e.err.Error()
	}
	return e.msg
}

// Unwrap returns the next error in the error chain.
// If there is no next error, Unwrap returns nil.
func (e ActionError) Unwrap() error {
	return e.err
}

// GetAction returns an action by its domain.
func (c *Client) GetAction(domain string, opts ...RequestOption) (*Action, error) {
	actions, err := c.ListActions(opts...)
	if err != nil {
		return nil, &APIError{err: &ActionError{err: err, msg: fmt.Sprintf("get action for domain %s", domain)}}
	}

	idx := slices.IndexFunc(actions, func(s *Action) bool {
		return s.Domain == domain
	})
	if idx < 0 {
		return nil, &APIError{err: &ActionError{err: fmt.Errorf("action for domain %s not found", domain)}}
	}

	return actions[idx], nil
}

// ListActions returns a list of all actions.
func (c *Client) ListActions(opts ...RequestOption) ([]*Action, error) {
	apiURL, err := url.JoinPath(apiPathSegmentActions)
	if err != nil {
		return nil, &APIError{err: &ActionError{err: err, msg: "build API URL to list actions"}}
	}

	req := c.prepareRequestWithOptions(opts...)
	var result []*Action
	resp, err := req.SetResult(&result).Get(apiURL)
	if err != nil {
		return nil, &APIError{err: &ActionError{err: err, msg: "call API to list actions"}}
	}

	if !resp.IsSuccess() {
		return nil, &APIError{err: &ActionError{err: err, msg: "list actions"}, status: resp.StatusCode()}
	}

	return result, err
}

// RunAction runs an action with the provided data.
func (c *Client) RunAction(domain, name string, data map[string]any) error {
	apiURL, err := url.JoinPath(apiPathSegmentActions, domain, name)
	if err != nil {
		return &APIError{err: &ActionError{err: err, msg: "build API URL to run action"}}
	}

	resp, err := c.rc.R().SetBody(data).Post(apiURL)
	if err != nil {
		return &APIError{err: &ActionError{err: err, msg: "call API to run action"}}
	}
	if !resp.IsSuccess() {
		return &APIError{
			err: &ActionError{
				err: err,
				msg: fmt.Sprintf("running action %s of domain %s not succeeded", name, domain),
			},
			status: resp.StatusCode(),
		}
	}

	return nil
}
