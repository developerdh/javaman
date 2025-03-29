//go:build windows
// +build windows

package detect

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// 常见的JDK安装路径
var commonJDKPaths = []string{
	`C:\Program Files\Java`,
	`C:\Program Files (x86)\Java`,
}

// DetectJDKs 检测系统中已安装的JDK
func DetectJDKs() (map[string]string, error) {
	// 使用map来临时存储每个主版本号对应的所有JDK路径
	tempVersions := make(map[string][]string)
	result := make(map[string]string)

	// 1. 检查常见安装目录
	for _, basePath := range commonJDKPaths {
		if entries, err := os.ReadDir(basePath); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					jdkPath := filepath.Join(basePath, entry.Name())
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

	// 2. 检查注册表
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\JavaSoft\Java Development Kit`, registry.READ)
	if err == nil {
		defer k.Close()
		if subKeys, err := k.ReadSubKeyNames(0); err == nil {
			for _, regVersion := range subKeys {
				subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\JavaSoft\Java Development Kit\`+regVersion, registry.READ)
				if err == nil {
					if path, _, err := subKey.GetStringValue("JavaHome"); err == nil && isValidJDKPath(path) {
						version := NormalizeVersion(regVersion)
						tempVersions[version] = append(tempVersions[version], path)
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
			for _, regVersion := range subKeys {
				subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\JavaSoft\JDK\`+regVersion, registry.READ)
				if err == nil {
					if path, _, err := subKey.GetStringValue("JavaHome"); err == nil && isValidJDKPath(path) {
						version := NormalizeVersion(regVersion)
						tempVersions[version] = append(tempVersions[version], path)
					}
					subKey.Close()
				}
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
	javac := filepath.Join(path, "bin", "javac.exe")
	_, err := os.Stat(javac)
	return err == nil
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
