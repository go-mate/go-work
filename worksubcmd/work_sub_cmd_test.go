package worksubcmd

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

	workSpace := workspace.NewWorkSpace([]string{projectPath})
	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig().WithDebug(), []*workspace.Workspace{workSpace})
	require.NoError(t, Sync(worksExec))
}

func TestTide(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workSpace := workspace.NewWorkSpace([]string{projectPath})
	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig().WithDebug(), []*workspace.Workspace{workSpace})
	require.NoError(t, Tide(worksExec))
}
