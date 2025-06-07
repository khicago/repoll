package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsGitRepository(t *testing.T) {
	// 测试非Git目录
	tempDir := t.TempDir()
	if isGitRepository(tempDir) {
		t.Error("Expected false for non-Git directory")
	}
	
	// 测试不存在的目录
	if isGitRepository("/nonexistent/path") {
		t.Error("Expected false for nonexistent directory")
	}
	
	// 测试有.git目录的情况
	gitDir := t.TempDir()
	err := os.Mkdir(filepath.Join(gitDir, ".git"), 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	if !isGitRepository(gitDir) {
		t.Error("Expected true for directory with .git folder")
	}
	
	// 测试.git文件的情况（子模块）
	submoduleDir := t.TempDir()
	gitFile := filepath.Join(submoduleDir, ".git")
	err = os.WriteFile(gitFile, []byte("gitdir: ../.git/modules/submodule"), 0644)
	if err != nil {
		t.Fatalf("Failed to create .git file: %v", err)
	}
	
	if !isGitRepository(submoduleDir) {
		t.Error("Expected true for directory with .git file")
	}
}

func TestGetCurrentBranch(t *testing.T) {
	// 测试非Git目录
	tempDir := t.TempDir()
	branch, err := GetCurrentBranch(tempDir)
	if err == nil {
		t.Error("Expected error for non-Git directory")
	}
	if branch != "" {
		t.Errorf("Expected empty branch name, got: %s", branch)
	}
	
	// 测试不存在的目录
	branch, err = GetCurrentBranch("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}
	if branch != "" {
		t.Errorf("Expected empty branch name, got: %s", branch)
	}
}

func TestGetRemoteURL(t *testing.T) {
	// 测试非Git目录
	tempDir := t.TempDir()
	url, err := GetRemoteURL(tempDir, "origin")
	if err == nil {
		t.Error("Expected error for non-Git directory")
	}
	if url != "" {
		t.Errorf("Expected empty URL, got: %s", url)
	}
	
	// 测试不存在的目录
	url, err = GetRemoteURL("/nonexistent/path", "origin")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}
	if url != "" {
		t.Errorf("Expected empty URL, got: %s", url)
	}
	
	// 测试空remote参数（应该默认为origin）
	url, err = GetRemoteURL(tempDir, "")
	if err == nil {
		t.Error("Expected error for non-Git directory")
	}
}

func TestClone_ErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name      string
		url       string
		targetDir string
		expectErr bool
	}{
		{
			name:      "Empty URL",
			url:       "",
			targetDir: filepath.Join(tempDir, "test1"),
			expectErr: true,
		},
		{
			name:      "Empty target directory",
			url:       "https://github.com/test/repo.git",
			targetDir: "",
			expectErr: true,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Clone(test.url, test.targetDir)
			
			if test.expectErr && err == nil {
				t.Errorf("Expected error for test '%s'", test.name)
			}
		})
	}
}

func TestClone_DirectoryCreation(t *testing.T) {
	tempDir := t.TempDir()
	
	// 测试父目录创建
	deepPath := filepath.Join(tempDir, "deep", "nested", "path", "repo")
	
	// Clone会失败（因为URL无效），但应该创建父目录
	err := Clone("invalid-url", deepPath)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
	
	// 检查父目录是否被创建
	parentDir := filepath.Dir(deepPath)
	if _, statErr := os.Stat(parentDir); os.IsNotExist(statErr) {
		t.Errorf("Parent directory was not created: %s", parentDir)
	}
}

func TestUpdate_ErrorHandling(t *testing.T) {
	// 测试非Git目录
	tempDir := t.TempDir()
	err := Update(tempDir)
	if err == nil {
		t.Error("Expected error for non-Git directory")
	}
	
	expectedMsg := "not a valid Git repository"
	if err != nil && !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message containing '%s', got: %v", expectedMsg, err)
	}
	
	// 测试不存在的目录
	err = Update("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}
}

func TestUpdate_NonExistentDirectory(t *testing.T) {
	err := Update("/non/existent/directory")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestUpdate_EmptyPath(t *testing.T) {
	// 测试空路径
	err := Update("")
	if err == nil {
		t.Error("Expected error for empty path")
	}
	
	if !strings.Contains(err.Error(), "not a valid Git repository") {
		t.Errorf("Expected error to mention 'not a valid Git repository', got: %v", err)
	}
}

func TestUpdate_NotGitRepository(t *testing.T) {
	// 创建临时目录但不是 Git 仓库
	tempDir := t.TempDir()
	
	err := Update(tempDir)
	if err == nil {
		t.Error("Expected error for non-Git directory")
	}
	
	if !strings.Contains(err.Error(), "not a valid Git repository") {
		t.Errorf("Expected error to mention 'not a valid Git repository', got: %v", err)
	}
}

func TestUpdate_MockGitRepository(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	// 创建 .git 目录来模拟 Git 仓库
	gitDir := filepath.Join(tempDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	// 测试 Update 函数（会失败，因为不是真正的 Git 仓库）
	err = Update(tempDir)
	if err != nil {
		t.Logf("Update failed as expected (not a real Git repo): %v", err)
		// 验证错误是来自 Git 命令执行，而不是目录检查
		if strings.Contains(err.Error(), "not a Git repository") {
			t.Errorf("Unexpected error type, should fail at Git command execution: %v", err)
		}
	}
}

func TestUpdate_InvalidGitDirectory(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	// 创建 .git 文件而不是目录（某些情况下 Git 会这样做）
	gitFile := filepath.Join(tempDir, ".git")
	err := os.WriteFile(gitFile, []byte("gitdir: /some/other/path"), 0644)
	if err != nil {
		t.Fatalf("Failed to create .git file: %v", err)
	}
	
	// 测试 Update 函数
	err = Update(tempDir)
	if err != nil {
		t.Logf("Update failed as expected: %v", err)
	}
}

func TestUpdate_ReadOnlyDirectory(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	// 创建 .git 目录
	gitDir := filepath.Join(tempDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	// 将目录设置为只读（在某些系统上可能不起作用）
	err = os.Chmod(tempDir, 0444)
	if err != nil {
		t.Logf("Failed to set directory as read-only: %v", err)
	}
	
	// 恢复权限以便清理
	defer func() {
		os.Chmod(tempDir, 0755)
	}()
	
	// 测试 Update 函数
	err = Update(tempDir)
	if err != nil {
		t.Logf("Update failed as expected: %v", err)
	}
}

func TestUpdate_RelativePath(t *testing.T) {
	// 测试相对路径
	err := Update("./non-existent-relative-path")
	if err == nil {
		t.Error("Expected error for non-existent relative path")
	}
	
	if !strings.Contains(err.Error(), "not a valid Git repository") {
		t.Errorf("Expected error to mention 'not a valid Git repository', got: %v", err)
	}
}

func TestUpdate_CurrentDirectory(t *testing.T) {
	// 测试当前目录（如果不是 Git 仓库）
	err := Update(".")
	if err != nil {
		t.Logf("Update failed as expected for current directory: %v", err)
		// 这个测试的结果取决于当前目录是否是 Git 仓库
		// 我们只验证函数不会 panic
	}
}

func TestUpdate_SymbolicLink(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	// 创建一个符号链接指向不存在的目录
	linkPath := filepath.Join(tempDir, "symlink")
	err := os.Symlink("/non/existent/target", linkPath)
	if err != nil {
		t.Skipf("Failed to create symbolic link (may not be supported): %v", err)
	}
	
	// 测试 Update 函数
	err = Update(linkPath)
	if err == nil {
		t.Error("Expected error for broken symbolic link")
	}
}

func TestUpdate_NestedGitRepository(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	nestedDir := filepath.Join(tempDir, "nested", "repo")
	err := os.MkdirAll(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}
	
	// 在嵌套目录中创建 .git 目录
	gitDir := filepath.Join(nestedDir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	// 测试 Update 函数
	err = Update(nestedDir)
	if err != nil {
		t.Logf("Update failed as expected for nested repo: %v", err)
	}
} 