package modulepath

import (
	"testing"

	"github.com/yyle88/neatjson/neatjsons"
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
}
