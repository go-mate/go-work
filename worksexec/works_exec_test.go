package worksexec_test

import (
	"testing"

	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestNewWorksExec(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workSpace := workspace.NewWorkSpace([]string{projectPath})
	must.Full(workSpace)
	t.Log(neatjsons.S(workSpace))

	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig(), []*workspace.Workspace{workSpace})
	must.Full(worksExec)
	t.Log(neatjsons.S(worksExec.GetWorkspaces()))
}
