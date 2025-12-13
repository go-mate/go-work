package workspath

// scanConfig holds internal scanning configuration
// scanConfig 保存内部扫描配置
type scanConfig struct {
	currentProject bool // Include root if it's a module // 包含根目录（如果是模块）
	currentPackage bool // Include current DIR // 包含当前目录
	scanDeep       bool // Include submodules // 包含子模块
	skipNoGo       bool // Skip modules without Go files // 跳过无 Go 文件的模块
	debugMode      bool // Enable debug logging // 启用调试日志
}

// Option configures scanning behavior
// Option 配置扫描行为
type Option func(*scanConfig)

// WithCurrentProject includes root DIR if it contains go.mod
// WithCurrentProject 包含根目录（如果包含 go.mod）
func WithCurrentProject() Option {
	return func(c *scanConfig) { c.currentProject = true }
}

// WithCurrentPackage includes the current DIR in results
// WithCurrentPackage 在结果中包含当前目录
func WithCurrentPackage() Option {
	return func(c *scanConfig) { c.currentPackage = true }
}

// ScanDeep includes submodules in scan results
// ScanDeep 在扫描结果中包含子模块
func ScanDeep() Option {
	return func(c *scanConfig) { c.scanDeep = true }
}

// SkipNoGo excludes modules without Go source files
// SkipNoGo 排除没有 Go 源文件的模块
func SkipNoGo() Option {
	return func(c *scanConfig) { c.skipNoGo = true }
}

// DebugMode enables detailed debug logging
// DebugMode 启用详细调试日志
func DebugMode() Option {
	return func(c *scanConfig) { c.debugMode = true }
}

// WithDebug sets debug logging based on the provided flag
// WithDebug 根据传入的标志设置调试日志
func WithDebug(debug bool) Option {
	return func(c *scanConfig) { c.debugMode = debug }
}
