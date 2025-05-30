// Package main contains tests for configuration handling functionality.
// This file provides comprehensive test coverage for TOML configuration
// reading, parsing, and repository URL/path calculation functions.
package main

import (
	"os"
	"testing"
)

func TestReadConfig_ValidFile(t *testing.T) {
	// Create a temporary TOML config file
	configContent := `
[[sites]]
remote = "https://github.com/user"
dir = "/home/user/projects"

[[sites.repos]]
repo = "test-repo"
rename = "custom-name"
warm_up = true
memo = "test memo"

[[sites]]
remote = "git@gitlab.com:org"
dir = "/home/user/gitlab"

[[sites.repos]]
repo = "another-repo"
`

	tmpFile, err := os.CreateTemp("", "test_config_*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test reading the config
	config, err := readConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("readConfig failed: %v", err)
	}

	// Verify the parsed configuration
	if len(config.Sites) != 2 {
		t.Errorf("Expected 2 sites, got %d", len(config.Sites))
	}

	site1 := config.Sites[0]
	if site1.RemotePrefix != "https://github.com/user" {
		t.Errorf("Expected remote 'https://github.com/user', got '%s'", site1.RemotePrefix)
	}

	if len(site1.Repos) != 1 {
		t.Errorf("Expected 1 repo in first site, got %d", len(site1.Repos))
	}

	repo := site1.Repos[0]
	if repo.Repo != "test-repo" {
		t.Errorf("Expected repo 'test-repo', got '%s'", repo.Repo)
	}
	if repo.Rename != "custom-name" {
		t.Errorf("Expected rename 'custom-name', got '%s'", repo.Rename)
	}
	if !repo.WarmUp {
		t.Error("Expected WarmUp to be true")
	}
	if repo.Memo != "test memo" {
		t.Errorf("Expected memo 'test memo', got '%s'", repo.Memo)
	}
}

func TestReadConfig_NonExistentFile(t *testing.T) {
	_, err := readConfig("/non/existent/file.toml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestReadConfig_InvalidTOML(t *testing.T) {
	// Create a file with invalid TOML content
	tmpFile, err := os.CreateTemp("", "invalid_*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidContent := `[invalid toml content`
	if _, err := tmpFile.Write([]byte(invalidContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	_, err = readConfig(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid TOML, got nil")
	}
}

func TestRepoUrl_HTTPSRemote(t *testing.T) {
	repo := Repo{Repo: "test-repo"}
	site := SiteConfig{RemotePrefix: "https://github.com/user"}

	url := repo.RepoUrl(site)
	expected := "https://github.com/user/test-repo.git"

	if url != expected {
		t.Errorf("Expected URL '%s', got '%s'", expected, url)
	}
}

func TestRepoUrl_SSHRemote(t *testing.T) {
	repo := Repo{Repo: "test-repo"}
	site := SiteConfig{RemotePrefix: "git@github.com:user"}

	url := repo.RepoUrl(site)
	expected := "git@github.com:usertest-repo.git"

	if url != expected {
		t.Errorf("Expected URL '%s', got '%s'", expected, url)
	}
}

func TestRepoUrl_WithWhitespace(t *testing.T) {
	repo := Repo{Repo: "  test-repo  "}
	site := SiteConfig{RemotePrefix: "https://github.com/user"}

	url := repo.RepoUrl(site)
	expected := "https://github.com/user/test-repo.git"

	if url != expected {
		t.Errorf("Expected URL '%s', got '%s'", expected, url)
	}
}

func TestFullPath_Default(t *testing.T) {
	repo := Repo{Repo: "test-repo"}
	site := SiteConfig{Dir: "/home/user/projects"}

	path := repo.FullPath(site)
	expected := "/home/user/projects/test-repo"

	if path != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, path)
	}
}

func TestFullPath_WithRename(t *testing.T) {
	repo := Repo{Repo: "test-repo", Rename: "custom-name"}
	site := SiteConfig{Dir: "/home/user/projects"}

	path := repo.FullPath(site)
	expected := "/home/user/projects/custom-name"

	if path != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, path)
	}
}

func TestFullPath_BaseRename(t *testing.T) {
	repo := Repo{Repo: "complex/path/test-repo", Rename: "{base}"}
	site := SiteConfig{Dir: "/home/user/projects"}

	path := repo.FullPath(site)
	expected := "/home/user/projects/test-repo"

	if path != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, path)
	}
}

func TestFullPath_WithWhitespace(t *testing.T) {
	repo := Repo{Repo: "  test-repo  ", Rename: "  custom-name  "}
	site := SiteConfig{Dir: "/home/user/projects"}

	path := repo.FullPath(site)
	expected := "/home/user/projects/custom-name"

	if path != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, path)
	}
}

func TestIsURL_HTTPS(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"https://github.com/user", true},
		{"http://example.com", true},
		{"git@github.com:user", false},
		{"ftp://example.com", false},
		{"", false},
		{"not-a-url", false},
	}

	for _, tt := range tests {
		result := isURL(tt.input)
		if result != tt.expected {
			t.Errorf("isURL(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestEnsureTrailingSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://github.com/user", "https://github.com/user/"},
		{"https://github.com/user/", "https://github.com/user/"},
		{"", "/"},
		{"/path", "/path/"},
		{"/path/", "/path/"},
	}

	for _, tt := range tests {
		result := ensureTrailingSlash(tt.input)
		if result != tt.expected {
			t.Errorf("ensureTrailingSlash(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
} 