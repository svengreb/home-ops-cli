// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

package archive

import (
	"github.com/mholt/archives"
)

const (
	// FilePathSeparator is the separator for folders/paths within an archive which is independent of the operating system.
	// A slash character is the standard, like e.g. defined in [the ZIP file format specification in section 4.4.17], even when an archive
	// file is created on operating systems like Windows where backslashes are used for file system paths.
	//
	// [the ZIP file format specification in section 4.4.17]: https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT
	FilePathSeparator = '/'
)

var (
	// ArchivalMapping maps archival formats.
	ArchivalMapping = map[string]archives.Archival{
		"tar": archives.Tar{},
		"zip": archives.Zip{},
	}

	// CompressionMapping compression formats.
	CompressionMapping = map[string]archives.Compression{
		"br":  archives.Brotli{},
		"bz2": archives.Bz2{},
		"gz":  archives.Gz{},
		"lz4": archives.Lz4{},
		"xz":  archives.Xz{},
		"zst": archives.Zstd{},
	}
)
