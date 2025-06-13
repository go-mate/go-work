package worksexec

import (
	"fmt"

	"github.com/emirpasic/gods/v2/sets/linkedhashset"
	"github.com/go-mate/go-work/workspace"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

type WorksExec struct {
	execConfig *osexec.ExecConfig
	workspaces []*workspace.Workspace
}

func NewWorksExec(execConfig *osexec.ExecConfig, workspaces []*workspace.Workspace) (wse *WorksExec) {
	return &WorksExec{
		execConfig: must.Nice(execConfig),
		workspaces: must.Have(workspaces),
	}
}

func (wse *WorksExec) Subprojects() []string {
	set := linkedhashset.New[string]()
	for _, wsp := range wse.workspaces {
		for _, path := range wsp.Projects {
			set.Add(path)
		}
	}
	return set.Values()
}

func (wse *WorksExec) GetWorkspaces() []*workspace.Workspace {
	return wse.workspaces
}

func (wse *WorksExec) GetNewCommand() *osexec.ExecConfig {
	return wse.execConfig.NewConfig()
}

func (wse *WorksExec) GetSubCommand(path string) *osexec.ExecConfig {
	return wse.GetNewCommand().WithPath(path)
}

func (wse *WorksExec) ForeachWorkRun(run func(execConfig *osexec.ExecConfig, workspace *workspace.Workspace) error) error {
	for idx, wsp := range wse.GetWorkspaces() {
		processMessage := fmt.Sprintf("(%d/%d)", idx, len(wse.GetWorkspaces()))

		if wsp.WorkRoot != "" {
			zaplog.SUG.Debugln("run", processMessage)

			if err := run(wse.GetSubCommand(wsp.WorkRoot), wsp); err != nil {
				return erero.Wro(err)
			}
		} else {
			zaplog.SUG.Debugln("run", processMessage, "no work root so nothing to do")
		}
	}
	return nil
}

func (wse *WorksExec) ForeachSubExec(run func(execConfig *osexec.ExecConfig, projectPath string) error) error {
	for idx, wsp := range wse.GetWorkspaces() {
		for num, projectPath := range wsp.Projects {
			process1Message := fmt.Sprintf("(%d/%d)", idx, len(wse.GetWorkspaces()))
			process2Message := fmt.Sprintf("(%d/%d)", num, len(wsp.Projects))

			zaplog.SUG.Debugln("run", process1Message, process2Message, projectPath)

			if err := run(wse.GetSubCommand(projectPath), projectPath); err != nil {
				return erero.Wro(err)
			}
		}
	}
	return nil
}
