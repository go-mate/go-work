package workspace_test

import (
	"testing"

	"github.com/go-mate/go-work/workspace"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/runpath"
)

func TestNewWorkspace(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	wsp := workspace.NewWorkspace("", []string{projectPath})
	must.Full(wsp)
	t.Log(neatjsons.S(wsp))
}
