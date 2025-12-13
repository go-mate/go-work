// Package main: go-work entrypoint
// Lists Go module paths and versions in the current workspace
//
// main: go-work 的程序入口
// 列举当前工作区中的 Go 模块路径和版本
package main

import (
	"os"
	"path/filepath"

	"github.com/go-mate/go-work/workspath"
	"github.com/spf13/cobra"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/tern"
	"github.com/yyle88/zaplog"
	"golang.org/x/mod/modfile"
)

func main() {
	workPath := rese.C1(os.Getwd())

	rootCmd := &cobra.Command{
		Use:   "go-work",
		Short: "List Go modules in workspace",
		Long:  "go-work: Lists Go module paths in the current workspace",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			showPathList(workPath)
		},
	}

	rootCmd.AddCommand(newVersionCmd(workPath))
	must.Done(rootCmd.Execute())
}

// showPathList lists all Go module paths in workspace
// showPathList 列举工作区中所有 Go 模块路径
func showPathList(workPath string) {
	type Result struct {
		Path   string `json:"path"`
		Module string `json:"module"`
	}
	var results []*Result
	for _, path := range getModulePaths(workPath) {
		modFile := parseModFile(path)
		results = append(results, &Result{
			Path:   path,
			Module: modFile.Module.Mod.Path,
		})
	}
	zaplog.SUG.Debugln(neatjsons.S(results))
}

// newVersionCmd creates version subcommand to show go versions
// newVersionCmd 创建 version 子命令来显示 go 版本
func newVersionCmd(workPath string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "List Go versions used in each module",
		Long:  "Shows the Go version specified in each module's go.mod file",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			showVersionList(workPath)
		},
	}
}

// showVersionList lists go versions from each module's go.mod
// showVersionList 列举每个模块 go.mod 中的 go 版本
func showVersionList(workPath string) {
	type Result struct {
		Path    string `json:"path"`
		Module  string `json:"module"`
		Version string `json:"version"`
	}
	var results []*Result
	for _, path := range getModulePaths(workPath) {
		modFile := parseModFile(path)
		goVersion := tern.BFV(modFile.Go != nil, func() string {
			return modFile.Go.Version
		}, "unknown")
		results = append(results, &Result{
			Path:    path,
			Module:  modFile.Module.Mod.Path,
			Version: goVersion,
		})
	}
	zaplog.SUG.Debugln(neatjsons.S(results))
}

// getModulePaths returns all Go module paths in workspace
// getModulePaths 返回工作区中所有 Go 模块路径
func getModulePaths(workPath string) []string {
	return workspath.GetModulePaths(
		workPath,
		workspath.WithCurrentProject(),
		workspath.ScanDeep(),
		workspath.SkipNoGo(),
	)
}

// parseModFile parses go.mod file and returns modfile.File
// parseModFile 解析 go.mod 文件并返回 modfile.File
func parseModFile(modulePath string) *modfile.File {
	modPath := filepath.Join(modulePath, "go.mod")
	content := rese.V1(os.ReadFile(modPath))
	return rese.P1(modfile.Parse(modPath, content, nil))
}
