---
layout: default
title: Examples
permalink: /examples/
---

# Examples

Real-world configuration examples for different use cases and development scenarios.

## 🎯 Quick Start Examples

### Personal Open Source Setup

```toml
# Basic setup for open source contributions
[[sites]]
remote = "https://github.com/"
dir = "./oss/"
warm_up_all = true

  [[sites.repos]]
  repo = "golang/go"
  memo = "Go language source"

  [[sites.repos]]
  repo = "microsoft/vscode"
  memo = "VS Code editor"

  [[sites.repos]]
  repo = "kubernetes/kubernetes"
  rename = "k8s"
  memo = "Container orchestration"
```

**Usage:**
```bash
repoll personal-oss.toml
```

**Result:**
```
✓ Cloned golang/go -> ./oss/go
🔥 Warmed up ./oss/go (Go project)
✓ Cloned microsoft/vscode -> ./oss/vscode
🔥 Warmed up ./oss/vscode (Node.js project)
✓ Cloned kubernetes/kubernetes -> ./oss/k8s
🔥 Warmed up ./oss/k8s (Go project)
```

### Team Development Environment

```toml
# Company projects setup
[[sites]]
remote = "git@company.com:"
dir = "./company/"
warm_up_all = true

  [[sites.repos]]
  repo = "frontend/web-app"
  rename = "webapp"
  memo = "Main web application"

  [[sites.repos]]
  repo = "backend/api-service"
  rename = "api"
  memo = "REST API service"

  [[sites.repos]]
  repo = "infrastructure/k8s-configs"
  rename = "k8s"
  warm_up = false
  memo = "Kubernetes configurations"
```

## 🏢 Enterprise Scenarios

### Microservices Architecture

```toml
# Backend services
[[sites]]
remote = "https://gitlab.company.com/"
dir = "./services/"
warm_up_all = true

  [[sites.repos]]
  repo = "platform/user-service"
  memo = "User authentication and management"

  [[sites.repos]]
  repo = "platform/payment-service"
  memo = "Payment processing"

  [[sites.repos]]
  repo = "platform/notification-service"
  memo = "Email and SMS notifications"

  [[sites.repos]]
  repo = "platform/api-gateway"
  memo = "Main API gateway"

# Frontend applications
[[sites]]
remote = "https://gitlab.company.com/"
dir = "./frontend/"
warm_up_all = true

  [[sites.repos]]
  repo = "ui/customer-portal"
  rename = "portal"
  memo = "Customer-facing web app"

  [[sites.repos]]
  repo = "ui/admin-dashboard"
  rename = "admin"
  memo = "Internal admin interface"

# Shared libraries
[[sites]]
remote = "https://gitlab.company.com/"
dir = "./shared/"
warm_up_all = true

  [[sites.repos]]
  repo = "shared/common-utils"
  rename = "utils"
  memo = "Shared utility functions"

  [[sites.repos]]
  repo = "shared/ui-components"
  rename = "components"
  memo = "Reusable UI components"
```

### Multi-Platform Development

```toml
# GitHub dependencies
[[sites]]
remote = "https://github.com/"
dir = "./deps/"
warm_up_all = true

  [[sites.repos]]
  repo = "gin-gonic/gin"
  memo = "Go HTTP framework"

  [[sites.repos]]
  repo = "facebook/react"
  memo = "UI library"

  [[sites.repos]]
  repo = "microsoft/TypeScript"
  memo = "TypeScript compiler"

# Internal GitLab
[[sites]]
remote = "https://gitlab.internal.com/"
dir = "./internal/"
warm_up_all = true

  [[sites.repos]]
  repo = "core/authentication"
  rename = "auth"
  memo = "Internal auth service"

  [[sites.repos]]
  repo = "core/logging"
  memo = "Centralized logging"

# Legacy Bitbucket
[[sites]]
remote = "https://bitbucket.org/legacy-team/"
dir = "./legacy/"
warm_up_all = false

  [[sites.repos]]
  repo = "old-system"
  warm_up = false
  memo = "Legacy system - manual setup"
```

## 🔧 Technology-Specific Setups

### Go Development Environment

```toml
# Go ecosystem projects
[[sites]]
remote = "https://github.com/"
dir = "./go-projects/"
warm_up_all = true

  [[sites.repos]]
  repo = "golang/go"
  memo = "Go standard library"

  [[sites.repos]]
  repo = "spf13/cobra"
  memo = "CLI library"

  [[sites.repos]]
  repo = "spf13/viper"
  memo = "Configuration management"

  [[sites.repos]]
  repo = "gorilla/mux"
  memo = "HTTP router"

  [[sites.repos]]
  repo = "stretchr/testify"
  memo = "Testing toolkit"

  [[sites.repos]]
  repo = "uber-go/zap"
  memo = "Structured logging"

# Personal Go projects
[[sites]]
remote = "https://github.com/yourusername/"
dir = "./my-go/"
warm_up_all = true

  [[sites.repos]]
  repo = "awesome-cli-tool"
  memo = "My CLI application"

  [[sites.repos]]
  repo = "microservice-template"
  memo = "Go microservice boilerplate"
```

### Node.js Full-Stack Setup

```toml
# Frontend frameworks
[[sites]]
remote = "https://github.com/"
dir = "./frontend/"
warm_up_all = true

  [[sites.repos]]
  repo = "facebook/react"
  memo = "React library"

  [[sites.repos]]
  repo = "vuejs/vue"
  memo = "Vue.js framework"

  [[sites.repos]]
  repo = "angular/angular"
  memo = "Angular framework"

  [[sites.repos]]
  repo = "vercel/next.js"
  rename = "nextjs"
  memo = "React framework"

# Backend frameworks
[[sites]]
remote = "https://github.com/"
dir = "./backend/"
warm_up_all = true

  [[sites.repos]]
  repo = "expressjs/express"
  memo = "Express.js framework"

  [[sites.repos]]
  repo = "fastify/fastify"
  memo = "Fast web framework"

  [[sites.repos]]
  repo = "nestjs/nest"
  memo = "NestJS framework"

  [[sites.repos]]
  repo = "koajs/koa"
  memo = "Koa.js framework"

# Tools and utilities
[[sites]]
remote = "https://github.com/"
dir = "./tools/"
warm_up_all = true

  [[sites.repos]]
  repo = "eslint/eslint"
  memo = "JavaScript linter"

  [[sites.repos]]
  repo = "prettier/prettier"
  memo = "Code formatter"

  [[sites.repos]]
  repo = "webpack/webpack"
  memo = "Module bundler"
```

### Python Data Science Setup

```toml
# Core Python projects
[[sites]]
remote = "https://github.com/"
dir = "./python-core/"
warm_up_all = true

  [[sites.repos]]
  repo = "python/cpython"
  memo = "Python interpreter"

  [[sites.repos]]
  repo = "pypa/pip"
  memo = "Package installer"

  [[sites.repos]]
  repo = "psf/requests"
  memo = "HTTP library"

# Data science frameworks
[[sites]]
remote = "https://github.com/"
dir = "./data-science/"
warm_up_all = true

  [[sites.repos]]
  repo = "numpy/numpy"
  memo = "Numerical computing"

  [[sites.repos]]
  repo = "pandas-dev/pandas"
  memo = "Data analysis"

  [[sites.repos]]
  repo = "matplotlib/matplotlib"
  memo = "Plotting library"

  [[sites.repos]]
  repo = "scikit-learn/scikit-learn"
  rename = "sklearn"
  memo = "Machine learning"

  [[sites.repos]]
  repo = "tensorflow/tensorflow"
  memo = "Deep learning"

# Web frameworks
[[sites]]
remote = "https://github.com/"
dir = "./python-web/"
warm_up_all = true

  [[sites.repos]]
  repo = "django/django"
  memo = "Web framework"

  [[sites.repos]]
  repo = "pallets/flask"
  memo = "Micro framework"

  [[sites.repos]]
  repo = "tiangolo/fastapi"
  memo = "Modern API framework"
```

## 🎨 Advanced Configuration Patterns

### Environment-Based Setup

**development.toml:**
```toml
# Development environment
[[sites]]
remote = "https://github.com/"
dir = "./dev/"
warm_up_all = true

  [[sites.repos]]
  repo = "company/main-app"
  memo = "Main application - dev branch"

[[sites]]
remote = "https://github.com/"
dir = "./experiments/"
warm_up_all = false

  [[sites.repos]]
  repo = "experimental/new-feature"
  warm_up = true
  memo = "Feature development"
```

**staging.toml:**
```toml
# Staging environment
[[sites]]
remote = "https://github.com/"
dir = "./staging/"
warm_up_all = false

  [[sites.repos]]
  repo = "company/main-app"
  warm_up = false
  memo = "Staging deployment"

  [[sites.repos]]
  repo = "company/config-staging"
  rename = "config"
  memo = "Staging configurations"
```

**production.toml:**
```toml
# Production environment
[[sites]]
remote = "https://github.com/"
dir = "./prod/"
warm_up_all = false

  [[sites.repos]]
  repo = "company/main-app"
  warm_up = false
  memo = "Production code"

  [[sites.repos]]
  repo = "company/config-prod"
  rename = "config"
  memo = "Production configurations"

  [[sites.repos]]
  repo = "company/monitoring"
  memo = "Production monitoring"
```

### Team-Based Organization

**frontend-team.toml:**
```toml
# Frontend team repositories
[[sites]]
remote = "https://github.com/company/"
dir = "./frontend/"
warm_up_all = true

  [[sites.repos]]
  repo = "web-app"
  memo = "Main web application"

  [[sites.repos]]
  repo = "mobile-app"
  memo = "React Native mobile app"

  [[sites.repos]]
  repo = "design-system"
  memo = "Shared UI components"

# External dependencies
[[sites]]
remote = "https://github.com/"
dir = "./deps/"
warm_up_all = true

  [[sites.repos]]
  repo = "facebook/react"
  memo = "React library"

  [[sites.repos]]
  repo = "vercel/next.js"
  rename = "nextjs"
  memo = "Next.js framework"
```

**backend-team.toml:**
```toml
# Backend team repositories
[[sites]]
remote = "https://github.com/company/"
dir = "./services/"
warm_up_all = true

  [[sites.repos]]
  repo = "user-service"
  memo = "User management API"

  [[sites.repos]]
  repo = "order-service"
  memo = "Order processing API"

  [[sites.repos]]
  repo = "payment-service"
  memo = "Payment gateway API"

# Infrastructure
[[sites]]
remote = "https://github.com/company/"
dir = "./infra/"
warm_up_all = false

  [[sites.repos]]
  repo = "docker-configs"
  memo = "Docker configurations"

  [[sites.repos]]
  repo = "k8s-manifests"
  memo = "Kubernetes manifests"
```

**devops-team.toml:**
```toml
# Infrastructure and deployment
[[sites]]
remote = "https://github.com/company/"
dir = "./infrastructure/"
warm_up_all = false

  [[sites.repos]]
  repo = "terraform-modules"
  rename = "terraform"
  memo = "Infrastructure as code"

  [[sites.repos]]
  repo = "ansible-playbooks"
  rename = "ansible"
  memo = "Configuration management"

  [[sites.repos]]
  repo = "helm-charts"
  rename = "helm"
  memo = "Kubernetes package manager"

# Monitoring and observability
[[sites]]
remote = "https://github.com/"
dir = "./monitoring/"
warm_up_all = true

  [[sites.repos]]
  repo = "prometheus/prometheus"
  memo = "Monitoring system"

  [[sites.repos]]
  repo = "grafana/grafana"
  memo = "Visualization platform"

  [[sites.repos]]
  repo = "jaegertracing/jaeger"
  memo = "Distributed tracing"
```

## 🔄 Dynamic Configuration Generation

### Generate from Existing Projects

```bash
# Generate configuration from current directory
repoll mkconf . > repos.toml

# Generate from specific directory
repoll mkconf ~/projects > my-projects.toml

# Generate from multiple directories
repoll mkconf ./frontend ./backend > full-stack.toml
```

### Example Generated Configuration

When running `repoll mkconf` on a directory with existing repositories:

```toml
# Generated by repoll mkconf on 2024-01-15T10:30:00Z

[[sites]]
remote = "https://github.com/"
dir = "./"

  [[sites.repos]]
  repo = "golang/go"
  memo = "Auto-detected from ./go/"

  [[sites.repos]]
  repo = "microsoft/vscode"
  memo = "Auto-detected from ./vscode/"

[[sites]]
remote = "git@gitlab.com:"
dir = "./"

  [[sites.repos]]
  repo = "company/internal-tool"
  memo = "Auto-detected from ./internal-tool/"
```

## 📊 Complex Multi-Project Setup

### Full Development Environment

```toml
# Core development tools
[[sites]]
remote = "https://github.com/"
dir = "./tools/"
warm_up_all = true

  [[sites.repos]]
  repo = "git/git"
  memo = "Git version control"

  [[sites.repos]]
  repo = "vim/vim"
  memo = "Text editor"

  [[sites.repos]]
  repo = "BurntSushi/ripgrep"
  memo = "Fast text search"

# Programming languages
[[sites]]
remote = "https://github.com/"
dir = "./languages/"
warm_up_all = false

  [[sites.repos]]
  repo = "golang/go"
  warm_up = true
  memo = "Go programming language"

  [[sites.repos]]
  repo = "rust-lang/rust"
  warm_up = false
  memo = "Rust programming language"

  [[sites.repos]]
  repo = "nodejs/node"
  warm_up = false
  memo = "Node.js runtime"

# Frameworks and libraries
[[sites]]
remote = "https://github.com/"
dir = "./frameworks/"
warm_up_all = true

  # Go frameworks
  [[sites.repos]]
  repo = "gin-gonic/gin"
  memo = "Go HTTP framework"

  [[sites.repos]]
  repo = "gorilla/mux"
  memo = "Go HTTP router"

  # JavaScript frameworks
  [[sites.repos]]
  repo = "facebook/react"
  memo = "React UI library"

  [[sites.repos]]
  repo = "vuejs/vue"
  memo = "Vue.js framework"

  # CSS frameworks
  [[sites.repos]]
  repo = "twbs/bootstrap"
  memo = "CSS framework"

  [[sites.repos]]
  repo = "tailwindlabs/tailwindcss"
  memo = "Utility-first CSS"

# DevOps and infrastructure
[[sites]]
remote = "https://github.com/"
dir = "./devops/"
warm_up_all = false

  [[sites.repos]]
  repo = "kubernetes/kubernetes"
  rename = "k8s"
  memo = "Container orchestration"

  [[sites.repos]]
  repo = "docker/cli"
  memo = "Docker command line"

  [[sites.repos]]
  repo = "hashicorp/terraform"
  memo = "Infrastructure as code"

  [[sites.repos]]
  repo = "ansible/ansible"
  memo = "Configuration management"

# Personal projects
[[sites]]
remote = "https://github.com/yourusername/"
dir = "./my-projects/"
warm_up_all = true

  [[sites.repos]]
  repo = "portfolio-website"
  memo = "Personal website"

  [[sites.repos]]
  repo = "awesome-tool"
  memo = "My CLI tool"

  [[sites.repos]]
  repo = "learning-rust"
  memo = "Rust learning project"

# Company projects
[[sites]]
remote = "git@company.com:"
dir = "./work/"
warm_up_all = true

  [[sites.repos]]
  repo = "main-product"
  memo = "Company's main product"

  [[sites.repos]]
  repo = "internal-tools"
  memo = "Internal development tools"

  [[sites.repos]]
  repo = "documentation"
  memo = "Technical documentation"
```

## 🚀 Usage Examples

### Basic Operations

```bash
# Process all configurations
repoll repos.toml

# Preview what will be done
repoll --dry-run repos.toml

# Verbose output for debugging
repoll --verbose repos.toml

# Process multiple configurations
repoll dev.toml staging.toml prod.toml
```

### Team Workflow Examples

```bash
# Frontend team daily setup
repoll frontend-team.toml

# Backend team environment
repoll backend-team.toml

# Full stack development
repoll frontend-team.toml backend-team.toml

# DevOps environment setup
repoll devops-team.toml
```

### Environment Management

```bash
# Set up development environment
repoll development.toml

# Deploy to staging
repoll staging.toml

# Prepare production
repoll production.toml

# Multi-environment setup
repoll development.toml staging.toml production.toml
```

## 🔧 Pro Tips

### 1. Organize by Purpose

```bash
# Create purpose-specific configurations
learning.toml          # Educational repositories
contributions.toml     # Open source contributions
experiments.toml       # Experimental projects
production.toml        # Production code
```

### 2. Use Descriptive Memos

```toml
[[sites.repos]]
repo = "complex/repository-name"
memo = "Critical service - handles user authentication and session management"
```

### 3. Strategic Warm-up Usage

```toml
# Enable warm-up for active development
[[sites.repos]]
repo = "active/project"
warm_up = true
memo = "Currently working on this"

# Disable for large or reference repositories
[[sites.repos]]
repo = "linux/kernel"
warm_up = false
memo = "Reference only - too large for auto-setup"
```

### 4. Environment Variables

```bash
# Use environment-specific configurations
export REPOLL_ENV=development
repoll ${REPOLL_ENV}.toml
```

### 5. Batch Processing

```bash
# Process team configurations in parallel
repoll team-*.toml

# Environment-specific processing
repoll *-${ENVIRONMENT}.toml
```

## 📞 Next Steps

- **[Configuration Guide](configuration.md)** - Detailed configuration options
- **[API Reference](api.md)** - Complete command reference
- **[Getting Started](getting-started.md)** - Basic usage tutorial

---

<div class="notice--info">
  <h4>💡 Pro Tip:</h4>
  <p>Start with a simple configuration and gradually add more repositories. Use <code>repoll mkconf</code> to generate initial configurations from existing projects, then customize as needed!</p>
</div> 