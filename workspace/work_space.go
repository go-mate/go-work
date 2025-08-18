package workspace

import (
	"path/filepath"

	"github.com/yyle88/must"
	"github.com/yyle88/osexistpath/osmustexist"
)

// Workspace represents a Go workspace DIR containing multiple subprojects.
// It provides structure for managing multiple Go modules within a single workspace.
// Workspace 代表一个 Go 工作目录，里面有多个子项目
type Workspace struct {
	WorkRoot string   // 工作区根目录
	Projects []string // 该 Workspace 各个子项目的路径
}

// NewWorkSpace creates a new workspace without a root DIR.
// This is an alias for NewWorkspace with empty workRoot for convenience.
func NewWorkSpace(projects []string) (wsp *Workspace) {
	return NewWorkspace("", projects)
}

// NewWorkspace creates a new workspace with the specified root DIR and projects.
// If workRoot is provided, it validates the existence of go.work file.
// All project paths are validated to ensure they contain go.mod files.
func NewWorkspace(workRoot string, projects []string) (wsp *Workspace) {
	if workRoot != "" {
		osmustexist.MustRoot(workRoot)
		osmustexist.MustFile(filepath.Join(workRoot, "go.work"))
	}
	for _, path := range must.Have(projects) {
		osmustexist.MustRoot(path)
		osmustexist.MustFile(filepath.Join(path, "go.mod"))
	}
	return &Workspace{
		WorkRoot: workRoot,
		Projects: projects,
	}
}
