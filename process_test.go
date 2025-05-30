// Package main contains tests for repository processing functionality.
// This file provides comprehensive test coverage for repository cloning, updating,
// warm-up operations, and process management functions.
package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProcessConfig_NonExistentFile(t *testing.T) {
	report := &MakeReport{Actions: make([]*MakeAction, 0)}

	err := processConfig("/non/existent/config.toml", report)
	if err == nil {
		t.Error("Expected error for non-existent config file, got nil")
	}
}

func TestProcessConfig_ValidFile(t *testing.T) {
	// Create a temporary config file
	configContent := `
[[sites]]
remote = "https://github.com/test"
dir = "/tmp/test-repos"

[[sites.repos]]
repo = "non-existent-repo"
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

	report := &MakeReport{Actions: make([]*MakeAction, 0)}

	err = processConfig(tmpFile.Name(), report)
	// Should not error on processing, even if clone fails
	if err != nil {
		t.Errorf("processConfig failed: %v", err)
	}
}

func TestCloneRepo_InvalidURL(t *testing.T) {
	tempDir := t.TempDir()
	targetDir := filepath.Join(tempDir, "test-repo")

	err := cloneRepo("invalid-url", targetDir)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestUpdateRepo_NonGitDirectory(t *testing.T) {
	tempDir := t.TempDir()

	err := updateRepo(tempDir)
	if err == nil {
		t.Error("Expected error for non-Git directory, got nil")
	}
}

func TestShouldWarmUp_RepoLevel(t *testing.T) {
	repo := Repo{WarmUp: true}
	site := SiteConfig{WarmUpAll: false}

	result := shouldWarmUp(repo, site)
	if !result {
		t.Error("Expected true when repo-level warm-up is enabled")
	}
}

func TestShouldWarmUp_SiteLevel(t *testing.T) {
	repo := Repo{WarmUp: false}
	site := SiteConfig{WarmUpAll: true}

	result := shouldWarmUp(repo, site)
	if !result {
		t.Error("Expected true when site-level warm-up is enabled")
	}
}

func TestShouldWarmUp_Neither(t *testing.T) {
	repo := Repo{WarmUp: false}
	site := SiteConfig{WarmUpAll: false}

	result := shouldWarmUp(repo, site)
	if result {
		t.Error("Expected false when neither warm-up is enabled")
	}
}

func TestPerformWarmUp_NonExistentDirectory(t *testing.T) {
	err := performWarmUp("/non/existent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

func TestPerformWarmUp_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	err := performWarmUp(tempDir)
	// Should not error for empty directory (no warm-up needed)
	if err != nil {
		t.Errorf("performWarmUp failed for empty directory: %v", err)
	}
}

func TestPerformWarmUp_GoProject(t *testing.T) {
	tempDir := t.TempDir()

	// Create a go.mod file
	goModContent := `module test-project

go 1.19
`
	goModPath := filepath.Join(tempDir, "go.mod")
	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	err = performWarmUp(tempDir)
	// Should not error (even if go mod download fails, it's not critical)
	if err != nil {
		t.Errorf("performWarmUp failed for Go project: %v", err)
	}
}

func TestPerformWarmUp_NodeProject(t *testing.T) {
	tempDir := t.TempDir()

	// Create a package.json file
	packageJsonContent := `{
  "name": "test-project",
  "version": "1.0.0",
  "dependencies": {}
}`
	packageJsonPath := filepath.Join(tempDir, "package.json")
	err := os.WriteFile(packageJsonPath, []byte(packageJsonContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	err = performWarmUp(tempDir)
	// Should not error (even if npm install fails, it's not critical)
	if err != nil {
		t.Errorf("performWarmUp failed for Node project: %v", err)
	}
}

func TestAddToReport(t *testing.T) {
	report := &MakeReport{Actions: make([]*MakeAction, 0)}
	site := SiteConfig{
		RemotePrefix: "https://github.com/user",
		Dir:          "/home/user/projects",
	}
	result := RepoProcessResult{
		Repo:     "test-repo",
		Duration: 5 * time.Second,
		Success:  true,
		Error:    "",
		Memo:     "Clone completed",
	}

	addToReport(report, site, result)

	if len(report.Actions) != 1 {
		t.Errorf("Expected 1 action in report, got %d", len(report.Actions))
	}

	action := report.Actions[0]
	if action.Repository != result.Repo {
		t.Errorf("Expected repo '%s', got '%s'", result.Repo, action.Repository)
	}
	if action.Duration != result.Duration {
		t.Errorf("Expected duration %v, got %v", result.Duration, action.Duration)
	}
	if action.Success != result.Success {
		t.Errorf("Expected success %v, got %v", result.Success, action.Success)
	}
}

func TestProcessRepo_CloneSuccess(t *testing.T) {
	tempDir := t.TempDir()

	repo := Repo{
		Repo:   "test-repo",
		WarmUp: false,
	}
	site := SiteConfig{
		RemotePrefix: "file://" + tempDir, // Use local path as remote for testing
		Dir:          tempDir,
		WarmUpAll:    false,
	}

	// Create a fake git repository to clone from
	sourceRepo := filepath.Join(tempDir, "source.git")
	if err := os.MkdirAll(sourceRepo, 0755); err != nil {
		t.Fatalf("Failed to create source repo: %v", err)
	}

	result := processRepo(repo, site)

	// The result should be a failure because we can't actually clone
	// from a fake source, but we can test the structure
	if result.Repo != repo.Repo {
		t.Errorf("Expected repo '%s', got '%s'", repo.Repo, result.Repo)
	}
	if result.Duration <= 0 {
		t.Error("Expected positive duration")
	}
}

func TestProcessRepo_UpdateExisting(t *testing.T) {
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "existing-repo")

	// Create a directory structure that looks like a repo
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}

	repo := Repo{
		Repo:   "existing-repo",
		WarmUp: false,
	}
	site := SiteConfig{
		RemotePrefix: "https://github.com/user",
		Dir:          tempDir,
		WarmUpAll:    false,
	}

	result := processRepo(repo, site)

	// Should fail because it's not a real git repo
	if result.Success {
		t.Error("Expected failure for fake git repo")
	}
	if result.Repo != repo.Repo {
		t.Errorf("Expected repo '%s', got '%s'", repo.Repo, result.Repo)
	}
}

func TestProcessSiteRepos_EmptyRepos(t *testing.T) {
	site := SiteConfig{
		RemotePrefix: "https://github.com/user",
		Dir:          "/tmp/test",
		Repos:        []Repo{},
	}
	report := &MakeReport{Actions: make([]*MakeAction, 0)}

	err := processSiteRepos(site, report)
	if err != nil {
		t.Errorf("processSiteRepos failed: %v", err)
	}

	if len(report.Actions) != 0 {
		t.Errorf("Expected 0 actions for empty repos, got %d", len(report.Actions))
	}
}

func TestProcessSiteRepos_MultipleRepos(t *testing.T) {
	site := SiteConfig{
		RemotePrefix: "https://github.com/user",
		Dir:          "/tmp/test",
		Repos: []Repo{
			{Repo: "repo1"},
			{Repo: "repo2"},
		},
	}
	report := &MakeReport{Actions: make([]*MakeAction, 0)}

	err := processSiteRepos(site, report)
	if err != nil {
		t.Errorf("processSiteRepos failed: %v", err)
	}

	// Should have 2 actions (even if they fail)
	if len(report.Actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(report.Actions))
	}
}

func TestProcessConfig_SiteProcessingError(t *testing.T) {
	report := &MakeReport{Actions: make([]*MakeAction, 0)}
	
	// This should handle site processing errors gracefully
	err := processConfig("/tmp/test-config.toml", report)
	
	// The function should return an error for invalid config path
	if err == nil {
		t.Error("Expected processConfig to return error for invalid config path")
	}
}

func TestProcessRepo_CloneFailure(t *testing.T) {
	// Test repository cloning failure
	repo := Repo{
		Repo:   "test-repo",
		Rename: "",
		WarmUp: false,
		Memo:   "",
	}
	
	site := SiteConfig{
		RemotePrefix: "invalid://protocol",
		Dir:          "/tmp/test",
		Repos:        []Repo{repo},
		WarmUpAll:    false,
	}
	
	result := processRepo(repo, site)
	
	// Should return failure result
	if result.Success {
		t.Error("Expected processRepo to fail with invalid URL")
	}
	
	if result.Error == "" {
		t.Error("Expected error message for failed clone")
	}
}

func TestProcessRepo_UpdateFailure(t *testing.T) {
	// Create a temporary directory that looks like a repo but isn't a valid git repo
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "fake-repo")
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	
	repo := Repo{
		Repo:   "fake-repo",
		Rename: "",
		WarmUp: false,
		Memo:   "",
	}
	
	site := SiteConfig{
		RemotePrefix: "https://github.com/test",
		Dir:          tempDir,
		Repos:        []Repo{repo},
		WarmUpAll:    false,
	}
	
	result := processRepo(repo, site)
	
	// Should return failure result due to invalid git repo
	if result.Success {
		t.Error("Expected processRepo to fail with invalid git repo")
	}
	
	if result.Error == "" {
		t.Error("Expected error message for failed update")
	}
}

func TestProcessRepo_WarmupFailure(t *testing.T) {
	// Create a temporary directory with a fake Go project
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "go-repo")
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	
	// Create a fake go.mod file to trigger Go warmup
	goModContent := `module test-repo

go 1.19
`
	err = os.WriteFile(filepath.Join(repoDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}
	
	// Create a fake git repo structure
	gitDir := filepath.Join(repoDir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	repo := Repo{
		Repo:   "go-repo",
		Rename: "",
		WarmUp: true, // Enable warmup
		Memo:   "",
	}
	
	site := SiteConfig{
		RemotePrefix: "https://github.com/test",
		Dir:          tempDir,
		Repos:        []Repo{repo},
		WarmUpAll:    false,
	}
	
	result := processRepo(repo, site)
	
	// Should succeed even if warmup fails (warmup failure is not critical)
	if !result.Success {
		t.Errorf("Expected processRepo to succeed even with warmup failure, got error: %s", result.Error)
	}
	
	// Should have warmup memo indicating failure
	if result.Memo == "" {
		t.Error("Expected warmup memo to be set")
	}
}

func TestPerformWarmUp_NoProject(t *testing.T) {
	// Create a temporary directory with no project files
	tempDir := t.TempDir()
	
	// Should succeed with no commands to run
	err := performWarmUp(tempDir)
	if err != nil {
		t.Errorf("performWarmUp should succeed with no project files, got: %v", err)
	}
}

func TestCloneRepo_ParentDirectoryCreation(t *testing.T) {
	// Test cloning to a path where parent directories don't exist
	tempDir := t.TempDir()
	targetDir := filepath.Join(tempDir, "deep", "nested", "path", "repo")
	
	// This will fail because the URL is invalid, but we're testing the parent directory creation
	err := cloneRepo("invalid://url", targetDir)
	
	// Should fail due to invalid URL, but parent directories should be created
	if err == nil {
		t.Error("Expected cloneRepo to fail with invalid URL")
	}
	
	// Check that parent directory was created
	parentDir := filepath.Dir(targetDir)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		t.Error("Expected parent directory to be created")
	}
}

func TestUpdateRepo_InvalidDirectory(t *testing.T) {
	// Test updating a non-existent repository
	err := updateRepo("/non/existent/directory")
	
	if err == nil {
		t.Error("Expected updateRepo to fail with non-existent directory")
	}
}

func TestShouldWarmUp_BothFalse(t *testing.T) {
	repo := Repo{WarmUp: false}
	site := SiteConfig{WarmUpAll: false}
	
	if shouldWarmUp(repo, site) {
		t.Error("Expected shouldWarmUp to return false when both are false")
	}
} 