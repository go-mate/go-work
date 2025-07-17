package workspace

import (
	"path/filepath"

	"github.com/yyle88/must"
	"github.com/yyle88/osexistpath/osmustexist"
)

// Workspace 代表一个 Go 工作目录，里面有多个子项目
type Workspace struct {
	WorkRoot string   // 工作区根目录
	Projects []string // 该 Workspace 各个子项目的路径
}

func NewWorkSpace(projects []string) (wsp *Workspace) {
	return NewWorkspace("", projects)
}

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
