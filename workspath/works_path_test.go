// Package workspath: Tests Go module path detection and workspace traversing
// Verifies module detection with different option configurations
// Tests filtering logic and path resolution
//
// workspath: 测试 Go 模块路径发现和工作区遍历
// 验证不同选项配置的模块检测
// 测试过滤逻辑和路径解析准确性
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

// TestGetModulePaths tests module path detection with various options
// Verifies that paths are detected and filtered as expected
//
// TestGetModulePaths 测试带各种选项的模块路径发现
// 验证路径被正确发现和过滤
func TestGetModulePaths(t *testing.T) {
	pkgPath := runpath.PARENT.Path()
	t.Log(pkgPath)

	options := NewOptions().
		WithIncludeCurrentProject(true).
		WithIncludeCurrentPackage(true).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(true)
	t.Log(neatjsons.S(options))

	paths := GetModulePaths(pkgPath, options)
	t.Log(neatjsons.S(paths))

	require.NotEmpty(t, paths, "should find at least one module path")
	require.Contains(t, paths, runpath.PARENT.Up(1), "should include project root")
}

// TestGetModulePathsOptions tests different option combinations
// Validates handling of different inclusion and exclusion settings
//
// TestGetModulePathsOptions 测试不同的选项组合
// 验证不同包含和排除设置的行为
func TestGetModulePathsOptions(t *testing.T) {
	tempDIR := setupTestGoProject(t)
	defer cleanupTestDIR(t, tempDIR)

	// Test different option combinations
	tests := []struct {
		name     string
		options  *Options
		expected int
	}{
		{
			name:     "include current project just",
			options:  NewOptions().WithIncludeCurrentProject(true),
			expected: 1,
		},
		{
			name:     "include submodules just",
			options:  NewOptions().WithIncludeSubModules(true),
			expected: 2, // main + sub
		},
		{
			name:     "exclude directories without Go files",
			options:  NewOptions().WithIncludeSubModules(true).WithExcludeNoGo(true),
			expected: 1, // just the one with Go files
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := GetModulePaths(tempDIR, tt.options)
			require.Len(t, paths, tt.expected, "should have expected count of paths")
		})
	}
}

// TestNewOptions verifies default Options initialization
// Ensures each flag starts as false when initialized
//
// TestNewOptions 验证默认 Options 初始化
// 确保所有标志默认从 false 开始
func TestNewOptions(t *testing.T) {
	options := NewOptions()
	require.False(t, options.IncludeCurrentProject)
	require.False(t, options.IncludeCurrentPackage)
	require.False(t, options.IncludeSubModules)
	require.False(t, options.ExcludeNoGo)
	require.False(t, options.DebugMode)
}

// TestOptionsBuilderPattern tests the option chaining implementation
// Validates that options can be chained and applied as expected
//
// TestOptionsBuilderPattern 测试构建器模式实现
// 验证选项可以被链式调用并正确应用
func TestOptionsBuilderPattern(t *testing.T) {
	options := NewOptions().
		WithIncludeCurrentProject(true).
		WithIncludeCurrentPackage(true).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(true)

	require.True(t, options.IncludeCurrentProject)
	require.True(t, options.IncludeCurrentPackage)
	require.True(t, options.IncludeSubModules)
	require.True(t, options.ExcludeNoGo)
	require.True(t, options.DebugMode)
}

// setupTestGoProject creates a temp test project structure
// Returns temp DIR path with main and submodule projects
//
// setupTestGoProject 创建临时测试项目结构
// 返回带主项目和子模块项目的临时 DIR 路径
func setupTestGoProject(t *testing.T) string {
	tempDIR := rese.V1(os.MkdirTemp("", "test-workspath-*"))

	// Create main project with go.mod and go files
	must.Done(os.WriteFile(filepath.Join(tempDIR, "go.mod"), []byte("module test\n\ngo 1.22.8\n"), 0644))
	must.Done(os.WriteFile(filepath.Join(tempDIR, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0644))

	// Create a submodule DIR with go.mod but no go files
	subDIR := filepath.Join(tempDIR, "submodule")
	must.Done(os.MkdirAll(subDIR, 0755))
	must.Done(os.WriteFile(filepath.Join(subDIR, "go.mod"), []byte("module test/sub\n\ngo 1.22.8\n"), 0644))

	return tempDIR
}

// cleanupTestDIR removes temp test DIR and its contents
// Called in cleanup to ensure test completes without issues
//
// cleanupTestDIR 删除临时测试 DIR 及所有内容
// 通过 defer 调用以确保测试清理
func cleanupTestDIR(t *testing.T, tempDIR string) {
	must.Done(os.RemoveAll(tempDIR))
}
