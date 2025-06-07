package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/khicago/repoll/internal/git"
	"github.com/khicago/repoll/internal/reporter"
)

// GenerateFromDirectory generates a configuration by scanning a directory for Git repositories
func GenerateFromDirectory(targetDir string, report *reporter.MkconfReport) (*Config, error) {
	config := &Config{
		Sites: make([]SiteConfig, 0),
	}

	siteMap := make(map[string]*SiteConfig)

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// Check if this directory is a Git repository
		gitDir := filepath.Join(path, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			return nil
		}

		// Discover repository information
		repoInfo, err := git.DiscoverRepository(path)
		if err != nil {
			// Log but continue with other repositories
			fmt.Printf("Warning: Failed to discover repository at %s: %v\n", path, err)
			return nil
		}

		// Add to report
		if report != nil {
			action := &reporter.MkconfAction{
				Time:        time.Now(),
				Path:        path,
				Origin:      repoInfo.Origin,
				HasOrigin:   repoInfo.HasOrigin,
				Uncommitted: repoInfo.Uncommitted,
				Unmerged:    repoInfo.Unmerged,
			}
			report.Actions = append(report.Actions, action)
		}

		if !repoInfo.HasOrigin {
			fmt.Printf("Skipping repository without origin: %s\n", path)
			return nil
		}

		// Extract repository name and remote prefix
		repoName := git.ExtractRepoNameFromURL(repoInfo.Origin)
		remotePrefix := git.ExtractRemotePrefix(repoInfo.Origin)

		if repoName == "" || remotePrefix == "" {
			fmt.Printf("Warning: Could not parse repository info for %s\n", path)
			return nil
		}

		// Get or create site configuration
		site, exists := siteMap[remotePrefix]
		if !exists {
			// Determine appropriate directory
			relPath, err := filepath.Rel(targetDir, filepath.Dir(path))
			if err != nil {
				relPath = filepath.Dir(path)
			}
			if relPath == "." {
				relPath = "./"
			} else {
				relPath = "./" + relPath + "/"
			}

			site = &SiteConfig{
				RemotePrefix: remotePrefix,
				Dir:          relPath,
				Repos:        make([]Repo, 0),
				WarmUpAll:    false,
			}
			siteMap[remotePrefix] = site
			config.Sites = append(config.Sites, *site)
		}

		// Create repository configuration
		repo := Repo{
			Repo:   repoName,
			WarmUp: shouldDefaultWarmUp(path),
			Memo:   generateMemoFromPath(path),
		}

		// Check if we need to use a custom name
		expectedPath := filepath.Join(site.Dir, strings.Split(repoName, "/")[len(strings.Split(repoName, "/"))-1])
		actualPath, _ := filepath.Rel(targetDir, path)
		if actualPath != expectedPath {
			repo.Rename = filepath.Base(path)
		}

		// Add to the correct site
		for i := range config.Sites {
			if config.Sites[i].RemotePrefix == remotePrefix {
				config.Sites[i].Repos = append(config.Sites[i].Repos, repo)
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(config.Sites) == 0 {
		return nil, fmt.Errorf("no Git repositories found in %s", targetDir)
	}

	return config, nil
}

// shouldDefaultWarmUp determines if warm-up should be enabled by default based on project type
func shouldDefaultWarmUp(path string) bool {
	// Check for common project files that indicate warm-up should be enabled
	warmupIndicators := []string{
		"go.mod",          // Go projects
		"package.json",    // Node.js projects
		"requirements.txt", // Python projects
		"Cargo.toml",      // Rust projects
		"pom.xml",         // Maven projects
		"build.gradle",    // Gradle projects
	}

	for _, indicator := range warmupIndicators {
		if _, err := os.Stat(filepath.Join(path, indicator)); err == nil {
			return true
		}
	}

	return false
}

// generateMemoFromPath generates a memo based on the repository path
func generateMemoFromPath(path string) string {
	// Try to read common documentation files for a description
	readmeFiles := []string{"README.md", "README.rst", "README.txt", "README"}
	
	for _, readme := range readmeFiles {
		readmePath := filepath.Join(path, readme)
		if content, err := os.ReadFile(readmePath); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if len(line) > 10 && len(line) < 100 && !strings.HasPrefix(line, "#") {
					// Try to find a descriptive line
					if strings.Contains(strings.ToLower(line), "description") ||
						strings.Contains(strings.ToLower(line), "about") {
						return line
					}
				}
			}
			// If no description found, use the first substantial line
			for _, line := range lines {
				line = strings.TrimSpace(line)
				line = strings.TrimPrefix(line, "# ")
				if len(line) > 10 && len(line) < 100 {
					return line
				}
			}
		}
	}

	return ""
} 