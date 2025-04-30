//go:build linux || darwin
// +build linux darwin

package detect

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 常见的JDK安装路径
var commonJDKPaths = map[string][]string{
	"linux": {
		"/usr/lib/jvm",
		"/usr/java",
		"/opt/java",
	},
	"darwin": {
		"/Library/Java/JavaVirtualMachines",
		"/System/Library/Java/JavaVirtualMachines",
	},
}

// DetectJDKs 检测系统中已安装的JDK
func DetectJDKs() (map[string]string, error) {
	// 使用map来临时存储每个主版本号对应的所有JDK路径
	tempVersions := make(map[string][]string)
	result := make(map[string]string)

	// 1. 检查常见安装目录
	paths := commonJDKPaths[runtime.GOOS]
	for _, basePath := range paths {
		if entries, err := os.ReadDir(basePath); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					var jdkPath string
					if runtime.GOOS == "darwin" {
						// macOS的JDK路径结构: /Library/Java/JavaVirtualMachines/jdk-17.jdk/Contents/Home
						jdkPath = filepath.Join(basePath, entry.Name(), "Contents", "Home")
					} else {
						jdkPath = filepath.Join(basePath, entry.Name())
					}

					if isValidJDKPath(jdkPath) {
						version := ExtractVersionFromDirName(entry.Name())
						if version != "" {
							tempVersions[version] = append(tempVersions[version], jdkPath)
						}
					}
				}
			}
		}
	}

	// 2. 检查环境变量中的Java路径
	if javaPath, err := exec.LookPath("java"); err == nil {
		if realPath, err := filepath.EvalSymlinks(javaPath); err == nil {
			// 获取JAVA_HOME路径（bin目录的父目录）
			javaHome := filepath.Dir(filepath.Dir(realPath))
			if isValidJDKPath(javaHome) {
				// 尝试获取版本信息
				if version := getJavaVersion(filepath.Join(javaHome, "bin", "java")); version != "" {
					tempVersions[version] = append(tempVersions[version], javaHome)
				}
			}
		}
	}

	// 3. 检查JAVA_HOME环境变量
	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		if isValidJDKPath(javaHome) {
			if version := getJavaVersion(filepath.Join(javaHome, "bin", "java")); version != "" {
				tempVersions[version] = append(tempVersions[version], javaHome)
			}
		}
	}

	// 选择每个版本的最佳路径
	for version, paths := range tempVersions {
		// 如果有多个路径，选择最后一个（通常是最新安装的）
		if len(paths) > 0 {
			result[version] = paths[len(paths)-1]
		}
	}

	return result, nil
}

// isValidJDKPath 验证路径是否包含有效的JDK
func isValidJDKPath(path string) bool {
	java := filepath.Join(path, "bin", "java")
	if _, err := os.Stat(java); err != nil {
		return false
	}

	// 验证java是否可以正常运行并返回版本信息
	cmd := exec.Command(filepath.Join(path, "bin", "java"), "-version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// getJavaVersion 通过运行java -version命令获取版本信息
func getJavaVersion(javaPath string) string {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	// 解析版本信息
	versionStr := string(output)
	lines := strings.Split(versionStr, "\n")
	if len(lines) > 0 {
		firstLine := strings.ToLower(lines[0])
		if strings.Contains(firstLine, "version") {
			parts := strings.Split(firstLine, `"`)
			if len(parts) > 1 {
				return NormalizeVersion(parts[1])
			}
		}
	}
	return ""
}

// GetLatestJDK 获取最新的JDK版本和路径
func GetLatestJDK() (version string, path string, err error) {
	jdks, err := DetectJDKs()
	if err != nil {
		return "", "", err
	}
	if len(jdks) == 0 {
		return "", "", fmt.Errorf("no JDK installations found")
	}

	// 找出最高版本
	var highestVersion string
	for v := range jdks {
		if highestVersion == "" || CompareVersions(v, highestVersion) > 0 {
			highestVersion = v
		}
	}

	return highestVersion, jdks[highestVersion], nil
}
