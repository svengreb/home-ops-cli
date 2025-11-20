// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

package extraction

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"

	"github.com/svengreb/home-ops-cli/support/archive"
)

// Extractor provides methods to unarchive and decompress archives.
// Their format is auto-detected and the extraction is safed by checking for the ["ZipSlip" vulnerability] and proper directory creation.
//
// ["ZipSlip" vulnerability]: https://security.snyk.io/research/zip-slip-vulnerability
type Extractor struct {
	// opts are the extractor options.
	opts *Options
}

// Decompress decompresses an archive.
func (e *Extractor) Decompress(archivePath, dstPath, dstFileName string) (err error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open archive file %s: %w\n", archivePath, err)
	}
	defer func(f *os.File) { err = errors.Join(err, f.Close()) }(file)

	format, formatReader, err := archives.Identify(context.Background(), archivePath, file)
	if err != nil {
		return fmt.Errorf("identify format of archive file %s: %w\n", archivePath, err)
	}

	// Handle compressed archives by assigning the required reader.
	decompressor, isCompressed := format.(archives.Decompressor)
	if !isCompressed {
		return fmt.Errorf("archive file %s is not a supported compression format", archivePath)
	}

	reader, err := decompressor.OpenReader(formatReader)
	if err != nil {
		return errors.Join(err, err)
	}
	defer func() { err = errors.Join(err, reader.Close()) }()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("read file information for archive file %s: %w\n", file.Name(), err)
	}

	decompressedFilePath := filepath.Join(dstPath, dstFileName)
	// Create the decompressed archive file on the target file system using the same permissions as the compressed archive.
	decompressedFile, err := os.OpenFile(decompressedFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("create decompressed archive file %s: %w", decompressedFilePath, err)
	}
	defer func() { err = errors.Join(err, decompressedFile.Close()) }()

	if _, err = io.Copy(decompressedFile, reader); err != nil {
		return fmt.Errorf("write decompressed archive file %s: %w", decompressedFilePath, err)
	}

	return err
}

// Unarchive unarchives and archive.
func (e *Extractor) Unarchive(archivePath, dstPath string) (err error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open archive file %s: %w", archivePath, err)
	}
	defer func() { err = errors.Join(err, file.Close()) }()

	format, stream, err := archives.Identify(context.Background(), archivePath, file)
	if err != nil {
		return fmt.Errorf("identify format of archive file %s: %w", archivePath, err)
	}

	extractor, isSupportedFormat := format.(archives.Extractor)
	if !isSupportedFormat {
		return errors.New("unsupported archive file format for extraction")
	}

	if e.opts.createDestinationDirectory {
		if err = os.MkdirAll(dstPath, os.ModeDir); err != nil {
			return fmt.Errorf("create destination directory %s: %w", dstPath, err)
		}
	}

	return extractor.Extract(context.Background(), io.NopCloser(stream), func(ctx context.Context, f archives.FileInfo) error {
		return e.unarchiveFile(f, dstPath)
	})
}

func (e *Extractor) unarchiveFile(f archives.FileInfo, dstPath string) (err error) {
	// Skip processing early when the file should be excluded.
	if e.isFileExcluded(f.NameInArchive) {
		return nil
	}

	fileDstPath := filepath.Join(dstPath, f.NameInArchive)

	// Check for the "ZipSlip" vulnerability (path traversal attack).
	// References:
	//   1. https://security.snyk.io/research/zip-slip-vulnerability
	//   2. https://github.com/snyk/zip-slip-vulnerability
	if !strings.HasPrefix(fileDstPath, filepath.Clean(dstPath)+string(os.PathSeparator)) {
		return fmt.Errorf(`"ZipSlip" vulnerability (path traversal attack) detected for file %q`, f.NameInArchive)
	}

	// Ensure the parent directory exists.
	parentDir := filepath.Dir(fileDstPath)
	if err = os.MkdirAll(parentDir, os.ModePerm); err != nil {
		return fmt.Errorf("create parent directory %s: %w", parentDir, err)
	}

	// Handle directories early by creating them with permissions from the archive.
	if f.IsDir() {
		if err = os.MkdirAll(fileDstPath, f.Mode()); err != nil {
			return fmt.Errorf("create directory %s: %w", f.NameInArchive, err)
		}
		return nil
	}

	if !e.opts.processLinks && f.LinkTarget != "" {
		return fmt.Errorf("file %s in archive is a (symbolic or hard) link", f.NameInArchive)
	}

	reader, err := f.Open()
	if err != nil {
		return fmt.Errorf("read archive file %s: %w", f.NameInArchive, err)
	}
	defer func() { err = errors.Join(err, reader.Close()) }()

	// Create the new file on the target file system using the permissions from the file within the archive.
	dstFile, err := os.OpenFile(fileDstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("create file %s to extraxt archive file %s: %w", fileDstPath, f.NameInArchive, err)
	}
	defer func() { err = errors.Join(err, dstFile.Close()) }()

	if _, err = io.Copy(dstFile, reader); err != nil {
		return fmt.Errorf("copy archive file %s to %s: %w", f.NameInArchive, fileDstPath, err)
	}

	return err
}

// indentifyFormat can identify unknown content based on the file name and/or the header (which peeks at the stream).
func (e *Extractor) indentifyFormat(ctx context.Context, filePath string, r io.Reader) (archives.Format, io.Reader, error) {
	// Unless the stream is an io.Seeker, the returned stream value must be used to ensure to re-read the bytes consumed during the method
	// call.
	format, reader, err := archives.Identify(ctx, filePath, r)
	if err != nil {
		return nil, nil, err
	}

	return format, reader, nil
}

// isFileExcluded checks if the file name is excluded according to the exclusion list, meaning it is in the exclusion list, its parent
// folder/path is in the list, or the list is empty.
func (e *Extractor) isFileExcluded(fileName string) bool {
	// Include all files if the list is empty.
	if len(e.opts.excludedFiles) == 0 {
		return false
	}

	for _, excludedFileName := range e.opts.excludedFiles {
		// Exclude exact matches and files with the same parent folder/path.
		if fileName == excludedFileName ||
			strings.HasPrefix(
				fileName,
				strings.TrimSuffix(excludedFileName, string(archive.FilePathSeparator))+string(archive.FilePathSeparator),
			) {
			return true
		}
	}
	return false
}

// NewExtractor creates a new Extractor with the given list of Option.
func NewExtractor(opts ...Option) *Extractor {
	opt := newDefaultOptions(opts...)
	return &Extractor{
		opts: opt,
	}
}
