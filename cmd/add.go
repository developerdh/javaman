package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"javaman/internal/config"
	"javaman/internal/env"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Add a new JDK installation",
	Long: `Add a new JDK installation to be managed by javaman.
The path should point to the JDK installation directory (containing the bin folder).

Examples:
  Windows: javaman add "C:\Program Files\Java\jdk-17"
  Linux:   javaman add /usr/lib/jvm/java-17-openjdk-amd64
  macOS:   javaman add /Library/Java/JavaVirtualMachines/jdk-17.jdk/Contents/Home

The version will be detected automatically from the JDK installation.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// 获取绝对路径
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		// 验证JDK路径
		if !env.IsValidJDKPath(absPath) {
			return fmt.Errorf("invalid JDK path: %s", absPath)
		}

		// 从java命令获取版本信息
		version, err := getJavaVersion(absPath)
		if err != nil {
			// 如果无法从命令获取版本，尝试从路径名获取
			version = extractVersionFromPath(absPath)
			if version == "" {
				return fmt.Errorf("could not determine JDK version")
			}
		}

		// 添加到配置
		if err := config.AddVersion(version, absPath); err != nil {
			return fmt.Errorf("failed to add version: %w", err)
		}

		fmt.Printf("Added JDK version %s\n", version)
		fmt.Printf("Path: %s\n", absPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

// getJavaVersion 通过运行java -version命令获取版本信息
func getJavaVersion(jdkPath string) (string, error) {
	javaExe := "java"
	if runtime.GOOS == "windows" {
		javaExe = "java.exe"
	}

	cmd := exec.Command(filepath.Join(jdkPath, "bin", javaExe), "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	versionStr := string(output)
	// 查找版本信息
	lines := strings.Split(versionStr, "\n")
	if len(lines) > 0 {
		firstLine := strings.ToLower(lines[0])
		if strings.Contains(firstLine, "version") {
			parts := strings.Split(firstLine, `"`)
			if len(parts) > 1 {
				version := parts[1]
				// 处理版本号格式
				if strings.HasPrefix(version, "1.") {
					// 对于1.8.0这样的格式，提取主版本号
					dotIndex := strings.Index(version[2:], ".")
					if dotIndex != -1 {
						return version[2 : 2+dotIndex], nil
					}
					return version[2:], nil
				} else {
					// 对于11.0.2这样的格式，提取主版本号
					dotIndex := strings.Index(version, ".")
					if dotIndex != -1 {
						return version[:dotIndex], nil
					}
					return version, nil
				}
			}
		}
	}
	return "", fmt.Errorf("could not parse version information")
}

// extractVersionFromPath 从路径中提取版本号（作为备用方法）
func extractVersionFromPath(path string) string {
	base := filepath.Base(path)
	base = strings.ToLower(base)

	// 处理常见的JDK目录名模式
	if strings.Contains(base, "jdk") {
		// 移除"jdk"和可能的连字符
		base = strings.TrimPrefix(base, "jdk")
		base = strings.TrimPrefix(base, "-")
		base = strings.TrimPrefix(base, "_")
		base = strings.TrimSpace(base)

		// 如果是 1.8 这样的格式，转换为 8
		if strings.HasPrefix(base, "1.") {
			base = base[2:]
		}

		// 提取版本号的第一部分
		parts := strings.Split(base, ".")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}
