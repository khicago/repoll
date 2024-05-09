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
		return runCommandWithTimer(cmd)
	} else {
		fmt.Printf("Updating repository in %s...\n", fullPath)
		cmd := exec.Command("git", "-C", fullPath, "pull")
		return runCommandWithTimer(cmd)
	}
}
