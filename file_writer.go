// Package main contains file writing functionality for repoll configuration files.
// This file implements the logic for saving configuration structures to TOML files
// with proper formatting and error handling.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bagaking/goulp/wlog"
)

// saveConfigToFile saves a configuration structure to a TOML file.
// The file is named with a timestamp to avoid conflicts.
//
// Parameters:
//   - config: Configuration structure to save
//
// Returns:
//   - error: Any error encountered during file writing
func saveConfigToFile(config Config) error {
	// Generate TOML content
	var content strings.Builder

	for _, site := range config.Sites {
		content.WriteString("[[sites]]\n")
		content.WriteString(fmt.Sprintf("remote = \"%s\"\n", site.RemotePrefix))
		content.WriteString(fmt.Sprintf("dir = \"%s\"\n", site.Dir))
		content.WriteString("\n")

		for _, repo := range site.Repos {
			content.WriteString("  [[sites.repos]]\n")
			content.WriteString(fmt.Sprintf("  repo = \"%s\"\n", repo.Repo))
			if repo.Rename != "" {
				content.WriteString(fmt.Sprintf("  rename = \"%s\"\n", repo.Rename))
			}
			if repo.WarmUp {
				content.WriteString("  warm_up = true\n")
			}
			if repo.Memo != "" {
				content.WriteString(fmt.Sprintf("  memo = \"%s\"\n", repo.Memo))
			}
			content.WriteString("\n")
		}

		content.WriteString("\n")
	}

	// Write to file with timestamp
	configFileName := time.Now().Format("20060102-150405") + "_conf.toml"
	if err := os.WriteFile(configFileName, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configFileName, err)
	}

	wlog.Common().Infof("Configuration saved to: %s", configFileName)
	return nil
}

// mkStatusMemo creates a status memo string based on conditions.
// This is a helper function for building status messages.
//
// Parameters:
//   - str: Current status string
//   - need: Whether to add the memo
//   - memo: Memo string to add
//
// Returns:
//   - string: Updated status string
func mkStatusMemo(str string, need bool, memo string) string {
	if !need {
		return str
	}
	if str != "" {
		str += " & "
	}
	str = memo
	return str
}
