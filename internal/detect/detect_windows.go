//go:build windows
// +build windows

package detect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// 常见的JDK安装路径
var commonJDKPaths = []string{
	`C:\Program Files\Java`,
	`C:\Program Files (x86)\Java`,
}

// DetectJDKs 检测系统中已安装的JDK
func DetectJDKs() (map[string]string, error) {
	result := make(map[string]string)

	// 1. 检查常见安装目录
	for _, basePath := range commonJDKPaths {
		if entries, err := os.ReadDir(basePath); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					jdkPath := filepath.Join(basePath, entry.Name())
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

	// 2. 检查注册表
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\JavaSoft\Java Development Kit`, registry.READ)
	if err == nil {
		defer k.Close()
		if subKeys, err := k.ReadSubKeyNames(0); err == nil {
			for _, version := range subKeys {
				subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\JavaSoft\Java Development Kit\`+version, registry.READ)
				if err == nil {
					if path, _, err := subKey.GetStringValue("JavaHome"); err == nil && isValidJDKPath(path) {
						result[version] = path
					}
					subKey.Close()
				}
			}
		}
	}

	// 3. 检查JDK11+的注册表位置
	k, err = registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\JavaSoft\JDK`, registry.READ)
	if err == nil {
		defer k.Close()
		if subKeys, err := k.ReadSubKeyNames(0); err == nil {
			for _, version := range subKeys {
				subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\JavaSoft\JDK\`+version, registry.READ)
				if err == nil {
					if path, _, err := subKey.GetStringValue("JavaHome"); err == nil && isValidJDKPath(path) {
						result[version] = path
					}
					subKey.Close()
				}
			}
		}
	}

	return result, nil
}

// isValidJDKPath 验证路径是否包含有效的JDK
func isValidJDKPath(path string) bool {
	javac := filepath.Join(path, "bin", "javac.exe")
	_, err := os.Stat(javac)
	return err == nil
}

// extractVersion 从目录名中提取版本号
func extractVersion(dirName string) string {
	dirName = strings.ToLower(dirName)
	if strings.Contains(dirName, "jdk") {
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
