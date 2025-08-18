# go-work

**Auto execute commands across multiple Go modules in workspace with smart path discovery**

go-work is a powerful workspace management tool that auto discovers all Go modules in your workspace and executes commands across them efficiently. Perfect for managing monorepos, multi-module projects, and complex Go workspaces with multiple dependencies.

## Features

- 🔍 **Auto Discovers**: Auto discovers all Go modules in your workspace
- 🎯 **Smart Filtering**: Excludes directories without Go source files by default to avoid empty modules
- 🏗️ **Flexible Options**: Configure current project, submodules, and filtering with fine-grained control
- 🚀 **Efficient Execution**: Runs commands sequentially across modules with proper error handling
- 🐛 **Debug Mode**: Detailed output for troubleshooting including execution paths and command outputs
- ⚡ **Shell Integration**: Supports different shell types (bash, zsh, etc.) with automatic detection
- 🛡️ **Error Resilience**: Continues execution on individual module failures with comprehensive logging

## Install

```bash
go install github.com/go-mate/go-work/cmd/go-work@latest
```

## Usage

### Basic Usage

```bash
# Auto run go mod tidy across all Go modules
cd awesome-path && go-work -c="go mod tidy -e"

# Auto check git status in all modules with debug output
cd awesome-path && go-work -c="git status" --debug

# Auto build all modules
cd awesome-path && go-work -c="go build ./..."

# Auto run tests across all modules
cd awesome-path && go-work -c="go test ./..."

# Run linting across all modules
cd awesome-path && go-work -c="golangci-lint run"
```

### Advanced Examples

```bash
# Auto format code across all modules
go-work -c="go fmt ./..."

# Auto update dependencies across all modules  
go-work -c="go get -u ./..."

# Auto run custom scripts across all modules
go-work -c="./scripts/custom-check.sh"

# Auto clean modules and build cache
go-work -c="go clean -cache -modcache"

# Run security scanning across all modules
go-work -c="gosec ./..."

# Generate coverage reports for all modules
go-work -c="go test -coverprofile=coverage.out ./..."
```

### Debug Mode

```bash
# See detailed execution information including paths and outputs
go-work -c="go mod tidy" --debug

# Debug mode shows:
# - Discovered module paths
# - Command execution details
# - Shell type detection
# - Command outputs and errors
# - Execution flow
```

## How It Works

go-work auto-run:

1. **Discovers** all directories containing `go.mod` files by walking the filesystem
2. **Filters** out directories without Go source files (configurable with `--exclude-no-go`)
3. **Executes** your command in each discovered module directory using the detected shell
4. **Reports** success/failure for each operation with detailed error context
5. **Continues** execution even if individual modules fail, logging all errors
6. **Summarizes** the overall operation status

This automation saves you from manually running commands in each module directory and provides comprehensive workspace management!

---

## Command Line Options

```
Usage:
  go-work [flags]

Flags:
  -c, --command string   command to run in each path (required)
      --debug           enable debug mode for detailed output
  -h, --help            help for go-work
```

## Shell Support

go-work auto-detects your shell environment (`$SHELL`) and supports:
- bash
- zsh  
- fish
- sh
- And other POSIX-compatible shells

## Error Handling

go-work provides robust error handling:
- Individual module failures don't stop the overall execution
- Detailed error context with `erero.Wro()` error wrapping
- Comprehensive logging for troubleshooting
- Debug mode for deep inspection of execution flow

## License

MIT License. See [LICENSE](LICENSE).

---

## Contributing

Contributions are welcome! To contribute:

1. Fork the repo on GitHub (using the webpage interface).
2. Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. Navigate to the cloned project (`cd repo-name`)
4. Create a feature branch (`git checkout -b feature/xxx`).
5. Stage changes (`git add .`)
6. Commit changes (`git commit -m "Add feature xxx"`).
7. Push to the branch (`git push origin feature/xxx`).
8. Open a pull request on GitHub (on the GitHub webpage).

Please ensure tests pass and include relevant documentation updates.

---

## Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

If you find this package valuable, give me some stars on GitHub! Thank you!!!

**Thank you for your support!**

**Happy Coding with this package!** 🎉

Give me stars. Thank you!!!

---

## GitHub Stars

[![starring](https://starchart.cc/go-mate/go-work.svg?variant=adaptive)](https://starchart.cc/go-mate/go-work)
