// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

// Package orchestration provides commands to interact with orchestration systems.
package orchestration

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/svengreb/home-ops-cli/cmd/config"
	cmdKubernetes "github.com/svengreb/home-ops-cli/cmd/orchestration/kubernetes"
)

const (
	// CmdName is the name of the command provided by this package.
	CmdName = "orchestration"
)

// New creates a new command to interact with orchestration systems.
func New(cfg *config.Config, logger *log.Logger) (*cobra.Command, error) {
	c := &cobra.Command{
		Use:   CmdName,
		Short: `Provides commands to interact with orchestration systems`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				err.Error()
			}
		},
	}

	kubernetesCommand, err := cmdKubernetes.New(cfg, logger)
	if err != nil {
		return nil, err
	}

	c.AddCommand(
		kubernetesCommand,
	)

	return c, nil
}
