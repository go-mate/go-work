package workcmd

import (
	"github.com/go-mate/go-work/workcfg"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

func NewWorkCmd(config *workcfg.WorksExec) *cobra.Command {
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
			must.Done(Sync(config))
		},
	})

	return cmd
}

func NewModCmd(config *workcfg.WorksExec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mod",
		Short: "go mod -->>",
		Long:  "go mod -->>",
	}
	cmd.AddCommand(NewTidyCmd(config))
	cmd.AddCommand(NewTideCmd(config))
	return cmd
}

func NewTidyCmd(config *workcfg.WorksExec) *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "go mod tidy",
		Long:  "go mod tidy",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Tidy(config))
		},
	}
}

func NewTideCmd(config *workcfg.WorksExec) *cobra.Command {
	return &cobra.Command{
		Use:   "tide",
		Short: "go mod tidy -e",
		Long:  "go mod tidy -e",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Tide(config))
		},
	}
}

func Sync(config *workcfg.WorksExec) error {
	return config.ForeachWorkRun(func(workspace *workcfg.Workspace, execConfig *osexec.ExecConfig) error {
		data, err := execConfig.Exec("go", "work", "sync")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func Tidy(config *workcfg.WorksExec) error {
	return config.ForeachSubExec(func(projectPath string, execConfig *osexec.ExecConfig) error {
		data, err := execConfig.Exec("go", "mod", "tidy")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func Tide(config *workcfg.WorksExec) error {
	return config.ForeachSubExec(func(projectPath string, execConfig *osexec.ExecConfig) error {
		data, err := execConfig.Exec("go", "mod", "tidy", "-e")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func UpdateGoWorkVersion(config *workcfg.WorksExec, versionNum string) error {
	return config.ForeachWorkRun(func(workspace *workcfg.Workspace, execConfig *osexec.ExecConfig) error {
		data, err := execConfig.Exec("go", "work", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func UpdateGoModuleVersion(config *workcfg.WorksExec, versionNum string) error {
	return config.ForeachSubExec(func(projectPath string, execConfig *osexec.ExecConfig) error {
		data, err := execConfig.Exec("go", "mod", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}
