---
layout: default
title: API Reference
permalink: /api/
---

# API Reference

Complete technical reference for repoll's command-line interface, configuration structure, and internal APIs.

## 🖥️ Command Line Interface

### Main Commands

#### `repoll [options] [config-files...]`

Process one or more configuration files to clone and manage Git repositories.

**Usage:**
```bash
# Process default configuration
repoll

# Process specific configuration file
repoll repos.toml

# Process multiple configuration files
repoll dev.toml staging.toml prod.toml

# Process with options
repoll --dry-run --verbose repos.toml
```

**Options:**
- `--dry-run, -n`: Show what would be done without executing
- `--verbose, -v`: Enable verbose output
- `--help, -h`: Show help message
- `--version`: Show version information

**Examples:**
```bash
# Basic usage
repoll repos.toml

# Dry run to preview actions
repoll --dry-run repos.toml

# Verbose output for debugging
repoll --verbose repos.toml

# Multiple configurations
repoll personal.toml work.toml
```

#### `repoll mkconf [directory]`

Generate a configuration file from existing repositories in a directory.

**Usage:**
```bash
# Generate from current directory
repoll mkconf

# Generate from specific directory
repoll mkconf ./projects

# Generate from multiple directories
repoll mkconf ./frontend ./backend
```

**Output:**
Creates a `repos.toml` file with discovered repositories.

**Examples:**
```bash
# Scan current directory
repoll mkconf .

# Scan specific directory
repoll mkconf ~/development/projects

# Generate with custom output
repoll mkconf ./src > custom-repos.toml
```

#### `repoll version`

Display version information and build details.

**Usage:**
```bash
repoll version
```

**Output:**
```
repoll version 1.0.0
Build: 2024-01-15T10:30:00Z
Go version: go1.21.0
```

#### `repoll help`

Show comprehensive help information.

**Usage:**
```bash
repoll help

# Show help for specific command
repoll help mkconf
```

### Global Options

All commands support these global options:

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--help` | `-h` | Show help message | |
| `--version` | | Show version information | |
| `--verbose` | `-v` | Enable verbose logging | `false` |
| `--dry-run` | `-n` | Preview actions without executing | `false` |

## 📋 Configuration API

### Configuration Structure

The configuration is loaded from TOML files with the following structure:

```go
type Config struct {
    Sites []SiteConfig `toml:"sites"`
}

type SiteConfig struct {
    Remote     string       `toml:"remote"`
    Dir        string       `toml:"dir"`
    WarmUpAll  bool         `toml:"warm_up_all"`
    Repos      []Repo       `toml:"repos"`
}

type Repo struct {
    Repo     string `toml:"repo"`
    Rename   string `toml:"rename"`
    WarmUp   *bool  `toml:"warm_up"`
    Memo     string `toml:"memo"`
}
```

### Configuration Functions

#### Loading Configuration

```go
func LoadConfig(filename string) (*Config, error)
```

Load and parse a TOML configuration file.

**Parameters:**
- `filename`: Path to the TOML configuration file

**Returns:**
- `*Config`: Parsed configuration structure
- `error`: Error if file cannot be loaded or parsed

**Example:**
```go
config, err := LoadConfig("repos.toml")
if err != nil {
    log.Fatal(err)
}
```

#### Saving Configuration

```go
func SaveConfig(config *Config, filename string) error
```

Save a configuration structure to a TOML file.

**Parameters:**
- `config`: Configuration structure to save
- `filename`: Output file path

**Returns:**
- `error`: Error if file cannot be written

**Example:**
```go
err := SaveConfig(config, "generated-repos.toml")
if err != nil {
    log.Fatal(err)
}
```

#### Validating Configuration

```go
func ValidateConfig(config *Config) error
```

Validate a configuration structure for correctness.

**Parameters:**
- `config`: Configuration to validate

**Returns:**
- `error`: Validation error, or nil if valid

**Validation Rules:**
- Each site must have a non-empty `remote` field
- Each site must have a non-empty `dir` field
- Repository names must be non-empty
- Directory paths should be relative

## 🗄️ Repository Management API

### Repository Functions

#### Generate Repository URL

```go
func GenerateRepoURL(remote, repo string) string
```

Generate a complete repository URL from remote prefix and repository path.

**Parameters:**
- `remote`: Remote URL prefix (e.g., "https://github.com/")
- `repo`: Repository path (e.g., "owner/repository")

**Returns:**
- `string`: Complete repository URL

**Examples:**
```go
url := GenerateRepoURL("https://github.com/", "golang/go")
// Returns: "https://github.com/golang/go"

url := GenerateRepoURL("git@github.com:", "user/repo")
// Returns: "git@github.com:user/repo"
```

#### Generate Local Path

```go
func GenerateLocalPath(baseDir, repo, rename string) string
```

Generate the local directory path for a repository.

**Parameters:**
- `baseDir`: Base directory for repositories
- `repo`: Repository path
- `rename`: Custom name (optional)

**Returns:**
- `string`: Local directory path

**Examples:**
```go
path := GenerateLocalPath("./projects", "golang/go", "")
// Returns: "./projects/go"

path := GenerateLocalPath("./projects", "owner/long-name", "short")
// Returns: "./projects/short"
```

### Git Operations API

#### Clone Repository

```go
func CloneRepository(url, localPath string) error
```

Clone a Git repository to a local directory.

**Parameters:**
- `url`: Repository URL
- `localPath`: Local directory path

**Returns:**
- `error`: Error if cloning fails

**Example:**
```go
err := CloneRepository("https://github.com/golang/go", "./projects/go")
if err != nil {
    log.Printf("Clone failed: %v", err)
}
```

#### Update Repository

```go
func UpdateRepository(localPath string) error
```

Update an existing Git repository by pulling latest changes.

**Parameters:**
- `localPath`: Local repository directory

**Returns:**
- `error`: Error if update fails

**Example:**
```go
err := UpdateRepository("./projects/go")
if err != nil {
    log.Printf("Update failed: %v", err)
}
```

#### Check Repository Status

```go
func CheckRepoStatus(localPath string) (*RepoStatus, error)

type RepoStatus struct {
    Exists    bool
    IsGitRepo bool
    HasChanges bool
    Branch    string
    Remote    string
}
```

Check the status of a local repository.

**Parameters:**
- `localPath`: Local repository directory

**Returns:**
- `*RepoStatus`: Repository status information
- `error`: Error if status check fails

**Example:**
```go
status, err := CheckRepoStatus("./projects/go")
if err != nil {
    log.Printf("Status check failed: %v", err)
} else {
    fmt.Printf("Repository exists: %v\n", status.Exists)
    fmt.Printf("Is Git repo: %v\n", status.IsGitRepo)
    fmt.Printf("Has changes: %v\n", status.HasChanges)
}
```

## 🔥 Warm-up API

### Warm-up Functions

#### Detect Project Type

```go
func DetectProjectType(path string) ProjectType

type ProjectType int

const (
    Unknown ProjectType = iota
    Go
    NodeJS
    Python
    Rust
    JavaMaven
    JavaGradle
)
```

Detect the project type based on indicator files.

**Parameters:**
- `path`: Project directory path

**Returns:**
- `ProjectType`: Detected project type

**Example:**
```go
projectType := DetectProjectType("./projects/myapp")
switch projectType {
case Go:
    fmt.Println("Go project detected")
case NodeJS:
    fmt.Println("Node.js project detected")
default:
    fmt.Println("Unknown project type")
}
```

#### Execute Warm-up

```go
func ExecuteWarmUp(path string, projectType ProjectType) error
```

Execute warm-up commands for a specific project type.

**Parameters:**
- `path`: Project directory path
- `projectType`: Project type

**Returns:**
- `error`: Error if warm-up fails

**Warm-up Commands by Project Type:**

| Project Type | Detection File | Commands |
|--------------|----------------|----------|
| Go | `go.mod` | `go mod download` |
| Node.js | `package.json` | `npm install` or `yarn install` |
| Python | `requirements.txt` | `pip install -r requirements.txt` |
| Rust | `Cargo.toml` | `cargo fetch` |
| Java (Maven) | `pom.xml` | `mvn dependency:resolve` |
| Java (Gradle) | `build.gradle` | `gradle dependencies` |

**Example:**
```go
projectType := DetectProjectType("./projects/myapp")
err := ExecuteWarmUp("./projects/myapp", projectType)
if err != nil {
    log.Printf("Warm-up failed: %v", err)
}
```

#### Should Warm Up

```go
func ShouldWarmUp(repo Repo, siteWarmUpAll bool) bool
```

Determine if warm-up should be executed for a repository.

**Parameters:**
- `repo`: Repository configuration
- `siteWarmUpAll`: Site-level warm-up setting

**Returns:**
- `bool`: True if warm-up should be executed

**Decision Logic:**
1. If `repo.WarmUp` is explicitly set, use that value
2. If `siteWarmUpAll` is true, return true
3. Otherwise, auto-detect based on project files

**Example:**
```go
repo := Repo{
    Repo: "golang/go",
    WarmUp: nil, // Auto-detect
}
shouldWarm := ShouldWarmUp(repo, true)
fmt.Printf("Should warm up: %v\n", shouldWarm)
```

## 🚨 Error Handling

### Error Types

repoll defines specific error types for different failure scenarios:

```go
type ConfigError struct {
    File    string
    Message string
}

type GitError struct {
    Repository string
    Operation  string
    Message    string
}

type WarmUpError struct {
    Project     string
    ProjectType ProjectType
    Message     string
}
```

### Error Examples

#### Configuration Errors
```go
// Invalid TOML syntax
ConfigError{
    File: "repos.toml",
    Message: "invalid TOML syntax at line 5",
}

// Missing required field
ConfigError{
    File: "repos.toml",
    Message: "site missing required field 'remote'",
}
```

#### Git Errors
```go
// Repository not found
GitError{
    Repository: "user/nonexistent",
    Operation:  "clone",
    Message:    "repository does not exist",
}

// Authentication failed
GitError{
    Repository: "private/repo",
    Operation:  "clone",
    Message:    "authentication required",
}
```

#### Warm-up Errors
```go
// Command failed
WarmUpError{
    Project:     "./projects/myapp",
    ProjectType: NodeJS,
    Message:     "npm install failed with exit code 1",
}
```

## 📊 Logging and Output

### Log Levels

repoll supports different verbosity levels:

- **Quiet**: Only errors and critical messages
- **Normal**: Standard operation messages
- **Verbose**: Detailed operation information
- **Debug**: All internal operations and state

### Output Formats

#### Normal Output
```
✓ Cloned golang/go -> ./projects/go
🔥 Warmed up ./projects/go (Go project)
✓ Updated microsoft/vscode -> ./projects/vscode
⚠ Skipped warm-up for ./projects/large-repo (disabled)
```

#### Verbose Output
```
[INFO] Loading configuration from repos.toml
[INFO] Found 3 sites with 15 repositories
[INFO] Processing site: https://github.com/
[INFO] Cloning https://github.com/golang/go to ./projects/go
[INFO] Clone completed in 2.3s
[INFO] Detected Go project at ./projects/go
[INFO] Executing warm-up: go mod download
[INFO] Warm-up completed in 1.8s
```

#### Dry Run Output
```
[DRY RUN] Would clone golang/go -> ./projects/go
[DRY RUN] Would warm up ./projects/go (Go project)
[DRY RUN] Would update microsoft/vscode -> ./projects/vscode
[DRY RUN] Would skip warm-up for ./projects/large-repo (disabled)
```

## 🔧 Advanced Usage

### Programmatic Usage

You can use repoll as a library in your Go programs:

```go
package main

import (
    "log"
    "github.com/khicago/repoll/internal/config"
    "github.com/khicago/repoll/internal/process"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("repos.toml")
    if err != nil {
        log.Fatal(err)
    }

    // Process repositories
    processor := process.NewProcessor()
    err = processor.ProcessConfig(cfg)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Custom Processing

```go
// Custom repository processing
for _, site := range cfg.Sites {
    for _, repo := range site.Repos {
        url := GenerateRepoURL(site.Remote, repo.Repo)
        localPath := GenerateLocalPath(site.Dir, repo.Repo, repo.Rename)
        
        // Custom logic here
        if needsProcessing(localPath) {
            err := CloneRepository(url, localPath)
            if err != nil {
                log.Printf("Failed to clone %s: %v", repo.Repo, err)
                continue
            }
        }
    }
}
```

## 📞 Integration Examples

### GitHub Actions

```yaml
name: Setup Development Environment
on: [push, pull_request]

jobs:
  setup:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install repoll
        run: |
          curl -fsSL https://raw.githubusercontent.com/khicago/repoll/main/install.sh | bash
      - name: Setup repositories
        run: |
          repoll --verbose dev-repos.toml
```

### Makefile Integration

```makefile
.PHONY: setup update clean

setup:
	@echo "Setting up development environment..."
	repoll repos.toml

update:
	@echo "Updating repositories..."
	repoll --verbose repos.toml

clean:
	@echo "Cleaning repositories..."
	rm -rf ./projects/*

dry-run:
	@echo "Preview setup actions..."
	repoll --dry-run repos.toml
```

### Shell Script Integration

```bash
#!/bin/bash
set -e

echo "🚀 Setting up development environment..."

# Install repoll if not exists
if ! command -v repoll &> /dev/null; then
    echo "Installing repoll..."
    curl -fsSL https://raw.githubusercontent.com/khicago/repoll/main/install.sh | bash
fi

# Generate configuration if not exists
if [[ ! -f "repos.toml" ]]; then
    echo "Generating configuration from current directory..."
    repoll mkconf . > repos.toml
fi

# Process repositories
echo "Processing repositories..."
repoll --verbose repos.toml

echo "✅ Setup completed!"
```

---

<div class="notice--success">
  <h4>🎯 Quick Reference:</h4>
  <ul>
    <li><strong>Basic usage:</strong> <code>repoll repos.toml</code></li>
    <li><strong>Dry run:</strong> <code>repoll --dry-run repos.toml</code></li>
    <li><strong>Generate config:</strong> <code>repoll mkconf .</code></li>
    <li><strong>Verbose output:</strong> <code>repoll --verbose repos.toml</code></li>
  </ul>
</div> 