package goworkcmd

import (
	"testing"

	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestSync(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	wsp := workspace.NewWorkspace("", []string{projectPath})

	wse := worksexec.NewWorksExec(osexec.NewCommandConfig().WithDebug(), []*workspace.Workspace{wsp})

	require.NoError(t, Sync(wse))
}

func TestTide(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	wsp := workspace.NewWorkspace("", []string{projectPath})

	wse := worksexec.NewWorksExec(osexec.NewCommandConfig().WithDebug(), []*workspace.Workspace{wsp})

	require.NoError(t, Tide(wse))
}
