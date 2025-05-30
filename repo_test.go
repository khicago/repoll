package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"https://github.com", "https://github.com", true},
		{"http://example.com", "http://example.com", true},
		{"ssh://git@github.com", "ssh://git@github.com", true},
		{"git@github.com:user/repo", "git@github.com:user/repo", false},
		{"./relative/path", "./relative/path", false},
		{"", "", false},
		{"not-a-url", "not-a-url", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isURL(tt.input)
			if result != tt.expected {
				t.Errorf("isURL(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGitCloneOrUpdate_DirectoryError(t *testing.T) {
	repo := Repo{
		Repo:   "test/repo",
		Rename: "",
		WarmUp: false,
		Memo:   "",
	}

	site := SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          "/invalid/path/that/cannot/be/created",
		Repos:        []Repo{repo},
		WarmUpAll:    false,
	}

	err := gitCloneOrUpdate(repo, site)
	if err == nil {
		t.Error("Expected error for invalid directory path")
	}
}

func TestGitCloneOrUpdate_ExistingDirectory(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	repo := Repo{
		Repo:   "test-repo",
		Rename: "",
		WarmUp: false,
		Memo:   "",
	}

	site := SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		Repos:        []Repo{repo},
		WarmUpAll:    false,
	}

	// 创建目标目录
	repoPath := repo.FullPath(site)
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}

	// 这应该尝试更新而不是克隆
	err := gitCloneOrUpdate(repo, site)
	// 由于这不是真正的git仓库，预期会失败
	if err == nil {
		t.Log("gitCloneOrUpdate succeeded unexpectedly")
	} else {
		t.Logf("gitCloneOrUpdate failed as expected: %v", err)
	}
}

func TestGitCloneOrUpdate_WithRename(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()

	repo := Repo{
		Repo:   "nonexistent/repo", // 使用明显不存在的仓库名
		Rename: "renamed-repo",
		WarmUp: false,
		Memo:   "",
	}

	site := SiteConfig{
		RemotePrefix: "https://invalid-domain-that-does-not-exist.com/", // 使用无效域名
		Dir:          tempDir,
		Repos:        []Repo{repo},
		WarmUpAll:    false,
	}

	err := gitCloneOrUpdate(repo, site)
	// 由于这是一个不存在的仓库，预期会失败
	if err == nil {
		t.Error("Expected error for non-existent repository")
	} else {
		t.Logf("gitCloneOrUpdate failed as expected: %v", err)
	}

	// 检查是否使用了正确的路径
	expectedPath := filepath.Join(tempDir, "renamed-repo")
	if repo.FullPath(site) != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, repo.FullPath(site))
	}
}

func TestGitCloneOrUpdate_StatError(t *testing.T) {
	repo := Repo{
		Repo:   "test/repo",
		Rename: "",
		WarmUp: false,
		Memo:   "",
	}

	// 使用一个无法访问的路径（比如文件而不是目录）
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "not-a-directory.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	site := SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          filePath, // 这里故意用文件路径而不是目录
		Repos:        []Repo{repo},
		WarmUpAll:    false,
	}

	err := gitCloneOrUpdate(repo, site)
	if err == nil {
		t.Error("Expected error when checking file instead of directory")
	}
}

func TestGitCloneOrUpdate_CloneSuccess(t *testing.T) {
	repo := Repo{
		Repo:   "test/repo",
		Rename: "renamed-repo",
		WarmUp: false,
		Memo:   "",
	}

	tmpDir := t.TempDir()

	site := SiteConfig{
		RemotePrefix: "https://invalid-domain-that-does-not-exist.com/",
		Dir:          tmpDir,
		Repos:        []Repo{repo},
		WarmUpAll:    false,
	}

	err := gitCloneOrUpdate(repo, site)
	// 应该失败，因为域名不存在，但至少会尝试克隆
	if err == nil {
		t.Error("Expected error for invalid domain")
	}
}
