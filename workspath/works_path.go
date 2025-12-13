// Package workspath: Go module path detection and workspace scanning
// Scans workspace structures and detects Go modules with configurable options
//
// workspath: Go 模块路径发现和工作区扫描
// 扫描工作区结构并使用可配置选项发现 Go 模块
package workspath

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/emirpasic/gods/v2/sets/linkedhashset"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/osexistpath/osomitexist"
	"github.com/yyle88/zaplog"
)

// ProjectPath contains discovered Go project information
// ProjectPath 包含已发现的 Go 项目信息
type ProjectPath struct {
	Root    string // Root path of the Go project // Go 项目的根路径
	SubPath string // Relative path from project root to current DIR // 从项目根到当前 DIR 的相对路径
}

// GetProjectPath locates Go project root by traversing up the DIR tree
// Returns ProjectPath with project root and middle path
//
// GetProjectPath 通过向上遍历 DIR 树定位 Go 项目根目录
// 返回包含项目根路径和中间路径的 ProjectPath
func GetProjectPath(path string) (*ProjectPath, bool) {
	root := path
	subs := ""
	for !osomitexist.IsFile(filepath.Join(root, "go.mod")) {
		subName := filepath.Base(root)
		parent := filepath.Dir(root)
		if parent == root {
			return nil, false
		}
		root = parent
		subs = filepath.Join(subName, subs)
	}
	return &ProjectPath{
		Root:    root,
		SubPath: subs,
	}, true
}

// GetProjectRoot locates Go project root by traversing up the DIR tree
// Returns project root path and true if go.mod is found
//
// GetProjectRoot 通过向上遍历 DIR 树定位 Go 项目根目录
// 找到 go.mod 时返回项目根路径和 true
func GetProjectRoot(path string) (string, bool) {
	root := path
	for !osomitexist.IsFile(filepath.Join(root, "go.mod")) {
		parent := filepath.Dir(root)
		if parent == root {
			return "", false
		}
		root = parent
	}
	return root, true
}

// GetModulePaths detects Go module paths starting from path
// Returns slice of module paths based on options
//
// GetModulePaths 从 path 开始发现 Go 模块路径
// 根据选项返回模块路径切片
func GetModulePaths(root string, opts ...Option) []string {
	cfg := &scanConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	set := linkedhashset.New[string]()

	// WithCurrentProject: find project root and add it
	// WithCurrentProject: 查找项目根并添加
	if cfg.currentProject {
		if projectRoot, ok := GetProjectRoot(root); ok {
			set.Add(projectRoot)
		}
	}

	// WithCurrentPackage: add current path itself
	// WithCurrentPackage: 添加当前路径本身
	if cfg.currentPackage {
		set.Add(root)
	}

	if cfg.debugMode {
		zaplog.SUG.Debugln("init:", neatjsons.S(set.Values()))
	}

	if cfg.scanDeep {
		must.Done(filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if isHidden(info) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if !info.IsDir() && info.Name() == "go.mod" {
				if moduleRoot := filepath.Dir(path); osmustexist.IsRoot(moduleRoot) {
					set.Add(moduleRoot)
				}
			}
			return nil
		}))
		if cfg.debugMode {
			zaplog.SUG.Debugln("subs:", neatjsons.S(set.Values()))
		}
	}

	if cfg.skipNoGo {
		set = set.Select(func(idx int, modulePath string) bool {
			return existsGoFiles(modulePath)
		})
		if cfg.debugMode {
			zaplog.SUG.Debugln("skip empty:", neatjsons.S(set.Values()))
		}
	}

	return set.Values()
}

// isHidden checks if path should be skipped (hidden files/dirs)
// isHidden 检查是否应跳过路径（隐藏文件/目录）
func isHidden(info fs.FileInfo) bool {
	return strings.HasPrefix(info.Name(), ".")
}

// existsGoFiles checks if DIR contains .go source files
// existsGoFiles 检查 DIR 是否包含 .go 源文件
func existsGoFiles(root string) bool {
	found := false
	must.Done(filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if isHidden(info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			if path != root && osomitexist.IsFile(filepath.Join(path, "go.mod")) {
				return filepath.SkipDir
			}
		} else if filepath.Ext(info.Name()) == ".go" {
			found = true
			return filepath.SkipAll
		}
		return nil
	}))
	return found
}
