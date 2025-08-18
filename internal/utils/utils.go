package utils

import (
	"path/filepath"

	"github.com/yyle88/osexistpath/osomitexist"
)

func GetProjectPath(currentPath string) (string, string, bool) {
	projectPath := currentPath
	shortMiddle := ""
	for !osomitexist.IsFile(filepath.Join(projectPath, "go.mod")) {
		subName := filepath.Base(projectPath)

		parentPath := filepath.Dir(projectPath)
		if parentPath == projectPath {
			return "", "", false
		}
		projectPath = parentPath
		shortMiddle = filepath.Join(subName, shortMiddle)
	}
	return projectPath, shortMiddle, true
}
