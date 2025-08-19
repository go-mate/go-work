// Package worksubcmd: Cobra command definitions for Go workspace and module operations
// Provides subcommands for go work sync, go mod tidy, and version management
// Integrates with worksexec for batch execution across multiple projects
//
// worksubcmd: Go 工作区和模块操作的 Cobra 命令定义
// 提供 go work sync、go mod tidy 和版本管理的子命令
// 与 worksexec 集成实现跨多个项目的批量执行
package worksubcmd

import (
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

// NewWorkCmd creates the main 'work' command with sync subcommand
// Provides workspace-level operations like go work sync
// Returns cobra command ready to be added to parent command
//
// 创建带有 sync 子命令的主 'work' 命令
// 提供工作区级别的操作，如 go work sync
// 返回准备好添加到父命令的 cobra 命令
func NewWorkCmd(worksExec *worksexec.WorksExec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "work",
		Short: "go work -->>",
		Long:  "go work -->>",
		Run: func(cmd *cobra.Command, args []string) {
			panic(erero.New("wrong"))
		},
	}

	// Add sync subcommand for workspace synchronization
	// 添加用于工作区同步的 sync 子命令
	cmd.AddCommand(&cobra.Command{
		Use:   "sync",
		Short: "go work sync",
		Long:  "go work sync",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Sync(worksExec))
		},
	})

	return cmd
}

// NewModCmd creates the main 'mod' command with tidy and tide subcommands
// Provides module-level operations for dependency management
// Returns cobra command with tidy and tide subcommands attached
//
// 创建带有 tidy 和 tide 子命令的主 'mod' 命令
// 提供用于依赖管理的模块级别操作
// 返回附加了 tidy 和 tide 子命令的 cobra 命令
func NewModCmd(worksExec *worksexec.WorksExec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mod",
		Short: "go mod -->>",
		Long:  "go mod -->>",
		Run: func(cmd *cobra.Command, args []string) {
			panic(erero.New("wrong"))
		},
	}
	// Add tidy and tide subcommands for module management
	// 添加用于模块管理的 tidy 和 tide 子命令
	cmd.AddCommand(NewTidyCmd(worksExec))
	cmd.AddCommand(NewTideCmd(worksExec))
	return cmd
}

// NewTidyCmd creates the 'tidy' subcommand for standard go mod tidy
// Executes go mod tidy across all subprojects in workspace
// Returns ready-to-use cobra command
//
// 为标准 go mod tidy 创建 'tidy' 子命令
// 在工作区的所有子项目中执行 go mod tidy
// 返回可用的 cobra 命令
func NewTidyCmd(worksExec *worksexec.WorksExec) *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "go mod tidy",
		Long:  "go mod tidy",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Tidy(worksExec))
		},
	}
}

// NewTideCmd creates the 'tide' subcommand for go mod tidy with -e flag
// Executes go mod tidy -e (error-tolerant) across all subprojects
// Returns ready-to-use cobra command for error-tolerant dependency cleanup
//
// 为带 -e 标志的 go mod tidy 创建 'tide' 子命令
// 在所有子项目中执行 go mod tidy -e（错误容忍）
// 返回用于错误容忍依赖清理的可用 cobra 命令
func NewTideCmd(worksExec *worksexec.WorksExec) *cobra.Command {
	return &cobra.Command{
		Use:   "tide",
		Short: "go mod tidy -e",
		Long:  "go mod tidy -e",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Tide(worksExec))
		},
	}
}

// Sync executes 'go work sync' command in each workspace root
// Synchronizes workspace go.work file with actual module structure
// Returns error if any workspace sync operation fails
//
// 在每个工作区根目录中执行 'go work sync' 命令
// 使工作区 go.work 文件与实际模块结构保持同步
// 如果任何工作区同步操作失败则返回错误
func Sync(worksExec *worksexec.WorksExec) error {
	return worksExec.ForeachWorkRun(func(execConfig *osexec.ExecConfig, workspace *workspace.Workspace) error {
		// Execute go work sync command
		// 执行 go work sync 命令
		data, err := execConfig.Exec("go", "work", "sync")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

// Tidy executes 'go mod tidy' command in each subproject DIR
// Cleans up module dependencies and removes unused ones
// Returns error if any module tidy operation fails
//
// 在每个子项目 DIR 中执行 'go mod tidy' 命令
// 清理模块依赖并移除未使用的依赖
// 如果任何模块 tidy 操作失败则返回错误
func Tidy(worksExec *worksexec.WorksExec) error {
	return worksExec.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		// Execute go mod tidy command
		// 执行 go mod tidy 命令
		data, err := execConfig.Exec("go", "mod", "tidy")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

// Tide executes 'go mod tidy -e' command in each subproject DIR
// Error-tolerant version of tidy that continues despite download failures
// Returns error if any critical module operation fails
//
// 在每个子项目 DIR 中执行 'go mod tidy -e' 命令
// 错误容忍版本的 tidy，在下载失败时继续执行
// 如果任何关键模块操作失败则返回错误
func Tide(worksExec *worksexec.WorksExec) error {
	return worksExec.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		// Execute go mod tidy -e command (error-tolerant)
		// 执行 go mod tidy -e 命令（错误容忍）
		data, err := execConfig.Exec("go", "mod", "tidy", "-e")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

// UpdateGoWorkGoVersion updates Go version in workspace go.work files
// Executes 'go work edit -go <version>' in each workspace root
// Returns error if any workspace version update fails
//
// 更新工作区 go.work 文件中的 Go 版本
// 在每个工作区根目录中执行 'go work edit -go <version>'
// 如果任何工作区版本更新失败则返回错误
func UpdateGoWorkGoVersion(worksExec *worksexec.WorksExec, versionNum string) error {
	return worksExec.ForeachWorkRun(func(execConfig *osexec.ExecConfig, workspace *workspace.Workspace) error {
		// Execute go work edit -go command
		// 执行 go work edit -go 命令
		data, err := execConfig.Exec("go", "work", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

// UpdateModuleGoVersion updates Go version in all module go.mod files
// Executes 'go mod edit -go <version>' in each subproject DIR
// Returns error if any module version update fails
//
// 更新所有模块 go.mod 文件中的 Go 版本
// 在每个子项目 DIR 中执行 'go mod edit -go <version>'
// 如果任何模块版本更新失败则返回错误
func UpdateModuleGoVersion(worksExec *worksexec.WorksExec, versionNum string) error {
	return worksExec.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		// Execute go mod edit -go command
		// 执行 go mod edit -go 命令
		data, err := execConfig.Exec("go", "mod", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}
