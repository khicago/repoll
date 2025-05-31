---
layout: default
title: API Reference
permalink: /api/
---

# API Reference

repoll provides both a command-line interface and a Go package API for repository management.

## Command Line Interface

### Main Commands

#### `repoll [config-files...]`

Process one or more configuration files to clone/update repositories.

**Usage:**
```bash
repoll config1.toml config2.toml
```

**Options:**
- `--report`: Generate detailed execution report
- `--dry-run`: Show what would be done without executing
- `--update-only`: Only update existing repositories, don't clone new ones
- `--verbose`: Enable verbose logging

**Examples:**
```bash
# Basic usage
repoll repos.toml

# With detailed reporting
repoll repos.toml --report

# Multiple configuration files
repoll dev.toml staging.toml --report
```

#### `repoll mkconf [directory]`

Generate configuration file from existing Git repositories in a directory.

**Usage:**
```bash
repoll mkconf ./projects/
```

**Options:**
- `--report`: Generate execution report
- `--output <file>`: Specify output file (default: auto-generated name)

**Examples:**
```bash
# Scan current directory
repoll mkconf .

# Scan specific directory with report
repoll mkconf ./my-projects/ --report
```

#### `repoll version`

Display version information.

**Output:**
```
repoll version 1.0.0
Built: 2023-12-01T10:30:00Z
Commit: abc123def456
```

#### `repoll help`

Show help information for commands.

## Configuration API

### Config Structure

```go
type Config struct {
    Sites []SiteConfig `toml:"sites"`
}

type SiteConfig struct {
    RemotePrefix string `toml:"remote"`
    Dir          string `toml:"dir"`
    Repos        []Repo `toml:"repos"`
    WarmUpAll    bool   `toml:"warm_up_all"`
}

type Repo struct {
    Repo   string `toml:"repo"`
    Rename string `toml:"rename"`
    WarmUp bool   `toml:"warm_up"`
    Memo   string `toml:"memo"`
}
```

### Configuration Functions

#### `readConfig(configPath string) (*Config, error)`

Read and parse a TOML configuration file.

**Parameters:**
- `configPath`: Path to the TOML configuration file

**Returns:**
- `*Config`: Parsed configuration structure
- `error`: Any error encountered during reading or parsing

**Example:**
```go
config, err := readConfig("repos.toml")
if err != nil {
    log.Fatal(err)
}
```

#### `saveConfigToFile(config Config) error`

Save a configuration structure to a timestamped TOML file.

**Parameters:**
- `config`: Configuration structure to save

**Returns:**
- `error`: Any error encountered during file writing

**Example:**
```go
config := Config{
    Sites: []SiteConfig{
        {
            RemotePrefix: "https://github.com/",
            Dir:          "./projects/",
            Repos: []Repo{
                {Repo: "golang/go", WarmUp: true},
            },
        },
    },
}

err := saveConfigToFile(config)
if err != nil {
    log.Fatal(err)
}
```

## Repository Management API

### Repository Operations

#### `(repo Repo) RepoUrl(site SiteConfig) string`

Generate the complete Git repository URL for cloning.

**Parameters:**
- `site`: Site configuration containing remote prefix

**Returns:**
- `string`: Complete repository URL

**Example:**
```go
repo := Repo{Repo: "golang/go"}
site := SiteConfig{RemotePrefix: "https://github.com/"}
url := repo.RepoUrl(site) // Returns: "https://github.com/golang/go.git"
```

#### `(repo Repo) FullPath(site SiteConfig) string`

Generate the complete local filesystem path for the repository.

**Parameters:**
- `site`: Site configuration containing base directory

**Returns:**
- `string`: Complete local path

**Example:**
```go
repo := Repo{Repo: "golang/go", Rename: "go-lang"}
site := SiteConfig{Dir: "./projects/"}
path := repo.FullPath(site) // Returns: "./projects/go-lang"
```

### Git Operations

#### `cloneRepo(url, targetDir string) error`

Clone a Git repository from URL to target directory.

**Parameters:**
- `url`: Git repository URL
- `targetDir`: Local directory path for cloning

**Returns:**
- `error`: Any error encountered during cloning

#### `updateRepo(repoDir string) error`

Update an existing Git repository by pulling latest changes.

**Parameters:**
- `repoDir`: Local directory path of the repository

**Returns:**
- `error`: Any error encountered during updating

### Warm-up API

#### `performWarmUp(repoDir string) error`

Execute warm-up operations based on detected project type.

**Parameters:**
- `repoDir`: Local directory path of the repository

**Returns:**
- `error`: Any error encountered during warm-up

**Supported Project Types:**
- **Go**: Runs `go mod download` and `go mod tidy`
- **Node.js**: Runs `npm install` or `yarn install`
- **Python**: Runs `pip install -r requirements.txt` (planned)
- **Rust**: Runs `cargo fetch` (planned)

#### `shouldWarmUp(repo Repo, site SiteConfig) bool`

Determine if warm-up should be performed for a repository.

**Parameters:**
- `repo`: Repository configuration
- `site`: Site configuration

**Returns:**
- `bool`: True if warm-up should be performed

## Discovery API

### Git Repository Discovery

#### `discoverGitRepo(path string) (*RepoDiscoveryResult, error)`

Discover Git repository information from a directory path.

**Returns:**
```go
type RepoDiscoveryResult struct {
    Path        string
    Origin      string
    HasOrigin   bool
    Uncommitted bool
    Unmerged    bool
}
```

#### `getRepoNameFromURL(url string) string`

Extract repository name from a Git URL.

**Example:**
```go
name := getRepoNameFromURL("https://github.com/golang/go.git")
// Returns: "go"
```

#### `getRemotePrefix(origin string) string`

Extract remote prefix from a Git origin URL.

**Example:**
```go
prefix := getRemotePrefix("https://github.com/golang/go.git")
// Returns: "https://github.com/golang"
```

## Reporting API

### Report Structures

```go
type MakeReport struct {
    Actions []*MakeAction
}

type MakeAction struct {
    Time       time.Time
    Repository string
    Duration   time.Duration
    Success    bool
    Error      string
    Memo       string
}

type MkconfReport struct {
    Actions []*MkconfAction
}

type MkconfAction struct {
    Time        time.Time
    Path        string
    Origin      string
    HasOrigin   bool
    Uncommitted bool
    Unmerged    bool
}
```

### Report Methods

#### `(mr *MakeReport) Report() string`

Generate a formatted string report of repository operations.

#### `(mr *MkconfReport) Report() string`

Generate a formatted string report of repository discovery.

## Error Handling

All API functions return errors following Go conventions. Common error types:

- **Configuration Errors**: Invalid TOML syntax, missing files
- **Git Errors**: Clone failures, network issues, authentication problems
- **Filesystem Errors**: Permission denied, disk space issues
- **Command Errors**: Missing dependencies (git, go, npm)

**Example Error Handling:**
```go
config, err := readConfig("repos.toml")
if err != nil {
    if os.IsNotExist(err) {
        log.Fatal("Configuration file not found")
    }
    log.Fatalf("Failed to read config: %v", err)
}
```

## Integration Examples

### Custom CLI Tool

```go
package main

import (
    "log"
    "github.com/khicago/repoll"
)

func main() {
    config, err := readConfig("my-repos.toml")
    if err != nil {
        log.Fatal(err)
    }
    
    report := &MakeReport{Actions: make([]*MakeAction, 0)}
    
    for _, site := range config.Sites {
        err := processSiteRepos(site, report)
        if err != nil {
            log.Printf("Site processing failed: %v", err)
        }
    }
    
    fmt.Println(report.Report())
}
```

### Web Service Integration

```go
func handleRepoSync(w http.ResponseWriter, r *http.Request) {
    configPath := r.FormValue("config")
    
    config, err := readConfig(configPath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Process repositories...
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
        "message": "Repositories synchronized",
    })
}
```

For more implementation examples, see the [Examples](examples.md) page. 