package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Clone clones a Git repository from URL to target directory
func Clone(url, targetDir string) error {
	// Ensure parent directory exists
	parentDir := filepath.Dir(targetDir)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	cmd := exec.Command("git", "clone", url, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// Update updates an existing Git repository by pulling latest changes
func Update(repoDir string) error {
	// Check if it's a valid Git repository
	if !isGitRepository(repoDir) {
		return fmt.Errorf("not a valid Git repository: %s", repoDir)
	}

	// Fetch latest changes
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = repoDir
	if output, err := fetchCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git fetch failed: %w\nOutput: %s", err, string(output))
	}

	// Get current branch
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchCmd.Dir = repoDir
	branchOutput, err := branchCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	currentBranch := strings.TrimSpace(string(branchOutput))

	// Pull changes
	pullCmd := exec.Command("git", "pull", "origin", currentBranch)
	pullCmd.Dir = repoDir
	if output, err := pullCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git pull failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// isGitRepository checks if a directory is a Git repository
func isGitRepository(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	if info, err := os.Stat(gitDir); err == nil {
		return info.IsDir() || info.Mode().IsRegular() // Handle both .git directory and .git file (submodules)
	}
	return false
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch(repoDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetRemoteURL returns the remote URL for the specified remote (default: origin)
func GetRemoteURL(repoDir, remote string) (string, error) {
	if remote == "" {
		remote = "origin"
	}
	
	cmd := exec.Command("git", "remote", "get-url", remote)
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
} 