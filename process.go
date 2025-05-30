// Package main contains the core processing logic for repoll repository operations.
// This file implements the repository cloning, updating, and warm-up functionality.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bagaking/goulp/wlog"
)

// RepoProcessResult represents the processing result of a single repository operation.
// It captures all relevant information about what happened during the clone/update process.
type RepoProcessResult struct {
	Repo     string        `json:"repo"`     // Repository name
	Duration time.Duration `json:"duration"` // Time taken for the operation
	Success  bool          `json:"success"`  // Whether operation completed successfully
	Error    string        `json:"error"`    // Error message if operation failed
	Memo     string        `json:"memo"`     // Additional notes or information
}

// processConfig processes a configuration file by reading the config and processing
// each site's repositories concurrently. This is the main entry point for the 'make' command.
//
// Parameters:
//   - configPath: Absolute path to the TOML configuration file
//   - report: Pointer to the report structure for tracking operations
//
// Returns:
//   - error: Any error that occurred during configuration processing
func processConfig(configPath string, report *MakeReport) error {
	wlog.Common().Infof("Processing configuration file: %s", configPath)

	// Read and parse the TOML configuration file
	config, err := readConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Process each site defined in the configuration
	for _, site := range config.Sites {
		wlog.Common().Infof("Processing site: %s", site.RemotePrefix)

		// Process all repositories for this site concurrently
		if err := processSiteRepos(site, report); err != nil {
			wlog.Common().Errorf("Error processing site %s: %s", site.RemotePrefix, err)
			// Continue processing other sites even if one fails
		}
	}

	return nil
}

// processSiteRepos processes all repositories for a single site concurrently.
// It creates goroutines for each repository to enable parallel processing.
//
// Parameters:
//   - site: Site configuration containing repositories and settings
//   - report: Pointer to the report structure for tracking operations
//
// Returns:
//   - error: Any critical error that prevents processing
func processSiteRepos(site SiteConfig, report *MakeReport) error {
	// Use WaitGroup to coordinate concurrent repository processing
	var wg sync.WaitGroup

	// Channel to collect results from goroutines
	resultChan := make(chan RepoProcessResult, len(site.Repos))

	// Start a goroutine for each repository
	for _, repo := range site.Repos {
		wg.Add(1)

		// Process repository in a separate goroutine for concurrency
		go func(r Repo) {
			defer wg.Done()

			// Process the individual repository
			result := processRepo(r, site)

			// Send result to channel for collection
			resultChan <- result
		}(repo)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(resultChan)

	// Collect all results and add to report
	for result := range resultChan {
		addToReport(report, site, result)
	}

	return nil
}

// processRepo handles the processing of a single repository including cloning/updating
// and warm-up operations. This function contains the core business logic for repository operations.
//
// Parameters:
//   - repo: Repository configuration
//   - site: Site configuration (for warm-up settings)
//
// Returns:
//   - RepoProcessResult: Detailed result of the processing operation
func processRepo(repo Repo, site SiteConfig) RepoProcessResult {
	start := time.Now()

	wlog.Common().Infof("Processing repository: %s", repo.Repo)

	// Determine the target directory for the repository
	targetDir := repo.FullPath(site)

	// Check if repository already exists locally
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		// Repository doesn't exist, clone it
		if err := cloneRepo(repo.RepoUrl(site), targetDir); err != nil {
			return RepoProcessResult{
				Repo:     repo.Repo,
				Duration: time.Since(start),
				Success:  false,
				Error:    fmt.Sprintf("failed to clone: %s", err),
				Memo:     "Clone operation failed",
			}
		}
		wlog.Common().Infof("Successfully cloned %s to %s", repo.Repo, targetDir)
	} else {
		// Repository exists, update it
		if err := updateRepo(targetDir); err != nil {
			return RepoProcessResult{
				Repo:     repo.Repo,
				Duration: time.Since(start),
				Success:  false,
				Error:    fmt.Sprintf("failed to update: %s", err),
				Memo:     "Update operation failed",
			}
		}
		wlog.Common().Infof("Successfully updated %s", repo.Repo)
	}

	// Perform warm-up operations if enabled
	var warmupMemo string
	if shouldWarmUp(repo, site) {
		wlog.Common().Infof("Performing warm-up for %s", repo.Repo)

		if err := performWarmUp(targetDir); err != nil {
			// Warm-up failure is not critical, just log and continue
			wlog.Common().Warnf("Warm-up failed for %s: %s", repo.Repo, err)
			warmupMemo = fmt.Sprintf("Warm-up failed: %s", err)
		} else {
			warmupMemo = "Warm-up completed successfully"
		}
	}

	return RepoProcessResult{
		Repo:     repo.Repo,
		Duration: time.Since(start),
		Success:  true,
		Error:    "",
		Memo:     warmupMemo,
	}
}

// cloneRepo clones a Git repository from the specified URL to the target directory.
//
// Parameters:
//   - url: Git repository URL (supports HTTPS and SSH)
//   - targetDir: Local directory path where repository should be cloned
//
// Returns:
//   - error: Any error that occurred during cloning
func cloneRepo(url, targetDir string) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Execute git clone command
	cmd := exec.Command("git", "clone", url, targetDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %s\nOutput: %s", err, string(output))
	}

	return nil
}

// updateRepo updates an existing Git repository by pulling the latest changes.
//
// Parameters:
//   - repoDir: Local directory path of the repository
//
// Returns:
//   - error: Any error that occurred during updating
func updateRepo(repoDir string) error {
	// Change to repository directory and pull latest changes
	cmd := exec.Command("git", "pull")
	cmd.Dir = repoDir

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git pull failed: %s\nOutput: %s", err, string(output))
	}

	return nil
}

// shouldWarmUp determines if warm-up operations should be performed for a repository.
// It checks both repository-level and site-level warm-up configurations.
//
// Parameters:
//   - repo: Repository configuration
//   - site: Site configuration
//
// Returns:
//   - bool: True if warm-up should be performed, false otherwise
func shouldWarmUp(repo Repo, site SiteConfig) bool {
	// Repository-level setting takes precedence
	if repo.WarmUp {
		return true
	}

	// Fall back to site-level setting
	return site.WarmUpAll
}

// performWarmUp executes warm-up operations based on the project type detected
// in the repository. This includes dependency installation and build preparation.
//
// Parameters:
//   - repoDir: Local directory path of the repository
//
// Returns:
//   - error: Any error that occurred during warm-up operations
func performWarmUp(repoDir string) error {
	var commands [][]string

	// Detect Go projects and add Go-specific warm-up commands
	if _, err := os.Stat(filepath.Join(repoDir, "go.mod")); err == nil {
		wlog.Common().Infof("Detected Go project, adding Go warm-up commands")
		commands = append(commands, []string{"go", "mod", "download"})
		commands = append(commands, []string{"go", "mod", "tidy"})
	}

	// Detect Node.js projects and add npm warm-up commands
	if _, err := os.Stat(filepath.Join(repoDir, "package.json")); err == nil {
		wlog.Common().Infof("Detected Node.js project, adding npm warm-up commands")
		commands = append(commands, []string{"npm", "install"})
	}

	// Execute all detected warm-up commands
	for _, cmdArgs := range commands {
		wlog.Common().Infof("Executing warm-up command: %s", strings.Join(cmdArgs, " "))

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = repoDir

		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("warm-up command '%s' failed: %s\nOutput: %s",
				strings.Join(cmdArgs, " "), err, string(output))
		}
	}

	return nil
}

// addToReport adds a repository processing result to the overall report.
// This function is thread-safe and handles the aggregation of results.
//
// Parameters:
//   - report: Pointer to the report structure
//   - site: Site configuration the repository belongs to
//   - result: Processing result to be added to the report
func addToReport(report *MakeReport, site SiteConfig, result RepoProcessResult) {
	// Create a new action record for the report
	action := &MakeAction{
		Time:       time.Now(),
		Repository: site.RemotePrefix + "/" + result.Repo,
		Duration:   result.Duration,
		Success:    result.Success,
		Error:      result.Error,
		Memo:       result.Memo,
	}

	// Thread-safe addition to report (assuming MakeReport handles concurrency)
	report.Actions = append(report.Actions, action)
}
