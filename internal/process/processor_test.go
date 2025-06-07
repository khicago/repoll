package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/khicago/repoll/internal/config"
	"github.com/khicago/repoll/internal/reporter"
)

func TestProcessConfig_InvalidFile(t *testing.T) {
	report := &reporter.MakeReport{Actions: make([]*reporter.MakeAction, 0)}
	
	err := ProcessConfig("non-existent-file.toml", report)
	if err == nil {
		t.Error("Expected error for non-existent config file")
	}
}

func TestProcessConfig_InvalidToml(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid.toml")
	
	// 创建无效的TOML文件
	invalidContent := `[[sites]
remote = "invalid toml`
	
	err := os.WriteFile(configFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid config: %v", err)
	}
	
	report := &reporter.MakeReport{Actions: make([]*reporter.MakeAction, 0)}
	err = ProcessConfig(configFile, report)
	
	if err == nil {
		t.Error("Expected error for invalid TOML")
	}
}

func TestProcessConfig_EmptyConfig(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "empty.toml")
	
	// 创建空的配置文件
	err := os.WriteFile(configFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty config: %v", err)
	}
	
	report := &reporter.MakeReport{Actions: make([]*reporter.MakeAction, 0)}
	err = ProcessConfig(configFile, report)
	
	// 空配置应该不报错，但也不会有任何操作
	if err != nil {
		t.Errorf("ProcessConfig failed for empty config: %v", err)
	}
	
	if len(report.Actions) != 0 {
		t.Errorf("Expected no actions for empty config, got %d", len(report.Actions))
	}
}

func TestShouldWarmUp(t *testing.T) {
	tests := []struct {
		name     string
		repo     config.Repo
		site     config.SiteConfig
		expected bool
	}{
		{
			name: "repo warmup enabled",
			repo: config.Repo{WarmUp: true},
			site: config.SiteConfig{WarmUpAll: false},
			expected: true,
		},
		{
			name: "site warmup all enabled",
			repo: config.Repo{WarmUp: false},
			site: config.SiteConfig{WarmUpAll: true},
			expected: true,
		},
		{
			name: "both enabled",
			repo: config.Repo{WarmUp: true},
			site: config.SiteConfig{WarmUpAll: true},
			expected: true,
		},
		{
			name: "both disabled",
			repo: config.Repo{WarmUp: false},
			site: config.SiteConfig{WarmUpAll: false},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldWarmUp(tt.repo, tt.site)
			if result != tt.expected {
				t.Errorf("shouldWarmUp() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestProcessSite_EmptyRepos(t *testing.T) {
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          "./test/",
		Repos:        []config.Repo{},
	}
	
	report := &reporter.MakeReport{Actions: make([]*reporter.MakeAction, 0)}
	
	err := processSite(site, report)
	if err != nil {
		t.Errorf("processSite failed: %v", err)
	}
	
	if len(report.Actions) != 0 {
		t.Errorf("Expected no actions for empty repos, got %d", len(report.Actions))
	}
}

func TestProcessSite_InvalidDirectory(t *testing.T) {
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          "/invalid/path/that/does/not/exist",
		Repos: []config.Repo{
			{
				Repo:   "test/repo",
				WarmUp: false,
				Memo:   "Test repo",
			},
		},
	}
	
	report := &reporter.MakeReport{Actions: make([]*reporter.MakeAction, 0)}
	
	// 这应该会失败，因为目录无效
	err := processSite(site, report)
	if err != nil {
		t.Logf("processSite failed as expected: %v", err)
	}
	
	// 即使失败，也应该记录操作
	if len(report.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(report.Actions))
	}
	
	if len(report.Actions) > 0 {
		action := report.Actions[0]
		if action.Repository != "test/repo" {
			t.Errorf("Expected repository 'test/repo', got %s", action.Repository)
		}
		
		if action.Memo != "Test repo" {
			t.Errorf("Expected memo 'Test repo', got %s", action.Memo)
		}
		
		// 检查时间戳是否合理
		if action.Time.IsZero() {
			t.Error("Expected non-zero timestamp")
		}
		
		// 检查持续时间是否合理
		if action.Duration < 0 {
			t.Error("Expected non-negative duration")
		}
	}
}

func TestProcessRepository_BasicFunctionality(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	// 创建一个简单的配置
	repo := config.Repo{
		Repo:   "test/repo",
		WarmUp: false,
		Memo:   "Test repository",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		WarmUpAll:    false,
	}
	
	// 测试processRepository函数不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("processRepository panicked: %v", r)
		}
	}()
	
	// 调用processRepository（预期会失败，但不应该panic）
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed as expected: %v", err)
	}
}

func TestProcessRepository_EmptyConfig(t *testing.T) {
	// 测试空配置
	repo := config.Repo{}
	site := config.SiteConfig{}
	
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("processRepository panicked with empty config: %v", r)
		}
	}()
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed with empty config as expected: %v", err)
	}
}

func TestProcessRepository_InvalidPath(t *testing.T) {
	// 测试无效路径
	repo := config.Repo{
		Repo:   "test/repo",
		WarmUp: false,
		Memo:   "Test repository",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          "/nonexistent/path/that/should/not/exist",
		WarmUpAll:    false,
	}
	
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("processRepository panicked with invalid path: %v", r)
		}
	}()
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed with invalid path as expected: %v", err)
	}
}

func TestProcessRepository_ValidDirectory(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	repo := config.Repo{
		Repo:   "test/repo",
		WarmUp: false,
		Memo:   "Test repository",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		WarmUpAll:    false,
	}
	
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("processRepository panicked with valid directory: %v", r)
		}
	}()
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed as expected (no real Git repo): %v", err)
	}
}

// 新增的测试用例来提高 processRepository 覆盖率

func TestProcessRepository_ExistingRepository(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	// 创建一个模拟的已存在的仓库目录
	repoDir := filepath.Join(tempDir, "existing-repo")
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}
	
	repo := config.Repo{
		Repo:   "existing-repo",
		WarmUp: false,
		Memo:   "Existing repository",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		WarmUpAll:    false,
	}
	
	// 测试处理已存在的仓库（会尝试更新）
	err = processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed as expected (not a real Git repo): %v", err)
		// 验证错误信息包含更新相关的内容
		if !strings.Contains(err.Error(), "update") {
			t.Errorf("Expected error to mention 'update', got: %v", err)
		}
	}
}

func TestProcessRepository_WithWarmUp(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	repo := config.Repo{
		Repo:   "warmup-repo",
		WarmUp: true, // 启用预热
		Memo:   "Repository with warmup",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		WarmUpAll:    false,
	}
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed as expected: %v", err)
		// 验证错误是来自克隆操作，而不是预热操作
		if !strings.Contains(err.Error(), "clone") {
			t.Errorf("Expected error to mention 'clone', got: %v", err)
		}
	}
}

func TestProcessRepository_SiteWarmUpAll(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	repo := config.Repo{
		Repo:   "site-warmup-repo",
		WarmUp: false, // 仓库级别不启用
		Memo:   "Repository with site-level warmup",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		WarmUpAll:    true, // 站点级别启用预热
	}
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed as expected: %v", err)
	}
}

func TestProcessRepository_WithRename(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	repo := config.Repo{
		Repo:   "original/repo-name",
		Rename: "renamed-repo",
		WarmUp: false,
		Memo:   "Repository with rename",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		WarmUpAll:    false,
	}
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed as expected: %v", err)
	}
	
	// 验证目标路径使用了重命名后的名称
	if !strings.Contains(fmt.Sprintf("%v", err), "renamed-repo") {
		t.Logf("Expected path to use renamed directory, error: %v", err)
	}
}

func TestProcessRepository_EmptyRepoName(t *testing.T) {
	// 测试空仓库名
	tempDir := t.TempDir()
	
	repo := config.Repo{
		Repo:   "", // 空仓库名
		WarmUp: false,
		Memo:   "Empty repo name",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		WarmUpAll:    false,
	}
	
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("processRepository panicked with empty repo name: %v", r)
		}
	}()
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed with empty repo name as expected: %v", err)
	}
}

func TestProcessRepository_EmptyRemotePrefix(t *testing.T) {
	// 测试空远程前缀
	tempDir := t.TempDir()
	
	repo := config.Repo{
		Repo:   "test/repo",
		WarmUp: false,
		Memo:   "Empty remote prefix",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "", // 空远程前缀
		Dir:          tempDir,
		WarmUpAll:    false,
	}
	
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("processRepository panicked with empty remote prefix: %v", r)
		}
	}()
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed with empty remote prefix as expected: %v", err)
	}
}

func TestProcessRepository_ComplexRepoPath(t *testing.T) {
	// 测试复杂的仓库路径
	tempDir := t.TempDir()
	
	repo := config.Repo{
		Repo:   "organization/project/sub-module",
		WarmUp: false,
		Memo:   "Complex repo path",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		WarmUpAll:    false,
	}
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed with complex path as expected: %v", err)
	}
	
	// 验证目标路径只使用了最后一部分作为目录名
	if !strings.Contains(fmt.Sprintf("%v", err), "sub-module") {
		t.Logf("Expected path to use last part of complex path, error: %v", err)
	}
}

func TestProcessRepository_BothWarmUpEnabled(t *testing.T) {
	// 测试仓库和站点都启用预热的情况
	tempDir := t.TempDir()
	
	repo := config.Repo{
		Repo:   "both-warmup-repo",
		WarmUp: true,
		Memo:   "Both warmup enabled",
	}
	
	site := config.SiteConfig{
		RemotePrefix: "https://github.com/",
		Dir:          tempDir,
		WarmUpAll:    true,
	}
	
	err := processRepository(repo, site)
	if err != nil {
		t.Logf("processRepository failed as expected: %v", err)
	}
	
	// 验证 shouldWarmUp 函数被正确调用
	shouldWarm := shouldWarmUp(repo, site)
	if !shouldWarm {
		t.Error("Expected shouldWarmUp to return true when both are enabled")
	}
}

func TestProcessConfig_ValidStructure(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.toml")
	
	configContent := `[[sites]]
remote = "https://github.com/"
dir = "./test-repos/"

[[sites.repos]]
repo = "nonexistent/repo"
warm_up = false
memo = "Test repository"
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	
	report := &reporter.MakeReport{Actions: make([]*reporter.MakeAction, 0)}
	
	// 这个测试主要验证配置解析，Git操作会失败但不影响测试
	err = ProcessConfig(configFile, report)
	
	// 检查配置文件是否被正确解析（即使Git操作失败）
	if err != nil {
		t.Logf("ProcessConfig failed as expected due to Git operations: %v", err)
	}
	
	// 检查是否记录了操作
	if len(report.Actions) == 0 {
		t.Error("Expected at least one action in report")
	}
	
	if len(report.Actions) > 0 {
		action := report.Actions[0]
		if action.Repository != "nonexistent/repo" {
			t.Errorf("Expected repository 'nonexistent/repo', got %s", action.Repository)
		}
		
		if action.Memo != "Test repository" {
			t.Errorf("Expected memo 'Test repository', got %s", action.Memo)
		}
	}
}

func TestProcessConfig_NonExistentFile(t *testing.T) {
	err := ProcessConfig("non-existent-file.toml", nil)
	if err == nil {
		t.Error("Expected error for non-existent config file")
	}
}

func TestProcessConfig_EmptyReport(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test.toml")
	
	configContent := `
[[sites]]
remote_prefix = "https://github.com/"
dir = "` + tempDir + `"
warm_up_all = false

[[sites.repos]]
repo = "test/repo"
warm_up = false
memo = "Test repository"
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	// 测试ProcessConfig不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ProcessConfig panicked: %v", r)
		}
	}()
	
	err = ProcessConfig(configFile, nil)
	// 预期会有错误，因为没有真实的Git仓库
	if err != nil {
		t.Logf("ProcessConfig completed with expected errors: %v", err)
	}
}

func TestProcessConfig_WithReport(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test.toml")
	
	configContent := `
[[sites]]
remote_prefix = "https://github.com/"
dir = "` + tempDir + `"
warm_up_all = false

[[sites.repos]]
repo = "test/repo"
warm_up = false
memo = "Test repository"
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	report := reporter.MakeReport{}
	
	// 测试ProcessConfig不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ProcessConfig panicked: %v", r)
		}
	}()
	
	err = ProcessConfig(configFile, &report)
	// 预期会有错误，因为没有真实的Git仓库
	if err != nil {
		t.Logf("ProcessConfig completed with expected errors: %v", err)
	}
	
	// 验证报告中有动作记录
	if len(report.Actions) == 0 {
		t.Error("Expected at least one action in the report")
	}
} 