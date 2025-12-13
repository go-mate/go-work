// Package workspath: Tests Go module path detection and workspace scanning
// Verifies module detection with different option configurations
//
// workspath: 测试 Go 模块路径发现和工作区扫描
// 验证不同选项配置的模块检测
package workspath

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

// TestGetProjectPath tests project root detection with ProjectPath
// TestGetProjectPath 测试项目根目录发现（返回 ProjectPath）
func TestGetProjectPath(t *testing.T) {
	pkgPath := runpath.PARENT.Path()
	t.Log("pkgPath:", pkgPath)

	info, ok := GetProjectPath(pkgPath)
	require.True(t, ok)
	require.NotEmpty(t, info.Root)
	require.Equal(t, runpath.PARENT.Up(1), info.Root)
	require.Equal(t, "workspath", info.SubPath)
	t.Log("info:", neatjsons.S(info))
}

// TestGetProjectPath_NotFound tests when go.mod is not found
// TestGetProjectPath_NotFound 测试找不到 go.mod 的情况
func TestGetProjectPath_NotFound(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "test-no-gomod-*"))
	defer must.Done(os.RemoveAll(tempDIR))

	info, ok := GetProjectPath(tempDIR)
	require.False(t, ok)
	require.Nil(t, info)
}

// TestGetProjectRoot tests simple project root detection
// TestGetProjectRoot 测试简单项目根目录发现
func TestGetProjectRoot(t *testing.T) {
	pkgPath := runpath.PARENT.Path()
	t.Log("pkgPath:", pkgPath)

	root, ok := GetProjectRoot(pkgPath)
	require.True(t, ok)
	require.NotEmpty(t, root)
	require.Equal(t, runpath.PARENT.Up(1), root)
	t.Log("root:", root)
}

// TestGetProjectRoot_NotFound tests when go.mod is not found
// TestGetProjectRoot_NotFound 测试找不到 go.mod 的情况
func TestGetProjectRoot_NotFound(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "test-no-gomod-*"))
	defer must.Done(os.RemoveAll(tempDIR))

	root, ok := GetProjectRoot(tempDIR)
	require.False(t, ok)
	require.Empty(t, root)
}

// TestGetModulePaths tests basic module scanning
// TestGetModulePaths 测试基本模块扫描
func TestGetModulePaths(t *testing.T) {
	// Use project root (go-work) not sub DIR (workspath)
	// 使用项目根目录 (go-work) 而不是子目录 (workspath)
	projectRoot := runpath.PARENT.Up(1)
	t.Log("projectRoot:", projectRoot)

	paths := GetModulePaths(projectRoot, WithCurrentProject(), ScanDeep(), SkipNoGo(), DebugMode())
	t.Log("paths:", neatjsons.S(paths))

	require.NotEmpty(t, paths)
	require.Contains(t, paths, projectRoot)
}

// TestGetModulePaths_RootIsModule tests scanning when root is a module
// TestGetModulePaths_RootIsModule 测试 root 本身是模块的情况
func TestGetModulePaths_RootIsModule(t *testing.T) {
	tempDIR := setupTestProject(t)
	defer cleanupDIR(t, tempDIR)

	paths := GetModulePaths(tempDIR, WithCurrentProject())
	require.Len(t, paths, 1)
	require.Equal(t, tempDIR, paths[0])
}

// TestGetModulePaths_ScanDeep tests submodule scanning
// TestGetModulePaths_ScanDeep 测试子模块扫描
func TestGetModulePaths_ScanDeep(t *testing.T) {
	tempDIR := setupTestProject(t)
	defer cleanupDIR(t, tempDIR)

	// Without ScanDeep - just root
	paths := GetModulePaths(tempDIR, WithCurrentProject())
	require.Len(t, paths, 1)

	// With ScanDeep - root + submodule
	paths = GetModulePaths(tempDIR, WithCurrentProject(), ScanDeep())
	require.Len(t, paths, 2)
}

// TestGetModulePaths_SkipNoGo tests empty module filtering
// TestGetModulePaths_SkipNoGo 测试空模块过滤
func TestGetModulePaths_SkipNoGo(t *testing.T) {
	tempDIR := setupTestProject(t)
	defer cleanupDIR(t, tempDIR)

	// With ScanDeep, no filter - includes empty submodule
	paths := GetModulePaths(tempDIR, WithCurrentProject(), ScanDeep())
	require.Len(t, paths, 2)

	// With ScanDeep + SkipNoGo - excludes empty submodule
	paths = GetModulePaths(tempDIR, WithCurrentProject(), ScanDeep(), SkipNoGo())
	require.Len(t, paths, 1)
	require.Equal(t, tempDIR, paths[0])
}

// TestGetModulePaths_NotModule tests scanning non-module DIR
// TestGetModulePaths_NotModule 测试扫描非模块目录
func TestGetModulePaths_NotModule(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "test-not-module-*"))
	defer must.Done(os.RemoveAll(tempDIR))

	paths := GetModulePaths(tempDIR, WithCurrentProject())
	require.Empty(t, paths)
}

// TestGetModulePaths_WithCurrentPackage tests current package inclusion
// TestGetModulePaths_WithCurrentPackage 测试当前包包含
func TestGetModulePaths_WithCurrentPackage(t *testing.T) {
	pkgPath := runpath.PARENT.Path()

	paths := GetModulePaths(pkgPath, WithCurrentPackage(), ScanDeep(), SkipNoGo())
	t.Log("paths:", neatjsons.S(paths))

	require.NotEmpty(t, paths)
	require.Contains(t, paths, pkgPath)
}

// =====================================================
// Test Helpers
// 测试辅助函数
// =====================================================

// setupTestProject creates temp test project structure
// Returns temp DIR with main project and empty submodule
//
// setupTestProject 创建临时测试项目结构
// 返回带主项目和空子模块的临时 DIR
func setupTestProject(t *testing.T) string {
	tempDIR := rese.V1(os.MkdirTemp("", "test-workspath-*"))

	// Main project with go.mod and .go files
	// 主项目带 go.mod 和 .go 文件
	must.Done(os.WriteFile(filepath.Join(tempDIR, "go.mod"), []byte("module test\n\ngo 1.22.8\n"), 0644))
	must.Done(os.WriteFile(filepath.Join(tempDIR, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0644))

	// Submodule with go.mod but no .go files
	// 子模块带 go.mod 但无 .go 文件
	subDIR := filepath.Join(tempDIR, "submodule")
	must.Done(os.MkdirAll(subDIR, 0755))
	must.Done(os.WriteFile(filepath.Join(subDIR, "go.mod"), []byte("module test/sub\n\ngo 1.22.8\n"), 0644))

	return tempDIR
}

// cleanupDIR removes temp test DIR
// cleanupDIR 删除临时测试 DIR
func cleanupDIR(t *testing.T, tempDIR string) {
	must.Done(os.RemoveAll(tempDIR))
}
