// Package workspace_test: Tests workspace structure management and validation
// Verifies workspace creation with and without workRoot configurations
// Tests project path validation and go.work file handling
//
// workspace_test: 测试工作区结构管理和验证
// 验证带和不带 workRoot 配置的工作区创建
// 测试项目路径验证和 go.work 文件处理
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

// TestNewWorkspace tests workspace creation without workRoot
// Verifies basic workspace initialization with project paths
//
// TestNewWorkspace 测试不带 workRoot 的工作区创建
// 验证带项目路径的基本工作区初始化
func TestNewWorkspace(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)
	t.Log(projectPath)

	ws := workspace.NewWorkSpace([]string{projectPath})
	must.Full(ws)
	t.Log(neatjsons.S(ws))

	require.NotNil(t, ws)
	require.Equal(t, "", ws.WorkRoot)
	require.Contains(t, ws.Projects, projectPath)
}

// TestNewWorkspaceWithWorkRoot tests workspace creation with workRoot and go.work file
// Creates temp test environment and validates workspace structure
//
// TestNewWorkspaceWithWorkRoot 测试带 workRoot 和 go.work 文件的工作区创建
// 创建临时测试环境并验证工作区结构完整性
func TestNewWorkspaceWithWorkRoot(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "test-workspace-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Create a go.work file
	workFilePath := filepath.Join(tempDIR, "go.work")
	must.Done(os.WriteFile(workFilePath, []byte("go 1.22.8\n"), 0644))

	// Create a test project DIR with go.mod
	projectDIR := filepath.Join(tempDIR, "testproject")
	must.Done(os.MkdirAll(projectDIR, 0755))
	modFilePath := filepath.Join(projectDIR, "go.mod")
	must.Done(os.WriteFile(modFilePath, []byte("module test\n\ngo 1.22.8\n"), 0644))

	ws := workspace.NewWorkspace(tempDIR, []string{projectDIR})
	require.NotNil(t, ws)
	require.Equal(t, tempDIR, ws.WorkRoot)
	require.Contains(t, ws.Projects, projectDIR)
	t.Log(neatjsons.S(ws))
}

// TestNewWorkSpaceAlias verifies that NewWorkSpace is an alias to NewWorkspace
// Ensures both functions produce matching workspace structures
//
// TestNewWorkSpaceAlias 验证 NewWorkSpace 是 NewWorkspace 的别名
// 确保两个函数产生相同的工作区结构
func TestNewWorkSpaceAlias(t *testing.T) {
	projectPath := runpath.PARENT.Up(1)

	// Test that NewWorkSpace is an alias to NewWorkspace with blank workRoot
	ws1 := workspace.NewWorkSpace([]string{projectPath})
	ws2 := workspace.NewWorkspace("", []string{projectPath})

	require.Equal(t, ws1.WorkRoot, ws2.WorkRoot)
	require.Equal(t, ws1.Projects, ws2.Projects)
}
