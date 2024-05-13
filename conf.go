package main

import (
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type (
	// Repo represents a single repository within a site
	Repo struct {
		Repo   string `toml:"repo"`
		Rename string `toml:"rename,omitempty"`
		WarmUp bool   `toml:"warm_up,omitempty"`
		Memo   string `toml:"memo,omitempty"`
	}

	// SiteConfig represents a remote location and the directory to clone the repos to
	SiteConfig struct {
		RemotePrefix string `toml:"remote"`
		Dir          string `toml:"dir"`
		Repos        []Repo `toml:"repos"`
		WarmUpAll    bool   `toml:"warm_up,omitempty"`
	}

	// Config represents the top level toml configuration file format
	Config struct {
		Sites []SiteConfig `toml:"sites"`
	}
)

func readConfig(configPath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (repo Repo) RepoUrl(site SiteConfig) string {
	repoURL := strings.TrimSpace(repo.Repo) + ".git"
	if isURL(site.RemotePrefix) { // "https://git@...../"
		repoURL = ensureTrailingSlash(site.RemotePrefix) + repoURL
	} else { // "git@....:"
		repoURL = site.RemotePrefix + repoURL
	}

	return repoURL
}

func (repo Repo) FullPath(site SiteConfig) string {
	repo.Repo = strings.TrimSpace(repo.Repo)
	if repo.Rename != "" {
		if strings.ToLower(repo.Rename) == "{base}" {
			return filepath.Join(site.Dir, filepath.Base(repo.Repo))
		}
		return filepath.Join(site.Dir, strings.TrimSpace(repo.Rename))
	}
	return filepath.Join(site.Dir, repo.Repo)
}
