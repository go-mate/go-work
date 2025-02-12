package worksubcmd

import (
	"testing"

	"github.com/go-mate/go-work/workconfig"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestSync(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := workconfig.NewWorkspace("", []string{projectPath})
	workspace.MustCheck()

	workspaces := workconfig.NewWorkspaces(workspace)
	workspaces.MustCheck()

	config := workconfig.NewWorkspacesExecConfig(workspaces, osexec.NewCommandConfig().WithDebugMode(true))
	config.MustCheck()

	require.NoError(t, Sync(config))
}

func TestTide(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := workconfig.NewWorkspace("", []string{projectPath})
	workspace.MustCheck()

	workspaces := workconfig.NewWorkspaces(workspace)
	workspaces.MustCheck()

	config := workconfig.NewWorkspacesExecConfig(workspaces, osexec.NewCommandConfig().WithDebugMode(true))
	config.MustCheck()

	require.NoError(t, Tide(config))
}
