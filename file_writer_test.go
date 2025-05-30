// Package main contains tests for file writing functionality.
// This file provides comprehensive test coverage for configuration file writing
// and status memo building functions.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestSaveConfigToFile(t *testing.T) {
	// Create a test configuration
	config := Config{
		Sites: []SiteConfig{
			{
				RemotePrefix: "https://github.com/user",
				Dir:          "/home/user/projects",
				Repos: []Repo{
					{
						Repo:   "repo1",
						Rename: "",
						WarmUp: false,
						Memo:   "",
					},
					{
						Repo:   "repo2",
						Rename: "custom-name",
						WarmUp: true,
						Memo:   "test memo",
					},
				},
				WarmUpAll: false,
			},
		},
	}

	// Save the configuration
	err := saveConfigToFile(config)
	if err != nil {
		t.Fatalf("saveConfigToFile failed: %v", err)
	}

	// Find the generated file (it has a timestamp in the name)
	files, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	var configFile string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_conf.toml") {
			configFile = file.Name()
			break
		}
	}

	if configFile == "" {
		t.Fatal("No configuration file found")
	}

	// Clean up
	defer os.Remove(configFile)

	// Read and verify the content
	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	contentStr := string(content)

	// Verify the content contains expected elements
	expectedElements := []string{
		"[[sites]]",
		`remote = "https://github.com/user"`,
		`dir = "/home/user/projects"`,
		"[[sites.repos]]",
		`repo = "repo1"`,
		`repo = "repo2"`,
		`rename = "custom-name"`,
		"warm_up = true",
		`memo = "test memo"`,
	}

	for _, element := range expectedElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("Expected content to contain %q, but it didn't.\nContent:\n%s", element, contentStr)
		}
	}
}

func TestSaveConfigToFile_EmptyConfig(t *testing.T) {
	// Test with empty configuration
	config := Config{
		Sites: []SiteConfig{},
	}

	err := saveConfigToFile(config)
	if err != nil {
		t.Fatalf("saveConfigToFile failed with empty config: %v", err)
	}

	// Find and clean up the generated file
	files, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_conf.toml") {
			os.Remove(file.Name())
			break
		}
	}
}

func TestMkStatusMemo_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		need     bool
		memo     string
		expected string
	}{
		{
			name:     "Add memo to empty string",
			str:      "",
			need:     true,
			memo:     "uncommitted",
			expected: "uncommitted",
		},
		{
			name:     "Don't add memo when not needed",
			str:      "",
			need:     false,
			memo:     "uncommitted",
			expected: "",
		},
		{
			name:     "Add memo to existing string",
			str:      "existing",
			need:     true,
			memo:     "uncommitted",
			expected: "uncommitted", // Note: function overwrites existing string
		},
		{
			name:     "Don't modify existing string when not needed",
			str:      "existing",
			need:     false,
			memo:     "uncommitted",
			expected: "existing",
		},
		{
			name:     "Add empty memo",
			str:      "",
			need:     true,
			memo:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mkStatusMemo(tt.str, tt.need, tt.memo)
			if result != tt.expected {
				t.Errorf("mkStatusMemo(%q, %v, %q) = %q, expected %q",
					tt.str, tt.need, tt.memo, result, tt.expected)
			}
		})
	}
}

func TestSaveConfigToFile_ComplexConfig(t *testing.T) {
	// Test with a more complex configuration
	config := Config{
		Sites: []SiteConfig{
			{
				RemotePrefix: "https://github.com/org1",
				Dir:          "/path/to/org1",
				Repos: []Repo{
					{
						Repo:   "repo1",
						Rename: "",
						WarmUp: false,
						Memo:   "",
					},
				},
				WarmUpAll: false,
			},
			{
				RemotePrefix: "git@gitlab.com:org2",
				Dir:          "/path/to/org2",
				Repos: []Repo{
					{
						Repo:   "repo2",
						Rename: "custom-repo2",
						WarmUp: true,
						Memo:   "important repo",
					},
					{
						Repo:   "repo3",
						Rename: "",
						WarmUp: false,
						Memo:   "test repo",
					},
				},
				WarmUpAll: true,
			},
		},
	}

	err := saveConfigToFile(config)
	if err != nil {
		t.Fatalf("saveConfigToFile failed: %v", err)
	}

	// Find and clean up the generated file
	files, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	var configFile string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_conf.toml") {
			configFile = file.Name()
			break
		}
	}

	if configFile != "" {
		defer os.Remove(configFile)

		// Read and verify the content
		content, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		contentStr := string(content)

		// Verify multiple sites are present
		siteCount := strings.Count(contentStr, "[[sites]]")
		if siteCount != 2 {
			t.Errorf("Expected 2 sites in config, found %d", siteCount)
		}

		// Verify specific content
		expectedElements := []string{
			`remote = "https://github.com/org1"`,
			`remote = "git@gitlab.com:org2"`,
			`rename = "custom-repo2"`,
			`memo = "important repo"`,
			`memo = "test repo"`,
		}

		for _, element := range expectedElements {
			if !strings.Contains(contentStr, element) {
				t.Errorf("Expected content to contain %q", element)
			}
		}
	}
}

func TestWriteConfig_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		setupFunc func(string) error
		expectErr bool
	}{
		{
			name: "Empty config",
			config: Config{
				Sites: []SiteConfig{
					{
						RemotePrefix: "",
						Dir:          "",
						Repos:        []Repo{},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "Config with empty site names",
			config: Config{
				Sites: []SiteConfig{
					{
						RemotePrefix: "",
						Dir:          "",
						Repos:        []Repo{},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "Very long config content",
			config: Config{
				Sites: func() []SiteConfig {
					var sites []SiteConfig
					for i := 0; i < 100; i++ {
						sites = append(sites, SiteConfig{
							RemotePrefix: fmt.Sprintf("https://github.com/org%d", i),
							Dir:          fmt.Sprintf("/path/to/site%d", i),
							Repos: []Repo{
								{
									Repo:    fmt.Sprintf("repo%d", i),
									Rename:  "",
									WarmUp:  i%2 == 0,
									Memo:    fmt.Sprintf("memo%d", i),
								},
							},
						})
					}
					return sites
				}(),
			},
			expectErr: false,
		},
		{
			name: "Config with special characters in paths",
			config: Config{
				Sites: []SiteConfig{
					{
						RemotePrefix: "https://github.com/user-name_123",
						Dir:          "/path/with spaces/and-dashes_underscores/directory",
						Repos: []Repo{
							{
								Repo:   "repo-with_special.chars",
								Rename: "",
								WarmUp: true,
								Memo:   "Special repo",
							},
						},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			if tt.setupFunc != nil {
				if err := tt.setupFunc(tempDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Note: saveConfigToFile doesn't take a file path parameter, it generates its own filename
			// We need to change directory to tempDir for this test
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			os.Chdir(tempDir)

			err := saveConfigToFile(tt.config)
			
			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// If write succeeded, verify a config file was created
			if !tt.expectErr && err == nil {
				files, err := filepath.Glob("*_conf.toml")
				if err != nil || len(files) == 0 {
					t.Error("Config file was not created")
				}

				if len(files) > 0 {
					// Try to read back the content
					content, err := os.ReadFile(files[0])
					if err != nil {
						t.Errorf("Could not read back config file: %v", err)
					}

					// Basic validation that it's valid TOML-like content
					if len(content) == 0 && len(tt.config.Sites) > 0 {
						t.Error("Config file is empty but should contain data")
					}
				}
			}
		})
	}
}

func TestSaveConfigToFile_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create a read-only directory
	readOnlyDir := filepath.Join(tempDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0444); err != nil {
		t.Skipf("Could not create read-only directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755) // Restore permissions for cleanup

	config := Config{
		Sites: []SiteConfig{
			{
				RemotePrefix: "https://github.com/user",
				Dir:          "/home/user/projects",
				Repos: []Repo{
					{
						Repo:   "test-repo",
						Rename: "",
						WarmUp: false,
						Memo:   "",
					},
				},
			},
		},
	}

	// Change to read-only directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(readOnlyDir)

	err := saveConfigToFile(config)
	
	// Should fail due to permission error
	if err == nil {
		t.Error("Expected permission error but write succeeded")
	}
}

func TestSaveConfigToFile_LargeContent(t *testing.T) {
	// Test with very large content to check memory and performance
	var repos []Repo
	for i := 0; i < 1000; i++ {
		repos = append(repos, Repo{
			Repo:   fmt.Sprintf("large-repo-%d", i),
			Rename: "",
			WarmUp: i%2 == 0,
			Memo:   fmt.Sprintf("Large repo number %d", i),
		})
	}

	config := Config{
		Sites: []SiteConfig{
			{
				RemotePrefix: "https://github.com/organization",
				Dir:          "/very/long/path/to/base/directory",
				Repos:        repos,
			},
		},
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	err := saveConfigToFile(config)
	if err != nil {
		t.Errorf("Failed to write large config: %v", err)
	}

	// Verify a config file was created and has content
	files, err := filepath.Glob("*_conf.toml")
	if err != nil || len(files) == 0 {
		t.Error("Large config file was not created")
		return
	}

	stat, err := os.Stat(files[0])
	if err != nil {
		t.Errorf("Could not stat config file: %v", err)
	}

	if stat.Size() == 0 {
		t.Error("Large config file is empty")
	}

	// Should be fairly large (>10KB for 1000 repos)
	if stat.Size() < 10000 {
		t.Errorf("Config file size %d seems too small for 1000 repos", stat.Size())
	}
}

func TestSaveConfigToFile_ConcurrentWrite(t *testing.T) {
	// Test concurrent writes to different directories
	config := Config{
		Sites: []SiteConfig{
			{
				RemotePrefix: "https://github.com/user",
				Dir:          "/home/user/projects",
				Repos: []Repo{
					{
						Repo:   "concurrent-test",
						Rename: "",
						WarmUp: false,
						Memo:   "",
					},
				},
			},
		},
	}

	tempDir := t.TempDir()
	
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// Start 10 concurrent writes to different subdirectories
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			// Create a subdirectory for this goroutine
			subDir := filepath.Join(tempDir, fmt.Sprintf("subdir-%d", index))
			if err := os.MkdirAll(subDir, 0755); err != nil {
				errors <- fmt.Errorf("failed to create subdir %d: %v", index, err)
				return
			}
			
			// Change to the subdirectory
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			os.Chdir(subDir)
			
			if err := saveConfigToFile(config); err != nil {
				errors <- fmt.Errorf("write %d failed: %v", index, err)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent write error: %v", err)
	}

	// Verify config files were created in each subdirectory
	for i := 0; i < 10; i++ {
		subDir := filepath.Join(tempDir, fmt.Sprintf("subdir-%d", i))
		files, err := filepath.Glob(filepath.Join(subDir, "*_conf.toml"))
		if err != nil || len(files) == 0 {
			t.Errorf("Config file was not created in subdirectory %d", i)
		}
	}
}
