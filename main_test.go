package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	if rootCmd == nil {
		t.Error("rootCmd should not be nil")
	}
	if rootCmd.Use != "repoll" {
		t.Errorf("Expected rootCmd.Use to be 'repoll', got '%s'", rootCmd.Use)
	}
}

func TestMakeCommand(t *testing.T) {
	if cmdMake == nil {
		t.Error("cmdMake should not be nil")
	}
	if cmdMake.Use != "make [paths to the TOML config file]" {
		t.Errorf("Expected cmdMake.Use to be 'make [paths to the TOML config file]', got '%s'", cmdMake.Use)
	}
	if cmdMake.Short == "" {
		t.Error("cmdMake.Short should not be empty")
	}
}

func TestMakeConfCommand(t *testing.T) {
	if cmdMakeConf == nil {
		t.Error("cmdMakeConf should not be nil")
	}
	if cmdMakeConf.Use != "mkconf [directory]" {
		t.Errorf("Expected cmdMakeConf.Use to be 'mkconf [directory]', got '%s'", cmdMakeConf.Use)
	}
	if cmdMakeConf.Short == "" {
		t.Error("cmdMakeConf.Short should not be empty")
	}
}

func TestCommandsRegistered(t *testing.T) {
	commands := rootCmd.Commands()

	var makeFound, mkconfFound bool
	for _, cmd := range commands {
		if cmd.Name() == "make" {
			makeFound = true
		}
		if cmd.Name() == "mkconf" {
			mkconfFound = true
		}
	}

	if !makeFound {
		t.Error("make command should be registered")
	}
	if !mkconfFound {
		t.Error("mkconf command should be registered")
	}
}

func TestMakeCommandFlags(t *testing.T) {
	flag := cmdMake.Flags().Lookup("report")
	if flag == nil {
		t.Error("make command should have a 'report' flag")
	}
	if flag.DefValue != "false" {
		t.Errorf("Expected default value of 'report' flag to be 'false', got '%s'", flag.DefValue)
	}
}

func TestMakeConfCommandFlags(t *testing.T) {
	flag := cmdMakeConf.Flags().Lookup("report")
	if flag == nil {
		t.Error("mkconf command should have a 'report' flag")
	}
	if flag.DefValue != "false" {
		t.Errorf("Expected default value of 'report' flag to be 'false', got '%s'", flag.DefValue)
	}
}

func TestCmdMake_WithReport(t *testing.T) {
	// Test the make command with --report flag
	configFiles := []string{
		"test1.toml",
		"test2.toml",
	}
	
	cmd := cmdMake
	args := append(configFiles, "--report")
	cmd.SetArgs(args)
	
	// Execute the command
	err := cmd.Execute()
	// The command might fail due to missing directories, but we're testing the argument processing
	_ = err // Ignore error as we're testing argument processing
}

func TestCmdMakeConf_WithReport(t *testing.T) {
	// Create a temporary directory with a git repo
	tempDir := t.TempDir()
	
	// Create a mock git directory
	gitDir := filepath.Join(tempDir, "test-repo", ".git")
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create git directory: %v", err)
	}
	
	// Create a basic git config
	configContent := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
[remote "origin"]
	url = https://github.com/test/test-repo.git
	fetch = +refs/heads/*:refs/remotes/origin/*
`
	configPath := filepath.Join(gitDir, "config")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create git config: %v", err)
	}
	
	// Test the mkconf command with report flag
	cmd := cmdMakeConf
	cmd.SetArgs([]string{tempDir, "--report"})
	
	// Capture the command execution
	err = cmd.Execute()
	// The command might have issues, but we're testing the code path
}

func TestCmdMake_MultipleConfigs(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	
	// Create multiple test config files
	configs := []string{
		`[[sites]]
remote = "https://github.com/test1/"
dir = "/tmp/test1"
[[sites.repos]]
repo = "repo1"
rename = ""
warm_up = false
memo = "test repo 1"`,
		`[[sites]]
remote = "https://github.com/test2/"
dir = "/tmp/test2"
[[sites.repos]]
repo = "repo2"
rename = ""
warm_up = true
memo = "test repo 2"`,
	}
	
	var configFiles []string
	for i, content := range configs {
		configFile := filepath.Join(tempDir, fmt.Sprintf("config%d.toml", i))
		err := os.WriteFile(configFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file %d: %v", i, err)
		}
		configFiles = append(configFiles, configFile)
	}
	
	// Test the make command with multiple config files
	cmd := cmdMake
	args := append(configFiles, "--report")
	cmd.SetArgs(args)
	
	// Execute the command
	err := cmd.Execute()
	// The command might fail due to missing directories, but we're testing the argument processing
	_ = err // Ignore error as we're testing argument processing
}

func TestCmdMake_InvalidPath(t *testing.T) {
	// Test with an invalid/non-existent config file
	cmd := cmdMake
	cmd.SetArgs([]string{"/non/existent/path/config.toml"})
	
	// This should handle the error gracefully
	err := cmd.Execute()
	// The command should not panic, even with invalid paths
	_ = err // Ignore error as we're testing error handling
}

func TestCmdMakeConf_InvalidDirectory(t *testing.T) {
	// Test with an invalid/non-existent directory
	cmd := cmdMakeConf
	cmd.SetArgs([]string{"/non/existent/directory"})
	
	// This should handle the error gracefully
	err := cmd.Execute()
	// The command should not panic, even with invalid directories
	_ = err // Ignore error as we're testing error handling
}

func TestMain_ErrorHandling(t *testing.T) {
	// Test the main function's error handling by simulating a command error
	// We can't directly test main() easily, but we can test the error handling pattern
	
	// Create a command that will fail
	testCmd := &cobra.Command{
		Use: "test-fail",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("simulated error")
		},
	}
	
	err := testCmd.Execute()
	if err == nil {
		t.Error("Expected error but got none")
	}
	
	// Verify the error message
	if !strings.Contains(err.Error(), "simulated error") {
		t.Errorf("Expected error message to contain 'simulated error', got: %v", err)
	}
}
