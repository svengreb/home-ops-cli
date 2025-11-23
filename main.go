// Copyright 2018 Sven Greb <development@svengreb.de>
// This source code is licensed under the Apache License 2.0 found in the license file.

// Package main provides the official HomeOps CLI.
package main

import (
	"os"
	// Embed the IANA timezone database into the binary to improve portability instead of relying on the host system to provide the data.
	_ "time/tzdata"

	// Automatically sets the "GOMAXPROCS" environment variable to match Linux container CPU quota for better performance and compatibility.
	_ "go.uber.org/automaxprocs"

	"github.com/svengreb/home-ops-cli/cmd"
)

// main is the main HomeOps CLI execution function.
func main() {
	os.Exit(cmd.Main())
}
