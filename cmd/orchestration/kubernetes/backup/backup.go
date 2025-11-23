// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

// Package backup provides commands to back up data in the context of Kubernetes.
package backup

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/svengreb/home-ops-cli/cmd/config"
	cmdWorkload "github.com/svengreb/home-ops-cli/cmd/orchestration/kubernetes/backup/workload"
)

const (
	// CmdName is the name of the command provided by this package.
	CmdName = "backup"
)

// New creates a new command to back up data in the context of Kubernetes.
func New(cfg *config.Config, logger *log.Logger) (*cobra.Command, error) {
	c := &cobra.Command{
		Use:   CmdName,
		Short: `Provides commands to create data backups in the context of Kubernetes`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				err.Error()
			}
		},
	}

	workloadCommand, err := cmdWorkload.New(cfg, logger)
	if err != nil {
		return nil, err
	}

	c.AddCommand(
		workloadCommand,
	)

	return c, nil
}
