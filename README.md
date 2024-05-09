# Repoll

Repoll (Repo Puller) is a command-line tool that simplifies the process 
of cloning and updating multiple Git repositories based on a 
TOML configuration file.

It also supports optional "warm-up" actions, such as running
`go mod download` for Go projects or `npm install`/`yarn` for 
Node.js projects, ensuring that the cloned repositories are 
ready for use.

## Getting Started

### Prerequisites

- Git must be installed on your machine.
- Go must be installed on your machine (if you plan to build the project from source).
- For warm-up features:
    - Go projects should have `go` installed.
    - Node.js projects should have `npm` or `yarn` installed, depending on the project.

### Installing

To get started with Repo Puller, you can either download the 
pre-built binaries from the [Releases page](#) or compile the 
project from source.

To build from source, clone this repository and run the following 
command from the root directory of the project:

`go build` or `go install`

This will produce an executable called `repoll`.

## Configuration

Create a TOML configuration file (e.g., `config.toml`) that specifies the remote hosting service, directory for local clones, and the repositories to clone or update:

```toml
[[sites]]
    remote="https://code.example.org/"
    dir="./my_repositories/"

    [[sites.repos]]
        repo="user/repo1"
        warm_up=true
    [[sites.repos]]
        repo="user/repo2"
```

## Usage

Run Repo Puller with the path to your configuration file:

`repoll /path/to/config.toml`

This will start cloning or updating repositories as specified in the configuration.

## Features

- Easily clone or update multiple repositories with a single command.
- Support for post-clone "warm-up" actions to prepare repositories for development.
- Real-time progress updates and execution duration.
- Console output with color-coded prefixes and spinner animations for better readability.

## Contributing

If you would like to contribute to the development of Repo Puller, 
please follow these steps:

- Fork the repository.
- Create a new feature branch (git checkout -b feature/amazing-feature).
- Make your changes.
- Commit your changes (git commit -m 'Add an amazing feature').
- Push to the branch (git push origin feature/amazing-feature).
- Open a new pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Repoll is provided "as is", without warranty of any kind. Use it at your own risk.