// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

// Package workload provides commands to back up data of Kubernetes workloads.
package workload

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/svengreb/home-ops-cli/cmd/config"
	cmdHomeAssistant "github.com/svengreb/home-ops-cli/cmd/orchestration/kubernetes/backup/workload/homeassistant"
)

const (
	// CmdName is the name of the command provided by this package.
	CmdName = "workload"
)

// New creates a new command to back up data of Kubernetes workloads.
func New(cfg *config.Config, logger *log.Logger) (*cobra.Command, error) {
	c := &cobra.Command{
		Use:   CmdName,
		Short: `Provides commands to create data backups of Kubernetes workloads`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				err.Error()
			}
		},
	}

	homeAssistantCommand, err := cmdHomeAssistant.New(cfg, logger)
	if err != nil {
		return nil, err
	}

	c.AddCommand(
		homeAssistantCommand,
	)

	return c, nil
}
