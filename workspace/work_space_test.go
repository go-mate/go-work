package workspace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-mate/go-work/workspace"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

func TestNewWorkspace(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	workSpace := workspace.NewWorkSpace([]string{projectPath})
	must.Full(workSpace)
	t.Log(neatjsons.S(workSpace))

	require.NotNil(t, workSpace)
	require.Equal(t, "", workSpace.WorkRoot)
	require.Contains(t, workSpace.Projects, projectPath)
}

func TestNewWorkspaceWithWorkRoot(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "test-workspace-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Create a go.work file
	goWorkFile := filepath.Join(tempDIR, "go.work")
	must.Done(os.WriteFile(goWorkFile, []byte("go 1.22.8\n"), 0644))

	// Create a test project DIR with go.mod
	projectDIR := filepath.Join(tempDIR, "testproject")
	must.Done(os.MkdirAll(projectDIR, 0755))
	goModFile := filepath.Join(projectDIR, "go.mod")
	must.Done(os.WriteFile(goModFile, []byte("module test\n\ngo 1.22.8\n"), 0644))

	workSpace := workspace.NewWorkspace(tempDIR, []string{projectDIR})
	require.NotNil(t, workSpace)
	require.Equal(t, tempDIR, workSpace.WorkRoot)
	require.Contains(t, workSpace.Projects, projectDIR)
	t.Log(neatjsons.S(workSpace))
}

func TestNewWorkSpaceAlias(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)

	// Test that NewWorkSpace is an alias for NewWorkspace with empty workRoot
	workSpace1 := workspace.NewWorkSpace([]string{projectPath})
	workSpace2 := workspace.NewWorkspace("", []string{projectPath})

	require.Equal(t, workSpace1.WorkRoot, workSpace2.WorkRoot)
	require.Equal(t, workSpace1.Projects, workSpace2.Projects)
}
