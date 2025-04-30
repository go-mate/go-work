package main

import (
	"github.com/go-mate/go-work/workcfg"
	"github.com/go-mate/go-work/workcmd"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
)

// go run main.go help
// go run main.go mod tidy
// go run main.go mod tide
// go run main.go work sync
func main() {
	projectPath := runpath.PARENT.Up(3)
	zaplog.SUG.Debugln(projectPath)

	workspace := workcfg.NewWorkspace("", []string{projectPath})

	commandConfig := osexec.NewCommandConfig()
	commandConfig.WithBash()
	commandConfig.WithDebugMode(true)

	worksExec := workcfg.NewWorksExec([]*workcfg.Workspace{workspace}, commandConfig)

	// 定义根命令
	var rootCmd = &cobra.Command{
		Use:   "go", // 根命令的名称
		Short: "run",
		Long:  "run",
		Run: func(cmd *cobra.Command, args []string) {
			zaplog.LOG.Info("run")
		},
	}
	rootCmd.AddCommand(workcmd.NewWorkCmd(worksExec))
	rootCmd.AddCommand(workcmd.NewModCmd(worksExec))

	must.Done(rootCmd.Execute())
}
