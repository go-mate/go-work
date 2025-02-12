package workconfig

import (
	"testing"

	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestWorkspace_MustCheck(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := NewWorkspace("", []string{projectPath})
	workspace.MustCheck()
}

func TestWorkspaces_MustCheck(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := NewWorkspace("", []string{projectPath})
	workspace.MustCheck()

	workspaces := NewWorkspaces(workspace)
	workspaces.MustCheck()
}

func TestWorkspacesExec_MustCheck(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := NewWorkspace("", []string{projectPath})
	workspace.MustCheck()

	workspaces := NewWorkspaces(workspace)
	workspaces.MustCheck()

	config := NewWorkspacesExecConfig(workspaces, osexec.NewCommandConfig())
	config.MustCheck()
}
