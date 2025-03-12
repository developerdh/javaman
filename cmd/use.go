package cmd

import (
	"fmt"

	"javaman/internal/config"
	"javaman/internal/env"

	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to specified JDK version",
	Long: `Switch to a specific JDK version that is managed by javaman.

This command will:
1. Set JAVA_HOME to the specified JDK installation directory
2. Update system PATH to include the JDK's bin directory
3. Update the last used version in configuration

Examples:
  javaman use 17    # Switch to JDK 17
  javaman use 8     # Switch to JDK 8
  javaman use lts   # Switch to version aliased as 'lts'

Note: On Windows, this command requires administrator privileges.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]

		// 获取版本对应的路径
		cfg := config.GetConfig()
		jdkPath, exists := cfg.Versions[version]
		if !exists {
			// 检查是否是别名
			if aliasPath, ok := cfg.Aliases[version]; ok {
				jdkPath = cfg.Versions[aliasPath]
				version = aliasPath
			} else {
				return fmt.Errorf("version %s not found. Use 'javaman list' to see available versions", version)
			}
		}

		// 验证JDK路径
		if !env.IsValidJDKPath(jdkPath) {
			return fmt.Errorf("invalid JDK path: %s", jdkPath)
		}

		// 设置环境变量
		if err := env.SetJavaHome(jdkPath); err != nil {
			return fmt.Errorf("failed to set JAVA_HOME: %w", err)
		}

		// 更新last_used
		cfg.Settings.LastUsed = version
		if err := config.SaveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Successfully switched to JDK %s\n", version)
		fmt.Printf("JAVA_HOME: %s\n", jdkPath)
		fmt.Printf("You may need to restart your terminal for changes to take effect.\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
