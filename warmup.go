package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func warmUpRepo(fullPath string) error {
	fmt.Printf("- warm-up > Executing warm-up operations for repository in %s...\n", fullPath)

	// Check if 'go.mod' exists to identify a Go project
	goModPath := filepath.Join(fullPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		fmt.Println("- warm-up > Detected Go project, running 'go mod download'...")
		cmd := exec.Command("go", "mod", "download")
		cmd.Dir = fullPath // Set the working directory
		return runCommandWithTimer(cmd)
	}

	// Check if 'package.json' exists to identify a Node.js project
	packageJSONPath := filepath.Join(fullPath, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		fmt.Println("- warm-up > Detected Node.js project, determining package manager...")

		// Use `yarn` if 'yarn.lock' is present, else use 'npm install'
		if _, err = os.Stat(filepath.Join(fullPath, "yarn.lock")); err == nil {
			fmt.Println("- warm-up > Running 'yarn install'...")
			cmd := exec.Command("yarn", "install")
			cmd.Dir = fullPath
			return runCommandWithTimer(cmd)
		} else {
			fmt.Println("- warm-up > Running 'npm install'...")
			cmd := exec.Command("npm", "install")
			cmd.Dir = fullPath
			return runCommandWithTimer(cmd)
		}
	}

	// Add more warm-up checks and commands for other project types here...

	return nil
}
