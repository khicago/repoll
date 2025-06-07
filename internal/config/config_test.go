package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFromFile_ValidConfig(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "valid.toml")
	
	configContent := `[[sites]]
remote = "https://github.com/"
dir = "./test-repos/"
warm_up_all = true

[[sites.repos]]
repo = "golang/example"
rename = "example-go"
warm_up = false
memo = "Test repository"

[[sites.repos]]
repo = "user/another"
warm_up = true
memo = "Another test repo"
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	
	config, err := ReadFromFile(configFile)
	if err != nil {
		t.Fatalf("ReadFromFile failed: %v", err)
	}
	
	// 验证基本结构
	if len(config.Sites) != 1 {
		t.Errorf("Expected 1 site, got %d", len(config.Sites))
	}
	
	site := config.Sites[0]
	if site.RemotePrefix != "https://github.com/" {
		t.Errorf("Expected remote prefix 'https://github.com/', got %s", site.RemotePrefix)
	}
	
	if site.Dir != "./test-repos/" {
		t.Errorf("Expected dir './test-repos/', got %s", site.Dir)
	}
	
	if !site.WarmUpAll {
		t.Error("Expected warm_up_all to be true")
	}
	
	// 验证仓库
	if len(site.Repos) != 2 {
		t.Errorf("Expected 2 repos, got %d", len(site.Repos))
	}
	
	repo1 := site.Repos[0]
	if repo1.Repo != "golang/example" {
		t.Errorf("Expected repo 'golang/example', got %s", repo1.Repo)
	}
	
	if repo1.Rename != "example-go" {
		t.Errorf("Expected rename 'example-go', got %s", repo1.Rename)
	}
	
	if repo1.WarmUp {
		t.Error("Expected warm_up to be false")
	}
	
	if repo1.Memo != "Test repository" {
		t.Errorf("Expected memo 'Test repository', got %s", repo1.Memo)
	}
}

func TestReadFromFile_NonExistentFile(t *testing.T) {
	_, err := ReadFromFile("non-existent-file.toml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestReadFromFile_InvalidToml(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid.toml")
	
	invalidContent := `[[sites]
remote = "missing quote`
	
	err := os.WriteFile(configFile, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid config: %v", err)
	}
	
	_, err = ReadFromFile(configFile)
	if err == nil {
		t.Error("Expected error for invalid TOML")
	}
}

func TestSaveToFile_ValidConfig(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "output.toml")
	
	config := Config{
		Sites: []SiteConfig{
			{
				RemotePrefix: "https://github.com/",
				Dir:          "./repos/",
				WarmUpAll:    true,
				Repos: []Repo{
					{
						Repo:   "test/repo",
						Rename: "renamed-repo",
						WarmUp: false,
						Memo:   "Test save",
					},
				},
			},
		},
	}
	
	err := SaveToFile(config, configFile)
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
	
	// 验证文件是否创建
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
	
	// 验证内容是否正确保存（通过重新读取）
	loadedConfig, err := ReadFromFile(configFile)
	if err != nil {
		t.Fatalf("Failed to reload saved config: %v", err)
	}
	
	if len(loadedConfig.Sites) != 1 {
		t.Errorf("Expected 1 site, got %d", len(loadedConfig.Sites))
	}
	
	site := loadedConfig.Sites[0]
	if site.RemotePrefix != "https://github.com/" {
		t.Errorf("Expected remote prefix preserved, got %s", site.RemotePrefix)
	}
	
	if len(site.Repos) != 1 {
		t.Errorf("Expected 1 repo, got %d", len(site.Repos))
	}
	
	repo := site.Repos[0]
	if repo.Repo != "test/repo" {
		t.Errorf("Expected repo 'test/repo', got %s", repo.Repo)
	}
}

func TestSaveToFile_InvalidPath(t *testing.T) {
	config := Config{}
	
	// 尝试保存到无效路径
	err := SaveToFile(config, "/invalid/path/config.toml")
	if err == nil {
		t.Error("Expected error for invalid file path")
	}
}

func TestRepo_RepoUrl(t *testing.T) {
	site := SiteConfig{
		RemotePrefix: "https://github.com/",
	}
	
	tests := []struct {
		name     string
		repo     Repo
		expected string
	}{
		{
			name:     "basic repo",
			repo:     Repo{Repo: "user/repo"},
			expected: "https://github.com/user/repo.git",
		},
		{
			name:     "repo already with .git",
			repo:     Repo{Repo: "user/repo.git"},
			expected: "https://github.com/user/repo.git",
		},
		{
			name:     "empty repo",
			repo:     Repo{Repo: ""},
			expected: "https://github.com/.git",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.repo.RepoUrl(site)
			if result != test.expected {
				t.Errorf("RepoUrl() = %s, expected %s", result, test.expected)
			}
		})
	}
}

func TestRepo_RepoUrl_WithTrailingSlash(t *testing.T) {
	site := SiteConfig{
		RemotePrefix: "https://github.com",  // 没有尾部斜杠
	}
	
	repo := Repo{Repo: "user/repo"}
	result := repo.RepoUrl(site)
	expected := "https://github.com/user/repo.git"
	
	if result != expected {
		t.Errorf("RepoUrl() = %s, expected %s", result, expected)
	}
}

func TestRepo_FullPath(t *testing.T) {
	site := SiteConfig{
		Dir: "/home/user/repos",
	}
	
	tests := []struct {
		name     string
		repo     Repo
		expected string
	}{
		{
			name:     "simple repo",
			repo:     Repo{Repo: "simple"},
			expected: "/home/user/repos/simple",
		},
		{
			name:     "repo with owner",
			repo:     Repo{Repo: "owner/repo"},
			expected: "/home/user/repos/repo",
		},
		{
			name:     "repo with rename",
			repo:     Repo{Repo: "owner/repo", Rename: "custom-name"},
			expected: "/home/user/repos/custom-name",
		},
		{
			name:     "repo with nested rename",
			repo:     Repo{Repo: "owner/repo", Rename: "group/custom-name"},
			expected: "/home/user/repos/custom-name",
		},
		{
			name:     "empty repo",
			repo:     Repo{Repo: ""},
			expected: "/home/user/repos",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.repo.FullPath(site)
			if result != test.expected {
				t.Errorf("FullPath() = %s, expected %s", result, test.expected)
			}
		})
	}
}

func TestRepo_DisplayName(t *testing.T) {
	tests := []struct {
		name     string
		repo     Repo
		expected string
	}{
		{
			name:     "repo without rename",
			repo:     Repo{Repo: "owner/repo"},
			expected: "owner/repo",
		},
		{
			name:     "repo with rename",
			repo:     Repo{Repo: "owner/repo", Rename: "custom-name"},
			expected: "custom-name",
		},
		{
			name:     "empty repo with rename",
			repo:     Repo{Repo: "", Rename: "custom-name"},
			expected: "custom-name",
		},
		{
			name:     "empty repo without rename",
			repo:     Repo{Repo: ""},
			expected: "",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.repo.DisplayName()
			if result != test.expected {
				t.Errorf("DisplayName() = %s, expected %s", result, test.expected)
			}
		})
	}
}

func TestConfig_MultiSites(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "multi.toml")
	
	configContent := `[[sites]]
remote = "https://github.com/"
dir = "./github-repos/"

[[sites.repos]]
repo = "golang/go"

[[sites]]
remote = "https://gitlab.com/"
dir = "./gitlab-repos/"

[[sites.repos]]
repo = "group/project"
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	
	config, err := ReadFromFile(configFile)
	if err != nil {
		t.Fatalf("ReadFromFile failed: %v", err)
	}
	
	if len(config.Sites) != 2 {
		t.Errorf("Expected 2 sites, got %d", len(config.Sites))
	}
	
	// 验证第一个站点
	site1 := config.Sites[0]
	if site1.RemotePrefix != "https://github.com/" {
		t.Errorf("Site 1: Expected remote 'https://github.com/', got %s", site1.RemotePrefix)
	}
	
	if len(site1.Repos) != 1 {
		t.Errorf("Site 1: Expected 1 repo, got %d", len(site1.Repos))
	}
	
	// 验证第二个站点
	site2 := config.Sites[1]
	if site2.RemotePrefix != "https://gitlab.com/" {
		t.Errorf("Site 2: Expected remote 'https://gitlab.com/', got %s", site2.RemotePrefix)
	}
}

func TestRepo_AllFields(t *testing.T) {
	repo := Repo{
		Repo:   "owner/repository",
		Rename: "custom-name",
		WarmUp: true,
		Memo:   "Important project",
	}
	
	// 验证所有字段都正确设置
	if repo.Repo != "owner/repository" {
		t.Errorf("Repo field mismatch: got %s", repo.Repo)
	}
	
	if repo.Rename != "custom-name" {
		t.Errorf("Rename field mismatch: got %s", repo.Rename)
	}
	
	if !repo.WarmUp {
		t.Error("WarmUp field should be true")
	}
	
	if repo.Memo != "Important project" {
		t.Errorf("Memo field mismatch: got %s", repo.Memo)
	}
} 