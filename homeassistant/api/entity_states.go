// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

package api

import (
	"fmt"
	"net/url"
	"time"
)

const (
	// apiPathSegmentEntityStates is the name of the API path segment for entity states.
	apiPathSegmentEntityStates = "states"
)

// EntityState represents the state of a Home Assistant entity.
// References:
//   1. https://www.home-assistant.io/docs/configuration/state_object/#about-the-state-object
//   2. https://data.home-assistant.io/docs/states
//   3. https://developers.home-assistant.io/docs/dev_101_states
type EntityState struct {
	// Attributes is a map to hold attributes that are specific to each state.
	// It can be mapped to a dedicated struct using the "github.com/go-viper/mapstructure" package.
	Attributes map[string]any `json:"attributes"`

	// EntityID is the identifier of the entity this state is mapped to.
	EntityID string `json:"entity_id"`

	// LastChanged is the time the state changed in the state machine in UTC time.
	// Note that this is not updated if only state attributes change!
	LastChanged time.Time `json:"last_changed"`

	// State is the current state of the entity.
	State string `json:"state"`
}

// EntityStateError represents an error for Home Assistant entity state.
type EntityStateError struct {
	err error
	msg string
}

func (e EntityStateError) Error() string {
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
func (e EntityStateError) Unwrap() error {
	return e.err
}

// GetEntityStateByID returns an entity state by its ID.
func (c *Client) GetEntityStateByID(id string, opts ...RequestOption) (*EntityState, error) {
	apiURL, err := url.JoinPath(apiPathSegmentEntityStates, id)
	if err != nil {
		return nil, &APIError{err: &EntityStateError{err: err, msg: "build API URL to get entity state by ID"}}
	}

	req := c.prepareRequestWithOptions(opts...)
	var result *EntityState
	resp, err := req.SetResult(&result).Get(apiURL)
	if err != nil {
		return nil, &APIError{err: &EntityStateError{err: err, msg: fmt.Sprintf("API call to get entity state for ID %s", id)}}
	}

	if !resp.IsSuccess() {
		return nil, &APIError{err: &EntityStateError{err: err, msg: fmt.Sprintf("get entity state for ID %s", id)}, status: resp.StatusCode()}
	}

	return result, err
}

// ListEntityStates returns a list of all entity states.
func (c *Client) ListEntityStates(opts ...RequestOption) ([]*EntityState, error) {
	apiURL, err := url.JoinPath(apiPathSegmentEntityStates)
	if err != nil {
		return nil, &APIError{err: &EntityStateError{err: err, msg: "build API URL to list entities"}}
	}

	req := c.prepareRequestWithOptions(opts...)
	var result []*EntityState
	resp, err := req.SetResult(&result).Get(apiURL)
	if err != nil {
		return nil, &APIError{err: &EntityStateError{err: err, msg: "API call to list entities"}}
	}

	if !resp.IsSuccess() {
		return nil, &APIError{err: &EntityStateError{err: err, msg: "list entity states"}, status: resp.StatusCode()}
	}

	return result, err
}
