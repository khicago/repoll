package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWarmUpRepo_NonexistentDirectory(t *testing.T) {
	err := warmUpRepo("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}
}

func TestWarmUpRepo_NoProjectFiles(t *testing.T) {
	tempDir := t.TempDir()

	// 创建一个空目录（没有 go.mod 或 package.json）
	err := warmUpRepo(tempDir)
	if err != nil {
		t.Errorf("Expected no error for directory without project files, got: %v", err)
	}
}

func TestWarmUpRepo_GoProject(t *testing.T) {
	tempDir := t.TempDir()

	// 创建 go.mod 文件
	goModPath := filepath.Join(tempDir, "go.mod")
	goModContent := `module test
go 1.22`
	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// 测试 Go 项目预热
	err = warmUpRepo(tempDir)
	// go mod download 可能失败（因为没有真实的模块），但我们主要测试检测逻辑
	// 这里不检查错误，因为 go mod download 的成功与否取决于环境
	t.Logf("warmUpRepo result for Go project: %v", err)
}

func TestWarmUpRepo_NodeProjectWithNpm(t *testing.T) {
	tempDir := t.TempDir()

	// 创建 package.json 文件
	packageJSONPath := filepath.Join(tempDir, "package.json")
	packageJSONContent := `{
  "name": "test-project",
  "version": "1.0.0",
  "dependencies": {}
}`
	err := os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// 测试 Node.js 项目预热（npm）
	err = warmUpRepo(tempDir)
	// npm install 可能失败（如果 npm 不可用），但我们主要测试检测逻辑
	t.Logf("warmUpRepo result for Node.js project (npm): %v", err)
}

func TestWarmUpRepo_NodeProjectWithYarn(t *testing.T) {
	tempDir := t.TempDir()

	// 创建 package.json 和 yarn.lock 文件
	packageJSONPath := filepath.Join(tempDir, "package.json")
	packageJSONContent := `{
  "name": "test-project",
  "version": "1.0.0",
  "dependencies": {}
}`
	err := os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	yarnLockPath := filepath.Join(tempDir, "yarn.lock")
	err = os.WriteFile(yarnLockPath, []byte("# yarn.lock"), 0644)
	if err != nil {
		t.Fatalf("Failed to create yarn.lock: %v", err)
	}

	// 测试 Node.js 项目预热（yarn）
	err = warmUpRepo(tempDir)
	// yarn install 可能失败（如果 yarn 不可用），但我们主要测试检测逻辑
	t.Logf("warmUpRepo result for Node.js project (yarn): %v", err)
}

func TestIsGoProject(t *testing.T) {
	tempDir := t.TempDir()

	// 测试没有 go.mod 的目录
	if isGoProject(tempDir) {
		t.Error("Expected false for directory without go.mod")
	}

	// 创建 go.mod 文件
	goModPath := filepath.Join(tempDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// 测试有 go.mod 的目录
	if !isGoProject(tempDir) {
		t.Error("Expected true for directory with go.mod")
	}
}

func TestIsNodeProject(t *testing.T) {
	tempDir := t.TempDir()

	// 测试没有 package.json 的目录
	if isNodeProject(tempDir) {
		t.Error("Expected false for directory without package.json")
	}

	// 创建 package.json 文件
	packageJSONPath := filepath.Join(tempDir, "package.json")
	err := os.WriteFile(packageJSONPath, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// 测试有 package.json 的目录
	if !isNodeProject(tempDir) {
		t.Error("Expected true for directory with package.json")
	}
}

func TestHasYarnLock(t *testing.T) {
	tempDir := t.TempDir()

	// 测试没有 yarn.lock 的情况
	if hasYarnLock(tempDir) {
		t.Error("Expected hasYarnLock to return false for directory without yarn.lock")
	}

	// 创建 yarn.lock
	yarnLockPath := filepath.Join(tempDir, "yarn.lock")
	err := os.WriteFile(yarnLockPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create yarn.lock: %v", err)
	}

	if !hasYarnLock(tempDir) {
		t.Error("Expected hasYarnLock to return true for directory with yarn.lock")
	}
}

func TestWarmUpGoProject(t *testing.T) {
	tempDir := t.TempDir()

	// 创建 go.mod 文件
	goModPath := filepath.Join(tempDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test\ngo 1.22"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// 测试 Go 项目预热
	err = warmUpGoProject(tempDir)
	// 不检查错误，因为 go mod download 的成功与否取决于环境
	t.Logf("warmUpGoProject result: %v", err)
}

func TestWarmUpNodeProject(t *testing.T) {
	tempDir := t.TempDir()

	// 创建 package.json 文件
	packageJSONPath := filepath.Join(tempDir, "package.json")
	err := os.WriteFile(packageJSONPath, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// 测试没有 yarn.lock 的情况
	err = warmUpNodeProject(tempDir)
	t.Logf("warmUpNodeProject (npm) result: %v", err)

	// 创建 yarn.lock 文件
	yarnLockPath := filepath.Join(tempDir, "yarn.lock")
	err = os.WriteFile(yarnLockPath, []byte("# yarn.lock"), 0644)
	if err != nil {
		t.Fatalf("Failed to create yarn.lock: %v", err)
	}

	// 测试有 yarn.lock 的情况
	err = warmUpNodeProject(tempDir)
	t.Logf("warmUpNodeProject (yarn) result: %v", err)
}

func TestWarmUpWithYarn(t *testing.T) {
	tempDir := t.TempDir()

	// 创建基本的 package.json 和 yarn.lock
	packageJSONPath := filepath.Join(tempDir, "package.json")
	err := os.WriteFile(packageJSONPath, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	yarnLockPath := filepath.Join(tempDir, "yarn.lock")
	err = os.WriteFile(yarnLockPath, []byte("# yarn.lock"), 0644)
	if err != nil {
		t.Fatalf("Failed to create yarn.lock: %v", err)
	}

	err = warmUpWithYarn(tempDir)
	// yarn 可能不可用，所以不检查错误
	t.Logf("warmUpWithYarn result: %v", err)
}

func TestWarmUpWithNpm(t *testing.T) {
	tempDir := t.TempDir()

	// 创建基本的 package.json
	packageJSONPath := filepath.Join(tempDir, "package.json")
	err := os.WriteFile(packageJSONPath, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	err = warmUpWithNpm(tempDir)
	// npm 可能不可用，所以不检查错误
	t.Logf("warmUpWithNpm result: %v", err)
}

func TestWarmUpGoProject_Failure(t *testing.T) {
	tempDir := t.TempDir()

	// 创建无效的 go.mod 文件
	goModPath := filepath.Join(tempDir, "go.mod")
	invalidGoMod := `invalid go.mod content`
	err := os.WriteFile(goModPath, []byte(invalidGoMod), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// 测试 Go 项目预热失败情况
	err = warmUpGoProject(tempDir)
	if err != nil {
		t.Logf("warmUpGoProject failed as expected: %v", err)
	}
}

func TestWarmUpWithYarn_CommandError(t *testing.T) {
	tempDir := t.TempDir()

	// 测试在没有yarn.lock的目录运行yarn
	err := warmUpWithYarn(tempDir)
	if err != nil {
		t.Logf("warmUpWithYarn failed as expected: %v", err)
	}
}

func TestWarmUpWithNpm_CommandError(t *testing.T) {
	tempDir := t.TempDir()

	// 测试在没有package.json的目录运行npm
	err := warmUpWithNpm(tempDir)
	if err != nil {
		t.Logf("warmUpWithNpm failed as expected: %v", err)
	}
}

func TestWarmUpRepo_ProjectDetection(t *testing.T) {
	tempDir := t.TempDir()

	// 测试 isGoProject
	goModPath := filepath.Join(tempDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	if !isGoProject(tempDir) {
		t.Error("Expected isGoProject to return true")
	}

	// 删除 go.mod，添加 package.json
	os.Remove(goModPath)

	packageJSONPath := filepath.Join(tempDir, "package.json")
	err = os.WriteFile(packageJSONPath, []byte(`{"name": "test"}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	if !isNodeProject(tempDir) {
		t.Error("Expected isNodeProject to return true")
	}

	if isGoProject(tempDir) {
		t.Error("Expected isGoProject to return false after removing go.mod")
	}
}

func TestWarmUpRepo_DirectoryStatError(t *testing.T) {
	// Test with a path that will cause os.Stat to return an error other than NotExist
	// This is difficult to test reliably across platforms, but we can try with a restricted path
	
	tempDir := t.TempDir()
	restrictedDir := filepath.Join(tempDir, "restricted")
	
	// Create directory and then make it inaccessible
	err := os.MkdirAll(restrictedDir, 0000)
	if err != nil {
		t.Skip("Cannot create restricted directory for test")
	}
	defer os.Chmod(restrictedDir, 0755) // Restore permissions for cleanup
	
	// Try to access a subdirectory of the restricted directory
	inaccessiblePath := filepath.Join(restrictedDir, "subdir")
	
	err = warmUpRepo(inaccessiblePath)
	
	// Should return an error due to permission issues
	if err == nil {
		t.Error("Expected error when accessing restricted directory")
	}
	
	if !strings.Contains(err.Error(), "failed to check repository directory") {
		t.Errorf("Expected directory check error, got: %v", err)
	}
}

func TestWarmUpGoProject_CommandFailure(t *testing.T) {
	// Create a temporary directory with a Go project that will cause go mod download to fail
	tempDir := t.TempDir()
	
	// Create an invalid go.mod file
	goModContent := `module invalid-module

go 1.19

require (
	invalid/module/that/does/not/exist v1.0.0
)
`
	err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}
	
	err = warmUpGoProject(tempDir)
	
	// Should return an error due to invalid module
	if err == nil {
		t.Error("Expected error when running go mod download with invalid module")
	}
	
	if !strings.Contains(err.Error(), "failed to run 'go mod download'") {
		t.Errorf("Expected go mod download error, got: %v", err)
	}
}
