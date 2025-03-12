//go:build linux || darwin
// +build linux darwin

package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SetJavaHome 设置JAVA_HOME环境变量
func SetJavaHome(jdkPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// 根据不同的shell修改配置文件
	shellRcFiles := []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".profile"),
	}

	// 读取现有PATH
	currentPath := os.Getenv("PATH")
	paths := strings.Split(currentPath, ":")
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
	newPath := strings.Join(newPaths, ":")

	// 更新shell配置文件
	for _, rcFile := range shellRcFiles {
		if _, err := os.Stat(rcFile); err == nil {
			content, err := os.ReadFile(rcFile)
			if err != nil {
				continue
			}

			lines := strings.Split(string(content), "\n")
			newLines := []string{}
			javaHomeFound := false
			pathFound := false

			for _, line := range lines {
				if strings.HasPrefix(line, "export JAVA_HOME=") {
					newLines = append(newLines, fmt.Sprintf("export JAVA_HOME=%s", jdkPath))
					javaHomeFound = true
				} else if strings.HasPrefix(line, "export PATH=") {
					newLines = append(newLines, fmt.Sprintf("export PATH=%s", newPath))
					pathFound = true
				} else {
					newLines = append(newLines, line)
				}
			}

			if !javaHomeFound {
				newLines = append(newLines, fmt.Sprintf("export JAVA_HOME=%s", jdkPath))
			}
			if !pathFound {
				newLines = append(newLines, fmt.Sprintf("export PATH=%s", newPath))
			}

			err = os.WriteFile(rcFile, []byte(strings.Join(newLines, "\n")), 0644)
			if err != nil {
				return fmt.Errorf("failed to update %s: %w", rcFile, err)
			}
		}
	}

	// 立即更新当前会话的环境变量
	os.Setenv("JAVA_HOME", jdkPath)
	os.Setenv("PATH", newPath)

	return nil
}

// GetJavaHome 获取当前JAVA_HOME环境变量
func GetJavaHome() (string, error) {
	return os.Getenv("JAVA_HOME"), nil
}

// IsValidJDKPath 验证JDK路径是否有效
func IsValidJDKPath(path string) bool {
	// 检查java是否存在
	javaExe := filepath.Join(path, "bin", "java")
	if _, err := os.Stat(javaExe); err != nil {
		return false
	}

	// 检查javac是否存在
	javacExe := filepath.Join(path, "bin", "javac")
	if _, err := os.Stat(javacExe); err != nil {
		return false
	}

	return true
}
