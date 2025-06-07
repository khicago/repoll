package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/khicago/repoll/internal/reporter"
)

func TestGenerateFromDirectory_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	
	_, err := GenerateFromDirectory(tempDir, nil)
	if err == nil {
		t.Error("Expected error for directory with no Git repositories")
	}
	
	expectedMsg := "no Git repositories found"
	if err != nil && err.Error() != "no Git repositories found in "+tempDir {
		t.Errorf("Expected error message containing '%s', got: %v", expectedMsg, err)
	}
}

func TestGenerateFromDirectory_WithGitRepo(t *testing.T) {
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
		t.Skipf("Git not available, skipping test: %v", err)
	}
	
	// 添加origin
	cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/test/repo.git")
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to add origin: %v", err)
	}
	
	report := &reporter.MkconfReport{Actions: make([]*reporter.MkconfAction, 0)}
	config, err := GenerateFromDirectory(tempDir, report)
	if err != nil {
		t.Fatalf("GenerateFromDirectory failed: %v", err)
	}
	
	// 验证生成的配置
	if len(config.Sites) != 1 {
		t.Errorf("Expected 1 site, got %d", len(config.Sites))
	}
	
	site := config.Sites[0]
	if site.RemotePrefix != "https://github.com/" {
		t.Errorf("Expected remote prefix 'https://github.com/', got %s", site.RemotePrefix)
	}
	
	if len(site.Repos) != 1 {
		t.Errorf("Expected 1 repo, got %d", len(site.Repos))
	}
	
	repo := site.Repos[0]
	if repo.Repo != "test/repo" {
		t.Errorf("Expected repo 'test/repo', got %s", repo.Repo)
	}
	
	// 验证报告
	if report == nil || len(report.Actions) != 1 {
		t.Errorf("Expected 1 action in report, got %d", len(report.Actions))
	}
}

func TestGenerateFromDirectory_NoOrigin(t *testing.T) {
	tempDir := t.TempDir()
	
	// 创建没有origin的Git仓库
	repoDir := filepath.Join(tempDir, "no-origin-repo")
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}
	
	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		t.Skipf("Git not available, skipping test: %v", err)
	}
	
	// 不添加origin，应该被跳过
	_, err = GenerateFromDirectory(tempDir, nil)
	if err == nil {
		t.Error("Expected error for directory with repositories without origin")
	}
}

func TestGenerateFromDirectory_MultipleRepos(t *testing.T) {
	tempDir := t.TempDir()
	
	// 创建两个不同的Git仓库
	repos := []struct {
		name   string
		origin string
		prefix string
	}{
		{"github-repo", "https://github.com/user/repo1.git", "https://github.com/"},
		{"gitlab-repo", "https://gitlab.com/user/repo2.git", "https://gitlab.com/"},
	}
	
	for _, repo := range repos {
		repoDir := filepath.Join(tempDir, repo.name)
		err := os.MkdirAll(repoDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create repo directory: %v", err)
		}
		
		cmd := exec.Command("git", "init")
		cmd.Dir = repoDir
		err = cmd.Run()
		if err != nil {
			t.Skipf("Git not available, skipping test: %v", err)
		}
		
		cmd = exec.Command("git", "remote", "add", "origin", repo.origin)
		cmd.Dir = repoDir
		err = cmd.Run()
		if err != nil {
			t.Fatalf("Failed to add origin: %v", err)
		}
	}
	
	config, err := GenerateFromDirectory(tempDir, nil)
	if err != nil {
		t.Fatalf("GenerateFromDirectory failed: %v", err)
	}
	
	// 应该有两个站点（不同的remote prefix）
	if len(config.Sites) != 2 {
		t.Errorf("Expected 2 sites, got %d", len(config.Sites))
	}
	
	// 验证站点配置
	foundGitHub := false
	foundGitLab := false
	for _, site := range config.Sites {
		if site.RemotePrefix == "https://github.com/" {
			foundGitHub = true
		}
		if site.RemotePrefix == "https://gitlab.com/" {
			foundGitLab = true
		}
	}
	
	if !foundGitHub {
		t.Error("GitHub site not found")
	}
	if !foundGitLab {
		t.Error("GitLab site not found")
	}
}

func TestShouldDefaultWarmUp(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name     string
		files    []string
		expected bool
	}{
		{
			name:     "Go project",
			files:    []string{"go.mod"},
			expected: true,
		},
		{
			name:     "Node.js project",
			files:    []string{"package.json"},
			expected: true,
		},
		{
			name:     "Python project",
			files:    []string{"requirements.txt"},
			expected: true,
		},
		{
			name:     "Rust project",
			files:    []string{"Cargo.toml"},
			expected: true,
		},
		{
			name:     "Maven project",
			files:    []string{"pom.xml"},
			expected: true,
		},
		{
			name:     "Gradle project",
			files:    []string{"build.gradle"},
			expected: true,
		},
		{
			name:     "Multiple indicators",
			files:    []string{"go.mod", "package.json"},
			expected: true,
		},
		{
			name:     "No indicators",
			files:    []string{"README.md", "LICENSE"},
			expected: false,
		},
		{
			name:     "Empty directory",
			files:    []string{},
			expected: false,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 为每个测试创建子目录
			testDir := filepath.Join(tempDir, test.name)
			err := os.MkdirAll(testDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}
			
			// 创建测试文件
			for _, file := range test.files {
				filePath := filepath.Join(testDir, file)
				err := os.WriteFile(filePath, []byte("test content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}
			
			result := shouldDefaultWarmUp(testDir)
			if result != test.expected {
				t.Errorf("shouldDefaultWarmUp(%s) = %v, expected %v", testDir, result, test.expected)
			}
		})
	}
}

func TestGenerateMemoFromPath(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name        string
		readmeFile  string
		content     string
		expectedLen int
		shouldFind  bool
	}{
		{
			name:       "README.md with description",
			readmeFile: "README.md",
			content: `# Test Project

This is a description about the project.

## Installation
...`,
			expectedLen: 10,
			shouldFind:  true,
		},
		{
			name:       "README with title only",
			readmeFile: "README.md",
			content: `# Simple Project Title

Some content here that is descriptive enough.`,
			expectedLen: 10,
			shouldFind:  true,
		},
		{
			name:       "README with short lines",
			readmeFile: "README.txt",
			content: `Short
Very short line
This is a longer line that should be picked up by the generator`,
			expectedLen: 10,
			shouldFind:  true,
		},
		{
			name:       "No README file",
			readmeFile: "",
			content:    "",
			expectedLen: 0,
			shouldFind:  false,
		},
		{
			name:       "README with only short lines",
			readmeFile: "README.rst",
			content: `# Title
Short
Brief
OK`,
			expectedLen: 0,
			shouldFind:  false,
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 为每个测试创建子目录
			testDir := filepath.Join(tempDir, test.name)
			err := os.MkdirAll(testDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}
			
			// 创建README文件（如果指定）
			if test.readmeFile != "" {
				readmePath := filepath.Join(testDir, test.readmeFile)
				err := os.WriteFile(readmePath, []byte(test.content), 0644)
				if err != nil {
					t.Fatalf("Failed to create README file: %v", err)
				}
			}
			
			result := generateMemoFromPath(testDir)
			
			if test.shouldFind {
				if len(result) < test.expectedLen {
					t.Errorf("generateMemoFromPath(%s) returned too short memo: '%s' (length: %d)", testDir, result, len(result))
				}
			} else {
				if result != "" {
					t.Errorf("generateMemoFromPath(%s) should return empty string, got: '%s'", testDir, result)
				}
			}
		})
	}
}

func TestGenerateFromDirectory_WithReport(t *testing.T) {
	tempDir := t.TempDir()
	
	// 创建Git仓库
	repoDir := filepath.Join(tempDir, "test-repo")
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}
	
	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		t.Skipf("Git not available, skipping test: %v", err)
	}
	
	cmd = exec.Command("git", "remote", "add", "origin", "https://github.com/test/repo.git")
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to add origin: %v", err)
	}
	
	// 创建一个测试文件并添加到暂存区（模拟未提交变更）
	testFile := filepath.Join(repoDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	report := &reporter.MkconfReport{Actions: make([]*reporter.MkconfAction, 0)}
	_, err = GenerateFromDirectory(tempDir, report)
	if err != nil {
		t.Fatalf("GenerateFromDirectory failed: %v", err)
	}
	
	// 验证报告内容
	if len(report.Actions) != 1 {
		t.Errorf("Expected 1 action in report, got %d", len(report.Actions))
	}
	
	action := report.Actions[0]
	if action.Path != repoDir {
		t.Errorf("Expected action path %s, got %s", repoDir, action.Path)
	}
	
	if action.Origin != "https://github.com/test/repo.git" {
		t.Errorf("Expected origin 'https://github.com/test/repo.git', got %s", action.Origin)
	}
	
	if !action.HasOrigin {
		t.Error("Expected HasOrigin to be true")
	}
} 