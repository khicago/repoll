// Package main contains tests for Git repository discovery functionality.
// This file provides comprehensive test coverage for Git repository discovery,
// status checking, and metadata extraction functions.
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuildStatusMemo_Comprehensive(t *testing.T) {
	tests := []struct {
		name        string
		uncommitted bool
		unmerged    bool
		expected    string
	}{
		{
			name:        "Clean repository",
			uncommitted: false,
			unmerged:    false,
			expected:    "",
		},
		{
			name:        "Uncommitted changes only",
			uncommitted: true,
			unmerged:    false,
			expected:    "uncommitted",
		},
		{
			name:        "Unmerged commits only",
			uncommitted: false,
			unmerged:    true,
			expected:    "unmerged",
		},
		{
			name:        "Both uncommitted and unmerged (unmerged takes priority)",
			uncommitted: true,
			unmerged:    true,
			expected:    "unmerged",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildStatusMemo(tt.uncommitted, tt.unmerged)
			if result != tt.expected {
				t.Errorf("buildStatusMemo(%v, %v) = %q, expected %q", tt.uncommitted, tt.unmerged, result, tt.expected)
			}
		})
	}
}

func TestCalculateRelativeLocalPath_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		repoPath string
		repoName string
		baseDir  string
		expected string
	}{
		{
			name:     "Simple relative path",
			repoPath: "/home/user/projects/myrepo",
			repoName: "myrepo",
			baseDir:  "/home/user/projects",
			expected: "/home/user/projects/",
		},
		{
			name:     "Nested directory structure",
			repoPath: "/home/user/projects/group/myrepo",
			repoName: "myrepo",
			baseDir:  "/home/user/projects",
			expected: "/home/user/projects/group/",
		},
		{
			name:     "Same directory as base",
			repoPath: "/home/user/projects",
			repoName: "projects",
			baseDir:  "/home/user/projects",
			expected: "/home/user/projects/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateRelativeLocalPath(tt.repoPath, tt.repoName, tt.baseDir)
			if result != tt.expected {
				t.Errorf("calculateRelativeLocalPath(%q, %q, %q) = %q, expected %q",
					tt.repoPath, tt.repoName, tt.baseDir, result, tt.expected)
			}
		})
	}
}

func TestAddToMkconfReport(t *testing.T) {
	report := &MkconfReport{
		Actions: make([]*MkconfAction, 0),
	}

	result := &RepoDiscoveryResult{
		Path:        "/test/path",
		Origin:      "https://github.com/user/repo.git",
		HasOrigin:   true,
		Uncommitted: true,
		Unmerged:    false,
	}

	addToMkconfReport(report, result)

	if len(report.Actions) != 1 {
		t.Errorf("Expected 1 action in report, got %d", len(report.Actions))
	}

	action := report.Actions[0]
	if action.Path != result.Path {
		t.Errorf("Expected action path %q, got %q", result.Path, action.Path)
	}
	if action.Origin != result.Origin {
		t.Errorf("Expected action origin %q, got %q", result.Origin, action.Origin)
	}
	if action.HasOrigin != result.HasOrigin {
		t.Errorf("Expected action HasOrigin %v, got %v", result.HasOrigin, action.HasOrigin)
	}
	if action.Uncommitted != result.Uncommitted {
		t.Errorf("Expected action Uncommitted %v, got %v", result.Uncommitted, action.Uncommitted)
	}
	if action.Unmerged != result.Unmerged {
		t.Errorf("Expected action Unmerged %v, got %v", result.Unmerged, action.Unmerged)
	}
}

func TestBuildSiteConfig_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		result   *RepoDiscoveryResult
		baseDir  string
		wantKey  string
		wantSite SiteConfig
	}{
		{
			name: "Valid repository with origin",
			result: &RepoDiscoveryResult{
				Path:        "/home/user/projects/myrepo",
				Origin:      "https://github.com/user/myrepo.git",
				HasOrigin:   true,
				Uncommitted: true,
				Unmerged:    false,
			},
			baseDir: "/home/user/projects",
			wantKey: "https://github.com/user @@ /home/user/projects",
			wantSite: SiteConfig{
				RemotePrefix: "https://github.com/user",
				Dir:          "/home/user/projects",
				Repos: []Repo{
					{Repo: "myrepo", Memo: "uncommitted"},
				},
			},
		},
		{
			name: "Repository without origin",
			result: &RepoDiscoveryResult{
				Path:      "/home/user/projects/myrepo",
				HasOrigin: false,
			},
			baseDir:  "/home/user/projects",
			wantKey:  "",
			wantSite: SiteConfig{},
		},
		{
			name: "Repository with unmerged commits",
			result: &RepoDiscoveryResult{
				Path:        "/home/user/projects/myrepo",
				Origin:      "git@github.com:user/myrepo.git",
				HasOrigin:   true,
				Uncommitted: false,
				Unmerged:    true,
			},
			baseDir: "/home/user/projects",
			wantKey: "git@github.com:user @@ /home/user/projects",
			wantSite: SiteConfig{
				RemotePrefix: "git@github.com:user",
				Dir:          "/home/user/projects",
				Repos: []Repo{
					{Repo: "myrepo", Memo: "unmerged"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotSite := buildSiteConfig(tt.result, tt.baseDir)

			if gotKey != tt.wantKey {
				t.Errorf("buildSiteConfig() key = %q, want %q", gotKey, tt.wantKey)
			}

			if gotSite.RemotePrefix != tt.wantSite.RemotePrefix {
				t.Errorf("buildSiteConfig() RemotePrefix = %q, want %q", gotSite.RemotePrefix, tt.wantSite.RemotePrefix)
			}

			if gotSite.Dir != tt.wantSite.Dir {
				t.Errorf("buildSiteConfig() Dir = %q, want %q", gotSite.Dir, tt.wantSite.Dir)
			}

			if len(gotSite.Repos) != len(tt.wantSite.Repos) {
				t.Errorf("buildSiteConfig() Repos length = %d, want %d", len(gotSite.Repos), len(tt.wantSite.Repos))
			}

			if len(gotSite.Repos) > 0 && len(tt.wantSite.Repos) > 0 {
				if gotSite.Repos[0].Repo != tt.wantSite.Repos[0].Repo {
					t.Errorf("buildSiteConfig() Repos[0].Repo = %q, want %q", gotSite.Repos[0].Repo, tt.wantSite.Repos[0].Repo)
				}
				if gotSite.Repos[0].Memo != tt.wantSite.Repos[0].Memo {
					t.Errorf("buildSiteConfig() Repos[0].Memo = %q, want %q", gotSite.Repos[0].Memo, tt.wantSite.Repos[0].Memo)
				}
			}
		})
	}
}

func TestDiscoverGitRepo_NonGitDirectory(t *testing.T) {
	// Create a temporary directory that is not a Git repository
	tempDir := t.TempDir()

	result, err := discoverGitRepo(tempDir)

	if err != nil {
		t.Errorf("Expected no error for non-Git directory, got: %v", err)
	}

	if result != nil {
		t.Error("Expected nil result for non-Git directory")
	}
}

func TestDiscoverGitRepo_GitRepository(t *testing.T) {
	// Create a temporary Git repository
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")

	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	result, err := discoverGitRepo(tempDir)

	if err != nil {
		t.Errorf("Expected no error for Git directory, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result for Git directory")
	}

	if result.Path != tempDir {
		t.Errorf("Expected path %q, got %q", tempDir, result.Path)
	}

	// Should not have origin since we didn't add one
	if result.HasOrigin {
		t.Error("Expected HasOrigin to be false for repo without origin")
	}
}

func TestGetGitRemoteOrigin_Alias(t *testing.T) {
	// Test that getGitRemoteOrigin is properly aliased to getGitOrigin
	tempDir := t.TempDir()

	// This should fail since it's not a git repository
	_, err1 := getGitOrigin(tempDir)
	_, err2 := getGitRemoteOrigin(tempDir)

	// Both should fail in the same way
	if (err1 == nil) != (err2 == nil) {
		t.Error("getGitOrigin and getGitRemoteOrigin should behave identically")
	}
}

func TestHasUncommittedChanges_EdgeCases(t *testing.T) {
	// Test with non-git directory
	tempDir := t.TempDir()
	
	result, err := hasUncommittedChanges(tempDir)
	if err == nil {
		t.Error("Expected error for non-git directory")
	}
	if result {
		t.Error("Expected false for non-git directory")
	}
}

func TestHasUnmergedCommits_EdgeCases(t *testing.T) {
	// Test with non-git directory
	tempDir := t.TempDir()
	
	result, err := hasUnmergedCommits(tempDir)
	if err == nil {
		t.Error("Expected error for non-git directory")
	}
	if result {
		t.Error("Expected false for non-git directory")
	}
}

func TestGetGitOrigin_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(string) error
		expectErr bool
	}{
		{
			name: "Empty directory",
			setupFunc: func(dir string) error {
				return nil // Do nothing
			},
			expectErr: true,
		},
		{
			name: "Git directory with no remotes",
			setupFunc: func(dir string) error {
				gitDir := filepath.Join(dir, ".git")
				if err := os.MkdirAll(gitDir, 0755); err != nil {
					return err
				}
				
				// Create basic git config without remotes
				configContent := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
`
				configPath := filepath.Join(gitDir, "config")
				return os.WriteFile(configPath, []byte(configContent), 0644)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			
			if err := tt.setupFunc(tempDir); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			
			_, err := getGitOrigin(tempDir)
			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestCalculateRelativeLocalPath_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		repoPath string
		repoName string
		baseDir  string
		expected string
	}{
		{
			name:     "Repository path equals base directory",
			repoPath: "/home/user/projects",
			repoName: "projects",
			baseDir:  "/home/user/projects",
			expected: "/home/user/projects/",
		},
		{
			name:     "Very long nested path",
			repoPath: "/home/user/projects/team/subteam/project/deep/nested/repo",
			repoName: "repo",
			baseDir:  "/home/user/projects",
			expected: "/home/user/projects/team/subteam/project/deep/nested/",
		},
		{
			name:     "Relative path that can't be calculated",
			repoPath: "/completely/different/path/repo",
			repoName: "repo",
			baseDir:  "/home/user/projects",
			expected: "/completely/different/path/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateRelativeLocalPath(tt.repoPath, tt.repoName, tt.baseDir)
			if result != tt.expected {
				t.Errorf("calculateRelativeLocalPath(%q, %q, %q) = %q, expected %q",
					tt.repoPath, tt.repoName, tt.baseDir, result, tt.expected)
			}
		})
	}
}

func TestDiscoverGitRepo_PermissionError(t *testing.T) {
	// This test might not work on all systems due to permission restrictions
	// Skip it if we can't create the test condition
	tempDir := t.TempDir()
	restrictedDir := filepath.Join(tempDir, "restricted")
	
	if err := os.MkdirAll(restrictedDir, 0000); err != nil {
		t.Skip("Cannot create restricted directory for test")
	}
	defer os.Chmod(restrictedDir, 0755) // Restore permissions for cleanup
	
	result, err := discoverGitRepo(restrictedDir)
	// Should not crash, should return nil result
	if result != nil {
		t.Error("Expected nil result for inaccessible directory")
	}
	// Error handling may vary by system, so we don't strictly check for error
	_ = err // Ignore error as it may vary by system
}

func TestBuildSiteConfig_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		result   *RepoDiscoveryResult
		baseDir  string
		wantKey  string
		validate func(*testing.T, SiteConfig)
	}{
		{
			name: "Empty origin URL",
			result: &RepoDiscoveryResult{
				Path:      "/home/user/projects/myrepo",
				Origin:    "",
				HasOrigin: true, // Inconsistent state for testing
			},
			baseDir: "/home/user/projects",
			wantKey: " @@ /home/user/projects",
			validate: func(t *testing.T, sc SiteConfig) {
				if sc.RemotePrefix != "" {
					t.Errorf("Expected empty RemotePrefix, got %q", sc.RemotePrefix)
				}
			},
		},
		{
			name: "Very long repository path",
			result: &RepoDiscoveryResult{
				Path:      "/very/long/path/to/deeply/nested/repository/structure/myrepo",
				Origin:    "https://github.com/user/myrepo.git",
				HasOrigin: true,
			},
			baseDir: "/very/long/path",
			wantKey: "https://github.com/user @@ /very/long/path",
			validate: func(t *testing.T, sc SiteConfig) {
				if len(sc.Repos) != 1 {
					t.Errorf("Expected 1 repo, got %d", len(sc.Repos))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotSite := buildSiteConfig(tt.result, tt.baseDir)
			
			if gotKey != tt.wantKey {
				t.Errorf("buildSiteConfig() key = %q, want %q", gotKey, tt.wantKey)
			}
			
			if tt.validate != nil {
				tt.validate(t, gotSite)
			}
		})
	}
}

func TestDiscoverGitRepo_UncommittedChangesError(t *testing.T) {
	// Create a temporary directory with a git repo that will cause status check to fail
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	// Create a corrupted git config that might cause git status to fail
	configContent := `[core]
	repositoryformatversion = 999999
	filemode = true
	bare = false
`
	configPath := filepath.Join(gitDir, "config")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create git config: %v", err)
	}
	
	result, err := discoverGitRepo(tempDir)
	
	// Should not fail even if git status fails
	if err != nil {
		t.Errorf("discoverGitRepo should handle git status errors gracefully: %v", err)
	}
	
	if result == nil {
		t.Error("Expected result even with git status error")
	}
}

func TestDiscoverGitRepo_UnmergedCommitsError(t *testing.T) {
	// Create a temporary directory with a git repo that will cause cherry check to fail
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	// Create a basic git config without upstream
	configContent := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
[remote "origin"]
	url = https://github.com/test/repo.git
	fetch = +refs/heads/*:refs/remotes/origin/*
`
	configPath := filepath.Join(gitDir, "config")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create git config: %v", err)
	}
	
	result, err := discoverGitRepo(tempDir)
	
	// Should not fail even if git cherry fails
	if err != nil {
		t.Errorf("discoverGitRepo should handle git cherry errors gracefully: %v", err)
	}
	
	if result == nil {
		t.Error("Expected result even with git cherry error")
	}
}

func TestHasUnmergedCommits_GitCherryFailure(t *testing.T) {
	// Test with a directory that will cause git cherry to fail
	tempDir := t.TempDir()
	
	// Create a fake git directory without proper setup
	gitDir := filepath.Join(tempDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	result, err := hasUnmergedCommits(tempDir)
	
	// Should return false and no error when git cherry fails
	if err != nil {
		t.Errorf("hasUnmergedCommits should handle git cherry failure gracefully: %v", err)
	}
	
	if result {
		t.Error("Expected false when git cherry fails")
	}
}

func TestGetGitOrigin_CommandFailure(t *testing.T) {
	// Test with a non-git directory
	tempDir := t.TempDir()
	
	_, err := getGitOrigin(tempDir)
	
	if err == nil {
		t.Error("Expected error when getting origin from non-git directory")
	}
}

func TestCalculateRelativeLocalPath_RelativePathError(t *testing.T) {
	// Test with paths that can't be made relative
	repoPath := "/completely/different/filesystem/repo"
	repoName := "repo"
	baseDir := "/home/user/projects"
	
	result := calculateRelativeLocalPath(repoPath, repoName, baseDir)
	
	// Should fallback to parent directory
	expected := "/completely/different/filesystem/"
	if result != expected {
		t.Errorf("Expected fallback path %q, got %q", expected, result)
	}
}
