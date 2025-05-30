// Package main contains URL parsing functionality for Git repository URLs.
// This file handles parsing and extracting information from various Git URL formats
// including HTTPS and SSH URLs from different hosting services.
package main

import (
	"strings"

	"github.com/bagaking/goulp/wlog"
)

// getRepoNameFromURL extracts the repository name from a Git URL.
// Supports both HTTPS and SSH URL formats and handles .git suffix removal.
//
// Examples:
//   - "https://github.com/user/repo.git" → "repo"
//   - "git@github.com:user/repo.git" → "repo"
//   - "https://github.com/user/repo" → "repo"
//
// Parameters:
//   - url: Git repository URL in any supported format
//
// Returns:
//   - string: Repository name extracted from URL
func getRepoNameFromURL(url string) string {
	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	// Extract the last component of the path
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	// Fallback for malformed URLs
	return url
}

// getRemotePrefix extracts the remote prefix (base URL) from a Git repository URL.
// This is used to group repositories by their hosting service or organization.
//
// Examples:
//   - "https://github.com/user/repo.git" → "https://github.com/user"
//   - "git@github.com:user/repo.git" → "git@github.com:user"
//   - "git@gitlab.example.com:group/repo.git" → "git@gitlab.example.com:group"
//
// Parameters:
//   - origin: Git remote origin URL
//
// Returns:
//   - string: Remote prefix for repository grouping
func getRemotePrefix(origin string) string {
	origin = strings.TrimSpace(origin)

	// Handle HTTPS URLs
	if strings.HasPrefix(origin, "https://") {
		return extractHTTPRemotePrefix(origin)
	}

	// Handle SSH URLs (git@host:path format)
	if strings.Contains(origin, "@") && strings.Contains(origin, ":") {
		// Find the last '/' or ':' before the repository name
		lastSlash := strings.LastIndex(origin, "/")
		lastColon := strings.LastIndex(origin, ":")

		// Use the later position (closer to repo name)
		cutPos := lastColon
		if lastSlash > lastColon {
			cutPos = lastSlash
		}

		if cutPos > 0 {
			return origin[:cutPos]
		}
	}

	// Fallback: return origin as-is for unknown formats
	wlog.Common().Warnf("Unknown URL format, using full origin as prefix: %s", origin)
	return origin
}

// extractHTTPRemotePrefix extracts the remote prefix from HTTPS Git URLs.
// This helper function specifically handles HTTP/HTTPS URL parsing.
//
// Parameters:
//   - origin: HTTPS Git repository URL
//
// Returns:
//   - string: Remote prefix for HTTPS URLs
func extractHTTPRemotePrefix(origin string) string {
	// Remove .git suffix and trim
	origin = strings.TrimSuffix(origin, ".git")
	origin = strings.TrimSpace(origin)

	// Find the last '/' to separate repo name from prefix
	lastSlash := strings.LastIndex(origin, "/")
	if lastSlash > 0 {
		return origin[:lastSlash]
	}

	// Fallback for malformed URLs
	return origin
}

// extractSSHRemotePrefix extracts the remote prefix from SSH Git URLs.
// This helper function specifically handles SSH URL parsing (git@host:path format).
//
// Parameters:
//   - origin: SSH Git repository URL
//
// Returns:
//   - string: Remote prefix for SSH URLs
func extractSSHRemotePrefix(origin string) string {
	// Handle SSH URLs (git@host:path format)
	if strings.Contains(origin, "@") && strings.Contains(origin, ":") {
		// Find the last '/' or ':' before the repository name
		lastSlash := strings.LastIndex(origin, "/")
		lastColon := strings.LastIndex(origin, ":")

		// Use the later position (closer to repo name)
		cutPos := lastColon
		if lastSlash > lastColon {
			cutPos = lastSlash
		}

		if cutPos > 0 {
			return origin[:cutPos]
		}
	}

	// Fallback for malformed URLs
	return origin
}
