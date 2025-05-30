// Package main contains configuration generation logic for repoll.
// This file implements the main configuration generation workflow,
// including repository discovery, grouping, and TOML structure creation.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bagaking/goulp/wlog"
)

// makeConfig is the main function for the mkconf command. It scans a directory
// for Git repositories and generates a TOML configuration file.
//
// The function performs the following operations:
// 1. Recursively discovers all Git repositories in the specified directory
// 2. Groups repositories by their remote prefix (hosting service/organization)
// 3. Generates a TOML configuration file with proper site structures
// 4. Reports statistics and any issues encountered
//
// Parameters:
//   - rootDir: Root directory to scan for Git repositories
//   - report: Report structure to track operations and results
//
// Returns:
//   - error: Any critical error that prevents configuration generation
func makeConfig(rootDir string, report *MkconfReport) error {
	wlog.Common().Infof("Starting repository discovery in: %s", rootDir)

	// Convert to absolute path for consistency
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", rootDir, err)
	}

	// Discover all Git repositories in the directory tree
	var repoResults []*RepoDiscoveryResult

	err = filepath.Walk(absRootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			wlog.Common().Warnf("Error accessing path %s: %s", path, err)
			return nil // Continue walking despite errors
		}

		// Only process directories
		if !info.IsDir() {
			return nil
		}

		// Analyze directory for Git repository
		if result, err := discoverGitRepo(path); err != nil {
			wlog.Common().Errorf("Error analyzing repository at %s: %s", path, err)
			return nil // Continue despite errors
		} else if result != nil {
			repoResults = append(repoResults, result)
			wlog.Common().Infof("Discovered repository: %s", path)

			// Add to report for tracking
			addToMkconfReport(report, result)

			// Skip subdirectories of Git repositories to avoid nested discovery
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory tree: %w", err)
	}

	// Generate configuration from discovered repositories
	if err := generateConfig(repoResults, absRootDir); err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	wlog.Common().Infof("Configuration generation completed. Found %d repositories.", len(repoResults))
	return nil
}

// generateConfigStruct creates a Config structure from repository discovery results.
// This function groups repositories by their remote prefix and creates a structured
// configuration that can be used for testing or further processing.
//
// Parameters:
//   - results: List of discovered repository results
//
// Returns:
//   - Config: Generated configuration structure
func generateConfigStruct(results []RepoDiscoveryResult) Config {
	// Group repositories by remote prefix
	siteGroups := make(map[string][]RepoDiscoveryResult)

	for _, repo := range results {
		if repo.Origin == "" {
			continue // Skip repositories without origin
		}

		prefix := getRemotePrefix(repo.Origin)
		siteGroups[prefix] = append(siteGroups[prefix], repo)
	}

	// Create site configurations
	var sites []SiteConfig
	for prefix, repos := range siteGroups {
		site := generateSiteConfigStruct(prefix, repos)
		sites = append(sites, site)
	}

	return Config{Sites: sites}
}

// generateSiteConfigStruct creates a SiteConfig from a group of repositories.
// This helper function processes repositories that share the same remote prefix.
//
// Parameters:
//   - remotePrefix: The remote prefix for this site
//   - repos: List of repositories for this site
//
// Returns:
//   - SiteConfig: Generated site configuration
func generateSiteConfigStruct(remotePrefix string, repos []RepoDiscoveryResult) SiteConfig {
	if len(repos) == 0 {
		return SiteConfig{RemotePrefix: remotePrefix}
	}

	// Calculate the common directory from the first repository
	baseDir := filepath.Dir(repos[0].Path)

	var repoConfigs []Repo
	for _, repo := range repos {
		repoName := getRepoNameFromURL(repo.Origin)
		memo := buildStatusMemo(repo.Uncommitted, repo.Unmerged)

		repoConfig := Repo{
			Repo:   repoName,
			Rename: "",
			WarmUp: false,
			Memo:   memo,
		}
		repoConfigs = append(repoConfigs, repoConfig)
	}

	return SiteConfig{
		RemotePrefix: remotePrefix,
		Dir:          baseDir,
		Repos:        repoConfigs,
		WarmUpAll:    false,
	}
}

// generateConfig creates and saves a TOML configuration file from repository discovery results.
// This function groups repositories by their remote prefix and generates a structured
// TOML configuration file with proper site organization.
//
// Parameters:
//   - repos: List of discovered repository results
//   - rootDir: Root directory being scanned (for relative path calculation)
//
// Returns:
//   - error: Any error encountered during config generation
func generateConfig(repos []*RepoDiscoveryResult, rootDir string) error {
	// Group repositories by remote prefix
	siteGroups := make(map[string][]*RepoDiscoveryResult)

	for _, repo := range repos {
		if !repo.HasOrigin {
			wlog.Common().Warnf("Skipping repository without origin: %s", repo.Path)
			continue
		}

		prefix := getRemotePrefix(repo.Origin)
		siteGroups[prefix] = append(siteGroups[prefix], repo)
	}

	// Generate configuration content
	var configContent strings.Builder
	configContent.WriteString("# Generated by repoll mkconf\n")
	configContent.WriteString(fmt.Sprintf("# Generated at: %s\n", time.Now().Format(time.RFC3339)))
	configContent.WriteString(fmt.Sprintf("# Root directory: %s\n\n", rootDir))

	// Sort site groups for consistent output
	var sortedPrefixes []string
	for prefix := range siteGroups {
		sortedPrefixes = append(sortedPrefixes, prefix)
	}
	sort.Strings(sortedPrefixes)

	// Generate site configurations
	for _, prefix := range sortedPrefixes {
		repos := siteGroups[prefix]
		if err := generateSiteConfig(&configContent, prefix, repos, rootDir); err != nil {
			return fmt.Errorf("failed to generate site config for %s: %w", prefix, err)
		}
	}

	// Write configuration to file
	configFileName := time.Now().Format("20060102-150405") + "_conf.toml"
	if err := os.WriteFile(configFileName, []byte(configContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configFileName, err)
	}

	wlog.Common().Infof("Configuration saved to: %s", configFileName)
	return nil
}

// generateSiteConfig generates a site section in the TOML configuration.
//
// Parameters:
//   - content: String builder for configuration content
//   - prefix: Remote prefix for this site
//   - repos: Repositories belonging to this site
//   - rootDir: Root directory for relative path calculation
//
// Returns:
//   - error: Any error encountered during site config generation
func generateSiteConfig(content *strings.Builder, prefix string, repos []*RepoDiscoveryResult, rootDir string) error {
	content.WriteString("[[sites]]\n")
	content.WriteString(fmt.Sprintf("remote = \"%s\"\n", prefix))

	// Calculate common directory for this site's repositories
	if len(repos) > 0 {
		commonDir := filepath.Dir(repos[0].Path)
		content.WriteString(fmt.Sprintf("dir = \"%s\"\n", commonDir))
	}

	content.WriteString("\n")

	// Generate repository entries
	for _, repo := range repos {
		repoName := getRepoNameFromURL(repo.Origin)
		content.WriteString(fmt.Sprintf("  [[sites.repos]]\n"))
		content.WriteString(fmt.Sprintf("  repo = \"%s\"\n", repoName))

		// Add status information as comments
		if repo.Uncommitted {
			content.WriteString("  # WARNING: Has uncommitted changes\n")
		}
		if repo.Unmerged {
			content.WriteString("  # WARNING: Has unmerged commits\n")
		}

		content.WriteString("\n")
	}

	content.WriteString("\n")
	return nil
}
