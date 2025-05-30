// Package main contains tests for URL parsing functionality.
// This file provides comprehensive test coverage for Git repository URL parsing
// and prefix extraction functions.
package main

import (
	"testing"
)

func TestGetRepoNameFromURL_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "HTTPS URL with .git suffix",
			url:      "https://github.com/user/repo.git",
			expected: "repo",
		},
		{
			name:     "HTTPS URL without .git suffix",
			url:      "https://github.com/user/repo",
			expected: "repo",
		},
		{
			name:     "SSH URL with .git suffix",
			url:      "git@github.com:user/repo.git",
			expected: "repo",
		},
		{
			name:     "SSH URL without .git suffix",
			url:      "git@github.com:user/repo",
			expected: "repo",
		},
		{
			name:     "Complex URL with multiple slashes",
			url:      "https://gitlab.example.com/group/subgroup/project.git",
			expected: "project",
		},
		{
			name:     "Empty string",
			url:      "",
			expected: "",
		},
		{
			name:     "Single word URL",
			url:      "simple",
			expected: "simple",
		},
		{
			name:     "URL with only .git",
			url:      ".git",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRepoNameFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("getRepoNameFromURL(%q) = %q, expected %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestGetRemotePrefix_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected string
	}{
		{
			name:     "HTTPS GitHub URL",
			origin:   "https://github.com/user/repo.git",
			expected: "https://github.com/user",
		},
		{
			name:     "HTTPS GitLab URL",
			origin:   "https://gitlab.com/group/repo.git",
			expected: "https://gitlab.com/group",
		},
		{
			name:     "SSH GitHub URL",
			origin:   "git@github.com:user/repo.git",
			expected: "git@github.com:user",
		},
		{
			name:     "SSH GitLab URL with slash",
			origin:   "git@gitlab.example.com:group/subgroup/repo.git",
			expected: "git@gitlab.example.com:group/subgroup",
		},
		{
			name:     "Unknown format",
			origin:   "unknown-format-url",
			expected: "unknown-format-url",
		},
		{
			name:     "URL with whitespace",
			origin:   "  https://github.com/user/repo.git  ",
			expected: "https://github.com/user",
		},
		{
			name:     "SSH URL without colon and slash",
			origin:   "git@server",
			expected: "git@server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRemotePrefix(tt.origin)
			if result != tt.expected {
				t.Errorf("getRemotePrefix(%q) = %q, expected %q", tt.origin, result, tt.expected)
			}
		})
	}
}

func TestExtractHTTPRemotePrefix_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected string
	}{
		{
			name:     "Standard HTTPS URL",
			origin:   "https://github.com/user/repo.git",
			expected: "https://github.com/user",
		},
		{
			name:     "HTTPS URL without .git",
			origin:   "https://github.com/user/repo",
			expected: "https://github.com/user",
		},
		{
			name:     "HTTPS URL with whitespace",
			origin:   "  https://gitlab.com/group/project.git  ",
			expected: "https://gitlab.com/group",
		},
		{
			name:     "URL without slash (fallback case)",
			origin:   "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "Empty URL",
			origin:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractHTTPRemotePrefix(tt.origin)
			if result != tt.expected {
				t.Errorf("extractHTTPRemotePrefix(%q) = %q, expected %q", tt.origin, result, tt.expected)
			}
		})
	}
}

func TestExtractSSHRemotePrefix_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected string
	}{
		{
			name:     "SSH URL with colon separator",
			origin:   "git@github.com:user/repo.git",
			expected: "git@github.com:user",
		},
		{
			name:     "SSH URL with slash and colon",
			origin:   "git@gitlab.example.com:group/subgroup/repo.git",
			expected: "git@gitlab.example.com:group/subgroup",
		},
		{
			name:     "SSH URL without @ or :",
			origin:   "invalid-ssh-format",
			expected: "invalid-ssh-format",
		},
		{
			name:     "SSH URL without colon",
			origin:   "git@server.com",
			expected: "git@server.com",
		},
		{
			name:     "SSH URL with only colon at end",
			origin:   "git@server:",
			expected: "git@server:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSSHRemotePrefix(tt.origin)
			if result != tt.expected {
				t.Errorf("extractSSHRemotePrefix(%q) = %q, expected %q", tt.origin, result, tt.expected)
			}
		})
	}
}

func TestGetRepoNameFromURL_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "URL with multiple .git occurrences",
			url:      "https://github.com/user/repo.git.git",
			expected: "repo.git",
		},
		{
			name:     "URL with query parameters",
			url:      "https://github.com/user/repo.git?ref=main",
			expected: "repo.git?ref=main",
		},
		{
			name:     "URL with fragment",
			url:      "https://github.com/user/repo.git#readme",
			expected: "repo.git#readme",
		},
		{
			name:     "URL with trailing slash",
			url:      "https://github.com/user/repo/",
			expected: "",
		},
		{
			name:     "Just a file name",
			url:      "repository.git",
			expected: "repository",
		},
		{
			name:     "Path with backslashes (Windows style)",
			url:      "C:\\path\\to\\repo.git",
			expected: "repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRepoNameFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("getRepoNameFromURL(%q) = %q, expected %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestGetRemotePrefix_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected string
	}{
		{
			name:     "URL with port number",
			origin:   "https://gitlab.example.com:8080/group/repo.git",
			expected: "https://gitlab.example.com:8080/group",
		},
		{
			name:     "SSH URL with port",
			origin:   "ssh://git@gitlab.com:2222/user/repo.git",
			expected: "ssh://git@gitlab.com:2222/user/repo.git",
		},
		{
			name:     "Very nested path",
			origin:   "https://github.com/org/team/subteam/project/repo.git",
			expected: "https://github.com/org/team/subteam/project",
		},
		{
			name:     "Empty string",
			origin:   "",
			expected: "",
		},
		{
			name:     "Just whitespace",
			origin:   "   \t\n   ",
			expected: "",
		},
		{
			name:     "SSH without colon",
			origin:   "git@server.com",
			expected: "git@server.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRemotePrefix(tt.origin)
			if result != tt.expected {
				t.Errorf("getRemotePrefix(%q) = %q, expected %q", tt.origin, result, tt.expected)
			}
		})
	}
}
