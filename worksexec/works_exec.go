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

func NewWorksExec(execConfig *osexec.ExecConfig, workspaces []*workspace.Workspace) (worksExec *WorksExec) {
	return &WorksExec{
		execConfig: must.Nice(execConfig),
		workspaces: must.Have(workspaces),
	}
}

func (W *WorksExec) Subprojects() []string {
	set := linkedhashset.New[string]()
	for _, wsp := range W.workspaces {
		for _, path := range wsp.Projects {
			set.Add(path)
		}
	}
	return set.Values()
}

func (W *WorksExec) GetWorkspaces() []*workspace.Workspace {
	return W.workspaces
}

func (W *WorksExec) GetNewCommand() *osexec.ExecConfig {
	return W.execConfig.NewConfig()
}

func (W *WorksExec) GetSubCommand(path string) *osexec.ExecConfig {
	return W.GetNewCommand().WithPath(path)
}

func (W *WorksExec) ForeachWorkRun(run func(execConfig *osexec.ExecConfig, workspace *workspace.Workspace) error) error {
	workspaces := W.GetWorkspaces()
	for idx, wsp := range workspaces {
		processMessage := fmt.Sprintf("(%d/%d)", idx, len(workspaces))

		if wsp.WorkRoot != "" {
			zaplog.SUG.Debugln("run", processMessage)

			if err := run(W.GetSubCommand(wsp.WorkRoot), wsp); err != nil {
				return erero.Wro(err)
			}
		} else {
			zaplog.SUG.Debugln("run", processMessage, "no work root so nothing to do")
		}
	}
	return nil
}

func (W *WorksExec) ForeachSubExec(run func(execConfig *osexec.ExecConfig, projectPath string) error) error {
	workspaces := W.GetWorkspaces()
	for idx, wsp := range workspaces {
		for num, projectPath := range wsp.Projects {
			process1Message := fmt.Sprintf("(%d/%d)", idx, len(workspaces))
			process2Message := fmt.Sprintf("(%d/%d)", num, len(wsp.Projects))

			zaplog.SUG.Debugln("run", process1Message, process2Message, projectPath)

			if err := run(W.GetSubCommand(projectPath), projectPath); err != nil {
				return erero.Wro(err)
			}
		}
	}
	return nil
}
