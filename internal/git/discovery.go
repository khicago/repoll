package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RepositoryInfo holds information about a discovered Git repository
type RepositoryInfo struct {
	Path        string
	Origin      string
	HasOrigin   bool
	Uncommitted bool
	Unmerged    bool
}

// DiscoverRepository discovers Git repository information from a directory path
func DiscoverRepository(path string) (*RepositoryInfo, error) {
	info := &RepositoryInfo{
		Path: path,
	}

	// Check if it's a Git repository
	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a Git repository: %s", path)
	}

	// Get origin URL
	origin, err := getGitOrigin(path)
	if err == nil && origin != "" {
		info.Origin = origin
		info.HasOrigin = true
	}

	// Check for uncommitted changes
	hasUncommitted, err := hasUncommittedChanges(path)
	if err == nil {
		info.Uncommitted = hasUncommitted
	}

	// Check for unmerged changes
	hasUnmerged, err := hasUnmergedChanges(path)
	if err == nil {
		info.Unmerged = hasUnmerged
	}

	return info, nil
}

// getGitOrigin gets the origin URL of a Git repository
func getGitOrigin(repoPath string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// hasUncommittedChanges checks if a repository has uncommitted changes
func hasUncommittedChanges(repoPath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// hasUnmergedChanges checks if a repository has unmerged changes
func hasUnmergedChanges(repoPath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "UU ") || strings.HasPrefix(line, "AA ") {
			return true, nil
		}
	}
	return false, scanner.Err()
}

// ExtractRepoNameFromURL extracts repository name from a Git URL
func ExtractRepoNameFromURL(url string) string {
	// Handle different URL formats
	url = strings.TrimSpace(url)
	
	// Remove .git suffix if present
	if strings.HasSuffix(url, ".git") {
		url = strings.TrimSuffix(url, ".git")
	}
	
	// Handle SSH format: git@github.com:owner/repo
	if strings.Contains(url, "@") && strings.Contains(url, ":") {
		parts := strings.Split(url, ":")
		if len(parts) >= 2 {
			return parts[len(parts)-1]
		}
	}
	
	// Handle HTTPS format: https://github.com/owner/repo
	if strings.Contains(url, "://") {
		parts := strings.Split(url, "/")
		if len(parts) >= 2 {
			return strings.Join(parts[len(parts)-2:], "/")
		}
	}
	
	// Fallback: try to extract from the end of the URL
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	
	return ""
}

// ExtractRemotePrefix extracts remote prefix from a Git origin URL
func ExtractRemotePrefix(origin string) string {
	origin = strings.TrimSpace(origin)
	
	// Handle SSH format: git@github.com:owner/repo -> https://github.com/
	if strings.Contains(origin, "@") && strings.Contains(origin, ":") {
		parts := strings.Split(origin, ":")
		if len(parts) >= 2 {
			hostPart := parts[0]
			if strings.Contains(hostPart, "@") {
				host := strings.Split(hostPart, "@")[1]
				return "https://" + host + "/"
			}
		}
	}
	
	// Handle HTTPS format: https://github.com/owner/repo -> https://github.com/
	if strings.Contains(origin, "://") {
		parts := strings.Split(origin, "/")
		if len(parts) >= 3 {
			return strings.Join(parts[:3], "/") + "/"
		}
	}
	
	return ""
}

// IsURL checks if a string is a URL
func IsURL(str string) bool {
	return strings.HasPrefix(str, "http://") || 
		   strings.HasPrefix(str, "https://") || 
		   strings.HasPrefix(str, "git@")
} 