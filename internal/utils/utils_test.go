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

func TestGetProjectPath(t *testing.T) {
	path := runpath.PARENT.Path()
	t.Log(path)
	projectPath, shortMiddle, isGoModule := GetProjectPath(path)
	require.True(t, isGoModule)
	t.Log(projectPath)
	t.Log(shortMiddle)

	require.Equal(t, path, filepath.Join(projectPath, shortMiddle))
}

func TestGetProjectPathFromRoot(t *testing.T) {
	// Get the actual project root by going up from utils DIR
	utilsPath := runpath.PARENT.Path() // This is internal/utils
	projectPath, shortMiddle, isGoModule := GetProjectPath(utilsPath)
	require.True(t, isGoModule)

	// Now test from the discovered project root
	rootProjectPath, rootShortMiddle, rootIsGoModule := GetProjectPath(projectPath)
	require.True(t, rootIsGoModule)
	require.Equal(t, projectPath, rootProjectPath)
	require.Empty(t, rootShortMiddle, "should be empty for project root")

	// Verify the utils path relationship
	require.Equal(t, "internal/utils", shortMiddle)
	require.Equal(t, utilsPath, filepath.Join(projectPath, shortMiddle))
}

func TestGetProjectPathNonGoModule(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "test-non-go-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Create a DIR without go.mod
	subDIR := filepath.Join(tempDIR, "subdir")
	must.Done(os.MkdirAll(subDIR, 0755))

	projectPath, shortMiddle, isGoModule := GetProjectPath(subDIR)
	require.False(t, isGoModule)
	require.Empty(t, projectPath)
	require.Empty(t, shortMiddle)
}

func TestGetProjectPathDeepNesting(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "test-deep-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Create go.mod at root
	must.Done(os.WriteFile(filepath.Join(tempDIR, "go.mod"), []byte("module test\n"), 0644))

	// Create deeply nested DIR
	deepPath := filepath.Join(tempDIR, "a", "b", "c", "d")
	must.Done(os.MkdirAll(deepPath, 0755))

	projectPath, shortMiddle, isGoModule := GetProjectPath(deepPath)
	require.True(t, isGoModule)
	require.Equal(t, tempDIR, projectPath)
	require.Equal(t, filepath.Join("a", "b", "c", "d"), shortMiddle)
	require.Equal(t, deepPath, filepath.Join(projectPath, shortMiddle))
}
