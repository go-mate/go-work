// demo1x: Demonstration application showing go-work package usage
// Shows how to build CLI tools using workspace execution capabilities
// Provides examples of mod tidy, tide, and work sync operations
//
// demo1x: 演示应用程序，展示 go-work 包的使用
// 展示如何使用工作区执行功能构建 CLI 工具
// 提供 mod tidy、tide 和 work sync 操作的示例
package main

import (
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/go-mate/go-work/worksubcmd"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

// Usage examples:
// go run main.go help
// go run main.go mod tidy
// go run main.go mod tide
// go run main.go work sync
//
// 使用示例：
// go run main.go help
// go run main.go mod tidy
// go run main.go mod tide
// go run main.go work sync
func main() {
	// Get project path by going up 3 levels from current location
	// 通过从当前位置向上 3 级获取项目路径
	projectPath := runpath.PARENT.Up(3)
	zaplog.SUG.Debugln(projectPath)

	// Create workspace with single project path
	// 使用单个项目路径创建工作区
	workSpace := workspace.NewWorkSpace([]string{projectPath})

	// Configure command execution with bash and debug mode
	// 配置使用 bash 和调试模式的命令执行
	commandConfig := osexec.NewCommandConfig()
	commandConfig.WithBash()
	commandConfig.WithDebug()

	// Create workspace exec
	// 创建工作区执行器
	worksExec := worksexec.NewWorksExec(commandConfig, []*workspace.Workspace{workSpace})

	// Define root command for CLI interface
	// 为 CLI 接口定义根命令
	var rootCmd = &cobra.Command{
		Use:   "go", // Root command name // 根命令名称
		Short: "run",
		Long:  "run",
		Run: func(cmd *cobra.Command, args []string) {
			zaplog.LOG.Info("run")
		},
	}

	// Add workspace and module management subcommands
	// 添加工作区和模块管理子命令
	rootCmd.AddCommand(worksubcmd.NewWorkCmd(worksExec))
	rootCmd.AddCommand(worksubcmd.NewModCmd(worksExec))

	// Execute the CLI application
	// 执行 CLI 应用程序
	must.Done(rootCmd.Execute())
}
