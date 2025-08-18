package workspath

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/emirpasic/gods/v2/sets/linkedhashset"
	"github.com/go-mate/go-work/internal/utils"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/zaplog"
)

// GetProjectPath finds the Go project root by traversing up the DIR tree.
// Returns project path, relative middle path, and whether it's a Go module.
func GetProjectPath(currentPath string) (string, string, bool) {
	return utils.GetProjectPath(currentPath)
}

// GetModulePaths discovers all Go module paths based on the provided options.
// It can include current project, submodules, and filter out empty directories.
func GetModulePaths(currentPath string, options *Options) []string {
	set := linkedhashset.New[string]()

	if options.IncludeCurrentProject || options.IncludeCurrentPackage {
		projectPath, shortMiddle, isGoModule := GetProjectPath(currentPath)
		if !isGoModule {
			must.None(projectPath)
			must.None(shortMiddle)
		} else {
			if options.IncludeCurrentProject {
				set.Add(projectPath) //把项目目录添加到结果里
			}

			if options.IncludeCurrentPackage {
				set.Add(currentPath) //把当前目录添加到结果里
			}
		}
		if options.DebugMode {
			zaplog.SUG.Debugln(neatjsons.S(set))
		}
	}

	//这里很有可能，当前目录下就是 go.mod 文件，就是把当前目录设置两次，因此使用 hash-set 去重复
	if options.IncludeSubModules {
		must.Done(filepath.Walk(currentPath, func(path string, info fs.FileInfo, err error) error {
			if exSkip, isHide := isHidePath(info); isHide {
				return exSkip
			}
			if !info.IsDir() && info.Name() == "go.mod" {
				if subRoot := filepath.Dir(path); osmustexist.IsRoot(subRoot) {
					set.Add(subRoot)
				}
				return nil
			}
			return nil
		}))
		if options.DebugMode {
			zaplog.SUG.Debugln(neatjsons.S(set))
		}
	}

	if options.ExcludeNoGo {
		//但是有些项目里是没有go文件的，比如空项目，或者大项目里只有子项目，而没有逻辑，因此需要去除
		set = set.Select(func(index int, value string) bool {
			if options.DebugMode {
				zaplog.SUG.Debugln(index, value)
			}
			return hasGoFiles(value)
		})
		if options.DebugMode {
			zaplog.SUG.Debugln(neatjsons.S(set))
		}
	}

	roots := set.Values()
	return roots
}

// hasGoFiles checks if a DIR contains any .go source files.
// It traverses the DIR tree but stops at nested go.mod boundaries.
func hasGoFiles(root string) bool {
	existGo := false
	must.Done(filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if exSkip, isHide := isHidePath(info); isHide {
			return exSkip
		}

		if info.IsDir() {
			//当遇到其它项目的 go.mod 时结束（由于传进来的就是项目根目录，因此要排除当前目录）
			if path != root && osmustexist.IsFile(filepath.Join(path, "go.mod")) {
				return filepath.SkipDir
			}
		} else {
			if filepath.Ext(info.Name()) == ".go" {
				existGo = true
				return filepath.SkipAll
			}
		}
		return nil
	}))
	return existGo
}

// isHidePath determines if a file or DIR should be skipped during traversal.
// Hidden files and directories (starting with '.') are skipped.
func isHidePath(info fs.FileInfo) (error, bool) {
	if info.IsDir() {
		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir, true
		}
	} else {
		if strings.HasPrefix(info.Name(), ".") {
			return nil, true
		}
	}
	return nil, false
}

// Options configures the behavior of module path discovery.
type Options struct {
	IncludeCurrentProject bool // 假如当前项目是go项目时，是否包含当前项目的根目录
	IncludeCurrentPackage bool // 假如当前项目是go项目时，是否包含当前目录
	IncludeSubModules     bool // 是否包含子模块的目录
	ExcludeNoGo           bool // 跳过不含go代码的项目
	DebugMode             bool // Enable detailed debug logging
}

// NewOptions creates a new Options instance with all flags set to false by default.
// Use the With* methods to configure the desired behavior.
func NewOptions() *Options {
	return &Options{
		IncludeCurrentProject: false, // Whether to include current project root DIR
		IncludeCurrentPackage: false, // Whether to include current DIR
		IncludeSubModules:     false, // Whether to include submodule directories
		ExcludeNoGo:           false, // Whether to exclude directories without Go files
		DebugMode:             false, // Whether to enable debug mode
	}
}

// WithIncludeCurrentProject sets whether to include the current project root DIR.
func (c *Options) WithIncludeCurrentProject(includeCurrentProject bool) *Options {
	c.IncludeCurrentProject = includeCurrentProject
	return c
}

// WithIncludeCurrentPackage sets whether to include the current DIR.
func (c *Options) WithIncludeCurrentPackage(includeCurrentPackage bool) *Options {
	c.IncludeCurrentPackage = includeCurrentPackage
	return c
}

// WithIncludeSubModules sets whether to include discovered submodule directories.
func (c *Options) WithIncludeSubModules(includeSubModules bool) *Options {
	c.IncludeSubModules = includeSubModules
	return c
}

// WithExcludeNoGo sets whether to exclude directories without Go source files.
func (c *Options) WithExcludeNoGo(excludeNoGo bool) *Options {
	c.ExcludeNoGo = excludeNoGo
	return c
}

// WithDebugMode sets whether to enable detailed debug logging.
func (c *Options) WithDebugMode(debugMode bool) *Options {
	c.DebugMode = debugMode
	return c
}
