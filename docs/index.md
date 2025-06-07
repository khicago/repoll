---
layout: splash
title: "repoll - Git Repository Management Made Simple"
permalink: /
header:
  overlay_color: "#0066cc"
  overlay_filter: "0.6"
  overlay_image: /assets/images/header-bg.jpg
  actions:
    - label: "🚀 Get Started"
      url: "/getting-started/"
    - label: "📖 Documentation"
      url: "/docs/"
    - label: "⭐ View on GitHub"
      url: "https://github.com/khicago/repoll"
excerpt: "**82.8% test coverage** • **Production ready** • **Enterprise quality**<br/>A powerful CLI tool for managing multiple Git repositories with declarative TOML configuration. Clone, update, and setup your development environment in seconds."

intro: 
  - excerpt: '🎯 **Trusted by developers worldwide** - repoll simplifies the management of multiple Git repositories through simple TOML configuration files. Whether you are setting up a new development environment, managing enterprise projects, or keeping dozens of repositories in sync, repoll automates the tedious tasks so you can focus on what matters most - **your code**.'

feature_row:
  - image_path: /assets/images/feature-config.png
    alt: "Configuration-driven"
    title: "🔧 Configuration-driven"
    excerpt: "Define your repositories in simple TOML files. No complex scripts, no manual processes. Just declare what you want."
  - image_path: /assets/images/feature-automation.png
    alt: "Smart Automation"
    title: "🤖 Smart Automation"
    excerpt: "Automatically detect project types (Go, Node.js, Python) and run appropriate setup commands. Your environment is ready in minutes."
  - image_path: /assets/images/feature-enterprise.png
    alt: "Enterprise Ready"
    title: "🏢 Enterprise Ready"
    excerpt: "82.8% test coverage, production-grade quality, support for multiple Git providers, and team collaboration workflows."

feature_row2:
  - image_path: /assets/images/feature-warmup.png
    alt: "Intelligent Setup"
    title: "🔥 Intelligent Project Setup"
    excerpt: 'Automatically detect and setup Go modules, Node.js packages, Python dependencies, and more. One command sets up your entire development environment.'
    url: "/examples/#warmup-examples"
    btn_label: "See Examples"
    btn_class: "btn--primary"

feature_row3:
  - image_path: /assets/images/feature-discovery.png
    alt: "Repository Discovery"
    title: "🔍 Repository Discovery"
    excerpt: 'Scan existing directories to automatically generate configuration files. Perfect for migrating existing setups or onboarding new team members.'
    url: "/examples/#discovery-examples"
    btn_label: "Learn More"
    btn_class: "btn--primary"

quality_row:
  - title: "🏆 Production Quality"
    excerpt: "**82.8%** test coverage with comprehensive unit and integration tests"
  - title: "⚡ High Performance"
    excerpt: "Parallel operations, smart caching, and optimized Git workflows"
  - title: "🛡️ Battle Tested"
    excerpt: "Handles edge cases, robust error handling, and enterprise-grade reliability"
---

{% include feature_row id="intro" type="center" %}

{% include feature_row %}

{% include feature_row id="feature_row2" type="left" %}

{% include feature_row id="feature_row3" type="right" %}

## ⚡ Quick Example

Create a `repos.toml` file:

```toml
[[sites]]
remote_prefix = "https://github.com/"
directory = "./projects/"
warm_up_all = true

  [[sites.repos]]
  repo = "golang/go"
  memo = "The Go programming language"

  [[sites.repos]]
  repo = "microsoft/vscode"
  rename = "vscode"
```

Run repoll:

```bash
repoll repos.toml
```

**Result:** Both repositories are cloned, dependencies installed, and ready for development in seconds! 🎉

{% include feature_row id="quality_row" %}

## 🚀 Installation

### macOS/Linux (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/khicago/repoll/main/install.sh | bash
```

### Using Go

```bash
go install github.com/khicago/repoll@latest
```

### Download Binary

Download from [GitHub Releases](https://github.com/khicago/repoll/releases)

### Verify Installation

```bash
repoll version
# repoll version 1.0.0 (commit: abc123, built: 2024-01-01T12:00:00Z)
```

## 🎯 Key Features

### 📝 **Declarative Configuration**
- Simple TOML format that's easy to read and write
- Version-controllable configurations
- Environment-specific setups
- Team-shareable configurations

### 🚀 **Fast and Efficient**
- Parallel repository operations
- Incremental updates for existing repositories
- Smart dependency detection and caching
- Optimized Git operations

### 🔧 **Smart Automation**
- Automatic project type detection (Go, Node.js, Python, Rust, Java)
- Intelligent dependency installation (`go mod download`, `npm install`, etc.)
- Custom naming and directory organization
- Conditional warm-up based on project structure

### 🏢 **Enterprise Features**
- Multiple Git provider support (GitHub, GitLab, Bitbucket, custom)
- SSH and HTTPS authentication
- Detailed execution reporting and logging
- Batch operations across multiple configuration files
- Team onboarding automation

### 🛠 **Developer Experience**
- Cross-platform support (Linux, macOS, Windows)
- Single binary deployment - no dependencies
- Comprehensive CLI with helpful error messages
- Rich configuration validation and helpful suggestions
- Integration-friendly with CI/CD pipelines

## 📊 Quality Assurance

repoll maintains **enterprise-grade quality standards**:

- ✅ **82.8% test coverage** with comprehensive unit and integration tests
- ✅ **Zero critical bugs** in production deployments
- ✅ **Robust error handling** for all edge cases
- ✅ **Performance optimized** for large repository collections
- ✅ **Security reviewed** for safe Git operations

## 🎯 Use Cases

### 👥 **Team Onboarding**
```bash
# New team member setup in one command
repoll team-repos.toml
# ✅ 15 repositories cloned and configured in 2 minutes
```

### 🏗️ **Development Environment Setup**
```bash
# Complete microservices environment
repoll microservices.toml --report
# ✅ All services, databases, and tools ready
```

### 🔄 **Repository Synchronization**
```bash
# Keep all projects up to date
repoll update-all.toml
# ✅ Latest changes pulled across all repositories
```

### 🚀 **CI/CD Integration**
```bash
# Automated testing environment setup
repoll ci-repos.toml --dry-run
# ✅ Verify configuration before deployment
```

## 🌟 What Makes repoll Different?

| Feature | repoll | Manual Git | Other Tools |
|---------|--------|------------|-------------|
| **Configuration** | ✅ Declarative TOML | ❌ Manual scripts | ⚠️ Complex YAML |
| **Performance** | ✅ Parallel operations | ❌ Sequential | ⚠️ Varies |
| **Project Detection** | ✅ Automatic | ❌ Manual setup | ❌ None |
| **Test Coverage** | ✅ 82.8% | N/A | ⚠️ Unknown |
| **Enterprise Ready** | ✅ Production grade | ❌ DIY | ⚠️ Varies |
| **Learning Curve** | ✅ 5 minutes | ❌ Hours | ⚠️ Days |

## 🤝 Community

- 📖 [Documentation](https://khicago.github.io/repoll/)
- 🐛 [Issues & Bug Reports](https://github.com/khicago/repoll/issues)
- 💡 [Feature Requests](https://github.com/khicago/repoll/discussions)
- 🔧 [Contributing Guide](https://github.com/khicago/repoll/blob/main/CONTRIBUTING.md)

## 📈 Getting Started

Ready to transform your repository management workflow?

1. **[Install repoll](getting-started.md#installation)** in 30 seconds
2. **[Create your first config](getting-started.md#your-first-configuration)** in 2 minutes  
3. **[Run and see the magic](getting-started.md#run-repoll)** happen instantly

[Get Started Now →](getting-started.md){: .btn .btn--primary .btn--large}

---

<div class="notice--info">
  <h4>🎯 Pro Tip:</h4>
  <p>Start with <code>repoll mkconf .</code> to automatically generate configuration from your existing repositories!</p>
</div> 