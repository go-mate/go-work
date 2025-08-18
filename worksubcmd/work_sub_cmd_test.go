package worksubcmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-mate/go-work/worksexec"
	"github.com/go-mate/go-work/workspace"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
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

func TestTidy(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)

	workSpace := workspace.NewWorkSpace([]string{projectPath})
	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig(), []*workspace.Workspace{workSpace})
	require.NoError(t, Tidy(worksExec))
}

func TestNewWorkCmd(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	workSpace := workspace.NewWorkSpace([]string{projectPath})
	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig(), []*workspace.Workspace{workSpace})

	cmd := NewWorkCmd(worksExec)
	require.NotNil(t, cmd)
	require.Equal(t, "work", cmd.Use)
	require.NotEmpty(t, cmd.Commands())

	// Test that sync subcommand exists
	found := false
	for _, subCmd := range cmd.Commands() {
		if subCmd.Use == "sync" {
			found = true
			break
		}
	}
	require.True(t, found, "sync subcommand should exist")
}

func TestNewModCmd(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	workSpace := workspace.NewWorkSpace([]string{projectPath})
	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig(), []*workspace.Workspace{workSpace})

	cmd := NewModCmd(worksExec)
	require.NotNil(t, cmd)
	require.Equal(t, "mod", cmd.Use)
	require.NotEmpty(t, cmd.Commands())

	// Test that tidy and tide subcommands exist
	subCommands := make(map[string]*cobra.Command)
	for _, subCmd := range cmd.Commands() {
		subCommands[subCmd.Use] = subCmd
	}
	require.Contains(t, subCommands, "tidy")
	require.Contains(t, subCommands, "tide")
}

func TestUpdateGoWorkGoVersion(t *testing.T) {
	tempDIR := setupTestWorkspaceWithGoWork(t)
	defer cleanupTestDIR(t, tempDIR)

	workSpace := workspace.NewWorkspace(tempDIR, []string{filepath.Join(tempDIR, "project1")})
	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig(), []*workspace.Workspace{workSpace})

	err := UpdateGoWorkGoVersion(worksExec, "1.22.8")
	require.NoError(t, err)
}

func TestUpdateModuleGoVersion(t *testing.T) {
	tempDIR := setupTestWorkspaceWithGoWork(t)
	defer cleanupTestDIR(t, tempDIR)

	projectPath := filepath.Join(tempDIR, "project1")
	workSpace := workspace.NewWorkSpace([]string{projectPath})
	worksExec := worksexec.NewWorksExec(osexec.NewCommandConfig(), []*workspace.Workspace{workSpace})

	err := UpdateModuleGoVersion(worksExec, "1.22.8")
	require.NoError(t, err)
}

func setupTestWorkspaceWithGoWork(t *testing.T) string {
	tempDIR := rese.V1(os.MkdirTemp("", "test-worksubcmd-*"))

	// Create go.work file
	goWorkContent := "go 1.22.8\n\nuse (\n\t./project1\n)\n"
	must.Done(os.WriteFile(filepath.Join(tempDIR, "go.work"), []byte(goWorkContent), 0644))

	// Create project1 with go.mod
	project1 := filepath.Join(tempDIR, "project1")
	must.Done(os.MkdirAll(project1, 0755))
	goModContent := "module test/project1\n\ngo 1.22.8\n"
	must.Done(os.WriteFile(filepath.Join(project1, "go.mod"), []byte(goModContent), 0644))

	return tempDIR
}

func cleanupTestDIR(t *testing.T, tempDIR string) {
	must.Done(os.RemoveAll(tempDIR))
}
