// go-work: Workspace management tool for executing commands across Go modules
// Auto discover and execute commands in multiple Go projects within a workspace
// Supports both standalone projects and Go workspace configurations
//
// go-work: Go 工作区管理工具，用于在多个 Go 模块中执行命令
// 自动发现并在工作区内的多个 Go 项目中执行命令
// 支持独立项目和 Go workspace 配置
package main

import (
	"os"

	"github.com/go-mate/go-work/workspath"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	// Get current working DIR for workspace analysis
	// 获取当前工作 DIR 用于工作区分析
	workPath := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(workPath))

	// Define command line parameters
	// 定义命令行参数
	var commandLine string
	var debugMode bool
	rootCmd := cobra.Command{
		Use:   "go-work",
		Short: "go-work",
		Long:  "go-work",
		Run: func(cmd *cobra.Command, args []string) {
			// Detect shell type from environment
			// 从环境中检测 shell 类型
			shellType := must.Nice(os.Getenv("SHELL"))
			if debugMode {
				zaplog.SUG.Debugln("current shell-type:", shellType)
			}
			run(workPath, shellType, commandLine, debugMode)
		},
	}
	rootCmd.Flags().StringVarP(&commandLine, "command", "c", "", "command to run in each path")
	rootCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "enable debug mode")
	must.Done(rootCmd.Execute())
}

// Execute command across all discovered Go modules in workspace
// Auto configure options to include current project and submodules
//
// 在工作区中发现的所有 Go 模块中执行命令
// 自动配置选项以包含当前项目和子模块
func run(workPath string, shellType string, commandLine string, debugMode bool) {
	// Configure module discovery options
	// 配置模块发现选项
	options := workspath.NewOptions().
		WithIncludeCurrentProject(true).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(debugMode)

	// Discover all Go module paths in workspace
	// 发现工作区中的所有 Go 模块路径
	modulePaths := workspath.GetModulePaths(workPath, options)
	if debugMode {
		zaplog.SUG.Debugln("module paths:", neatjsons.S(modulePaths))
	}

	// Execute command in each discovered module
	// 在每个发现的模块中执行命令
	for _, modulePath := range modulePaths {
		zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("--"))
		executeInSinglePath(modulePath, shellType, commandLine, debugMode)
		zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("--"))
	}
	eroticgo.GREEN.ShowMessage("SUCCESS")
}

// Execute command in single module path with specified shell
// Shows execution progress and handles output based on debug mode
//
// 在单个模块路径中使用指定 shell 执行命令
// 显示执行进度并根据调试模式处理输出
func executeInSinglePath(modulePath string, shellType string, commandLine string, debugMode bool) {
	// Build command message for display
	// 构建用于显示的命令消息
	commandMessage := eroticgo.AMBER.Sprint("cd", modulePath, "&&", commandLine)
	if debugMode {
		zaplog.SUG.Debugln("executing:", commandMessage)
	}

	// Execute command with path context
	// 在路径上下文中执行命令
	config := osexec.NewExecConfig().WithPath(modulePath)
	output := rese.V1(config.WithShell(shellType, "-c").Exec(commandLine))

	// Show detailed output in debug mode
	// 在调试模式下显示详细输出
	if debugMode {
		zaplog.SUG.Debugln("executing:", commandMessage, "output:", eroticgo.GREEN.Sprint(string(output)))
	}
	zaplog.SUG.Debugln("executing:", commandMessage, "->:", "success")
}
