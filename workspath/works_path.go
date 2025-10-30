// Package workspath: Go module path detection and workspace traversing engine
// Auto detects Go modules in workspace structures with configurable options
// Supports both single projects and complex multi-module workspace setups
//
// workspath: Go 模块路径发现和工作区遍历引擎
// 使用可配置选项自动发现工作区结构中的 Go 模块
// 支持单项目和复杂的多模块工作区设置
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

// PathInfo contains information about a discovered Go project
// Encapsulates project root path, relative middle path, and module status
//
// PathInfo 包含已发现的 Go 项目的信息
// 封装项目根路径、相对中间路径和模块状态
type PathInfo struct {
	ProjectPath string // Root path of the Go project // Go 项目的根路径
	ShortMiddle string // Relative path from project root to current DIR // 从项目根到当前 DIR 的相对路径
}

// GetProjectPath finds the Go project root through traversing up the DIR tree
// Returns PathInfo struct with project path, relative middle path, and module status
// Delegates to utils package with the implementation
//
// 通过向上遍历 DIR 树查找 Go 项目根目录
// 返回包含项目路径、相对中间路径和模块状态的 PathInfo 结构体
// 委托给 utils 包进行实际实现
func GetProjectPath(currentPath string) (*PathInfo, bool) {
	projectPath, shortMiddle, ok := utils.GetProjectPath(currentPath)
	return &PathInfo{
		ProjectPath: projectPath,
		ShortMiddle: shortMiddle,
	}, ok
}

// GetModulePaths detects Go module paths based on the provided options
// Can include current project, submodules, and exclude blank directories
// Uses linked hash set to maintain detection sequence and eliminate duplicates
//
// 根据提供的选项发现所有 Go 模块路径
// 可以包含当前项目、子模块，并过滤掉空目录
// 使用链式哈希集合维护发现顺序并消除重复项
func GetModulePaths(currentPath string, options *Options) []string {
	set := linkedhashset.New[string]()

	// Handle current project and package inclusion
	// 处理当前项目和包的包含
	if options.IncludeCurrentProject || options.IncludeCurrentPackage {
		info, ok := GetProjectPath(currentPath)
		if !ok {
			must.None(info.ProjectPath)
			must.None(info.ShortMiddle)
		} else {
			if options.IncludeCurrentProject {
				set.Add(info.ProjectPath) // Add project DIR to results // 把项目 DIR 添加到结果里
			}

			if options.IncludeCurrentPackage {
				set.Add(currentPath) // Add current DIR to results // 把当前 DIR 添加到结果里
			}
		}
		if options.DebugMode {
			zaplog.SUG.Debugln(neatjsons.S(set))
		}
	}

	// Find submodules through walking file tree
	// Current DIR might contain go.mod file, hash-set eliminates duplicates
	//
	// 通过遍历文件树发现子模块
	// 当前 DIR 可能包含 go.mod 文件，哈希集合消除重复项
	if options.IncludeSubModules {
		must.Done(filepath.Walk(currentPath, func(path string, info fs.FileInfo, err error) error {
			if skipErr, shouldSkip := skipHiddenPath(info); shouldSkip {
				return skipErr
			}
			if !info.IsDir() && info.Name() == "go.mod" {
				if moduleRoot := filepath.Dir(path); osmustexist.IsRoot(moduleRoot) {
					set.Add(moduleRoot)
				}
				return nil
			}
			return nil
		}))
		if options.DebugMode {
			zaplog.SUG.Debugln(neatjsons.S(set))
		}
	}

	// Exclude projects without Go source files
	// Some projects have no Go files (workspace roots, organizing directories)
	//
	// 过滤掉没有写 Go 源文件的项目
	// 有的项目没有 Go 文件（工作区根目录、组织性目录）
	if options.ExcludeNoGo {
		set = set.Select(func(idx int, modulePath string) bool {
			if options.DebugMode {
				zaplog.SUG.Debugln(idx, modulePath)
			}
			return containsGoFiles(modulePath)
		})
		if options.DebugMode {
			zaplog.SUG.Debugln(neatjsons.S(set))
		}
	}

	modulePaths := set.Values()
	return modulePaths
}

// containsGoFiles checks if a DIR contains .go source files
// Traverses the DIR tree but stops at nested go.mod boundaries
// Returns true when first .go file is found
//
// 检查 DIR 是否包含任何 .go 源文件
// 遍历 DIR 树但在嵌套的 go.mod 边界处停止
// 找到第一个 .go 文件时返回 true
func containsGoFiles(rootPath string) bool {
	found := false
	must.Done(filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if skipErr, shouldSkip := skipHiddenPath(info); shouldSkip {
			return skipErr
		}

		if info.IsDir() {
			// Stop when encountering nested project's go.mod (exclude current DIR)
			// 遇到其他项目的 go.mod 时停止（排除当前 DIR）
			if path != rootPath && osmustexist.IsFile(filepath.Join(path, "go.mod")) {
				return filepath.SkipDir
			}
		} else {
			if filepath.Ext(info.Name()) == ".go" {
				found = true
				return filepath.SkipAll
			}
		}
		return nil
	}))
	return found
}

// skipHiddenPath checks if a file should be skipped when traversing
// Hidden files and directories (starting with '.') are skipped
// Returns skip status and boolean flag
//
// 检查文件或 DIR 在遍历过程中是否应被跳过
// 跳过隐藏的文件和目录（以 '.' 开头）
// 返回跳过错误和布尔标志
func skipHiddenPath(info fs.FileInfo) (error, bool) {
	if strings.HasPrefix(info.Name(), ".") {
		if info.IsDir() {
			return filepath.SkipDir, true
		}
		return nil, true
	}
	return nil, false
}

// Options configures the settings of module path detection
// Provides fine-grained selection of which modules to include
// Supports debug mode with detailed detection process logging
//
// Options 配置模块路径发现行为
// 提供对包含哪些模块的细粒度控制
// 支持调试模式以进行详细的发现过程日志记录
type Options struct {
	IncludeCurrentProject bool // Include current project root DIR if it's a Go project // 如果当前项目是 Go 项目，是否包含当前项目的根 DIR
	IncludeCurrentPackage bool // Include current DIR if it's within a Go project // 如果当前 DIR 在 Go 项目内，是否包含当前 DIR
	IncludeSubModules     bool // Include discovered submodule directories // 是否包含发现的子模块 DIR
	ExcludeNoGo           bool // Skip projects without Go source files // 跳过不含 Go 源文件的项目
	DebugMode             bool // Enable detailed debug logging // 启用详细调试日志
}

// NewOptions creates a new Options instance with flags defaulting to false
// Use the With* methods to set the needed options
// Provides a clean starting point when building custom detection setups
//
// 创建一个新的 Options 实例，默认所有标志都设置为 false
// 使用 With* 方法配置所需的行为
// 为构建自定义发现配置提供干净的起点
func NewOptions() *Options {
	return &Options{
		IncludeCurrentProject: false, // If true, includes current project root DIR // 是否包含当前项目根 DIR
		IncludeCurrentPackage: false, // If true, includes current DIR // 是否包含当前 DIR
		IncludeSubModules:     false, // If true, includes submodule directories // 是否包含子模块目录
		ExcludeNoGo:           false, // If true, excludes directories without Go files // 是否排除没有 Go 文件的目录
		DebugMode:             false, // If true, enables debug mode // 是否启用调试模式
	}
}

// WithIncludeCurrentProject sets if the current project root DIR should be included
// Returns the Options instance to support method chaining
//
// 设置是否包含当前项目根 DIR
// 返回 Options 实例以支持方法链
func (o *Options) WithIncludeCurrentProject(includeCurrentProject bool) *Options {
	o.IncludeCurrentProject = includeCurrentProject
	return o
}

// WithIncludeCurrentPackage sets if the current DIR should be included
// Returns the Options instance to support method chaining
//
// 设置是否包含当前 DIR
// 返回 Options 实例以支持方法链
func (o *Options) WithIncludeCurrentPackage(includeCurrentPackage bool) *Options {
	o.IncludeCurrentPackage = includeCurrentPackage
	return o
}

// WithIncludeSubModules sets if detected submodule directories should be included
// Returns the Options instance to support method chaining
//
// 设置是否包含发现的子模块目录
// 返回 Options 实例以支持方法链
func (o *Options) WithIncludeSubModules(includeSubModules bool) *Options {
	o.IncludeSubModules = includeSubModules
	return o
}

// WithExcludeNoGo sets if directories without Go source files should be excluded
// Returns the Options instance to support method chaining
//
// 设置是否排除没有 Go 源文件的目录
// 返回 Options 实例以支持方法链
func (o *Options) WithExcludeNoGo(excludeNoGo bool) *Options {
	o.ExcludeNoGo = excludeNoGo
	return o
}

// WithDebugMode sets if detailed debug logging should be enabled
// Returns the Options instance to support method chaining
//
// 设置是否启用详细调试日志
// 返回 Options 实例以支持方法链
func (o *Options) WithDebugMode(debugMode bool) *Options {
	o.DebugMode = debugMode
	return o
}
