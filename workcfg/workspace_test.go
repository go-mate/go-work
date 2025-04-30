package workcfg

import (
	"testing"

	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/runpath"
)

func TestNewWorkspace(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := NewWorkspace("", []string{projectPath})
	must.Full(workspace)
}

func TestNewWorksExec(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workspace := NewWorkspace("", []string{projectPath})
	must.Full(workspace)

	worksExec := NewWorksExec([]*Workspace{workspace}, osexec.NewCommandConfig())
	must.Full(worksExec)
}
