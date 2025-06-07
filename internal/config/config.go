package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config represents the complete configuration structure
type Config struct {
	Sites []SiteConfig `toml:"sites"`
}

// SiteConfig represents a site configuration with repositories
type SiteConfig struct {
	RemotePrefix string `toml:"remote"`
	Dir          string `toml:"dir"`
	Repos        []Repo `toml:"repos"`
	WarmUpAll    bool   `toml:"warm_up_all"`
}

// Repo represents a single repository configuration
type Repo struct {
	Repo   string `toml:"repo"`
	Rename string `toml:"rename"`
	WarmUp bool   `toml:"warm_up"`
	Memo   string `toml:"memo"`
}

// ReadFromFile reads and parses a TOML configuration file
func ReadFromFile(configPath string) (*Config, error) {
	var config Config

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	return &config, nil
}

// SaveToFile saves a configuration structure to a TOML file
func SaveToFile(config Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("failed to encode TOML: %w", err)
	}

	return nil
}

// RepoUrl generates the complete Git repository URL for cloning
func (repo Repo) RepoUrl(site SiteConfig) string {
	url := strings.TrimSuffix(site.RemotePrefix, "/") + "/" + repo.Repo
	if !strings.HasSuffix(url, ".git") {
		url += ".git"
	}
	return url
}

// FullPath generates the complete local filesystem path for the repository
func (repo Repo) FullPath(site SiteConfig) string {
	dirname := repo.Repo
	if repo.Rename != "" {
		dirname = repo.Rename
	}
	
	// Extract just the repository name if it contains "/"
	if strings.Contains(dirname, "/") {
		parts := strings.Split(dirname, "/")
		dirname = parts[len(parts)-1]
	}
	
	return filepath.Join(site.Dir, dirname)
}

// DisplayName returns the display name for the repository
func (repo Repo) DisplayName() string {
	if repo.Rename != "" {
		return repo.Rename
	}
	return repo.Repo
}

// ToTOML converts a Config struct to TOML string format
func ToTOML(cfg *Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("config is nil")
	}
	
	var builder strings.Builder
	
	for i, site := range cfg.Sites {
		if i > 0 {
			builder.WriteString("\n")
		}
		
		builder.WriteString("[[sites]]\n")
		builder.WriteString(fmt.Sprintf("    remote_prefix = %q\n", site.RemotePrefix))
		builder.WriteString(fmt.Sprintf("    dir = %q\n", site.Dir))
		
		if site.WarmUpAll {
			builder.WriteString("    warm_up_all = true\n")
		}
		
		builder.WriteString("\n")
		
		for _, repo := range site.Repos {
			builder.WriteString("    [[sites.repos]]\n")
			builder.WriteString(fmt.Sprintf("        repo = %q\n", repo.Repo))
			
			if repo.Rename != "" {
				builder.WriteString(fmt.Sprintf("        rename = %q\n", repo.Rename))
			}
			
			if repo.WarmUp {
				builder.WriteString("        warm_up = true\n")
			}
			
			if repo.Memo != "" {
				builder.WriteString(fmt.Sprintf("        memo = %q\n", repo.Memo))
			}
			
			builder.WriteString("\n")
		}
	}
	
	return builder.String(), nil
} 