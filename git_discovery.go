// Package main contains Git repository discovery and status checking functionality.
// This file implements the logic for finding Git repositories, checking their status,
// and extracting relevant information for configuration generation.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/bagaking/goulp/wlog"
)

// RepoDiscoveryResult represents the result of discovering and analyzing a Git repository.
// It contains all the information needed to generate configuration entries and reports.
type RepoDiscoveryResult struct {
	Path        string // Local filesystem path to the repository
	Origin      string // Git remote origin URL
	Uncommitted bool   // Whether repository has uncommitted changes
	Unmerged    bool   // Whether repository has unmerged commits
	HasOrigin   bool   // Whether repository has a remote origin configured
}

// discoverGitRepo analyzes a directory to determine if it's a Git repository
// and extracts relevant information including origin, status, and metadata.
//
// Parameters:
//   - dirPath: Filesystem path to analyze
//
// Returns:
//   - *RepoDiscoveryResult: Repository information if Git repo found, nil otherwise
//   - error: Any error encountered during analysis
func discoverGitRepo(dirPath string) (*RepoDiscoveryResult, error) {
	// Check if directory contains a .git subdirectory
	gitDir := filepath.Join(dirPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return nil, nil // Not a Git repository
	}

	result := &RepoDiscoveryResult{
		Path: dirPath,
	}

	// Get remote origin URL
	if origin, err := getGitOrigin(dirPath); err != nil {
		wlog.Common().Warnf("Failed to get origin for %s: %s", dirPath, err)
		result.HasOrigin = false
	} else if origin != "" {
		result.Origin = origin
		result.HasOrigin = true
	}

	// Check for uncommitted changes
	if hasUncommitted, err := hasUncommittedChanges(dirPath); err != nil {
		wlog.Common().Warnf("Failed to check uncommitted changes for %s: %s", dirPath, err)
	} else {
		result.Uncommitted = hasUncommitted
	}

	// Check for unmerged commits
	if hasUnmerged, err := hasUnmergedCommits(dirPath); err != nil {
		wlog.Common().Warnf("Failed to check unmerged commits for %s: %s", dirPath, err)
	} else {
		result.Unmerged = hasUnmerged
	}

	return result, nil
}

// getGitOrigin retrieves the remote origin URL from a Git repository.
//
// Parameters:
//   - repoPath: Path to the Git repository
//
// Returns:
//   - string: Remote origin URL, empty if not found
//   - error: Any error encountered while querying Git
func getGitOrigin(repoPath string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git origin: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// getGitRemoteOrigin is an alias for getGitOrigin for backward compatibility.
// This function maintains compatibility with existing test code.
//
// Parameters:
//   - repoPath: Path to the Git repository
//
// Returns:
//   - string: Remote origin URL, empty if not found
//   - error: Any error encountered while querying Git
func getGitRemoteOrigin(repoPath string) (string, error) {
	return getGitOrigin(repoPath)
}

// hasUncommittedChanges checks if a Git repository has uncommitted changes.
// This includes both staged and unstaged modifications.
//
// Parameters:
//   - repoPath: Path to the Git repository
//
// Returns:
//   - bool: True if there are uncommitted changes
//   - error: Any error encountered while checking status
func hasUncommittedChanges(repoPath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	// If output is not empty, there are uncommitted changes
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// hasUnmergedCommits checks if a Git repository has commits that haven't been
// merged with the remote tracking branch.
//
// Parameters:
//   - repoPath: Path to the Git repository
//
// Returns:
//   - bool: True if there are unmerged commits
//   - error: Any error encountered while checking
func hasUnmergedCommits(repoPath string) (bool, error) {
	// Check if there are commits ahead of origin
	cmd := exec.Command("git", "cherry", "-v")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		// git cherry might fail if no upstream is set, which is normal
		wlog.Common().Warnf("Failed to check git cherry for %s (this may be normal): %s", repoPath, err)
		return false, nil
	}

	// If output is not empty, there are unmerged commits
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// buildStatusMemo creates a status memo string based on repository conditions.
// It prioritizes unmerged commits over uncommitted changes in the memo.
//
// Parameters:
//   - uncommitted: Whether repository has uncommitted changes
//   - unmerged: Whether repository has unmerged commits
//
// Returns:
//   - string: Status memo describing repository state
func buildStatusMemo(uncommitted, unmerged bool) string {
	if unmerged {
		return "unmerged"
	}
	if uncommitted {
		return "uncommitted"
	}
	return ""
}

// calculateRelativeLocalPath calculates the relative local path for a repository
// based on its filesystem path, repository name, and base directory.
//
// Parameters:
//   - repoPath: Full filesystem path to the repository
//   - repoName: Name of the repository
//   - baseDir: Base directory for path calculation
//
// Returns:
//   - string: Calculated relative local path
func calculateRelativeLocalPath(repoPath, repoName, baseDir string) string {
	// Remove the base directory from the repo path
	relPath, err := filepath.Rel(baseDir, repoPath)
	if err != nil {
		// Fallback: use the parent directory of the repo
		return filepath.Dir(repoPath) + "/"
	}

	// Remove the repository name from the end to get the parent directory
	parentDir := filepath.Dir(relPath)
	if parentDir == "." {
		return baseDir + "/"
	}

	return filepath.Join(baseDir, parentDir) + "/"
}

// addToMkconfReport adds a repository discovery result to the mkconf report.
//
// Parameters:
//   - report: Report structure to update
//   - result: Repository discovery result to add
func addToMkconfReport(report *MkconfReport, result *RepoDiscoveryResult) {
	action := &MkconfAction{
		Time:        time.Now(),
		Path:        result.Path,
		Origin:      result.Origin,
		HasOrigin:   result.HasOrigin,
		Uncommitted: result.Uncommitted,
		Unmerged:    result.Unmerged,
	}

	report.Actions = append(report.Actions, action)
}

// buildSiteConfig builds a site configuration from a repository discovery result.
// This function is used for testing and configuration generation.
//
// Parameters:
//   - result: Repository discovery result
//   - baseDir: Base directory for path calculation
//
// Returns:
//   - string: Site configuration key for grouping
//   - SiteConfig: Generated site configuration
func buildSiteConfig(result *RepoDiscoveryResult, baseDir string) (string, SiteConfig) {
	if !result.HasOrigin {
		return "", SiteConfig{}
	}

	remotePrefix := getRemotePrefix(result.Origin)
	repoName := getRepoNameFromURL(result.Origin)

	// Build status memo
	statusMemo := buildStatusMemo(result.Uncommitted, result.Unmerged)

	// Create site config
	siteConfig := SiteConfig{
		RemotePrefix: remotePrefix,
		Dir:          baseDir,
		Repos: []Repo{
			{
				Repo: repoName,
				Memo: statusMemo,
			},
		},
	}

	// Generate key for grouping
	key := fmt.Sprintf("%s @@ %s", remotePrefix, baseDir)

	return key, siteConfig
}
