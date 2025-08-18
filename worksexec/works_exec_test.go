package worksexec_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
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

	require.NotNil(t, worksExec)
	require.Len(t, worksExec.GetWorkspaces(), 1)
	require.NotEmpty(t, worksExec.Subprojects())
}

func TestWorksExecSubprojects(t *testing.T) {
	tempDIR := setupTestWorkspace(t)
	defer cleanupTestWorkspace(t, tempDIR)

	project1 := filepath.Join(tempDIR, "project1")
	project2 := filepath.Join(tempDIR, "project2")

	workSpace1 := workspace.NewWorkSpace([]string{project1})
	workSpace2 := workspace.NewWorkSpace([]string{project2})

	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig(), []*workspace.Workspace{workSpace1, workSpace2})

	subprojects := worksExec.Subprojects()
	require.Len(t, subprojects, 2)
	require.Contains(t, subprojects, project1)
	require.Contains(t, subprojects, project2)
}

func TestWorksExecCommands(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	workSpace := workspace.NewWorkSpace([]string{projectPath})
	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig(), []*workspace.Workspace{workSpace})

	// Test GetNewCommand
	newCmd := worksExec.GetNewCommand()
	require.NotNil(t, newCmd)

	// Test GetSubCommand
	subCmd := worksExec.GetSubCommand(projectPath)
	require.NotNil(t, subCmd)
}

func TestForeachSubExec(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	workSpace := workspace.NewWorkSpace([]string{projectPath})
	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig(), []*workspace.Workspace{workSpace})

	executedPaths := make([]string, 0)
	err := worksExec.ForeachSubExec(func(execConfig *osexec.ExecConfig, projectPath string) error {
		executedPaths = append(executedPaths, projectPath)
		return nil
	})

	require.NoError(t, err)
	require.NotEmpty(t, executedPaths)
	require.Contains(t, executedPaths, projectPath)
}

func setupTestWorkspace(t *testing.T) string {
	tempDIR := rese.V1(os.MkdirTemp("", "test-workspace-*"))

	// Create project1
	project1 := filepath.Join(tempDIR, "project1")
	must.Done(os.MkdirAll(project1, 0755))
	must.Done(os.WriteFile(filepath.Join(project1, "go.mod"), []byte("module test1\n\ngo 1.22.8\n"), 0644))

	// Create project2
	project2 := filepath.Join(tempDIR, "project2")
	must.Done(os.MkdirAll(project2, 0755))
	must.Done(os.WriteFile(filepath.Join(project2, "go.mod"), []byte("module test2\n\ngo 1.22.8\n"), 0644))

	return tempDIR
}

func cleanupTestWorkspace(t *testing.T, tempDIR string) {
	must.Done(os.RemoveAll(tempDIR))
}
