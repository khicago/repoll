package warmup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Perform executes warm-up operations based on detected project type
func Perform(repoDir string) error {
	// Detect project type and run appropriate warm-up commands
	if isGoProject(repoDir) {
		return warmUpGo(repoDir)
	}

	if isNodeProject(repoDir) {
		return warmUpNode(repoDir)
	}

	if isPythonProject(repoDir) {
		return warmUpPython(repoDir)
	}

	if isRustProject(repoDir) {
		return warmUpRust(repoDir)
	}

	// No specific project type detected, but that's okay
	return nil
}

// isGoProject checks if the directory contains a Go project
func isGoProject(dir string) bool {
	goModPath := filepath.Join(dir, "go.mod")
	_, err := os.Stat(goModPath)
	return err == nil
}

// isNodeProject checks if the directory contains a Node.js project
func isNodeProject(dir string) bool {
	packagePath := filepath.Join(dir, "package.json")
	_, err := os.Stat(packagePath)
	return err == nil
}

// isPythonProject checks if the directory contains a Python project
func isPythonProject(dir string) bool {
	requirementsPath := filepath.Join(dir, "requirements.txt")
	_, err := os.Stat(requirementsPath)
	return err == nil
}

// isRustProject checks if the directory contains a Rust project
func isRustProject(dir string) bool {
	cargoPath := filepath.Join(dir, "Cargo.toml")
	_, err := os.Stat(cargoPath)
	return err == nil
}

// warmUpGo performs warm-up for Go projects
func warmUpGo(repoDir string) error {
	// Run go mod download
	cmd := exec.Command("go", "mod", "download")
	cmd.Dir = repoDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod download failed: %w\nOutput: %s", err, string(output))
	}

	// Run go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = repoDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// warmUpNode performs warm-up for Node.js projects
func warmUpNode(repoDir string) error {
	// Check if yarn.lock exists, prefer yarn over npm
	yarnLockPath := filepath.Join(repoDir, "yarn.lock")
	if _, err := os.Stat(yarnLockPath); err == nil {
		return warmUpWithYarn(repoDir)
	}

	return warmUpWithNpm(repoDir)
}

// warmUpWithYarn performs warm-up using Yarn
func warmUpWithYarn(repoDir string) error {
	cmd := exec.Command("yarn", "install")
	cmd.Dir = repoDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("yarn install failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// warmUpWithNpm performs warm-up using npm
func warmUpWithNpm(repoDir string) error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = repoDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("npm install failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// warmUpPython performs warm-up for Python projects
func warmUpPython(repoDir string) error {
	cmd := exec.Command("pip", "install", "-r", "requirements.txt")
	cmd.Dir = repoDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pip install failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// warmUpRust performs warm-up for Rust projects
func warmUpRust(repoDir string) error {
	cmd := exec.Command("cargo", "fetch")
	cmd.Dir = repoDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cargo fetch failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}
