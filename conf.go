// Package main contains configuration handling functionality for repoll.
// This file defines the TOML configuration file format and provides utilities
// for reading, parsing, and working with repository configurations.
//
// Configuration Structure:
// - Config: Top-level structure containing multiple sites
// - SiteConfig: Represents a remote location (GitHub, GitLab, etc.)
// - Repo: Individual repository within a site
//
// The TOML format allows users to organize repositories by hosting service
// and configure batch operations like cloning, updating, and warm-up.
package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Repo represents a single repository within a site configuration.
// It contains the repository name and optional settings for customization.
type (
	Repo struct {
		Repo   string `toml:"repo"`              // Repository name (required)
		Rename string `toml:"rename,omitempty"`  // Optional custom local directory name
		WarmUp bool   `toml:"warm_up,omitempty"` // Enable warm-up operations for this repo
		Memo   string `toml:"memo,omitempty"`    // Optional notes or comments
	}

	// SiteConfig represents a remote location and the directory to clone repos to.
	// It groups repositories from the same hosting service or organization.
	SiteConfig struct {
		RemotePrefix string `toml:"remote"`            // Remote URL prefix (e.g., "https://github.com/user")
		Dir          string `toml:"dir"`               // Local directory for cloning repositories
		Repos        []Repo `toml:"repos"`             // List of repositories in this site
		WarmUpAll    bool   `toml:"warm_up,omitempty"` // Enable warm-up for all repos in this site
	}

	// Config represents the top-level TOML configuration file format.
	// It contains multiple site configurations for organizing repositories.
	Config struct {
		Sites []SiteConfig `toml:"sites"` // List of site configurations
	}
)

// readConfig reads and parses a TOML configuration file.
// It validates the file format and returns a structured configuration object.
//
// Parameters:
//   - configPath: Filesystem path to the TOML configuration file
//
// Returns:
//   - *Config: Parsed configuration structure
//   - error: Any error encountered during reading or parsing
func readConfig(configPath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to decode TOML config file %s: %w", configPath, err)
	}
	return &config, nil
}

// RepoUrl constructs the full Git repository URL for cloning operations.
// It combines the site's remote prefix with the repository name and handles
// different URL formats (HTTPS and SSH).
//
// Parameters:
//   - site: Site configuration containing the remote prefix
//
// Returns:
//   - string: Complete Git repository URL ready for cloning
func (repo Repo) RepoUrl(site SiteConfig) string {
	repoURL := strings.TrimSpace(repo.Repo) + ".git"

	if isURL(site.RemotePrefix) {
		// Handle HTTPS URLs: "https://github.com/user/" + "repo.git"
		repoURL = ensureTrailingSlash(site.RemotePrefix) + repoURL
	} else {
		// Handle SSH URLs: "git@github.com:user/" + "repo.git"
		repoURL = site.RemotePrefix + repoURL
	}

	return repoURL
}

// FullPath determines the complete local filesystem path where the repository
// should be cloned. It handles custom rename settings and path resolution.
//
// Parameters:
//   - site: Site configuration containing the base directory
//
// Returns:
//   - string: Complete local path for the repository
func (repo Repo) FullPath(site SiteConfig) string {
	repo.Repo = strings.TrimSpace(repo.Repo)

	// Handle custom rename settings
	if repo.Rename != "" {
		if strings.ToLower(repo.Rename) == "{base}" {
			// Use base name of the repository
			return filepath.Join(site.Dir, filepath.Base(repo.Repo))
		}
		// Use custom rename as directory name
		return filepath.Join(site.Dir, strings.TrimSpace(repo.Rename))
	}

	// Use repository name as directory name (default behavior)
	return filepath.Join(site.Dir, repo.Repo)
}
