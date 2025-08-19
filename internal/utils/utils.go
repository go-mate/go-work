// Package utils: Internal utilities for Go project path discovery and analysis
// Provides core functionality for traversing file system and finding Go modules
// Used by workspath package for module discovery operations
//
// utils: Go 项目路径发现和分析的内部工具
// 提供遍历文件系统和查找 Go 模块的核心功能
// 由 workspath 包用于模块发现操作
package utils

import (
	"path/filepath"

	"github.com/yyle88/osexistpath/osomitexist"
)

// GetProjectPath finds the Go project root by traversing up the DIR tree
// Searches for go.mod file starting from current path and moving up
// Returns project root path, relative middle path, and whether it's a Go module
//
// 通过向上遍历 DIR 树查找 Go 项目根目录
// 从当前路径开始搜索 go.mod 文件并向上移动
// 返回项目根路径、相对中间路径，以及是否为 Go 模块
func GetProjectPath(currentPath string) (string, string, bool) {
	// Start from current path and build middle path as we traverse up
	// 从当前路径开始，在向上遍历时构建中间路径
	projectPath := currentPath
	shortMiddle := ""

	// Keep searching until we find go.mod or reach filesystem root
	// 继续搜索直到找到 go.mod 或达到文件系统根目录
	for !osomitexist.IsFile(filepath.Join(projectPath, "go.mod")) {
		// Get current DIR name for building relative path
		// 获取当前 DIR 名用于构建相对路径
		subName := filepath.Base(projectPath)

		// Move up one level
		// 向上移动一级
		parentPath := filepath.Dir(projectPath)
		if parentPath == projectPath {
			// Reached filesystem root, no go.mod found
			// 达到文件系统根目录，未找到 go.mod
			return "", "", false
		}

		// Update paths for next iteration
		// 更新路径以进行下一次迭代
		projectPath = parentPath
		shortMiddle = filepath.Join(subName, shortMiddle)
	}

	// Found go.mod, return project info
	// 找到 go.mod，返回项目信息
	return projectPath, shortMiddle, true
}
