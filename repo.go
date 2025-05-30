package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func isURL(path string) bool {
	u, err := url.Parse(path)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func ensureTrailingSlash(s string) string {
	if !strings.HasSuffix(s, "/") {
		return s + "/"
	}
	return s
}

func gitCloneOrUpdate(repo Repo, site SiteConfig) error {
	fullPath := repo.FullPath(site)
	repoURL := repo.RepoUrl(site)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("Cloning %s into %s...\n", repoURL, fullPath)
		cmd := exec.Command("git", "clone", repoURL, fullPath)
		if err := runCommandWithTimer(cmd); err != nil {
			return fmt.Errorf("failed to clone repository %s to %s: %w", repoURL, fullPath, err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check if directory %s exists: %w", fullPath, err)
	} else {
		fmt.Printf("Updating repository in %s...\n", fullPath)
		cmd := exec.Command("git", "-C", fullPath, "pull")
		if err := runCommandWithTimer(cmd); err != nil {
			return fmt.Errorf("failed to update repository in %s: %w", fullPath, err)
		}
		return nil
	}
}
