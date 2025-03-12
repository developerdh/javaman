//go:build windows
// +build windows

package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

// SetJavaHome 设置JAVA_HOME环境变量
func SetJavaHome(jdkPath string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Session Manager\Environment`, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open registry key with ALL_ACCESS: %w", err)
	}
	defer key.Close()

	err = key.SetStringValue("JAVA_HOME", jdkPath)
	if err != nil {
		return fmt.Errorf("failed to set JAVA_HOME: %w", err)
	}

	// 更新PATH环境变量
	path, _, err := key.GetStringValue("Path")
	if err != nil {
		return fmt.Errorf("failed to get PATH: %w", err)
	}

	// 分割PATH
	paths := strings.Split(path, ";")
	newPaths := []string{}

	// 移除旧的Java相关路径
	for _, p := range paths {
		if !strings.Contains(strings.ToLower(p), "java") && !strings.Contains(strings.ToLower(p), "jdk") {
			newPaths = append(newPaths, p)
		}
	}

	// 添加新的Java路径
	binPath := filepath.Join(jdkPath, "bin")
	newPaths = append(newPaths, binPath)

	// 更新PATH
	newPath := strings.Join(newPaths, ";")
	err = key.SetStringValue("Path", newPath)
	if err != nil {
		return fmt.Errorf("failed to set PATH: %w", err)
	}

	// 广播环境变量更改消息
	broadcastEnvChange()

	return nil
}

// GetJavaHome 获取当前JAVA_HOME环境变量
func GetJavaHome() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Session Manager\Environment`, registry.READ)
	if err != nil {
		return "", fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	javaHome, _, err := key.GetStringValue("JAVA_HOME")
	if err != nil {
		if err == registry.ErrNotExist {
			return "", nil
		}
		return "", fmt.Errorf("failed to get JAVA_HOME: %w", err)
	}

	return javaHome, nil
}

// IsValidJDKPath 验证JDK路径是否有效
func IsValidJDKPath(path string) bool {
	// 检查java.exe是否存在
	javaExe := filepath.Join(path, "bin", "java.exe")
	if _, err := os.Stat(javaExe); err != nil {
		return false
	}

	// 检查javac.exe是否存在
	javacExe := filepath.Join(path, "bin", "javac.exe")
	if _, err := os.Stat(javacExe); err != nil {
		return false
	}

	return true
}

// 广播环境变量更改消息
func broadcastEnvChange() {
	// 加载user32.dll
	user32, err := syscall.LoadDLL("user32.dll")
	if err != nil {
		return
	}

	// 获取SendMessageTimeout函数
	sendMessageTimeout, err := user32.FindProc("SendMessageTimeoutW")
	if err != nil {
		return
	}

	// 广播WM_SETTINGCHANGE消息
	const (
		HWND_BROADCAST   = 0xFFFF
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Environment"))),
		uintptr(SMTO_ABORTIFHUNG),
		uintptr(1000),
		0,
	)
}
