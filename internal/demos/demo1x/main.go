package main

import (
	"github.com/go-mate/go-work/goworkcmd"
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
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

	wsp := workspace.NewWorkspace("", []string{projectPath})

	commandConfig := osexec.NewCommandConfig()
	commandConfig.WithBash()
	commandConfig.WithDebug()

	wse := worksexec.NewWorksExec(commandConfig, []*workspace.Workspace{wsp})

	// 定义根命令
	var rootCmd = &cobra.Command{
		Use:   "go", // 根命令的名称
		Short: "run",
		Long:  "run",
		Run: func(cmd *cobra.Command, args []string) {
			zaplog.LOG.Info("run")
		},
	}
	rootCmd.AddCommand(goworkcmd.NewWorkCmd(wse))
	rootCmd.AddCommand(goworkcmd.NewModCmd(wse))

	must.Done(rootCmd.Execute())
}
