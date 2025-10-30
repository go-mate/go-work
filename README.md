[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-mate/go-work/release.yml?branch=main&label=BUILD)](https://github.com/go-mate/go-work/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-mate/go-work)](https://pkg.go.dev/github.com/go-mate/go-work)
[![Coverage Status](https://img.shields.io/coveralls/github/go-mate/go-work/main.svg)](https://coveralls.io/github/go-mate/go-work?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.22--1.25-lightgrey.svg)](https://github.com/go-mate/go-work)
[![GitHub Release](https://img.shields.io/github/release/go-mate/go-work.svg)](https://github.com/go-mate/go-work/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mate/go-work)](https://goreportcard.com/report/github.com/go-mate/go-work)

# go-work

**Auto execute commands across multiple Go modules in workspace with smart path detection**

go-work is a workspace management package that auto discovers Go modules in the workspace and executes commands across them. Perfect package to manage monorepos, multi-module projects, and complex Go workspaces with multiple dependencies.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Features

- 🔍 **Auto Discovers**: Auto discovers Go modules in the workspace
- 🎯 **Smart Filtering**: Excludes paths without Go source files
- 🏗️ **Flexible Options**: Configure project and submodules
- ⚡ **Batch Execution**: Execute commands across multiple modules
- 🏢 **Monorepo Support**: Perfect fit with monorepo architecture

## Installation

```bash
go install github.com/go-mate/go-work/cmd/go-work@latest
```

## Usage

### Basic Usage

```bash
# Auto run go mod tidy across Go modules
cd awesome-path && go-work exec -c="go mod tidy -e"

# Auto check git status in each module with debug output
cd awesome-path && go-work exec -c="git status" --debug

# Auto build each module
cd awesome-path && go-work exec -c="go build ./..."

# Auto run tests across modules
cd awesome-path && go-work exec -c="go test ./..."

# Run linting across modules
cd awesome-path && go-work exec -c="golangci-lint run"
```

## Command Line Options

```
Usage:
  go-work exec [flags]

Flags:
  -c, --command string   command to run in each module path
      --debug            enable debug mode
  -h, --help             show help message
```

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![starring](https://starchart.cc/go-mate/go-work.svg?variant=adaptive)](https://starchart.cc/go-mate/go-work)
