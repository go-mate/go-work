[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-mate/go-work/release.yml?branch=main&label=BUILD)](https://github.com/go-mate/go-work/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-mate/go-work)](https://pkg.go.dev/github.com/go-mate/go-work)
[![Coverage Status](https://img.shields.io/coveralls/github/go-mate/go-work/main.svg)](https://coveralls.io/github/go-mate/go-work?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.22--1.25-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-mate/go-work.svg)](https://github.com/go-mate/go-work/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mate/go-work)](https://goreportcard.com/report/github.com/go-mate/go-work)

# go-work

**通过智能路径发现，列举工作区中的 Go 模块**

go-work 是一个工作区管理应用，它能自动发现工作区中的 Go 模块，并列举模块路径和版本。这是管理单体仓库、多模块项目以及复杂 Go 工作区的绝佳方案。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 功能

- 🔍 **自动发现**: 自动发现工作区中的 Go 模块
- 🎯 **智能过滤**: 排除不含 Go 源文件的路径
- 🏗️ **灵活选项**: 配置项目和子模块扫描
- 📋 **JSON 输出**: 清晰的 JSON 格式输出
- 🏢 **Monorepo 支持**: 完美适配 monorepo 架构

## 安装方式

```bash
go install github.com/go-mate/go-work/cmd/go-work@latest
```

## 用法

### 列举模块路径

```bash
# 列举当前工作区中的所有 Go 模块
cd awesome-path && go-work
```

输出:
```json
[
  {
    "path": "/Users/admin/awesome-path",
    "module": "github.com/example/awesome"
  }
]
```

### 列举模块版本

```bash
# 列举每个模块使用的 Go 版本
cd awesome-path && go-work version
```

输出:
```json
[
  {
    "path": "/Users/admin/awesome-path",
    "module": "github.com/example/awesome",
    "version": "1.22.8"
  }
]
```

## 命令行选项

```
用法:
  go-work [command]

可用命令:
  version     列举每个模块使用的 Go 版本
  help        关于任何命令的帮助

标志:
  -h, --help  go-work 的帮助信息
```

## 包用法

```go
import "github.com/go-mate/go-work/workspath"

// 获取项目根路径
root, ok := workspath.GetProjectRoot("/path/to/sub")

// 获取项目路径详情
info, ok := workspath.GetProjectPath("/path/to/sub")
// info.Root = "/path/to/project"
// info.SubPath = "sub"

// 使用选项扫描模块
paths := workspath.GetModulePaths(
    "/path/to/workspace",
    workspath.WithCurrentProject(),
    workspath.ScanDeep(),
    workspath.SkipNoGo(),
)
```

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 许可证类型

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

## 💬 联系与反馈

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **问题报告？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **新颖思路？** 创建 issue 讨论
- 📖 **文档疑惑？** 报告问题，帮助我们完善文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，协助解决性能问题
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：面向用户的更改需要更新文档
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来贡献此项目。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![标星点赞](https://starchart.cc/go-mate/go-work.svg?variant=adaptive)](https://starchart.cc/go-mate/go-work)
