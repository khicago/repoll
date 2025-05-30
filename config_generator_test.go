// Package main contains tests for configuration generation functionality.
// This file provides comprehensive test coverage for configuration generation
// and TOML file creation functions.
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateConfig_EmptyResults(t *testing.T) {
	// Test with empty discovery results
	results := []RepoDiscoveryResult{}

	config := generateConfigStruct(results)

	if len(config.Sites) != 0 {
		t.Errorf("Expected 0 sites for empty results, got %d", len(config.Sites))
	}
}

func TestGenerateConfig_SingleRepo(t *testing.T) {
	// Test with a single repository
	results := []RepoDiscoveryResult{
		{
			Path:       "/home/user/projects/repo1",
			Origin:     "https://github.com/user/repo1.git",
			Uncommitted: false,
			Unmerged:   false,
		},
	}

	config := generateConfigStruct(results)

	if len(config.Sites) != 1 {
		t.Fatalf("Expected 1 site, got %d", len(config.Sites))
	}

	site := config.Sites[0]
	if site.RemotePrefix != "https://github.com/user" {
		t.Errorf("Expected remote prefix 'https://github.com/user', got '%s'", site.RemotePrefix)
	}

	if len(site.Repos) != 1 {
		t.Fatalf("Expected 1 repo, got %d", len(site.Repos))
	}

	repo := site.Repos[0]
	if repo.Repo != "repo1" {
		t.Errorf("Expected repo name 'repo1', got '%s'", repo.Repo)
	}
}

func TestGenerateConfig_MultipleRepos(t *testing.T) {
	// Test with multiple repositories from different remotes
	results := []RepoDiscoveryResult{
		{
			Path:       "/home/user/projects/repo1",
			Origin:     "https://github.com/user/repo1.git",
			Uncommitted: false,
			Unmerged:   false,
		},
		{
			Path:       "/home/user/projects/repo2",
			Origin:     "https://github.com/user/repo2.git",
			Uncommitted: true,
			Unmerged:   false,
		},
		{
			Path:       "/home/user/projects/other-repo",
			Origin:     "git@gitlab.com:org/other-repo.git",
			Uncommitted: false,
			Unmerged:   true,
		},
	}

	config := generateConfigStruct(results)

	if len(config.Sites) != 2 {
		t.Fatalf("Expected 2 sites, got %d", len(config.Sites))
	}

	// Find GitHub site
	var githubSite *SiteConfig
	var gitlabSite *SiteConfig

	for i := range config.Sites {
		if strings.Contains(config.Sites[i].RemotePrefix, "github.com") {
			githubSite = &config.Sites[i]
		} else if strings.Contains(config.Sites[i].RemotePrefix, "gitlab.com") {
			gitlabSite = &config.Sites[i]
		}
	}

	if githubSite == nil {
		t.Fatal("GitHub site not found")
	}
	if gitlabSite == nil {
		t.Fatal("GitLab site not found")
	}

	if len(githubSite.Repos) != 2 {
		t.Errorf("Expected 2 repos in GitHub site, got %d", len(githubSite.Repos))
	}

	if len(gitlabSite.Repos) != 1 {
		t.Errorf("Expected 1 repo in GitLab site, got %d", len(gitlabSite.Repos))
	}
}

func TestGenerateSiteConfig_Basic(t *testing.T) {
	// Test basic site configuration generation
	remotePrefix := "https://github.com/user"
	repos := []RepoDiscoveryResult{
		{
			Path:       "/home/user/projects/repo1",
			Origin:     "https://github.com/user/repo1.git",
			Uncommitted: false,
			Unmerged:   false,
		},
	}
	
	site := generateSiteConfigStruct(remotePrefix, repos)
	
	if site.RemotePrefix != remotePrefix {
		t.Errorf("Expected remote prefix '%s', got '%s'", remotePrefix, site.RemotePrefix)
	}
	
	if site.Dir != "/home/user/projects" {
		t.Errorf("Expected dir '/home/user/projects', got '%s'", site.Dir)
	}
	
	if len(site.Repos) != 1 {
		t.Fatalf("Expected 1 repo, got %d", len(site.Repos))
	}
	
	repo := site.Repos[0]
	if repo.Repo != "repo1" {
		t.Errorf("Expected repo name 'repo1', got '%s'", repo.Repo)
	}
	
	if repo.Memo != "" {
		t.Errorf("Expected empty memo for clean repo, got '%s'", repo.Memo)
	}
}

func TestGenerateSiteConfig_WithStatus(t *testing.T) {
	// Test site configuration with repository status
	remotePrefix := "https://github.com/user"
	repos := []RepoDiscoveryResult{
		{
			Path:       "/home/user/projects/repo1",
			Origin:     "https://github.com/user/repo1.git",
			Uncommitted: true,
			Unmerged:   false,
		},
		{
			Path:       "/home/user/projects/repo2",
			Origin:     "https://github.com/user/repo2.git",
			Uncommitted: false,
			Unmerged:   true,
		},
	}
	
	site := generateSiteConfigStruct(remotePrefix, repos)
	
	if len(site.Repos) != 2 {
		t.Fatalf("Expected 2 repos, got %d", len(site.Repos))
	}
	
	// Check first repo with uncommitted changes
	repo1 := site.Repos[0]
	if repo1.Memo != "uncommitted" {
		t.Errorf("Expected memo 'uncommitted' for repo with uncommitted changes, got '%s'", repo1.Memo)
	}
	
	// Check second repo with unmerged commits
	repo2 := site.Repos[1]
	if repo2.Memo != "unmerged" {
		t.Errorf("Expected memo 'unmerged' for repo with unmerged commits, got '%s'", repo2.Memo)
	}
}

func TestMakeConfig_NonExistentDirectory(t *testing.T) {
	// Test with non-existent directory
	report := &MkconfReport{}

	err := makeConfig("/non/existent/directory", report)

	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

func TestMakeConfig_EmptyDirectory(t *testing.T) {
	// Create a temporary empty directory
	tempDir, err := os.MkdirTemp("", "test_empty_dir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	report := &MkconfReport{}

	err = makeConfig(tempDir, report)
	if err != nil {
		t.Fatalf("makeConfig failed: %v", err)
	}

	// Should have no repositories discovered
	if len(report.Actions) != 0 {
		t.Errorf("Expected 0 repositories in empty directory, got %d", len(report.Actions))
	}
}

func TestMakeConfig_WithGitRepo(t *testing.T) {
	// Create a temporary directory with a git repository
	tempDir, err := os.MkdirTemp("", "test_git_repo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a subdirectory for the git repo
	repoDir := filepath.Join(tempDir, "test-repo")
	err = os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}

	// Initialize git repository
	gitDir := filepath.Join(repoDir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	// Create a basic git config file
	configContent := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = https://github.com/user/test-repo.git
	fetch = +refs/heads/*:refs/remotes/origin/*
`
	configPath := filepath.Join(gitDir, "config")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write git config: %v", err)
	}

	report := &MkconfReport{}

	err = makeConfig(tempDir, report)
	if err != nil {
		t.Fatalf("makeConfig failed: %v", err)
	}

	// Should have discovered one repository
	if len(report.Actions) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(report.Actions))
	}

	if len(report.Actions) > 0 {
		action := report.Actions[0]
		if !strings.Contains(action.Path, "test-repo") {
			t.Errorf("Expected repo path to contain 'test-repo', got '%s'", action.Path)
		}

		if action.Origin != "https://github.com/user/test-repo.git" {
			t.Errorf("Expected origin 'https://github.com/user/test-repo.git', got '%s'", action.Origin)
		}
	}
}
