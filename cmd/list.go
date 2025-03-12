package cmd

import (
	"fmt"
	"sort"

	"javaman/internal/config"
	"javaman/internal/env"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all installed JDK versions",
	Long: `List all JDK versions that are currently managed by javaman.

This command shows:
- All installed JDK versions and their paths
- Currently active version (marked with *)
- Default version (if set)
- Version aliases (if any)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetConfig()
		currentJavaHome, err := env.GetJavaHome()
		if err != nil {
			return fmt.Errorf("failed to get current JAVA_HOME: %w", err)
		}

		if len(cfg.Versions) == 0 {
			fmt.Println("No JDK versions found.")
			fmt.Println("Use 'javaman add <path>' to add a JDK installation.")
			return nil
		}

		fmt.Println("Available JDK versions:")
		fmt.Println("---------------------")

		// 获取所有版本并排序
		versions := make([]string, 0, len(cfg.Versions))
		for version := range cfg.Versions {
			versions = append(versions, version)
		}
		sort.Strings(versions)

		// 显示所有版本
		for _, version := range versions {
			path := cfg.Versions[version]
			prefix := "  "
			if path == currentJavaHome {
				prefix = "* "
			}
			if version == cfg.Settings.Default {
				prefix = prefix + "[Default] "
			}
			fmt.Printf("%s%-10s -> %s\n", prefix, version, path)
		}

		// 显示别名
		if len(cfg.Aliases) > 0 {
			fmt.Println("\nAliases:")
			fmt.Println("--------")
			aliases := make([]string, 0, len(cfg.Aliases))
			for alias := range cfg.Aliases {
				aliases = append(aliases, alias)
			}
			sort.Strings(aliases)

			for _, alias := range aliases {
				target := cfg.Aliases[alias]
				fmt.Printf("  %-10s -> %s\n", alias, target)
			}
		}

		// 显示当前使用信息
		if currentJavaHome != "" {
			var currentVersion string
			for version, path := range cfg.Versions {
				if path == currentJavaHome {
					currentVersion = version
					break
				}
			}
			fmt.Printf("\nCurrent version: %s\n", currentVersion)
		}

		// 显示默认版本
		if cfg.Settings.Default != "" {
			fmt.Printf("Default version: %s\n", cfg.Settings.Default)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
