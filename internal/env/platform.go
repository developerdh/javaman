package env

import "runtime"

// IsWindows 检查当前是否是Windows系统
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsDarwin 检查当前是否是macOS系统
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux 检查当前是否是Linux系统
func IsLinux() bool {
	return runtime.GOOS == "linux"
}
