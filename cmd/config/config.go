// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

// Package config provides the configuration for the HomeOps CLI.
package config

import (
	"github.com/caarlos0/env/v11"
)

const (
	// DefaultOutputFormat is the default format for CLI output.
	DefaultOutputFormat = "text"

	// EnvVarPrefix is the prefix used when parsing any HomeOps CLI specific environment variables.
	// This prefix must be used when parsing environment variables with the [github.com/caarlos0/env/v11.ParseWithOptions] function.
	EnvVarPrefix = "HOMEOPSCTL_"
)

// Config is the configuration for the HomeOps CLI.
type Config struct {
	// DebugMode indicates whether the CLI runs in debug mode.
	DebugMode bool `env:"DEBUG"`

	// OutputFormat is the format of the CLI output.
	OutputFormat string `env:"OUTPUT_FORMAT"`

	// PrintFinalError indicates whether the final error should be printed when exiting with a non-zero code.
	PrintFinalError bool `env:"PRINT_FINAL_ERROR"`
}

// New creates a new HomeOps CLI application configuration.
func New() (*Config, error) {
	c := &Config{
		DebugMode:       false,
		OutputFormat:    DefaultOutputFormat,
		PrintFinalError: false,
	}
	// Parse configuration values from environment variables with the main prefix.
	if err := env.ParseWithOptions(c, env.Options{Prefix: EnvVarPrefix}); err != nil {
		return nil, err
	}

	return c, nil
}
