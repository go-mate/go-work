// Package workspace: Go workspace structure management in multi-module projects
// Defines workspace configuration and validates project structure
// Supports both classic projects and Go workspace mode with go.work files
//
// workspace: Go 工作区结构管理，用于多模块项目
// 定义工作区配置并验证项目结构完整性
// 支持传统项目和带有 go.work 文件的 Go 工作区模式
package workspace

import (
	"path/filepath"

	"github.com/yyle88/must"
	"github.com/yyle88/osexistpath/osmustexist"
)

// Workspace represents a Go workspace DIR containing multiple subprojects
// Provides structure to manage multiple Go modules within a single workspace
// Supports go.work file configuration when present
//
// Workspace 代表包含多个子项目的 Go 工作区 DIR
// 提供在单个工作区内管理多个 Go 模块的结构
// 可以在有或没有 go.work 文件配置的情况下运行
type Workspace struct {
	WorkRoot string   // Root DIR of the workspace // 工作区根目录
	Projects []string // Project paths within this workspace // 该工作区内的项目路径
}

// NewWorkSpace creates a new workspace without a root DIR
// This is an alias to NewWorkspace with blank workRoot
// Helps when managing standalone projects without go.work file
//
// 创建不带根目录的新工作区
// 这是 NewWorkspace 的别名，使用空 workRoot
// 用于管理没有 go.work 文件的独立项目
func NewWorkSpace(projects []string) *Workspace {
	return NewWorkspace("", projects)
}

// NewWorkspace creates a new workspace with the specified root DIR and projects
// If workRoot is provided, validates the existence of go.work file in root
// Each project path is validated to ensure it contains a go.mod file
// Enforces workspace structure through path validation
//
// 使用指定的根目录和项目创建新的工作区
// 如果提供了 workRoot，则验证根目录中 go.work 文件的存在
// 所有项目路径都会被验证以确保它们包含 go.mod 文件
// 通过路径验证强制执行工作区结构完整性
func NewWorkspace(workRoot string, projects []string) *Workspace {
	// Validate workspace root if provided
	// 如果提供了工作区根目录则进行验证
	if workRoot != "" {
		osmustexist.MustRoot(workRoot)
		osmustexist.MustFile(filepath.Join(workRoot, "go.work"))
	}

	// Validate each project path contains go.mod
	// 验证所有项目路径都包含 go.mod
	for _, projectPath := range must.Have(projects) {
		osmustexist.MustRoot(projectPath)
		osmustexist.MustFile(filepath.Join(projectPath, "go.mod"))
	}

	return &Workspace{
		WorkRoot: workRoot,
		Projects: projects,
	}
}
