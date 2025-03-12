package cmd

import (
	"fmt"

	"javaman/internal/config"
	"javaman/internal/env"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Display current JDK version",
	Long: `Display the currently active JDK version and its installation path.
This command shows:
- Current JAVA_HOME path
- Current JDK version
- Last explicitly selected version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentJavaHome, err := env.GetJavaHome()
		if err != nil {
			return fmt.Errorf("failed to get current JAVA_HOME: %w", err)
		}

		if currentJavaHome == "" {
			fmt.Println("No JDK currently active (JAVA_HOME not set)")
			return nil
		}

		// 获取配置信息
		cfg := config.GetConfig()

		// 查找当前JAVA_HOME对应的版本
		var currentVersion string
		for version, path := range cfg.Versions {
			if path == currentJavaHome {
				currentVersion = version
				break
			}
		}

		fmt.Println("Current Java Environment:")
		fmt.Println("------------------------")
		fmt.Printf("JAVA_HOME: %s\n", currentJavaHome)
		if currentVersion != "" {
			fmt.Printf("Version:    %s\n", currentVersion)
		} else {
			fmt.Printf("Version:    Unknown (path not managed by javaman)\n")
		}

		if cfg.Settings.LastUsed != "" {
			fmt.Printf("\nLast selected version: %s\n", cfg.Settings.LastUsed)
		}

		if cfg.Settings.Default != "" {
			fmt.Printf("Default version:      %s\n", cfg.Settings.Default)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
