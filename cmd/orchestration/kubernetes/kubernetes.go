// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

// Package kubernetes provides commands to interact with Kubernetes and its workloads.
package kubernetes

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	cmdBackup "github.com/svengreb/home-ops-cli/cmd/orchestration/kubernetes/backup"

	"github.com/svengreb/home-ops-cli/cmd/config"
)

const (
	// CmdName is the name of the command provided by this package.
	CmdName = "kubernetes"
)

// New creates a new command to interact with Kubernetes and its workloads.
func New(cfg *config.Config, logger *log.Logger) (*cobra.Command, error) {
	c := &cobra.Command{
		Use:   CmdName,
		Short: `Provides commands to interact with Kubernetes and its resources like workloads`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				err.Error()
			}
		},
	}

	backupCommand, err := cmdBackup.New(cfg, logger)
	if err != nil {
		return nil, err
	}

	c.AddCommand(
		backupCommand,
	)

	return c, nil
}
