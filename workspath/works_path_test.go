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

func TestGetModulePaths(t *testing.T) {
	parentPath := runpath.PARENT.Path()
	t.Log(parentPath)

	options := NewOptions().
		WithIncludeCurrentProject(true).
		WithIncludeCurrentPackage(true).
		WithIncludeSubModules(true).
		WithExcludeNoGo(true).
		WithDebugMode(true)
	t.Log(neatjsons.S(options))

	paths := GetModulePaths(parentPath, options)
	t.Log(neatjsons.S(paths))

	require.NotEmpty(t, paths, "should find at least one module path")
	require.Contains(t, paths, runpath.PARENT.Up(1), "should include project root")
}

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
			require.Len(t, paths, tt.expected, "should have expected number of paths")
		})
	}
}

func TestNewOptions(t *testing.T) {
	options := NewOptions()
	require.False(t, options.IncludeCurrentProject)
	require.False(t, options.IncludeCurrentPackage)
	require.False(t, options.IncludeSubModules)
	require.False(t, options.ExcludeNoGo)
	require.False(t, options.DebugMode)
}

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

func cleanupTestDIR(t *testing.T, tempDIR string) {
	must.Done(os.RemoveAll(tempDIR))
}
