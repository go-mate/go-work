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

func NewWorkCmd(worksExec *worksexec.WorksExec) *cobra.Command {
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
			must.Done(Sync(worksExec))
		},
	})

	return cmd
}

func NewModCmd(worksExec *worksexec.WorksExec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mod",
		Short: "go mod -->>",
		Long:  "go mod -->>",
		Run: func(cmd *cobra.Command, args []string) {
			panic(erero.New("wrong"))
		},
	}
	cmd.AddCommand(NewTidyCmd(worksExec))
	cmd.AddCommand(NewTideCmd(worksExec))
	return cmd
}

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

func Sync(worksExec *worksexec.WorksExec) error {
	return worksExec.ForeachWorkRun(func(execConfig *osexec.ExecConfig, workspace *workspace.Workspace) error {
		data, err := execConfig.Exec("go", "work", "sync")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func Tidy(worksExec *worksexec.WorksExec) error {
	return worksExec.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		data, err := execConfig.Exec("go", "mod", "tidy")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func Tide(worksExec *worksexec.WorksExec) error {
	return worksExec.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		data, err := execConfig.Exec("go", "mod", "tidy", "-e")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func UpdateGoWorkGoVersion(worksExec *worksexec.WorksExec, versionNum string) error {
	return worksExec.ForeachWorkRun(func(execConfig *osexec.ExecConfig, workspace *workspace.Workspace) error {
		data, err := execConfig.Exec("go", "work", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func UpdateModuleGoVersion(worksExec *worksexec.WorksExec, versionNum string) error {
	return worksExec.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		data, err := execConfig.Exec("go", "mod", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}
