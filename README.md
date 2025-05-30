# 🚀 repoll

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/github/release/khicago/repoll?style=for-the-badge)](https://github.com/khicago/repoll/releases)
[![Tests](https://img.shields.io/badge/Tests-82.9%25_Coverage-brightgreen?style=for-the-badge)](https://github.com/khicago/repoll/actions)

**🎯 The Ultimate Git Repository Management Tool**

*Clone, update, and warm-up multiple repositories with zero configuration overhead*

[📖 Documentation](https://khicago.github.io/repoll) • [🚀 Quick Start](#quick-start) • [💡 Examples](#examples) • [🤝 Contributing](#contributing)

</div>

---

## ✨ What is repoll?

**repoll** (Repository Puller) is a lightning-fast, developer-friendly CLI tool that revolutionizes how you manage multiple Git repositories. Whether you're working with microservices, managing open-source contributions, or handling complex multi-repo projects, repoll makes it effortless.

### 🎯 Why repoll?

- **⚡ Lightning Fast**: Concurrent operations with intelligent dependency management
- **🧠 Smart Warm-up**: Automatically prepares Go, Node.js, and other projects for development
- **📋 Simple Configuration**: One TOML file to rule them all
- **🔄 Flexible Workflows**: Clone, update, or sync with customizable strategies
- **📊 Rich Reporting**: Beautiful progress indicators and detailed execution reports
- **🛡️ Production Ready**: 82.9% test coverage with robust error handling

## 🚀 Quick Start

### Installation

**Option 1: Download Binary (Recommended)**
```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/khicago/repoll/main/install.sh | bash

# Windows (PowerShell)
iwr -useb https://raw.githubusercontent.com/khicago/repoll/main/install.ps1 | iex
```

**Option 2: Go Install**
```bash
go install github.com/khicago/repoll@latest
```

**Option 3: Build from Source**
```bash
git clone https://github.com/khicago/repoll.git
cd repoll
go build -o repoll
```

### Basic Usage

1. **Create a configuration file** (`repos.toml`):
```toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./projects/"
    warm_up_all = true

    [[sites.repos]]
        repo = "microsoft/vscode"
        warm_up = true
    
    [[sites.repos]]
        repo = "golang/go"
        
    [[sites.repos]]
        repo = "facebook/react"
```

2. **Run repoll**:
```bash
repoll repos.toml
```

3. **Watch the magic happen** ✨

## 💡 Examples

### 🏢 Enterprise Multi-Service Setup
```toml
# microservices.toml
[[sites]]
    remote_prefix = "https://git.company.com/"
    dir = "./microservices/"
    warm_up_all = true

    [[sites.repos]]
        repo = "team/user-service"
        
    [[sites.repos]]
        repo = "team/payment-service"
        
    [[sites.repos]]
        repo = "team/notification-service"

[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./tools/"
    
    [[sites.repos]]
        repo = "hashicorp/terraform"
        
    [[sites.repos]]
        repo = "kubernetes/kubernetes"
```

### 🌟 Open Source Contributor Setup
```toml
# opensource.toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./contributions/"
    
    [[sites.repos]]
        repo = "microsoft/TypeScript"
        warm_up = true
        
    [[sites.repos]]
        repo = "rust-lang/rust"
        
    [[sites.repos]]
        repo = "golang/go"
        warm_up = true
```

### 🚀 Advanced Configuration
```bash
# Generate configuration from existing directories
repoll mkconf ./my-projects/

# Update existing repositories only
repoll --update-only repos.toml

# Dry run to see what would happen
repoll --dry-run repos.toml

# Custom output format
repoll --format json repos.toml
```

## 🎮 Commands

| Command | Description | Example |
|---------|-------------|---------|
| `repoll <config>` | Clone/update repositories | `repoll repos.toml` |
| `repoll mkconf <dir>` | Generate config from directory | `repoll mkconf ./projects/` |
| `repoll version` | Show version info | `repoll version` |
| `repoll help` | Show help | `repoll help` |

## ⚙️ Configuration Reference

### Site Configuration
```toml
[[sites]]
    remote_prefix = "https://github.com/"  # Repository URL prefix
    dir = "./local-dir/"                   # Local directory for clones
    warm_up_all = false                    # Enable warm-up for all repos
```

### Repository Configuration
```toml
[[sites.repos]]
    repo = "owner/repository"              # Repository path
    warm_up = true                         # Enable project warm-up
    rename = "custom-name"                 # Custom local directory name
    memo = "Development version"           # Description/memo
```

### Warm-up Features

repoll intelligently detects project types and runs appropriate setup commands:

- **Go Projects**: `go mod download`
- **Node.js Projects**: `npm install` or `yarn install` (auto-detected)
- **Python Projects**: `pip install -r requirements.txt` (coming soon)
- **Rust Projects**: `cargo fetch` (coming soon)

## 📊 Performance & Quality

- **🏃‍♂️ Fast**: Concurrent operations with intelligent batching
- **🧪 Tested**: 82.9% test coverage with comprehensive edge case handling
- **🛡️ Reliable**: Production-ready error handling and recovery
- **📈 Scalable**: Handles hundreds of repositories efficiently
- **💾 Lightweight**: Single binary with minimal dependencies

## 🤝 Contributing

We love contributions! repoll is built with ❤️ by developers, for developers.

### 🌟 Ways to Contribute

- 🐛 **Bug Reports**: Found an issue? [Open an issue](https://github.com/khicago/repoll/issues)
- 💡 **Feature Requests**: Have an idea? [Start a discussion](https://github.com/khicago/repoll/discussions)
- 🔧 **Code Contributions**: [Submit a PR](https://github.com/khicago/repoll/pulls)
- 📖 **Documentation**: Help improve our docs
- 🌍 **Spread the Word**: Star the repo, share with friends!

### 🛠️ Development Setup

```bash
git clone https://github.com/khicago/repoll.git
cd repoll
go mod download
go test -v -cover
go build -o repoll
```

### 📋 Code Quality Standards

- ✅ **Tests Required**: Maintain >80% coverage
- 🧹 **Linting**: `golangci-lint run`
- 📖 **Documentation**: Update docs for new features
- 🎯 **Performance**: Benchmark critical paths

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/) and love ❤️
- CLI framework by [Cobra](https://github.com/spf13/cobra)
- Configuration parsing by [TOML](https://github.com/BurntSushi/toml)

---

<div align="center">

**⭐ If repoll helps you, please consider giving it a star! ⭐**

Made with ❤️ by the [repoll community](https://github.com/khicago/repoll/graphs/contributors)

</div>