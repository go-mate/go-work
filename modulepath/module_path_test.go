package modulepath

import (
	"testing"

	"github.com/yyle88/runpath"
)

func TestGetModulePaths(t *testing.T) {
	paths := GetModulePaths(runpath.PARENT.Path(), NewOptions().WithDebugMode(true))
	t.Log(paths)
}
