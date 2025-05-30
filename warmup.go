package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func warmUpRepo(fullPath string) error {
	fmt.Printf("- warm-up > Executing warm-up operations for repository in %s...\n", fullPath)

	// 检查目录是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("repository directory %s does not exist", fullPath)
	} else if err != nil {
		return fmt.Errorf("failed to check repository directory %s: %w", fullPath, err)
	}

	// 尝试 Go 项目预热
	if isGoProject(fullPath) {
		return warmUpGoProject(fullPath)
	}

	// 尝试 Node.js 项目预热
	if isNodeProject(fullPath) {
		return warmUpNodeProject(fullPath)
	}

	// 如果都不是，返回 nil（不是错误）
	return nil
}

// isGoProject 检查是否为 Go 项目
func isGoProject(fullPath string) bool {
	goModPath := filepath.Join(fullPath, "go.mod")
	_, err := os.Stat(goModPath)
	return err == nil
}

// isNodeProject 检查是否为 Node.js 项目
func isNodeProject(fullPath string) bool {
	packageJSONPath := filepath.Join(fullPath, "package.json")
	_, err := os.Stat(packageJSONPath)
	return err == nil
}

// warmUpGoProject 预热 Go 项目
func warmUpGoProject(fullPath string) error {
	fmt.Println("- warm-up > Detected Go project, running 'go mod download'...")
	cmd := exec.Command("go", "mod", "download")
	cmd.Dir = fullPath
	if err := runCommandWithTimer(cmd); err != nil {
		return fmt.Errorf("failed to run 'go mod download' in %s: %w", fullPath, err)
	}
	return nil
}

// warmUpNodeProject 预热 Node.js 项目
func warmUpNodeProject(fullPath string) error {
	fmt.Println("- warm-up > Detected Node.js project, determining package manager...")

	if hasYarnLock(fullPath) {
		return warmUpWithYarn(fullPath)
	}
	return warmUpWithNpm(fullPath)
}

// hasYarnLock 检查是否有 yarn.lock 文件
func hasYarnLock(fullPath string) bool {
	yarnLockPath := filepath.Join(fullPath, "yarn.lock")
	_, err := os.Stat(yarnLockPath)
	return err == nil
}

// warmUpWithYarn 使用 yarn 预热
func warmUpWithYarn(fullPath string) error {
	fmt.Println("- warm-up > Running 'yarn install'...")
	cmd := exec.Command("yarn", "install")
	cmd.Dir = fullPath
	if err := runCommandWithTimer(cmd); err != nil {
		return fmt.Errorf("failed to run 'yarn install' in %s: %w", fullPath, err)
	}
	return nil
}

// warmUpWithNpm 使用 npm 预热
func warmUpWithNpm(fullPath string) error {
	fmt.Println("- warm-up > Running 'npm install'...")
	cmd := exec.Command("npm", "install")
	cmd.Dir = fullPath
	if err := runCommandWithTimer(cmd); err != nil {
		return fmt.Errorf("failed to run 'npm install' in %s: %w", fullPath, err)
	}
	return nil
}
