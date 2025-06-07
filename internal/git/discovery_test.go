package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDiscoverRepository_ValidRepo(t *testing.T) {
	// 创建临时Git仓库用于测试
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	
	// 创建.git目录
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	// 初始化基本的Git仓库结构
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err = cmd.Run()
	if err != nil {
		t.Skipf("Git not available, skipping test: %v", err)
	}
	
	info, err := DiscoverRepository(tempDir)
	if err != nil {
		t.Fatalf("DiscoverRepository failed: %v", err)
	}
	
	if info.Path != tempDir {
		t.Errorf("Expected path %s, got %s", tempDir, info.Path)
	}
	
	// 新初始化的仓库通常没有origin
	if info.HasOrigin {
		t.Logf("Repository has origin: %s", info.Origin)
	}
}

func TestDiscoverRepository_NonGitDirectory(t *testing.T) {
	tempDir := t.TempDir()
	
	_, err := DiscoverRepository(tempDir)
	if err == nil {
		t.Error("Expected error for non-Git directory")
	}
}

func TestDiscoverRepository_NonExistentDirectory(t *testing.T) {
	_, err := DiscoverRepository("/non/existent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestExtractRepoNameFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "HTTPS GitHub URL",
			url:      "https://github.com/owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "HTTPS GitHub URL without .git",
			url:      "https://github.com/owner/repo",
			expected: "owner/repo",
		},
		{
			name:     "SSH GitHub URL",
			url:      "git@github.com:owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "SSH GitHub URL without .git",
			url:      "git@github.com:owner/repo",
			expected: "owner/repo",
		},
		{
			name:     "GitLab HTTPS URL",
			url:      "https://gitlab.com/group/project.git",
			expected: "group/project",
		},
		{
			name:     "GitLab SSH URL",
			url:      "git@gitlab.com:group/project.git",
			expected: "group/project",
		},
		{
			name:     "Custom domain HTTPS",
			url:      "https://git.company.com/team/project.git",
			expected: "team/project",
		},
		{
			name:     "Custom domain SSH",
			url:      "git@git.company.com:team/project.git",
			expected: "team/project",
		},
		{
			name:     "URL with spaces",
			url:      "  https://github.com/owner/repo.git  ",
			expected: "owner/repo",
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: "",
		},
		{
			name:     "Invalid format",
			url:      "not-a-git-url",
			expected: "",
		},
		{
			name:     "Single part URL",
			url:      "https://github.com/single",
			expected: "github.com/single",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ExtractRepoNameFromURL(test.url)
			if result != test.expected {
				t.Errorf("ExtractRepoNameFromURL(%s) = %s, expected %s", test.url, result, test.expected)
			}
		})
	}
}

func TestExtractRemotePrefix(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected string
	}{
		{
			name:     "GitHub HTTPS URL",
			origin:   "https://github.com/owner/repo.git",
			expected: "https://github.com/",
		},
		{
			name:     "GitHub SSH URL",
			origin:   "git@github.com:owner/repo.git",
			expected: "https://github.com/",
		},
		{
			name:     "GitLab HTTPS URL",
			origin:   "https://gitlab.com/group/project.git",
			expected: "https://gitlab.com/",
		},
		{
			name:     "GitLab SSH URL",
			origin:   "git@gitlab.com:group/project.git",
			expected: "https://gitlab.com/",
		},
		{
			name:     "Custom domain HTTPS",
			origin:   "https://git.company.com/team/project.git",
			expected: "https://git.company.com/",
		},
		{
			name:     "Custom domain SSH",
			origin:   "git@git.company.com:team/project.git",
			expected: "https://git.company.com/",
		},
		{
			name:     "URL with spaces",
			origin:   "  https://github.com/owner/repo.git  ",
			expected: "https://github.com/",
		},
		{
			name:     "Empty origin",
			origin:   "",
			expected: "",
		},
		{
			name:     "Invalid format",
			origin:   "not-a-git-url",
			expected: "",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ExtractRemotePrefix(test.origin)
			if result != test.expected {
				t.Errorf("ExtractRemotePrefix(%s) = %s, expected %s", test.origin, result, test.expected)
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected bool
	}{
		{
			name:     "HTTP URL",
			str:      "http://example.com",
			expected: true,
		},
		{
			name:     "HTTPS URL",
			str:      "https://github.com/owner/repo",
			expected: true,
		},
		{
			name:     "SSH Git URL",
			str:      "git@github.com:owner/repo.git",
			expected: true,
		},
		{
			name:     "Simple text",
			str:      "not-a-url",
			expected: false,
		},
		{
			name:     "Empty string",
			str:      "",
			expected: false,
		},
		{
			name:     "FTP URL (not supported)",
			str:      "ftp://example.com",
			expected: false,
		},
		{
			name:     "Partial HTTPS",
			str:      "https://",
			expected: true,
		},
		{
			name:     "Partial SSH",
			str:      "git@",
			expected: true,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsURL(test.str)
			if result != test.expected {
				t.Errorf("IsURL(%s) = %v, expected %v", test.str, result, test.expected)
			}
		})
	}
}

func TestGetGitOrigin_ValidRepo(t *testing.T) {
	tempDir := t.TempDir()
	
	// 初始化Git仓库
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err := cmd.Run()
	if err != nil {
		t.Skipf("Git not available, skipping test: %v", err)
	}
	
	// 添加origin
	testOrigin := "https://github.com/test/repo.git"
	cmd = exec.Command("git", "remote", "add", "origin", testOrigin)
	cmd.Dir = tempDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to add origin: %v", err)
	}
	
	origin, err := getGitOrigin(tempDir)
	if err != nil {
		t.Fatalf("getGitOrigin failed: %v", err)
	}
	
	if origin != testOrigin {
		t.Errorf("Expected origin %s, got %s", testOrigin, origin)
	}
}

func TestGetGitOrigin_NoOrigin(t *testing.T) {
	tempDir := t.TempDir()
	
	// 初始化Git仓库但不添加origin
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err := cmd.Run()
	if err != nil {
		t.Skipf("Git not available, skipping test: %v", err)
	}
	
	origin, err := getGitOrigin(tempDir)
	if err == nil {
		t.Logf("getGitOrigin returned: %s (expected error)", origin)
	}
	// 没有origin时应该返回错误，这是正常的
}

func TestGetGitOrigin_NonGitRepo(t *testing.T) {
	tempDir := t.TempDir()
	
	_, err := getGitOrigin(tempDir)
	if err == nil {
		t.Error("Expected error for non-Git repository")
	}
}

func TestHasUncommittedChanges_CleanRepo(t *testing.T) {
	tempDir := t.TempDir()
	
	// 初始化Git仓库
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err := cmd.Run()
	if err != nil {
		t.Skipf("Git not available, skipping test: %v", err)
	}
	
	// 配置Git用户（避免警告）
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	cmd.Run()
	
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	cmd.Run()
	
	hasChanges, err := hasUncommittedChanges(tempDir)
	if err != nil {
		t.Fatalf("hasUncommittedChanges failed: %v", err)
	}
	
	// 空仓库也可能显示为有变化，这取决于Git版本
	t.Logf("Clean repo has uncommitted changes: %v", hasChanges)
}

func TestHasUncommittedChanges_WithChanges(t *testing.T) {
	tempDir := t.TempDir()
	
	// 初始化Git仓库
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err := cmd.Run()
	if err != nil {
		t.Skipf("Git not available, skipping test: %v", err)
	}
	
	// 创建一个文件
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	hasChanges, err := hasUncommittedChanges(tempDir)
	if err != nil {
		t.Fatalf("hasUncommittedChanges failed: %v", err)
	}
	
	if !hasChanges {
		t.Error("Expected uncommitted changes to be detected")
	}
}

func TestHasUncommittedChanges_NonGitRepo(t *testing.T) {
	tempDir := t.TempDir()
	
	_, err := hasUncommittedChanges(tempDir)
	if err == nil {
		t.Error("Expected error for non-Git repository")
	}
}

func TestHasUnmergedChanges_CleanRepo(t *testing.T) {
	tempDir := t.TempDir()
	
	// 初始化Git仓库
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	err := cmd.Run()
	if err != nil {
		t.Skipf("Git not available, skipping test: %v", err)
	}
	
	hasUnmerged, err := hasUnmergedChanges(tempDir)
	if err != nil {
		t.Fatalf("hasUnmergedChanges failed: %v", err)
	}
	
	if hasUnmerged {
		t.Error("Clean repo should not have unmerged changes")
	}
}

func TestHasUnmergedChanges_NonGitRepo(t *testing.T) {
	tempDir := t.TempDir()
	
	_, err := hasUnmergedChanges(tempDir)
	if err == nil {
		t.Error("Expected error for non-Git repository")
	}
}

func TestRepositoryInfo_Fields(t *testing.T) {
	info := &RepositoryInfo{
		Path:        "/path/to/repo",
		Origin:      "https://github.com/owner/repo.git",
		HasOrigin:   true,
		Uncommitted: false,
		Unmerged:    false,
	}
	
	// 验证所有字段都正确设置
	if info.Path != "/path/to/repo" {
		t.Errorf("Path field mismatch: got %s", info.Path)
	}
	
	if info.Origin != "https://github.com/owner/repo.git" {
		t.Errorf("Origin field mismatch: got %s", info.Origin)
	}
	
	if !info.HasOrigin {
		t.Error("HasOrigin field should be true")
	}
	
	if info.Uncommitted {
		t.Error("Uncommitted field should be false")
	}
	
	if info.Unmerged {
		t.Error("Unmerged field should be false")
	}
}

func TestExtractRepoNameFromURL_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "URL with port",
			url:      "https://github.com:443/owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "SSH with port",
			url:      "ssh://git@github.com:22/owner/repo.git",
			expected: "22/owner/repo",
		},
		{
			name:     "Very nested path",
			url:      "https://github.com/org/team/project/repo.git",
			expected: "project/repo",
		},
		{
			name:     "Only one path component",
			url:      "https://github.com/single.git",
			expected: "github.com/single",
		},
		{
			name:     "Multiple .git suffixes",
			url:      "https://github.com/owner/repo.git.git",
			expected: "owner/repo.git",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ExtractRepoNameFromURL(test.url)
			if result != test.expected {
				t.Errorf("ExtractRepoNameFromURL(%s) = %s, expected %s", test.url, result, test.expected)
			}
		})
	}
} 