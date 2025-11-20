// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

package extraction

// Option is an option for an archive Extractor.
type Option func(*Options)

// Options are options for an archive Extractor.
type Options struct {
	// createDestinationDirectory indicates whether the destination directory should be created.
	// This option is useful when the destination already exists or the permissions to create it are not granted.
	createDestinationDirectory bool

	// excludedFiles is a list of files that should be excluded when processing files within an archive.
	excludedFiles []string

	// processLinks indicates whether symbolic or hard links should be processed.
	// This option is useful for security reasons and because not every archive format supports links.
	processLinks bool
}

// newDefaultOptions creates new default options an archive Extractor.
func newDefaultOptions(opts ...Option) *Options {
	opt := &Options{
		createDestinationDirectory: false,
		excludedFiles:              make([]string, 0),
	}
	for _, o := range opts {
		o(opt)
	}

	return opt
}

// WithCreateDestinationDirectory indicates whether the destination directory should be created.
// This option is useful when the destination already exists or the permissions to create it are not granted.
func WithCreateDestinationDirectory(create bool) Option {
	return func(o *Options) {
		o.createDestinationDirectory = create
	}
}

// WithExcludedFiles defines a list of files that should be excluded when processing files within an archive.
func WithExcludedFiles(files []string) Option {
	return func(o *Options) {
		o.excludedFiles = files
	}
}

// WithProcessLinks indicates whether symbolic or hard links should be processed.
// This option is useful for security reasons and because not every archive format supports links.
func WithProcessLinks(process bool) Option {
	return func(o *Options) {
		o.processLinks = process
	}
}
