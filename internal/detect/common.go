package detect

import (
	"strings"
)

// NormalizeVersion 标准化版本号，只保留主版本号
func NormalizeVersion(version string) string {
	// 移除开头的"jdk"、"1."等
	version = strings.ToLower(version)
	version = strings.TrimPrefix(version, "jdk")
	version = strings.TrimPrefix(version, "-")
	version = strings.TrimPrefix(version, "1.")
	version = strings.TrimSpace(version)

	// 只取主版本号
	if idx := strings.Index(version, "."); idx != -1 {
		version = version[:idx]
	}
	if idx := strings.Index(version, "_"); idx != -1 {
		version = version[:idx]
	}

	return version
}

// ExtractVersionFromDirName 从目录名中提取版本号
func ExtractVersionFromDirName(dirName string) string {
	dirName = strings.ToLower(dirName)
	if strings.Contains(dirName, "jdk") {
		// macOS的.jdk后缀处理
		dirName = strings.TrimSuffix(dirName, ".jdk")
		return NormalizeVersion(dirName)
	}
	return ""
}

// CompareVersions 比较两个版本号
func CompareVersions(v1, v2 string) int {
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
