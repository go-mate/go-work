// Package worksexec: Workspace execution engine for running commands across multiple Go projects
// Provides unified interface for executing commands in both workspace root and subprojects
// Supports batch execution with progress tracking and exception handling
//
// worksexec: 工作区执行引擎，用于在多个 Go 项目中运行命令
// 提供统一接口在工作区根目录和子项目中执行命令
// 支持批量执行，带有进度跟踪和错误处理
package worksexec

import (
	"fmt"

	"github.com/emirpasic/gods/v2/sets/linkedhashset"
	"github.com/go-mate/go-work/workspace"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

// WorksExec provides execution context for running commands across multiple workspaces
// Maintains execution configuration and workspace collection for batch operations
//
// WorksExec 提供在多个工作区中运行命令的执行上下文
// 维护执行配置和工作区集合以进行批量操作
type WorksExec struct {
	execConfig *osexec.ExecConfig     // Base execution configuration // 基础执行配置
	workspaces []*workspace.Workspace // Collection of workspace definitions // 工作区定义集合
}

// NewWorksExec creates a new workspace executor with provided configuration and workspaces
// Validates input parameters and initializes the execution context
//
// 使用提供的配置和工作区创建新的工作区执行器
// 验证输入参数并初始化执行上下文
func NewWorksExec(execConfig *osexec.ExecConfig, workspaces []*workspace.Workspace) (worksExec *WorksExec) {
	return &WorksExec{
		execConfig: must.Nice(execConfig),
		workspaces: must.Have(workspaces),
	}
}

// Subprojects returns all unique subproject paths from all workspaces
// Uses linked hash set to maintain sequence and eliminate duplicates
//
// 返回所有工作区中的唯一子项目路径
// 使用链式哈希集合维护顺序并去除重复项
func (W *WorksExec) Subprojects() []string {
	set := linkedhashset.New[string]()
	for _, wsp := range W.workspaces {
		for _, path := range wsp.Projects {
			set.Add(path)
		}
	}
	return set.Values()
}

// GetWorkspaces returns the collection of managed workspaces
// Provides access to all workspace definitions for iteration
//
// 返回管理的工作区集合
// 提供访问所有工作区定义以进行迭代
func (W *WorksExec) GetWorkspaces() []*workspace.Workspace {
	return W.workspaces
}

// GetNewCommand creates a new execution configuration based on the base config
// Returns a fresh config instance for independent command execution
//
// 基于基础配置创建新的执行配置
// 返回用于独立命令执行的新配置实例
func (W *WorksExec) GetNewCommand() *osexec.ExecConfig {
	return W.execConfig.NewConfig()
}

// GetSubCommand creates a new execution configuration with specific path context
// Useful for executing commands in specific project directories
//
// 创建带有特定路径上下文的新执行配置
// 用于在特定项目目录中执行命令
func (W *WorksExec) GetSubCommand(path string) *osexec.ExecConfig {
	return W.GetNewCommand().WithPath(path)
}

// ForeachWorkRun executes provided function on each workspace root DIR
// Skips workspaces without work root and provides progress tracking
// Only processes workspaces that have a valid WorkRoot configuration
//
// 在每个工作区根目录中执行提供的函数
// 跳过没有工作根目录的工作区并提供进度跟踪
// 只处理具有有效 WorkRoot 配置的工作区
func (W *WorksExec) ForeachWorkRun(run func(execConfig *osexec.ExecConfig, workspace *workspace.Workspace) error) error {
	// Cache expensive method call for performance
	// 缓存昂贵的方法调用以提高性能
	workspaces := W.GetWorkspaces()
	for idx, wsp := range workspaces {
		processMessage := fmt.Sprintf("(%d/%d)", idx, len(workspaces))

		if wsp.WorkRoot != "" {
			zaplog.SUG.Debugln("run", processMessage)

			if err := run(W.GetSubCommand(wsp.WorkRoot), wsp); err != nil {
				return erero.Wro(err)
			}
		} else {
			zaplog.SUG.Debugln("run", processMessage, "no work root so nothing to do")
		}
	}
	return nil
}

// ForeachSubExec executes provided function on each subproject path
// Provides detailed progress tracking with workspace and project indices
// Processes all projects within all workspaces sequentially
//
// 在每个子项目路径中执行提供的函数
// 提供带有工作区和项目索引的详细进度跟踪
// 按顺序处理所有工作区中的所有项目
func (W *WorksExec) ForeachSubExec(run func(execConfig *osexec.ExecConfig, projectPath string) error) error {
	// Cache expensive method call for performance
	// 缓存昂贵的方法调用以提高性能
	workspaces := W.GetWorkspaces()
	for idx, wsp := range workspaces {
		for num, projectPath := range wsp.Projects {
			process1Message := fmt.Sprintf("(%d/%d)", idx, len(workspaces))
			process2Message := fmt.Sprintf("(%d/%d)", num, len(wsp.Projects))

			zaplog.SUG.Debugln("run", process1Message, process2Message, projectPath)

			if err := run(W.GetSubCommand(projectPath), projectPath); err != nil {
				return erero.Wro(err)
			}
		}
	}
	return nil
}
