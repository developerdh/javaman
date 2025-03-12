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
						version := extractVersion(entry.Name())
						if version != "" {
							result[version] = jdkPath
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
					result[version] = javaHome
				}
			}
		}
	}

	// 3. 检查JAVA_HOME环境变量
	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		if isValidJDKPath(javaHome) {
			if version := getJavaVersion(filepath.Join(javaHome, "bin", "java")); version != "" {
				result[version] = javaHome
			}
		}
	}

	return result, nil
}

// isValidJDKPath 验证路径是否包含有效的JDK
func isValidJDKPath(path string) bool {
	javac := filepath.Join(path, "bin", "javac")
	_, err := os.Stat(javac)
	return err == nil
}

// extractVersion 从目录名中提取版本号
func extractVersion(dirName string) string {
	dirName = strings.ToLower(dirName)
	if strings.Contains(dirName, "jdk") {
		// 移除文件扩展名（针对macOS的.jdk后缀）
		dirName = strings.TrimSuffix(dirName, ".jdk")
		// 移除"jdk"并清理版本号
		version := strings.TrimPrefix(dirName, "jdk")
		version = strings.TrimPrefix(version, "-")
		version = strings.TrimSpace(version)
		// 处理一些常见的版本号格式
		if strings.HasPrefix(version, "1.") {
			version = version[2:] // 将"1.8"转换为"8"
		}
		return version
	}
	return ""
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
				version := parts[1]
				// 处理版本号格式
				var ver string
				if strings.HasPrefix(version, "1.") {
					if dotIndex := strings.Index(version[2:], "."); dotIndex != -1 {
						ver = version[2 : 2+dotIndex]
					} else {
						ver = version[2:]
					}
				} else {
					if dotIndex := strings.Index(version, "."); dotIndex != -1 {
						ver = version[:dotIndex]
					} else {
						ver = version
					}
				}
				if ver != "" {
					return ver
				}
				return version
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
		if highestVersion == "" || compareVersions(v, highestVersion) > 0 {
			highestVersion = v
		}
	}

	return highestVersion, jdks[highestVersion], nil
}

// compareVersions 比较两个版本号
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		if parts1[i] > parts2[i] {
			return 1
		} else if parts1[i] < parts2[i] {
			return -1
		}
	}

	if len(parts1) > len(parts2) {
		return 1
	} else if len(parts1) < len(parts2) {
		return -1
	}
	return 0
}
