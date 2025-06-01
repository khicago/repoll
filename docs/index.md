---
layout: splash
title: "repoll - Repository Management Made Simple"
permalink: /
header:
  overlay_color: "#000"
  overlay_filter: "0.5"
  overlay_image: /assets/images/header-bg.jpg
  actions:
    - label: "Get Started"
      url: "#installation"
    - label: "View on GitHub"
      url: "https://github.com/khicago/repoll"
excerpt: "A powerful tool for managing multiple Git repositories with configuration-driven automation. Clone, update, and warm-up your development environment with ease."
intro: 
  - excerpt: 'repoll simplifies the management of multiple Git repositories through declarative TOML configuration files. Whether you are setting up a new development environment, managing enterprise projects, or keeping your repositories in sync, repoll automates the tedious tasks so you can focus on what matters most - your code.'
feature_row:
  - image_path: /assets/images/feature-config.png
    alt: "Configuration-driven"
    title: "Configuration-driven"
    excerpt: "Define your repositories in simple TOML files. No complex scripts or manual processes."
  - image_path: /assets/images/feature-automation.png
    alt: "Smart Automation"
    title: "Smart Automation"
    excerpt: "Automatically clone, update, and warm-up repositories based on project type detection."
  - image_path: /assets/images/feature-enterprise.png
    alt: "Enterprise Ready"
    title: "Enterprise Ready"
    excerpt: "Support for multiple Git providers, custom naming, and team collaboration workflows."
feature_row2:
  - image_path: /assets/images/feature-warmup.png
    alt: "Intelligent Warm-up"
    title: "Intelligent Warm-up"
    excerpt: 'Automatically detect project types (Go, Node.js, Python) and run appropriate setup commands like `go mod download`, `npm install`, or `pip install`.'
    url: "/examples/#warm-up-examples"
    btn_label: "Learn More"
    btn_class: "btn--primary"
feature_row3:
  - image_path: /assets/images/feature-discovery.png
    alt: "Repository Discovery"
    title: "Repository Discovery"
    excerpt: 'Scan existing directories to automatically generate configuration files. Perfect for migrating existing setups or onboarding new team members.'
    url: "/examples/#discovery-examples"
    btn_label: "See Examples"
    btn_class: "btn--primary"
---

{% include feature_row id="intro" type="center" %}

{% include feature_row %}

{% include feature_row id="feature_row2" type="left" %}

{% include feature_row id="feature_row3" type="right" %}

## Installation

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/khicago/repoll/releases):

```bash
# Linux/macOS
curl -L https://github.com/khicago/repoll/releases/latest/download/repoll-linux-amd64 -o repoll
chmod +x repoll
sudo mv repoll /usr/local/bin/

# Or using wget
wget https://github.com/khicago/repoll/releases/latest/download/repoll-linux-amd64 -O repoll
```

### Build from Source

```bash
git clone https://github.com/khicago/repoll.git
cd repoll
go build -o repoll
```

### Verify Installation

```bash
repoll version
```

## Quick Start

### 1. Create a Configuration File

Create a `repos.toml` file:

```toml
[[sites]]
remote = "https://github.com/"
dir = "./projects/"

[[sites.repos]]
repo = "golang/go"
warm_up = true

[[sites.repos]]
repo = "kubernetes/kubernetes"
rename = "k8s"
memo = "Kubernetes main repository"
```

### 2. Run repoll

```bash
# Clone/update repositories
repoll repos.toml

# Generate detailed report
repoll repos.toml --report
```

### 3. Generate Configuration from Existing Repositories

```bash
# Scan current directory
repoll mkconf .

# Scan specific directory
repoll mkconf ./my-projects/
```

## Key Features

### 🚀 **Fast and Efficient**
- Parallel repository operations
- Incremental updates for existing repositories
- Smart caching and optimization

### 📝 **Declarative Configuration**
- Simple TOML format
- Version-controllable configurations
- Environment-specific setups

### 🔧 **Smart Automation**
- Automatic project type detection
- Intelligent warm-up processes
- Dependency installation

### 🏢 **Enterprise Features**
- Multiple Git provider support
- Custom repository naming
- Team collaboration workflows
- Detailed reporting and logging

### 🛠 **Developer Friendly**
- Cross-platform support (Linux, macOS, Windows)
- Single binary deployment
- Comprehensive CLI interface
- Rich API for integration

## Configuration Format

repoll uses TOML configuration files to define repository management rules:

```toml
# Global settings can be defined here

[[sites]]
remote = "https://github.com/"  # Git provider URL prefix
dir = "./projects/"             # Local base directory
warm_up_all = false            # Global warm-up setting

  [[sites.repos]]
  repo = "owner/repository"     # Repository path
  rename = "custom-name"        # Optional: custom local name
  warm_up = true               # Optional: enable warm-up
  memo = "Description"         # Optional: documentation

[[sites]]
remote = "https://gitlab.com/"
dir = "./gitlab-projects/"

  [[sites.repos]]
  repo = "group/project"
```

### Site Configuration

- **`remote`**: Git provider URL prefix (e.g., `https://github.com/`, `https://gitlab.com/`)
- **`dir`**: Local base directory where repositories will be cloned
- **`warm_up_all`**: Global setting to enable warm-up for all repositories in this site

### Repository Configuration

- **`repo`**: Repository path (e.g., `owner/repository`)
- **`rename`**: Optional custom local directory name
- **`warm_up`**: Enable automatic dependency installation and setup
- **`memo`**: Optional description or notes

## Use Cases

### Development Environment Setup
Quickly set up a complete development environment with all necessary repositories.

### Team Onboarding
New team members can get up and running with a single command.

### CI/CD Integration
Automate repository management in continuous integration pipelines.

### Project Migration
Easily migrate projects between different Git providers or directory structures.

### Dependency Management
Keep all project dependencies up-to-date across multiple repositories.

## Next Steps

- [View Examples](examples.md) - See practical usage examples
- [API Reference](api.md) - Detailed API documentation
- [GitHub Repository](https://github.com/khicago/repoll) - Source code and issues

---

*repoll is open source software released under the MIT License.* 