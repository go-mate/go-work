// Package utils: Tests core utilities in Go project path detection
// Validates project root detection through traversing DIR trees
// Tests edge cases such as non-Go directories and deep nesting
//
// utils: 测试 Go 项目路径发现的内部工具
// 验证通过遍历 DIR 树检测项目根目录
// 测试边界情况如非 Go 目录和深层嵌套
package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

// TestGetProjectPath tests project root detection from nested package
// Verifies correct path resolution and middle path calculation
//
// TestGetProjectPath 测试从子目录发现项目根目录
// 验证正确的路径解析和中间路径计算
func TestGetProjectPath(t *testing.T) {
	path := runpath.PARENT.Path()
	t.Log(path)
	projectPath, shortMiddle, ok := GetProjectPath(path)
	require.True(t, ok)
	t.Log(projectPath)
	t.Log(shortMiddle)

	require.Equal(t, path, filepath.Join(projectPath, shortMiddle))
}

// TestGetProjectPathFromRoot tests project path detection from project root
// Verifies that calling from root returns blank middle path
//
// TestGetProjectPathFromRoot 测试从项目根目录发现项目路径
// 验证从根目录调用返回空的中间路径
func TestGetProjectPathFromRoot(t *testing.T) {
	// Get the project root going up from utils DIR
	utilsPath := runpath.PARENT.Path() // This is core/utils
	projectPath, shortMiddle, ok := GetProjectPath(utilsPath)
	require.True(t, ok)

	{ // Now test from the discovered project root
		rootPath, rootMiddle, rootOk := GetProjectPath(projectPath)
		require.True(t, rootOk)
		require.Equal(t, projectPath, rootPath)
		require.Empty(t, rootMiddle, "should be blank when called from project root")
	}

	// Check the utils path relationship
	require.Equal(t, "internal/utils", shortMiddle)
	require.Equal(t, utilsPath, filepath.Join(projectPath, shortMiddle))
}

// TestGetProjectPathNonGoModule tests handling of non-Go directories
// Verifies that function returns false when no go.mod is found
//
// TestGetProjectPathNonGoModule 测试非 Go 目录的行为
// 验证当找不到 go.mod 时函数返回 false
func TestGetProjectPathNonGoModule(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "test-non-go-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Create a DIR without go.mod
	subPath := filepath.Join(tempDIR, "sub")
	must.Done(os.MkdirAll(subPath, 0755))

	projectPath, shortMiddle, ok := GetProjectPath(subPath)
	require.False(t, ok)
	require.Empty(t, projectPath)
	require.Empty(t, shortMiddle)
}

// TestGetProjectPathDeepNesting tests detection from multi-depth nested directories
// Verifies correct middle path calculation across multiple depths
//
// TestGetProjectPathDeepNesting 测试从深层嵌套目录的发现
// 验证跨多个层级的正确中间路径计算
func TestGetProjectPathDeepNesting(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "test-deep-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Create go.mod at root
	must.Done(os.WriteFile(filepath.Join(tempDIR, "go.mod"), []byte("module test\n"), 0644))

	// Create multi-depth nested DIR
	deepPath := filepath.Join(tempDIR, "a", "b", "c", "d")
	must.Done(os.MkdirAll(deepPath, 0755))

	projectPath, shortMiddle, ok := GetProjectPath(deepPath)
	require.True(t, ok)
	require.Equal(t, tempDIR, projectPath)
	require.Equal(t, filepath.Join("a", "b", "c", "d"), shortMiddle)
	require.Equal(t, deepPath, filepath.Join(projectPath, shortMiddle))
}
