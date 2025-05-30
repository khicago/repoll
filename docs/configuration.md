# Configuration Guide

repoll uses TOML configuration files to define repositories, directories, and behavior. This guide covers all configuration options in detail.

## Basic Structure

```toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./projects/"
    warm_up_all = false

    [[sites.repos]]
        repo = "owner/repository"
        warm_up = true
        rename = "custom-name"
        memo = "Description"
```

## Site Configuration

The `[[sites]]` section defines a hosting provider and local directory configuration.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `remote_prefix` | string | ✅ | URL prefix for repositories (e.g., `https://github.com/`) |
| `dir` | string | ✅ | Local directory for cloning repositories |
| `warm_up_all` | boolean | ❌ | Enable warm-up for all repositories in this site |

### Examples

#### GitHub
```toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./github-projects/"
```

#### GitLab
```toml
[[sites]]
    remote_prefix = "https://gitlab.com/"
    dir = "./gitlab-projects/"
```

#### Private Git Server
```toml
[[sites]]
    remote_prefix = "https://git.company.com/"
    dir = "./company-projects/"
```

#### SSH Access
```toml
[[sites]]
    remote_prefix = "git@github.com:"
    dir = "./ssh-projects/"
```

## Repository Configuration

The `[[sites.repos]]` section defines individual repositories within a site.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `repo` | string | ✅ | Repository path (e.g., `owner/repository`) |
| `warm_up` | boolean | ❌ | Enable warm-up for this repository |
| `rename` | string | ❌ | Custom directory name (defaults to repository name) |
| `memo` | string | ❌ | Description or note about the repository |

### Examples

#### Basic Repository
```toml
[[sites.repos]]
    repo = "golang/go"
```

#### Repository with Warm-up
```toml
[[sites.repos]]
    repo = "microsoft/vscode"
    warm_up = true
    memo = "VS Code development"
```

#### Repository with Custom Name
```toml
[[sites.repos]]
    repo = "very-long-repository-name/with-complex-structure"
    rename = "simple-name"
```

## Warm-up Features

Warm-up automatically prepares projects for development by running appropriate setup commands.

### Supported Project Types

| Project Type | Detection | Commands |
|--------------|-----------|----------|
| **Go** | `go.mod` file | `go mod download` |
| **Node.js** | `package.json` file | `npm install` or `yarn install` |
| **Python** | `requirements.txt` | `pip install -r requirements.txt` |
| **Rust** | `Cargo.toml` | `cargo fetch` |

### Warm-up Configuration

#### Site-level Warm-up
Enable warm-up for all repositories in a site:
```toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./projects/"
    warm_up_all = true  # Enables warm-up for all repos
```

#### Repository-level Warm-up
Enable warm-up for specific repositories:
```toml
[[sites.repos]]
    repo = "golang/go"
    warm_up = true  # Overrides site setting
```

## Complete Examples

### Multi-Platform Development
```toml
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

    [[sites.repos]]
        repo = "rust-lang/rust"
        memo = "Rust language"

# GitLab repositories
[[sites]]
    remote_prefix = "https://gitlab.com/"
    dir = "./gitlab/"

    [[sites.repos]]
        repo = "gitlab-org/gitlab"
        warm_up = true

# Company repositories
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

### Microservices Setup
```toml
[[sites]]
    remote_prefix = "https://github.com/company/"
    dir = "./microservices/"
    warm_up_all = true

    [[sites.repos]]
        repo = "user-service"
        memo = "User management microservice"

    [[sites.repos]]
        repo = "payment-service"
        memo = "Payment processing"

    [[sites.repos]]
        repo = "notification-service"
        memo = "Notification system"

    [[sites.repos]]
        repo = "api-gateway"
        memo = "Main API gateway"
```

### Open Source Contributions
```toml
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

## Best Practices

### Directory Structure
- Use descriptive directory names: `./microservices/`, `./tools/`, `./experiments/`
- Group related projects: separate work, personal, and open-source projects
- Use relative paths for portability

### Repository Organization
- Add meaningful `memo` descriptions for easy identification
- Use `rename` for long or unclear repository names
- Group repositories by team, purpose, or technology

### Warm-up Strategy
- Enable `warm_up_all` for development environments
- Disable warm-up for large repositories that don't need immediate setup
- Use repository-level warm-up for fine-grained control

### Performance Tips
- Limit the number of repositories per site to avoid overwhelming the system
- Use SSH keys for private repositories to avoid authentication prompts
- Consider network bandwidth when cloning large repositories

## Configuration Generation

### From Existing Directories
Generate configuration from existing Git repositories:
```bash
repoll mkconf ./my-projects/
```

This scans the directory and creates a configuration file automatically.

### Manual Configuration
1. Start with a simple configuration
2. Add repositories incrementally
3. Test with a few repositories before scaling up
4. Use `--dry-run` to preview changes

## Troubleshooting

### Common Issues

**Repository Not Found**
- Verify the `remote_prefix` and `repo` combination
- Check repository permissions and access rights

**Warm-up Failures**
- Ensure required tools are installed (Go, Node.js, etc.)
- Check project dependencies and requirements

**Permission Denied**
- Configure SSH keys for private repositories
- Verify Git credentials and access tokens

### Validation
Test your configuration:
```bash
# Dry run to see what would happen
repoll --dry-run config.toml

# Verbose output for debugging
repoll --verbose config.toml
```

## Advanced Configuration

### Environment Variables
Use environment variables in configuration:
```toml
[[sites]]
    remote_prefix = "${GITHUB_URL}"
    dir = "${PROJECT_DIR}/github/"
```

### Conditional Configuration
Different configurations for different environments:
```toml
# development.toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./dev/"
    warm_up_all = true

# production.toml
[[sites]]
    remote_prefix = "git@github.com:"
    dir = "/opt/repos/"
    warm_up_all = false
``` 