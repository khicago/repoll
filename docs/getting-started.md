---
layout: default
title: Getting Started
permalink: /getting-started/
---

# 🚀 Getting Started with repoll

`repoll` is a lightning-fast, developer-friendly CLI tool that revolutionizes how you manage multiple Git repositories. Whether you're managing personal projects or coordinating team development, repoll streamlines your workflow with intelligent automation.

## 📦 Installation

### Quick Install (macOS/Linux)
```bash
# Using curl (recommended)
curl -fsSL https://raw.githubusercontent.com/khicago/repoll/main/install.sh | bash

# Or using wget
wget -qO- https://raw.githubusercontent.com/khicago/repoll/main/install.sh | bash
```

### Windows (PowerShell)
```powershell
# Using PowerShell
iwr -useb https://raw.githubusercontent.com/khicago/repoll/main/install.ps1 | iex
```

### Alternative Installation Methods

#### Using Go
```bash
go install github.com/khicago/repoll/cmd/repoll@latest
```

#### Download Binary
Visit our [releases page](https://github.com/khicago/repoll/releases) and download the appropriate binary for your platform.

#### Build from Source
```bash
git clone https://github.com/khicago/repoll.git
cd repoll
go build -o repoll ./cmd/repoll
```

## ✅ Verify Installation

```bash
repoll --version
# Output: repoll version dev (commit: none, built: 2025-06-02)
```

## 🎯 Quick Start

### 1. Create Your First Configuration

Create a `repos.toml` file:

```toml
[[sites]]
    remote = "https://github.com/"
    dir = "./projects/"
    warm_up_all = false

    [[sites.repos]]
        repo = "golang/example"
        warm_up = true
        memo = "Go example projects"

    [[sites.repos]]
        repo = "microsoft/vscode-docs"
        memo = "VS Code documentation"
```

### 2. Run repoll

```bash
# Clone/update repositories
repoll run repos.toml

# With verbose output
repoll run --verbose repos.toml

# Dry run (see what would happen)
repoll run --dry-run --verbose repos.toml
```

**Expected Output:**
```
🚀 repoll - Git Repository Management Tool
ℹ Loading configuration from repos.toml
ℹ Found 1 site(s) with 2 repositories
  Processing site: https://github.com/
⟳ Cloning golang/example...
  Cloning https://github.com/golang/example.git to projects/example
  Starting warm-up for projects/example
  Warm-up completed for projects/example
✓ golang/example (2.9s)
⟳ Cloning microsoft/vscode-docs...
  Cloning https://github.com/microsoft/vscode-docs.git to projects/vscode-docs
✓ microsoft/vscode-docs (1m39s)

┌───────────┐
│ SUMMARY │
└───────────┘
Total repositories: 2
Total time: 1m42s

🎉 All repositories processed successfully!
```

## 🔧 Essential Commands

### Generate Configuration from Existing Projects
```bash
# Scan directory and generate repos.toml
repoll mkconf ./my-projects/

# With verbose output
repoll mkconf --verbose ./my-projects/
```

### Update Existing Repositories
```bash
# Run again to update existing repos
repoll run repos.toml
```

### Batch Process Multiple Configurations
```bash
# Process multiple config files
repoll run config1.toml config2.toml config3.toml
```

## 🔥 Warm-Up Feature

The warm-up feature prepares your projects for development by running common setup commands:

**Supported Project Types:**
- **Go projects** (`go.mod`): `go mod download`, `go mod tidy`
- **Node.js projects** (`package.json`): `npm install` or `yarn install`
- **Python projects** (`requirements.txt`): `pip install -r requirements.txt`
- **Rust projects** (`Cargo.toml`): `cargo fetch`
- **Maven projects** (`pom.xml`): `mvn dependency:resolve`
- **Gradle projects** (`build.gradle`): `gradle dependencies`

## 📋 Configuration Structure

### Complete Configuration Example

```toml
# Multiple sites configuration
[[sites]]
    remote = "https://github.com/"
    dir = "./github-projects/"
    warm_up_all = false

    [[sites.repos]]
        repo = "golang/go"
        rename = "go-lang"  # Custom local name
        warm_up = true
        memo = "The Go programming language"

    [[sites.repos]]
        repo = "microsoft/vscode"
        warm_up = false
        memo = "Visual Studio Code"

[[sites]]
    remote = "https://gitlab.com/"
    dir = "./gitlab-projects/"
    warm_up_all = true  # Enable warm-up for all repos in this site

    [[sites.repos]]
        repo = "gitlab-org/gitlab"
        memo = "GitLab Community Edition"
```

### Configuration Fields

#### Site Configuration (`[[sites]]`)
- `remote`: Git remote URL prefix (e.g., `"https://github.com/"`)
- `dir`: Local directory for repositories
- `warm_up_all`: Enable warm-up for all repositories in this site

#### Repository Configuration (`[[sites.repos]]`)
- `repo`: Repository path (e.g., `"owner/repository"`)
- `rename`: Custom local directory name (optional)
- `warm_up`: Enable warm-up for this specific repository
- `memo`: Description or notes about the repository

## 🛠️ Troubleshooting

### Common Issues

**Issue: "Repository not found"**
```bash
# Check if the repository URL is correct
git ls-remote https://github.com/owner/repo.git
```

**Issue: "Permission denied"**
```bash
# Set up SSH keys or use HTTPS with token
git config --global credential.helper store
```

**Issue: "Directory already exists"**
```bash
# repoll will update existing repositories automatically
# Use --dry-run to see what would happen
repoll run --dry-run repos.toml
```

## 💡 Pro Tips

### Speed Up Operations
```bash
# Use parallel processing (default behavior)
repoll run --verbose repos.toml

# Quiet mode for scripts
repoll run --quiet repos.toml
```

### Configuration Management
```bash
# Generate config from existing projects
repoll mkconf ./existing-projects/

# Combine with your existing config
cat generated-repos.toml >> my-repos.toml
```

### Team Collaboration
```bash
# Share configuration files in your team repository
git add repos.toml
git commit -m "Add repoll configuration"

# Team members can then run:
repoll run repos.toml
```

## 🎯 Next Steps

- 📖 Read the [Configuration Guide](./configuration.md) for advanced setups
- 🔧 Check out [CLI Reference](./cli-reference.md) for all available commands
- 🏗️ Learn about [Project Structure](./project-structure.md) for development
- 🤝 See [Contributing Guide](./contributing.md) to help improve repoll

---

**Need help?** Open an issue on [GitHub](https://github.com/khicago/repoll/issues) or check our [FAQ](./faq.md). 