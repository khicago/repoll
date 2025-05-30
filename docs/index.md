# repoll Documentation

Welcome to the official documentation for **repoll** - the ultimate Git repository management tool.

## Quick Navigation

- [Getting Started](getting-started.md)
- [Configuration Guide](configuration.md)
- [Command Reference](commands.md)
- [Examples](examples.md)
- [API Reference](api.md)
- [Troubleshooting](troubleshooting.md)

## What is repoll?

repoll (Repository Puller) is a lightning-fast, developer-friendly CLI tool that revolutionizes how you manage multiple Git repositories. Whether you're working with microservices, managing open-source contributions, or handling complex multi-repo projects, repoll makes it effortless.

## Key Features

- ⚡ **Lightning Fast**: Concurrent operations with intelligent dependency management
- 🧠 **Smart Warm-up**: Automatically prepares Go, Node.js, and other projects for development
- 📋 **Simple Configuration**: One TOML file to rule them all
- 🔄 **Flexible Workflows**: Clone, update, or sync with customizable strategies
- 📊 **Rich Reporting**: Beautiful progress indicators and detailed execution reports
- 🛡️ **Production Ready**: 82.9% test coverage with robust error handling

## Quick Start

1. Install repoll:
```bash
go install github.com/khicago/repoll@latest
```

2. Create a configuration file (`repos.toml`):
```toml
[[sites]]
    remote_prefix = "https://github.com/"
    dir = "./projects/"
    warm_up_all = true

    [[sites.repos]]
        repo = "microsoft/vscode"
        warm_up = true
```

3. Run repoll:
```bash
repoll repos.toml
```

## Community

- [GitHub Repository](https://github.com/khicago/repoll)
- [Issues & Bug Reports](https://github.com/khicago/repoll/issues)
- [Discussions](https://github.com/khicago/repoll/discussions)

---

**Made with ❤️ by the repoll community** 