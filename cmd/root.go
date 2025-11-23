// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

// Package cmd provides the HomeOps CLI application.
package cmd

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/svengreb/home-ops-cli/cmd/config"
	cmdOrchestration "github.com/svengreb/home-ops-cli/cmd/orchestration"
)

const (
	// FlagNameDebug is the name of the persistent command flag that indicates whether the CLI message printer runs in debug mode to output
	// messages with debug scope.
	FlagNameDebug = "debug"
)

// The values are populated during build-time through the `-X` linker flag.
var (
	// version is the application version.
	version = "development"
)

// app is the HomeOps CLI application.
type app struct {
	cfg    *config.Config
	cmd    *cobra.Command
	logger *log.Logger
}

// initCmd runs actions to initialize the application.
func (a *app) initCmd() error {
	a.cmd = &cobra.Command{
		// Define the application command name and documentations.
		Use:   "homeopsctl",
		Short: "HomeOps CLI",

		// Define the function that will be run before any subcommand.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return a.preRun()
		},
		// Print the command help and usage information when called without any subcommand.
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},

		// Disable cobra's verbose output on errors to use the custom CLI message printer for user-facing messages instead.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Define global (persistent) and local flags.
	a.cmd.PersistentFlags().BoolVar(&a.cfg.DebugMode, FlagNameDebug, false, "enables output with debug scope")

	// Set the version information for the automatically generated "version" flag.
	a.cmd.Version = version
	a.cmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)

	orchestrationCommand, err := cmdOrchestration.New(a.cfg, a.logger)
	if err != nil {
		return err
	}

	// Register all subcommands.
	a.cmd.AddCommand(
		orchestrationCommand,
	)

	return nil
}

// initConfig runs actions to initialize the application configuration.
func (a *app) initConfig() error {
	// Adjust the verbosity level when the debug mode has been enabled.
	if a.cfg.DebugMode {
		a.logger.SetLevel(log.DebugLevel)
	}

	return nil
}

// initLogger initializes the CLI message printer for user-facing messages.
func (a *app) initLogger() {
	// Set custom styles to improve the visual appearance for human CLI interactions.
	// The colors are part of the "Nord" color theme.
	// References:
	//   1. https://www.nordtheme.com
	styles := log.DefaultStyles()
	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FATAL").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#bf616a")).
		Foreground(lipgloss.Color("#eceff4"))
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#bf616a")).
		Foreground(lipgloss.Color("#2e3440"))
	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARNING").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#ebcb8b")).
		Foreground(lipgloss.Color("#2e3440"))
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#88c0d0")).
		Foreground(lipgloss.Color("#2e3440"))
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#b48ead")).
		Foreground(lipgloss.Color("#2e3440"))

	// Add a custom style for the `err` key to make it stand out.
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("#bf616a"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	a.logger.SetStyles(styles)
}

// preRun runs actions before the application starts.
func (a *app) preRun() error {
	// Initialize the CLI message printer.
	a.initLogger()

	// Populate the application configuration.
	if err := a.initConfig(); err != nil {
		a.logger.Error("Failed to initialize HomeOps CLI configuration", "err", err)
		return err
	}

	return nil
}

// Main is the main HomeOps CLI execution function.
func Main() int {
	a := &app{}

	a.logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
	})

	cfg, err := config.New()
	if err != nil {
		log.Printf("Failed to create HomeOps CLI configuration: %v", err)
		return 1
	}
	a.cfg = cfg

	err = a.initCmd()
	if err != nil {
		log.Printf("Failed to initialize HomeOps CLI: %v", err)
		return 1
	}

	// Execute the root command to parse and validate commands, flags and arguments and exit on any (downstream) error.
	if err = a.cmd.Execute(); err != nil {
		if a.cfg.PrintFinalError {
			a.logger.Error("Failed to execute HomeOps CLI", "err", err)
		}
		return 1
	}

	return 0
}
