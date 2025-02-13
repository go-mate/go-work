package workconfig

import (
	"fmt"
	"path/filepath"

	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/zaplog"
)

// Workspace 代表一个 Go 工作区
type Workspace struct {
	WorkRoot string   // 工作区根目录
	Projects []string // 该 Workspace 中的项目路径
}

// NewWorkspace 创建一个新的 Workspace
func NewWorkspace(workRoot string, projects []string) *Workspace {
	return &Workspace{
		WorkRoot: workRoot,
		Projects: projects,
	}
}

// MustCheck 确保 workspace 及其项目的必要文件存在
func (w *Workspace) MustCheck() {
	if w.WorkRoot != "" {
		osmustexist.MustRoot(w.WorkRoot)
		osmustexist.MustFile(filepath.Join(w.WorkRoot, "go.work"))
	}
	for _, project := range w.Projects {
		osmustexist.MustRoot(project)
		osmustexist.MustFile(filepath.Join(project, "go.mod"))
	}
}

// Workspaces 代表多个 Go 工作区
type Workspaces []*Workspace

// NewWorkspaces 创建多个 Workspace
func NewWorkspaces(workspaces ...*Workspace) Workspaces {
	return workspaces
}

// MustCheck 确保所有 Workspaces 的必要文件都存在
func (ws Workspaces) MustCheck() {
	for _, workspace := range ws {
		workspace.MustCheck()
	}
}

type WorkspacesExecConfig struct {
	workspaces Workspaces
	execConfig *osexec.ExecConfig
}

func NewWorkspacesExecConfig(workspaces Workspaces, execConfig *osexec.ExecConfig) *WorkspacesExecConfig {
	return &WorkspacesExecConfig{
		workspaces: workspaces,
		execConfig: execConfig,
	}
}

func (wc *WorkspacesExecConfig) MustCheck() {
	wc.workspaces.MustCheck()
	must.Nice(wc.execConfig)
}

func (wc *WorkspacesExecConfig) CollectSubprojectPaths() []string {
	var paths []string
	var mp = map[string]bool{}
	for _, workspace := range wc.workspaces {
		for _, project := range workspace.Projects {
			if !mp[project] {
				mp[project] = true
				paths = append(paths, project)
			}
		}
	}
	return paths
}

func (wc *WorkspacesExecConfig) GetWorkspaces() Workspaces {
	return wc.workspaces
}

func (wc *WorkspacesExecConfig) GetSubCommand(path string) *osexec.ExecConfig {
	return wc.execConfig.ShallowClone().WithPath(path)
}

func (wc *WorkspacesExecConfig) ForeachWorkRootRun(run func(workspace *Workspace, execConfig *osexec.ExecConfig) error) error {
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

func (wc *WorkspacesExecConfig) ForeachProjectExec(run func(projectPath string, execConfig *osexec.ExecConfig) error) error {
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
