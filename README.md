[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-mate/go-work/release.yml?branch=main&label=BUILD)](https://github.com/go-mate/go-work/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-mate/go-work)](https://pkg.go.dev/github.com/go-mate/go-work)
[![Coverage Status](https://img.shields.io/coveralls/github/go-mate/go-work/main.svg)](https://coveralls.io/github/go-mate/go-work?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.22--1.25-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-mate/go-work.svg)](https://github.com/go-mate/go-work/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mate/go-work)](https://goreportcard.com/report/github.com/go-mate/go-work)

# go-work

**List Go modules in workspace with smart path detection**

go-work is a workspace management application that auto discovers Go modules in the workspace and lists module paths and versions. Perfect fit with monorepos, multi-module projects, and complex Go workspaces.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->

## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Features

- 🔍 **Auto Discovers**: Auto discovers Go modules in the workspace
- 🎯 **Smart Filtering**: Excludes paths without Go source files
- 🏗️ **Flexible Options**: Configure project and submodules scanning
- 📋 **JSON Output**: Clean JSON format output
- 🏢 **Monorepo Support**: Perfect fit with monorepo architecture

## Installation

```bash
go install github.com/go-mate/go-work/cmd/go-work@latest
```

## Usage

### List Module Paths

```bash
# List all Go modules in current workspace
cd awesome-path && go-work
```

Output:
```json
[
  {
    "path": "/Users/admin/awesome-path",
    "module": "github.com/example/awesome"
  }
]
```

### List Module Versions

```bash
# List Go versions used in each module
cd awesome-path && go-work version
```

Output:
```json
[
  {
    "path": "/Users/admin/awesome-path",
    "module": "github.com/example/awesome",
    "version": "1.22.8"
  }
]
```

## Command Line Options

```
Usage:
  go-work [command]

Available Commands:
  version     List Go versions used in each module
  help        Help about any command

Flags:
  -h, --help  help for go-work
```

## Package Usage

```go
import "github.com/go-mate/go-work/workspath"

// Get project root path
root, ok := workspath.GetProjectRoot("/path/to/sub")

// Get project path with details
info, ok := workspath.GetProjectPath("/path/to/sub")
// info.Root = "/path/to/project"
// info.SubPath = "sub"

// Scan modules with options
paths := workspath.GetModulePaths(
    "/path/to/workspace",
    workspath.WithCurrentProject(),
    workspath.ScanDeep(),
    workspath.SkipNoGo(),
)
```

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 License

MIT License - see [LICENSE](LICENSE).

---

## 💬 Contact & Feedback

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Mistake reports?** Open an issue on GitHub with reproduction steps
- 💡 **Fresh ideas?** Create an issue to discuss
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
7. **Documentation**: Update documentation to support client-facing changes
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

[![Stargazers](https://starchart.cc/go-mate/go-work.svg?variant=adaptive)](https://starchart.cc/go-mate/go-work)
