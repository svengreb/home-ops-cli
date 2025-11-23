// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

// Package homeassistant provides commands to back up data of Home Assistant that runs as Kubernetes workload.
package homeassistant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
	"github.com/invopop/validation"
	"github.com/mholt/archives"
	"github.com/spf13/cobra"

	"github.com/svengreb/home-ops-cli/cmd/config"
	"github.com/svengreb/home-ops-cli/homeassistant/api"
	"github.com/svengreb/home-ops-cli/homeassistant/integrations/backup"
	"github.com/svengreb/home-ops-cli/support/archive/extraction"
)

const (
	// CmdName is the name of the command provided by this package.
	CmdName = "home-assistant"

	// EnvVarPrefix is the prefix used when parsing the subcommand-specific HomeOps CLI environment variables.
	// This prefix must be used, together with [github.com/svengreb/home-ops-cli/cmd/config.EnvVarPrefix],
	// when parsing environment variables with the [github.com/caarlos0/env/v11.ParseWithOptions] function.
	EnvVarPrefix = "ORCHESTRATION_KUBERNETES_BACKUP_WORKLOAD_HOME_ASSISTANT_"

	// temporaryWorkingDirectoryNamePattern is the name pattern for the creation of temporary working directories.
	temporaryWorkingDirectoryNamePattern = "homeopsctl-*"
)

var (
	apiAddressRegExp              = regexp.MustCompile("^https?://.*(:[0-9]+)?$")
	backupDirectoryRegExp         = regexp.MustCompile("^/.*$")
	dataExtractionDirectoryRegExp = backupDirectoryRegExp
	temporaryDirectoryRegExp      = backupDirectoryRegExp
)

// command is the command to back up data of Home Assistant that runs as Kubernetes workload.
type command struct {
	cfg     *commandConfig
	cfgRoot *config.Config
	cmd     *cobra.Command
	logger  *log.Logger
}

// commandConfig is the subcommand-specific configuration.
type commandConfig struct {
	ActionDomain            string        `env:"ACTION_DOMAIN"`
	ActionName              string        `env:"ACTION_NAME"`
	APIAddress              string        `env:"API_ADDRESS"`
	APIToken                string        `env:"API_TOKEN"`
	BackupDirectory         string        `env:"BACKUP_DIRECTORY"`
	BackupJSONInfoFileName  string        `env:"BACKUP_JSON_INFO_FILE_NAME"`
	DataExtractionDirectory string        `env:"DATA_EXTRACTION_DIRECTORY"`
	DownloadCopyWaitTime    time.Duration `env:"DOWNLOAD_COPY_WAIT_TIME"`
	ExcludedFiles           []string      `env:"EXCLUDED_FILES"`
	KeepDays                float64       `env:"KEEP_DAYS"`
	NamePrefix              string        `env:"NAME_PREFIX"`
	NameSuffix              string        `env:"NAME_SUFFIX"`
	TemporaryDirectory      string        `env:"TEMPORARY_DIRECTORY"`
}

func (c *command) validateConfig() validation.Errors {
	return validation.Errors{
		"api-address": validation.Validate(
			c.cfg.APIAddress,
			validation.Match(apiAddressRegExp).Error("API address must be a valid HTTP(S) URL"),
		),
		"backup-directory": validation.Validate(
			c.cfg.BackupDirectory,
			validation.Match(backupDirectoryRegExp).Error("Backup directory must be an absolute path"),
			validation.Required.Error("Backup directory path must not be empty"),
		),
		"data-extraction-directory": validation.Validate(
			c.cfg.DataExtractionDirectory,
			validation.Match(dataExtractionDirectoryRegExp).Error("Data extraction directory must be an absolute path"),
			validation.Required.Error("Data extraction directory path must not be empty"),
		),
		"temporary-directory": validation.Validate(
			c.cfg.TemporaryDirectory,
			validation.Match(temporaryDirectoryRegExp).Error("Temporary directory must be an absolute path"),
		),
	}
}

// initConfig runs actions to initialize the subcommand configuration.
func (c *command) initConfig() error {
	// Parse configuration values from environment variables…
	err := env.ParseWithOptions(c.cfg, env.Options{
		Prefix: config.EnvVarPrefix + EnvVarPrefix,
	})
	if err != nil {
		return err
	}

	// …and CLI flags afterward to ensure that they take precedence.
	if err = c.cmd.ParseFlags(os.Args[1:]); err != nil {
		return err
	}

	// Ensure that excluded files are valid path that do not include leading or trailing whitespaces!
	for idx, excludedFile := range c.cfg.ExcludedFiles {
		c.cfg.ExcludedFiles[idx] = strings.TrimSpace(excludedFile)
	}

	return nil
}

func (c *command) copyBackupJSONFile(filePath, dstDir string) error {
	srcFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open backup JSON information file %s: %w", filePath, err)
	}
	defer func(f *os.File) { err = errors.Join(err, f.Close()) }(srcFile)

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("read backup JSON information file %s: %w", filePath, err)
	}

	dstFilePath := filepath.Join(dstDir, fileInfo.Name())
	dstFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("create backup JSON information file %s: %w", dstFilePath, err)
	}
	defer func(f *os.File) { err = errors.Join(err, f.Close()) }(dstFile)

	c.logger.Debug("Copying backup archive JSON information file", "from", filePath, "to", dstFilePath)
	if _, copyErr := io.Copy(dstFile, srcFile); copyErr != nil {
		return fmt.Errorf("copy backup JSON information file %s to %s: %w", filePath, dstFilePath, err)
	}
	c.logger.Info("Copied the backup JSON information file", "from", dstFilePath, "to", dstFilePath)

	return nil
}

func (c *command) decompressBackupDataArchive(filePath, dstPath string) error {
	c.logger.Debug("Decompressing backup data archive file", "path", filePath)

	extractor := extraction.NewExtractor()
	fileName := filepath.Base(filePath)
	dstFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	if err := extractor.Decompress(filePath, dstPath, dstFileName); err != nil {
		return fmt.Errorf("decompress backup data archive file %s to %s: %w", filePath, dstFileName, err)
	}
	c.logger.Info("Decompressed backup data archive file", "path", filePath)

	return nil
}

func (c *command) extractBackupArchive(filePath, dstPath string) error {
	c.logger.Debug("Extracting backup data archive file", "path", filePath, "excluded", strings.Join(c.cfg.ExcludedFiles, ","))
	extractor := extraction.NewExtractor(extraction.WithExcludedFiles(c.cfg.ExcludedFiles))
	if err := extractor.Unarchive(filePath, dstPath); err != nil {
		return fmt.Errorf("extract backup archive file %s to %s: %w", filePath, dstPath, err)
	}
	c.logger.Info("Extracted backup data archive file", "path", filePath, "excluded", strings.Join(c.cfg.ExcludedFiles, ","))
	return nil
}

func (c *command) findLatestBackup() (*backup.ArchiveMetadata, error) {
	files, err := os.ReadDir(c.cfg.BackupDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("backup directory %s does not exist", c.cfg.BackupDirectory)
		}
		return nil, fmt.Errorf("determine the latest backup in directory %s: %w", c.cfg.BackupDirectory, err)
	}

	latest := &backup.ArchiveMetadata{}
	for _, entry := range files {
		if entry.IsDir() {
			c.logger.Debug("Skipped directory in backup directory", "name", entry.Name())
			continue
		}

		format, _, identifyErr := archives.Identify(context.Background(), filepath.Join(c.cfg.BackupDirectory, entry.Name()), nil)
		if identifyErr != nil {
			if errors.Is(identifyErr, archives.NoMatch) {
				c.logger.Warn("Skipped unsupported backup archive file", "name", entry.Name(), "reason", identifyErr)
				continue
			}
			return nil, fmt.Errorf("identify format of backup archive file %s: %w", entry.Name(), identifyErr)
		}
		c.logger.Debug(
			"Identified format of backup archive file",
			"name", entry.Name(),
			"media_type", format.MediaType(),
			"extension", format.Extension(),
		)

		fileSystem, fsErr := archives.FileSystem(context.Background(), filepath.Join(c.cfg.BackupDirectory, entry.Name()), nil)
		if fsErr != nil {
			return nil, fsErr
		}

		data, readErr := fs.ReadFile(fileSystem, backup.InfoFileName)
		if readErr != nil {
			if os.IsNotExist(readErr) {
				c.logger.Debug("file not found in archive, skipping processing", "name", backup.InfoFileName, "archive", entry.Name())
				continue
			}
			return nil, fmt.Errorf("read backup information file %s in archive %s: %w", backup.InfoFileName, entry.Name(), readErr)
		}

		var metadata *backup.ArchiveInfo
		decodeErr := json.Unmarshal(data, &metadata)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode backup information file %s in archive %s: %w", backup.InfoFileName, entry.Name(), decodeErr)
		}

		if latest.Name == "" || metadata.Date.After(latest.Info.Date) {
			latest.Info = metadata
			latest.Name = entry.Name()
		}
	}
	return latest, nil
}

// run runs the main functionality of the subcommand.
func (c *command) run() error {
	latestBackup, err := c.findLatestBackup()
	if err != nil {
		return fmt.Errorf("find archive file of latest backup: %w", err)
	}
	c.logger.Info(
		"Found archive file of latest backup",
		"name", latestBackup.Name,
		"directory", c.cfg.BackupDirectory,
		"date", latestBackup.Info.Date.Format(time.RFC3339),
	)

	workingDirectory := c.cfg.TemporaryDirectory
	if workingDirectory == "" {
		workingDirectory, err = os.MkdirTemp("", temporaryWorkingDirectoryNamePattern)
		if err != nil {
			return fmt.Errorf("create temporary directory for data extraction: %w", err)
		}
		c.logger.Debug("Created temporary directory for data extraction", "path", workingDirectory)
	} else {
		workingDirectory, err = os.MkdirTemp(workingDirectory, temporaryWorkingDirectoryNamePattern)
		if err != nil {
			return fmt.Errorf("create working directory for data extraction in temporary directory %s: %w", c.cfg.TemporaryDirectory, err)
		}
		c.logger.Debug("Created working directory for data extraction ", "path", workingDirectory)
	}

	backupFilePath := filepath.Join(c.cfg.BackupDirectory, latestBackup.Name)
	c.logger.Info(
		"Extracting backup archive file",
		"path", backupFilePath,
		"name", latestBackup.Info.Name,
		"date", latestBackup.Info.Date.Format(time.RFC3339),
		"home_assistant_version", latestBackup.Info.HomeAssistant.Version,
		"is_automatic_backup", latestBackup.Info.Extra.WithAutomaticSettings,
		"is_database_excluded", latestBackup.Info.HomeAssistant.ExcludeDatabase,
		"type", latestBackup.Info.Type,
		"slug", latestBackup.Info.Slug,
		"backup_format_version", latestBackup.Info.Version,
	)
	if err = c.extractBackupArchive(backupFilePath, workingDirectory); err != nil {
		return err
	}

	compressedArchivePath := filepath.Join(workingDirectory, backup.CompressedArchiveName)
	err = c.decompressBackupDataArchive(compressedArchivePath, workingDirectory)
	if err != nil {
		return err
	}

	err = c.extractBackupArchive(
		filepath.Join(
			workingDirectory,
			strings.TrimSuffix(backup.CompressedArchiveName, filepath.Ext(backup.CompressedArchiveName)),
		),
		c.cfg.DataExtractionDirectory,
	)
	if err != nil {
		return err
	}

	if err = c.copyBackupJSONFile(filepath.Join(workingDirectory, c.cfg.BackupJSONInfoFileName), c.cfg.DataExtractionDirectory); err != nil {
		return err
	}

	err = os.RemoveAll(workingDirectory)
	if err != nil {
		return fmt.Errorf("remove temporary data extraction working directory %s: %w", workingDirectory, err)
	}

	return nil
}

// New creates a new command to back up data of Home Assistant that runs as Kubernetes workload.
func New(cfg *config.Config, logger *log.Logger) (*cobra.Command, error) {
	c := &command{
		cfgRoot: cfg,
		cfg:     &commandConfig{},
		logger:  logger,
		cmd: &cobra.Command{
			Use:   CmdName,
			Short: `Provides commands to back up data of Home Assistant that runs as Kubernetes workload`,
			Run: func(cmd *cobra.Command, args []string) {
				if err := cmd.Help(); err != nil {
					err.Error()
				}
			},
		},
	}

	c.cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := c.initConfig(); err != nil {
			c.logger.Error("Failed to initialize command configuration", "err", err)
			return err
		}

		valErrs := c.validateConfig()
		// Remove any error object that is nil.
		maps.DeleteFunc(valErrs, func(s string, err error) bool { return err == nil })
		if len(valErrs) > 0 {
			l := c.logger.WithPrefix("Configuration Validation")
			for _, err := range valErrs {
				l.Error(err.Error())
			}
			l.Errorf("Failed with %d error(s)", len(valErrs))
			return fmt.Errorf("configuration validation failed with %d errors", len(valErrs))
		}

		return nil
	}

	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := c.run(); err != nil {
			c.logger.Error(err)
			return err
		}

		return nil
	}

	c.cmd.Flags().StringVar(
		&c.cfg.APIAddress,
		"api-address",
		api.DefaultHomeAssistantAPIAddress,
		"The Home Assistant API address that consists of the scheme, hostname and port",
	)
	c.cmd.Flags().StringVar(
		&c.cfg.APIToken,
		"api-token",
		"",
		"The Home Assistant API token",
	)
	c.cmd.Flags().StringVar(
		&c.cfg.BackupDirectory,
		"backup-directory",
		"",
		"The path to be directory that contains all Home Assistant data to backup",
	)
	c.cmd.Flags().StringVar(
		&c.cfg.BackupJSONInfoFileName,
		"backup-info-file-name",
		backup.InfoFileName,
		"The name of the JSON file that provides information about the backup",
	)
	c.cmd.Flags().StringVar(
		&c.cfg.DataExtractionDirectory,
		"data-extraction-directory",
		"",
		"The path to the directory to extract the backup data to",
	)
	c.cmd.Flags().StringVar(
		&c.cfg.NameSuffix,
		"name-suffix",
		"",
		"The suffix for the name of the backup archive file that is appended to the creation date",
	)
	c.cmd.Flags().StringVar(
		&c.cfg.TemporaryDirectory,
		"temporary-directory",
		"",
		"The directory used to store temporary application data",
	)

	return c.cmd, nil
}
