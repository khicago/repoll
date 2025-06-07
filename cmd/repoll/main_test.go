package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunProcessConfigs_EmptyConfigs(t *testing.T) {
	err := runProcessConfigs([]string{}, false)
	if err != nil {
		t.Errorf("runProcessConfigs with empty configs should not fail: %v", err)
	}
}

func TestRunProcessConfigs_NonExistentFile(t *testing.T) {
	// 测试不存在的配置文件
	err := runProcessConfigs([]string{"non-existent-file.toml"}, false)
	// 函数不应该返回错误，但会输出错误信息
	if err != nil {
		t.Errorf("runProcessConfigs should handle non-existent files gracefully: %v", err)
	}
}

func TestRunProcessConfigs_ValidConfig(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test.toml")
	
	configContent := `[[sites]]
remote_prefix = "https://github.com/"
directory = "` + tempDir + `/repos/"

[[sites.repos]]
repo = "user/test-repo"
warm_up = false
memo = "Test repository"
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	
	// 创建目标目录
	err = os.MkdirAll(filepath.Join(tempDir, "repos"), 0755)
	if err != nil {
		t.Fatalf("Failed to create repos directory: %v", err)
	}
	
	err = runProcessConfigs([]string{configFile}, false)
	if err != nil {
		t.Errorf("runProcessConfigs failed: %v", err)
	}
}

func TestRunProcessConfigs_WithReport(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-report.toml")
	
	configContent := `[[sites]]
remote_prefix = "https://github.com/"
directory = "` + tempDir + `/repos/"

[[sites.repos]]
repo = "user/test-repo"
warm_up = false
memo = "Test for report"
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	
	// 创建目标目录
	err = os.MkdirAll(filepath.Join(tempDir, "repos"), 0755)
	if err != nil {
		t.Fatalf("Failed to create repos directory: %v", err)
	}
	
	err = runProcessConfigs([]string{configFile}, true)
	if err != nil {
		t.Errorf("runProcessConfigs with report failed: %v", err)
	}
}

func TestRunProcessConfigs_MultipleConfigs(t *testing.T) {
	tempDir := t.TempDir()

	// 创建两个轻量级配置文件，避免网络调用
	configContent1 := `[[sites]]
remote_prefix = "https://github.com/"
directory = "` + tempDir + `/repos1/"
`

	configContent2 := `[[sites]]
remote_prefix = "https://gitlab.com/"
directory = "` + tempDir + `/repos2/"
`

	// 创建配置文件
	configFile1 := filepath.Join(tempDir, "config1.toml")
	configFile2 := filepath.Join(tempDir, "config2.toml")

	err := os.WriteFile(configFile1, []byte(configContent1), 0644)
	if err != nil {
		t.Fatalf("Failed to create config1.toml: %v", err)
	}

	err = os.WriteFile(configFile2, []byte(configContent2), 0644)
	if err != nil {
		t.Fatalf("Failed to create config2.toml: %v", err)
	}

	// 创建必要的目录
	err = os.MkdirAll(filepath.Join(tempDir, "repos1"), 0755)
	if err != nil {
		t.Fatalf("Failed to create repos1 directory: %v", err)
	}

	err = os.MkdirAll(filepath.Join(tempDir, "repos2"), 0755)
	if err != nil {
		t.Fatalf("Failed to create repos2 directory: %v", err)
	}

	// 测试处理多个配置文件
	err = runProcessConfigs([]string{configFile1, configFile2}, false)
	if err != nil {
		t.Fatalf("runProcessConfigs failed: %v", err)
	}
}

func TestRunProcessConfigs_InvalidConfig(t *testing.T) {
	tempDir := t.TempDir()
	
	// 创建无效的配置文件
	invalidContent := `invalid toml content [[[`
	configFile := filepath.Join(tempDir, "invalid.toml")
	err := os.WriteFile(configFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid config file: %v", err)
	}
	
	// 应该继续处理，不会返回错误（只是打印错误）
	err = runProcessConfigs([]string{configFile}, false)
	if err != nil {
		t.Fatalf("runProcessConfigs should not fail for invalid config: %v", err)
	}
}

func TestRunMakeConfig_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// 在空目录中运行应该返回错误
	err := runMakeConfig(tempDir, false)
	if err == nil {
		t.Error("Expected error when running runMakeConfig on empty directory without Git repositories")
	}
}

func TestRunMakeConfig_NonExistentDirectory(t *testing.T) {
	err := runMakeConfig("/non/existent/directory", false)
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestRunMakeConfig_WithGitRepo(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available, skipping test")
	}
	
	tempDir := t.TempDir()
	
	// 创建一个Git仓库
	repoDir := filepath.Join(tempDir, "test-repo")
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}
	
	// 初始化Git仓库
	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	// 添加origin
	cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/test/repo.git")
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to add origin: %v", err)
	}
	
	// 保存当前目录并切换到临时目录
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	// 测试runMakeConfig
	err = runMakeConfig(".", false)
	if err != nil {
		t.Fatalf("runMakeConfig failed: %v", err)
	}
	
	// 验证配置文件是否创建
	files, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}
	
	configFound := false
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "repoll-config-") && strings.HasSuffix(file.Name(), ".toml") {
			configFound = true
			break
		}
	}
	
	if !configFound {
		t.Error("Configuration file was not created")
	}
}

func TestRunMakeConfig_WithReport(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available, skipping test")
	}
	
	tempDir := t.TempDir()
	
	// 创建一个Git仓库
	repoDir := filepath.Join(tempDir, "test-repo")
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}
	
	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/test/repo.git")
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to add origin: %v", err)
	}
	
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	// 测试带报告的runMakeConfig
	err = runMakeConfig(".", true)
	if err != nil {
		t.Fatalf("runMakeConfig with report failed: %v", err)
	}
}

func TestGlobalVariables(t *testing.T) {
	// 验证全局变量不为空
	if version == "" {
		t.Error("version should not be empty")
	}
	if commit == "" {
		t.Error("commit should not be empty") 
	}
	if date == "" {
		t.Error("date should not be empty")
	}
}

func TestVersionString(t *testing.T) {
	// 测试版本字符串格式
	expectedFormat := fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
	
	// 通过创建一个临时的rootCmd来测试版本字符串
	rootCmd := &cobra.Command{
		Use:     "repoll",
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}
	
	if rootCmd.Version != expectedFormat {
		t.Errorf("Version string format incorrect. Expected: %s, Got: %s", expectedFormat, rootCmd.Version)
	}
	
	// 检查版本字符串是否包含预期的组件
	if !strings.Contains(rootCmd.Version, version) {
		t.Error("Version string should contain version")
	}
	if !strings.Contains(rootCmd.Version, commit) {
		t.Error("Version string should contain commit")
	}
	if !strings.Contains(rootCmd.Version, date) {
		t.Error("Version string should contain date")
	}
}

// 辅助函数：检查Git是否可用
func isGitAvailable() bool {
	cmd := exec.Command("git", "version")
	err := cmd.Run()
	return err == nil
} 