// Package workspath: Go module path discovery and workspace traversal engine
// Auto discovers Go modules in workspace structures with configurable options
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

// GetProjectPath finds the Go project root by traversing up the DIR tree
// Returns project path, relative middle path, and whether it's a Go module
// Delegates to utils package for actual implementation
//
// 通过向上遍历 DIR 树查找 Go 项目根目录
// 返回项目路径、相对中间路径，以及是否为 Go 模块
// 委托给 utils 包进行实际实现
func GetProjectPath(currentPath string) (string, string, bool) {
	return utils.GetProjectPath(currentPath)
}

// GetModulePaths discovers all Go module paths based on the provided options
// Can include current project, submodules, and filter out empty directories
// Uses linked hash set to maintain discovery order and eliminate duplicates
//
// 根据提供的选项发现所有 Go 模块路径
// 可以包含当前项目、子模块，并过滤掉空目录
// 使用链式哈希集合维护发现顺序并消除重复项
func GetModulePaths(currentPath string, options *Options) []string {
	set := linkedhashset.New[string]()

	// Handle current project and package inclusion
	// 处理当前项目和包的包含
	if options.IncludeCurrentProject || options.IncludeCurrentPackage {
		projectPath, shortMiddle, isGoModule := GetProjectPath(currentPath)
		if !isGoModule {
			must.None(projectPath)
			must.None(shortMiddle)
		} else {
			if options.IncludeCurrentProject {
				set.Add(projectPath) // Add project DIR to results // 把项目 DIR 添加到结果里
			}

			if options.IncludeCurrentPackage {
				set.Add(currentPath) // Add current DIR to results // 把当前 DIR 添加到结果里
			}
		}
		if options.DebugMode {
			zaplog.SUG.Debugln(neatjsons.S(set))
		}
	}

	// Discover submodules by walking file tree
	// Current DIR might contain go.mod file, hash-set eliminates duplicates
	//
	// 通过遍历文件树发现子模块
	// 当前 DIR 可能包含 go.mod 文件，哈希集合消除重复项
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

	// Filter out projects without Go source files
	// Some projects have no Go files (empty projects or just containers)
	//
	// 过滤掉没有写 Go 源文件的项目
	// 有的项目没有 Go 文件（空项目或仅为容器）
	if options.ExcludeNoGo {
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

// hasGoFiles checks if a DIR contains any .go source files
// Traverses the DIR tree but stops at nested go.mod boundaries
// Returns true immediately when first .go file is found for efficiency
//
// 检查 DIR 是否包含任何 .go 源文件
// 遍历 DIR 树但在嵌套的 go.mod 边界处停止
// 为提高效率，在找到第一个 .go 文件时立即返回 true
func hasGoFiles(root string) bool {
	existGo := false
	must.Done(filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if exSkip, isHide := isHidePath(info); isHide {
			return exSkip
		}

		if info.IsDir() {
			// Stop when encountering other project's go.mod (exclude current DIR)
			// 遇到其他项目的 go.mod 时停止（排除当前 DIR）
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

// isHidePath determines if a file or DIR should be skipped during traversal
// Hidden files and directories (starting with '.') are skipped automatically
// Returns appropriate skip error for directories vs files
//
// 确定文件或 DIR 在遍历过程中是否应被跳过
// 自动跳过隐藏的文件和目录（以 '.' 开头）
// 为目录与文件返回适当的跳过错误
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

// Options configures the behavior of module path discovery
// Provides fine-grained control over which modules to include
// Supports debug mode for detailed discovery process logging
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

// NewOptions creates a new Options instance with all flags set to false by default
// Use the With* methods to configure the desired behavior
// Provides a clean starting point for building custom discovery configurations
//
// 创建一个新的 Options 实例，默认所有标志都设置为 false
// 使用 With* 方法配置所需的行为
// 为构建自定义发现配置提供干净的起点
func NewOptions() *Options {
	return &Options{
		IncludeCurrentProject: false, // Whether to include current project root DIR // 是否包含当前项目根 DIR
		IncludeCurrentPackage: false, // Whether to include current DIR // 是否包含当前 DIR
		IncludeSubModules:     false, // Whether to include submodule directories // 是否包含子模块目录
		ExcludeNoGo:           false, // Whether to exclude directories without Go files // 是否排除没有 Go 文件的目录
		DebugMode:             false, // Whether to enable debug mode // 是否启用调试模式
	}
}

// WithIncludeCurrentProject sets whether to include the current project root DIR
// Returns the Options instance for method chaining
//
// 设置是否包含当前项目根 DIR
// 返回 Options 实例以支持方法链
func (c *Options) WithIncludeCurrentProject(includeCurrentProject bool) *Options {
	c.IncludeCurrentProject = includeCurrentProject
	return c
}

// WithIncludeCurrentPackage sets whether to include the current DIR
// Returns the Options instance for method chaining
//
// 设置是否包含当前 DIR
// 返回 Options 实例以支持方法链
func (c *Options) WithIncludeCurrentPackage(includeCurrentPackage bool) *Options {
	c.IncludeCurrentPackage = includeCurrentPackage
	return c
}

// WithIncludeSubModules sets whether to include discovered submodule directories
// Returns the Options instance for method chaining
//
// 设置是否包含发现的子模块目录
// 返回 Options 实例以支持方法链
func (c *Options) WithIncludeSubModules(includeSubModules bool) *Options {
	c.IncludeSubModules = includeSubModules
	return c
}

// WithExcludeNoGo sets whether to exclude directories without Go source files
// Returns the Options instance for method chaining
//
// 设置是否排除没有 Go 源文件的目录
// 返回 Options 实例以支持方法链
func (c *Options) WithExcludeNoGo(excludeNoGo bool) *Options {
	c.ExcludeNoGo = excludeNoGo
	return c
}

// WithDebugMode sets whether to enable detailed debug logging
// Returns the Options instance for method chaining
//
// 设置是否启用详细调试日志
// 返回 Options 实例以支持方法链
func (c *Options) WithDebugMode(debugMode bool) *Options {
	c.DebugMode = debugMode
	return c
}
