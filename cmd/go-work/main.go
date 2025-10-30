// Package main: go-work entrypoint
// go-work: Workspace management package that executes commands across Go modules
// Auto detects and executes commands in multiple Go projects within workspace
// Supports both standalone projects and Go workspace configurations
//
// main: go-work 的程序入口
// go-work: Go 工作区管理工具，在多个 Go 模块中执行命令
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
	workPath := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(workPath))

	rootCmd := &cobra.Command{
		Use:   "go-work",
		Short: "Workspace management in Go modules",
		Long:  "go-work: Auto detects and executes commands across multiple Go projects in workspace",
	}

	rootCmd.AddCommand(newExecCommand(workPath))
	must.Done(rootCmd.Execute())
}

// Config encapsulates execution configuration and parameters
// Centralizes settings like command text, debug mode, and paths
//
// Config 封装执行配置和参数
// 集中管理命令文本、调试模式和路径等设置
type Config struct {
	workPath    string // Base path to search in modules // 模块搜索的基础路径
	shellType   string // Runtime shell name used in command execution // 命令执行使用的 shell 类型
	commandLine string // Command text to execute in each module // 在每个模块中执行的命令行
	debugMode   bool   // Enable detailed logging output // 启用详细的日志输出
}

// newExecCommand creates the exec subcommand with workspace execution logic
// Returns configured cobra command that detects modules and runs commands
//
// newExecCommand 创建带有工作区执行逻辑的 exec 子命令
// 返回配置好的 cobra 命令，用于发现模块并执行命令
func newExecCommand(workPath string) *cobra.Command {
	cfg := &Config{workPath: workPath}

	execCmd := &cobra.Command{
		Use:   "exec",
		Short: "Execute command across Go modules in workspace",
		Long:  "Execute specified command in each detected Go module within the workspace",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.shellType = must.Nice(os.Getenv("SHELL"))
			cfg.run()
		},
	}

	execCmd.Flags().StringVarP(&cfg.commandLine, "command", "c", "", "command to run in each module path")
	execCmd.Flags().BoolVarP(&cfg.debugMode, "debug", "", false, "enable debug mode")

	return execCmd
}

// run executes the command across each module in workspace
// Detects modules and runs the command in each detected path
//
// run 在工作区的每个模块中执行命令
// 检测模块并在每个检测到的路径中运行命令
func (cfg *Config) run() {
	cfg.debugLog("current sh-type:", cfg.shellType)

	modulePaths := cfg.detectModules()
	cfg.debugLog("module paths:", neatjsons.S(modulePaths))

	cfg.executes(modulePaths)
	eroticgo.GREEN.ShowMessage("SUCCESS")
}

// detectModules finds each Go module in the workspace
// Returns slice containing paths to detected modules
//
// detectModules 检测工作区中的每个 Go 模块
// 返回包含检测到的模块路径的切片
func (cfg *Config) detectModules() []string {
	options := workspath.NewOptions().
		WithIncludeCurrentProject(true).
		WithIncludeCurrentPackage(false).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(cfg.debugMode)

	return workspath.GetModulePaths(cfg.workPath, options)
}

// executes runs the command in each module path
// Iterates through modules and executes the configured command
//
// executes 在每个模块路径中运行命令
// 遍历模块并执行配置的命令
func (cfg *Config) executes(modulePaths []string) {
	for _, modulePath := range modulePaths {
		zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("--"))
		cfg.execute(modulePath)
		zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("--"))
	}
}

// execute runs the command in a single module path
// Executes command with specified runtime and handles output
//
// execute 在单个模块路径中运行命令
// 使用指定的 shell 执行命令并处理输出
func (cfg *Config) execute(modulePath string) {
	commandMessage := eroticgo.AMBER.Sprint("cd", modulePath, "&&", cfg.commandLine)
	cfg.debugLog("executing:", commandMessage)

	execConfig := osexec.NewExecConfig().WithPath(modulePath)
	output := rese.V1(execConfig.WithShell(cfg.shellType, "-c").Exec(cfg.commandLine))

	if cfg.debugMode {
		zaplog.SUG.Debugln("executing:", commandMessage, "output:", eroticgo.GREEN.Sprint(string(output)))
	}
	zaplog.SUG.Debugln("executing:", commandMessage, "->:", "success")
}

// debugLog outputs log message when debug mode is active
// Reduces repetitive debug mode checks
//
// debugLog 在调试模式激活时输出日志消息
// 辅助函数用于减少重复的调试模式检查
func (cfg *Config) debugLog(args ...interface{}) {
	if cfg.debugMode {
		zaplog.ZAPS.Skip(1).SUG.Debugln(args...)
	}
}
