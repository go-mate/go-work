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

	var commandLine string
	var debugMode bool
	rootCmd := cobra.Command{
		Use:   "go-work",
		Short: "go-work",
		Long:  "go-work",
		Run: func(cmd *cobra.Command, args []string) {
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

func run(workPath string, shellType string, commandLine string, debugMode bool) {
	options := workspath.NewOptions().
		WithIncludeCurrentProject(true).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(debugMode)
	modulePaths := workspath.GetModulePaths(workPath, options)
	if debugMode {
		zaplog.SUG.Debugln("module paths:", neatjsons.S(modulePaths))
	}
	for _, modulePath := range modulePaths {
		zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("--"))
		executeInSinglePath(modulePath, shellType, commandLine, debugMode)
		zaplog.SUG.Debugln(eroticgo.BLUE.Sprint("--"))
	}
	eroticgo.GREEN.ShowMessage("SUCCESS")
}

func executeInSinglePath(modulePath string, shellType string, commandLine string, debugMode bool) {
	commandMessage := eroticgo.AMBER.Sprint("cd", modulePath, "&&", commandLine)
	if debugMode {
		zaplog.SUG.Debugln("executing:", commandMessage)
	}

	config := osexec.NewExecConfig().WithPath(modulePath)
	output := rese.V1(config.WithShell(shellType, "-c").Exec(commandLine))

	if debugMode {
		zaplog.SUG.Debugln("executing:", commandMessage, "output:", eroticgo.GREEN.Sprint(string(output)))
	}
	zaplog.SUG.Debugln("executing:", commandMessage, "->:", "success")
}
