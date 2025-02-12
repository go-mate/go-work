package worksubcmd

import (
	"github.com/go-mate/go-work/workconfig"
	"github.com/spf13/cobra"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

func NewWorkCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
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

func NewModCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mod",
		Short: "go mod -->>",
		Long:  "go mod -->>",
	}
	cmd.AddCommand(NewTidyCmd(config))
	cmd.AddCommand(NewTideCmd(config))
	return cmd
}

func NewTidyCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "go mod tidy",
		Long:  "go mod tidy",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Tidy(config))
		},
	}
}

func NewTideCmd(config *workconfig.WorkspacesExecConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "tide",
		Short: "go mod tidy -e",
		Long:  "go mod tidy -e",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(Tide(config))
		},
	}
}

func Sync(config *workconfig.WorkspacesExecConfig) error {
	return config.ForeachWorkRootRun(func(workspace *workconfig.Workspace, command *osexec.CommandConfig) error {
		data, err := command.Exec("go", "work", "sync")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func Tidy(config *workconfig.WorkspacesExecConfig) error {
	return config.ForeachProjectExec(func(projectPath string, command *osexec.CommandConfig) error {
		data, err := command.Exec("go", "mod", "tidy")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func Tide(config *workconfig.WorkspacesExecConfig) error {
	return config.ForeachProjectExec(func(projectPath string, command *osexec.CommandConfig) error {
		data, err := command.Exec("go", "mod", "tidy", "-e")
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func UpdateGoWorkVersion(config *workconfig.WorkspacesExecConfig, versionNum string) error {
	return config.ForeachWorkRootRun(func(workspace *workconfig.Workspace, command *osexec.CommandConfig) error {
		data, err := command.Exec("go", "work", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}

func UpdateGoModuleVersion(config *workconfig.WorkspacesExecConfig, versionNum string) error {
	return config.ForeachProjectExec(func(projectPath string, command *osexec.CommandConfig) error {
		data, err := command.Exec("go", "mod", "edit", "-go", versionNum)
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(string(data))
		return nil
	})
}
