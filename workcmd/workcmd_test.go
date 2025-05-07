package workcmd

import (
	"testing"

	"github.com/go-mate/go-work/workcfg"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestSync(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := workcfg.NewWorkspace("", []string{projectPath})

	worksExec := workcfg.NewWorksExec(osexec.NewCommandConfig().WithDebugMode(true), []*workcfg.Workspace{workspace})

	require.NoError(t, Sync(worksExec))
}

func TestTide(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := workcfg.NewWorkspace("", []string{projectPath})

	worksExec := workcfg.NewWorksExec(osexec.NewCommandConfig().WithDebugMode(true), []*workcfg.Workspace{workspace})

	require.NoError(t, Tide(worksExec))
}
