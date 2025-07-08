package modulepath

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

func GetProjectPath(currentPath string) (string, string, bool) {
	return utils.GetProjectPath(currentPath)
}

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

type Options struct {
	IncludeCurrentProject bool // 假如当前项目是go项目时，是否包含当前项目的根目录
	IncludeCurrentPackage bool // 假如当前项目是go项目时，是否包含当前目录
	IncludeSubModules     bool // 是否包含子模块的目录
	ExcludeNoGo           bool
	DebugMode             bool
}

func NewOptions() *Options {
	return &Options{
		IncludeCurrentProject: false, //是否包含当前项目的根目录
		IncludeCurrentPackage: false, //是否包含当前目录
		IncludeSubModules:     false, //是否包含子模块的目录
		ExcludeNoGo:           false, //是否包含没有 go 文件的目录
		DebugMode:             false,
	}
}

func (c *Options) WithIncludeCurrentProject(includeCurrentProject bool) *Options {
	c.IncludeCurrentProject = includeCurrentProject
	return c
}

func (c *Options) WithIncludeCurrentPackage(includeCurrentPackage bool) *Options {
	c.IncludeCurrentPackage = includeCurrentPackage
	return c
}

func (c *Options) WithIncludeSubModules(includeSubModules bool) *Options {
	c.IncludeSubModules = includeSubModules
	return c
}

func (c *Options) WithExcludeNoGo(excludeNoGo bool) *Options {
	c.ExcludeNoGo = excludeNoGo
	return c
}

func (c *Options) WithDebugMode(debugMode bool) *Options {
	c.DebugMode = debugMode
	return c
}
