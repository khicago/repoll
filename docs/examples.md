---
layout: default
title: Examples
permalink: /examples/
---

# Examples

Here are some practical examples of using repoll in different scenarios.

## Basic Usage

### Simple Repository Management

```toml
# repos.toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./projects/"
    warm_up_all = true

    [[sites.repos]]
        repo = "golang/go"
        warm_up = true

    [[sites.repos]]
        repo = "microsoft/vscode"
        memo = "VS Code editor"
```

Run with:
```bash
repoll repos.toml
```

## Enterprise Setup

### Microservices Development

```toml
# microservices.toml
[[sites]]
    remote_prefix = "https://git.company.com/"
    dir = "./microservices/"
    warm_up_all = true

    [[sites.repos]]
        repo = "team/user-service"
        memo = "User management microservice"

    [[sites.repos]]
        repo = "team/payment-service"
        memo = "Payment processing"

    [[sites.repos]]
        repo = "team/notification-service"
        memo = "Notification system"

    [[sites.repos]]
        repo = "team/api-gateway"
        memo = "Main API gateway"
```

### Multi-Platform Development

```toml
# development.toml
# GitHub repositories
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./github/"
    warm_up_all = true

    [[sites.repos]]
        repo = "golang/go"
        memo = "Go language"

    [[sites.repos]]
        repo = "microsoft/TypeScript"
        memo = "TypeScript language"

# GitLab repositories
[[sites]]
    remote_prefix = "https://gitlab.com/"
    dir = "./gitlab/"

    [[sites.repos]]
        repo = "gitlab-org/gitlab"
        warm_up = true

# Company repositories (SSH)
[[sites]]
    remote_prefix = "git@company.com:"
    dir = "./company/"
    warm_up_all = true

    [[sites.repos]]
        repo = "team/backend-service"
        rename = "backend"

    [[sites.repos]]
        repo = "team/frontend-app"
        rename = "frontend"
```

## Advanced Configuration

### Custom Project Structure

```toml
# advanced.toml
[[sites]]
    remote_prefix = "https://github.com/company/"
    dir = "./work/"
    warm_up_all = true

    [[sites.repos]]
        repo = "user-service"
        rename = "users"
        memo = "User management microservice"

    [[sites.repos]]
        repo = "very-long-repository-name"
        rename = "short-name"
        memo = "Shortened for convenience"

    [[sites.repos]]
        repo = "experimental-project"
        warm_up = false
        memo = "Don't warm up experimental code"
```

### Open Source Contributions

```toml
# opensource.toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./oss-contributions/"

    [[sites.repos]]
        repo = "kubernetes/kubernetes"
        warm_up = true
        memo = "Container orchestration"

    [[sites.repos]]
        repo = "hashicorp/terraform"
        warm_up = true
        memo = "Infrastructure as code"

    [[sites.repos]]
        repo = "prometheus/prometheus"
        warm_up = true
        memo = "Monitoring system"
```

## Command Usage Examples

### Configuration Generation

Generate configuration from existing directories:
```bash
# Scan current directory
repoll mkconf .

# Scan specific directory
repoll mkconf ./my-projects/

# Generate with report
repoll mkconf ./projects/ --report
```

### Repository Management

```bash
# Basic usage
repoll repos.toml

# With detailed reporting
repoll repos.toml --report

# Dry run to see what would happen
repoll --dry-run repos.toml

# Update only (no new clones)
repoll --update-only repos.toml
```

### Multiple Configurations

```bash
# Process multiple config files
repoll config1.toml config2.toml config3.toml

# Use different configurations for different environments
repoll development.toml --report
repoll production.toml --update-only
```

## Real-World Scenarios

### Team Onboarding

Create a configuration file for new team members:

```toml
# team-onboarding.toml
[[sites]]
    remote_prefix = "git@github.com:company/"
    dir = "./company-projects/"
    warm_up_all = true

    [[sites.repos]]
        repo = "main-application"
        memo = "Primary application"

    [[sites.repos]]
        repo = "shared-libraries"
        memo = "Common utilities"

    [[sites.repos]]
        repo = "deployment-scripts"
        memo = "Infrastructure and deployment"

    [[sites.repos]]
        repo = "documentation"
        warm_up = false
        memo = "Project documentation"
```

New team members just run:
```bash
repoll team-onboarding.toml
```

### Project Migration

Moving from one Git provider to another:

```toml
# migration.toml
# Old repositories (for reference)
[[sites]]
    remote_prefix = "https://old-git.company.com/"
    dir = "./old-repos/"

    [[sites.repos]]
        repo = "legacy-project"

# New repositories (active development)
[[sites]]
    remote_prefix = "https://github.com/company/"
    dir = "./new-repos/"
    warm_up_all = true

    [[sites.repos]]
        repo = "migrated-project"
```

### Continuous Integration

Use in CI/CD pipelines:

```bash
# In your CI script
repoll ci-dependencies.toml --update-only
# Run your tests
go test ./...
```

## Tips and Best Practices

1. **Group Related Projects**: Use separate sites for different teams or purposes
2. **Use Meaningful Names**: Add `memo` fields to document repository purposes
3. **Optimize Warm-up**: Only enable warm-up for repositories you actively develop
4. **Version Control Configs**: Keep your configuration files in version control
5. **Environment-Specific Configs**: Use different configurations for dev/staging/prod

## Troubleshooting Examples

### Common Issues

**Permission Denied**:
```bash
# Make sure SSH keys are set up
ssh-add ~/.ssh/id_rsa
repoll repos.toml
```

**Network Issues**:
```bash
# Use HTTPS instead of SSH
# Change from: git@github.com:user/repo
# To: https://github.com/user/repo
```

**Large Repositories**:
```bash
# Disable warm-up for large repos
warm_up = false
```

For more detailed configuration options, see the [Configuration Guide](configuration.md). 