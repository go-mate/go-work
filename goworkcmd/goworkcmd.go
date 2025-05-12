package goworkcmd

import (
	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

func NewWorkCmd(wse *worksexec.WorksExec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "work",
		Short: "go work -->>",
		Long:  "go work -->>",
		Run: func(cmd *cobra.Command, args []string) {
			panic(erero.New("wrong"))
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "sync",
		Short: "go work sync",
		Long:  "go work sync",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Sync(wse))
		},
	})

	return cmd
}

func NewModCmd(wse *worksexec.WorksExec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mod",
		Short: "go mod -->>",
		Long:  "go mod -->>",
	}
	cmd.AddCommand(NewTidyCmd(wse))
	cmd.AddCommand(NewTideCmd(wse))
	return cmd
}

func NewTidyCmd(wse *worksexec.WorksExec) *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "go mod tidy",
		Long:  "go mod tidy",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Tidy(wse))
		},
	}
}

func NewTideCmd(wse *worksexec.WorksExec) *cobra.Command {
	return &cobra.Command{
		Use:   "tide",
		Short: "go mod tidy -e",
		Long:  "go mod tidy -e",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Tide(wse))
		},
	}
}

func Sync(wse *worksexec.WorksExec) error {
	return wse.ForeachWorkRun(func(execConfig *osexec.ExecConfig, workspace *workspace.Workspace) error {
		data, err := execConfig.Exec("go", "work", "sync")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func Tidy(wse *worksexec.WorksExec) error {
	return wse.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		data, err := execConfig.Exec("go", "mod", "tidy")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func Tide(wse *worksexec.WorksExec) error {
	return wse.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		data, err := execConfig.Exec("go", "mod", "tidy", "-e")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func UpdateGoWorkVersion(wse *worksexec.WorksExec, versionNum string) error {
	return wse.ForeachWorkRun(func(execConfig *osexec.ExecConfig, workspace *workspace.Workspace) error {
		data, err := execConfig.Exec("go", "work", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func UpdateGoModuleVersion(wse *worksexec.WorksExec, versionNum string) error {
	return wse.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		data, err := execConfig.Exec("go", "mod", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}
