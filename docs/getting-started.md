# Getting Started

Welcome to repoll! This guide will help you get up and running in minutes.

## Installation

### Option 1: Go Install (Recommended)

```bash
go install github.com/khicago/repoll@latest
```

### Option 2: Download Binary

Download the latest release from [GitHub Releases](https://github.com/khicago/repoll/releases):

```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/khicago/repoll/main/install.sh | bash

# Windows (PowerShell)
iwr -useb https://raw.githubusercontent.com/khicago/repoll/main/install.ps1 | iex
```

### Option 3: Build from Source

```bash
git clone https://github.com/khicago/repoll.git
cd repoll
go build -o repoll
```

## Your First Configuration

Create a file named `repos.toml`:

```toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./my-projects/"
    warm_up_all = true

    [[sites.repos]]
        repo = "golang/go"
        warm_up = true

    [[sites.repos]]
        repo = "microsoft/vscode"
        memo = "VS Code editor"
```

## Run repoll

```bash
repoll repos.toml
```

You'll see output like:

```
🚀 Processing configuration: repos.toml
📁 Site: https://github.com/ -> ./my-projects/

⬇️  Cloning golang/go...
✅ Cloned successfully: ./my-projects/go/
🔥 Running warm-up: go mod download
✅ Warm-up completed

⬇️  Cloning microsoft/vscode...
✅ Cloned successfully: ./my-projects/vscode/
🔥 Running warm-up: npm install
✅ Warm-up completed

🎉 All repositories processed successfully!
Total time: 45.2s
```

## What Happened?

1. **Configuration Parsed**: repoll read your `repos.toml` file
2. **Repositories Cloned**: Each repository was cloned to the specified directory
3. **Warm-up Executed**: Project dependencies were automatically installed
4. **Progress Reported**: Real-time feedback on all operations

## Next Steps

- [Learn about Configuration](configuration.md) - Master the TOML configuration format
- [Explore Examples](examples.md) - See real-world use cases
- [Command Reference](commands.md) - Discover all available commands

## Quick Tips

### Generate Configuration from Existing Projects

If you already have Git repositories in a directory:

```bash
repoll mkconf ./existing-projects/
```

This will scan the directory and generate a configuration file automatically.

### Update Only Mode

To update existing repositories without cloning new ones:

```bash
repoll --update-only repos.toml
```

### Dry Run

See what repoll would do without actually doing it:

```bash
repoll --dry-run repos.toml
```

## Prerequisites

- **Git**: Must be installed and accessible in your PATH
- **Go** (optional): Required for Go project warm-up
- **Node.js/npm** (optional): Required for Node.js project warm-up

## Troubleshooting

Having issues? Check the [Troubleshooting Guide](troubleshooting.md) or [open an issue](https://github.com/khicago/repoll/issues). 