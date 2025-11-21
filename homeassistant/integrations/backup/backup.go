// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

package backup

import (
	"time"
)

const (
	// InfoFileName is the name of the (JSON) file that provides information about a backup TAR archive file.
	InfoFileName = "backup.json"

	// CompressedArchiveName is the name of the compressed archive file that is nested inside a backup TAR archive file.
	CompressedArchiveName = "homeassistant.tar.gz"
)

// ArchiveMetadata represents metadata of a "Home Assistant" backup archive.
type ArchiveMetadata struct {
	Info *ArchiveInfo
	Name string
}

// ArchiveInfo are information of a "Home Assistant" backup archive.
type ArchiveInfo struct {
	Compressed    bool                      `json:"compressed"`
	Date          time.Time                 `json:"date"`
	Extra         *ArchiveInfoExtra         `json:"extra"`
	HomeAssistant *ArchiveInfoHomeAssistant `json:"homeassistant"`
	Name          string                    `json:"name"`
	Protected     bool                      `json:"protected"`
	Slug          string                    `json:"slug"`
	Type          string                    `json:"type"`
	Version       uint64                    `json:"version"`
}

// ArchiveInfoExtra are extra information of a "Home Assistant" backup archive.
type ArchiveInfoExtra struct {
	InstanceID            string `json:"instance_id"`
	WithAutomaticSettings bool   `json:"with_automatic_settings"`
}

// ArchiveInfoHomeAssistant are "Home Assistant" specific information about a backup archive.
type ArchiveInfoHomeAssistant struct {
	ExcludeDatabase bool   `json:"exclude_database"`
	Version         string `json:"version"`
}
