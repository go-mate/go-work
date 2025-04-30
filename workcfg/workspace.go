package workcfg

import (
	"fmt"
	"path/filepath"

	"github.com/emirpasic/gods/v2/sets/linkedhashset"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/zaplog"
)

// Workspace 代表一个 Go 工作目录，里面有多个子项目
type Workspace struct {
	WorkRoot string   // 工作区根目录
	Projects []string // 该 Workspace 中的项目路径
}

func NewWorkspace(workRoot string, projects []string) *Workspace {
	if workRoot != "" {
		osmustexist.MustRoot(workRoot)
		osmustexist.MustFile(filepath.Join(workRoot, "go.work"))
	}
	for _, project := range must.Have(projects) {
		osmustexist.MustRoot(project)
		osmustexist.MustFile(filepath.Join(project, "go.mod"))
	}
	return &Workspace{
		WorkRoot: workRoot,
		Projects: projects,
	}
}

type WorksExec struct {
	workspaces []*Workspace
	execConfig *osexec.ExecConfig
}

func NewWorksExec(workspaces []*Workspace, execConfig *osexec.ExecConfig) *WorksExec {
	return &WorksExec{
		workspaces: must.Have(workspaces),
		execConfig: must.Nice(execConfig),
	}
}

func (wc *WorksExec) Subprojects() []string {
	paths := linkedhashset.New[string]()
	for _, workspace := range wc.workspaces {
		for _, project := range workspace.Projects {
			paths.Add(project)
		}
	}
	return paths.Values()
}

func (wc *WorksExec) GetWorkspaces() []*Workspace {
	return wc.workspaces
}

func (wc *WorksExec) GetNewCommand() *osexec.ExecConfig {
	return wc.execConfig.ShallowClone()
}

func (wc *WorksExec) GetSubCommand(path string) *osexec.ExecConfig {
	return wc.GetNewCommand().WithPath(path)
}

func (wc *WorksExec) ForeachWorkRun(run func(workspace *Workspace, execConfig *osexec.ExecConfig) error) error {
	for idx, workspace := range wc.GetWorkspaces() {
		processMessage := fmt.Sprintf("(%d/%d)", idx, len(wc.GetWorkspaces()))

		if workspace.WorkRoot != "" {
			zaplog.SUG.Debugln("run", processMessage)

			if err := run(workspace, wc.GetSubCommand(workspace.WorkRoot)); err != nil {
				return erero.Wro(err)
			}
		} else {
			zaplog.SUG.Debugln("run", processMessage, "no work root so nothing to do")
		}
	}
	return nil
}

func (wc *WorksExec) ForeachSubExec(run func(projectPath string, execConfig *osexec.ExecConfig) error) error {
	for idx, workspace := range wc.GetWorkspaces() {
		for num, projectPath := range workspace.Projects {
			process1Message := fmt.Sprintf("(%d/%d)", idx, len(wc.GetWorkspaces()))
			process2Message := fmt.Sprintf("(%d/%d)", num, len(workspace.Projects))

			zaplog.SUG.Debugln("run", process1Message, process2Message, projectPath)

			if err := run(projectPath, wc.GetSubCommand(projectPath)); err != nil {
				return erero.Wro(err)
			}
		}
	}
	return nil
}
