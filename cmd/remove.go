package cmd

import (
	"fmt"

	"javaman/internal/config"
	"javaman/internal/env"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [version]",
	Aliases: []string{"rm"},
	Short:   "Remove a JDK version",
	Long: `Remove a JDK version from javaman's management.
This command only removes the version from javaman's configuration,
it does not delete the actual JDK installation from your system.

Note: If you remove the currently active version, you'll need to
switch to another version using 'javaman use <version>'.

Examples:
  javaman remove 17   # Remove JDK 17
  javaman rm 8       # Remove JDK 8`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one version argument")
		}

		version := args[0]
		cfg := config.GetConfig()

		// 检查版本是否存在
		path, exists := cfg.Versions[version]
		if !exists {
			return fmt.Errorf("version %s not found", version)
		}

		// 检查是否是当前使用的版本
		currentJavaHome, err := env.GetJavaHome()
		if err == nil && currentJavaHome == path {
			fmt.Printf("Warning: Removing currently active version %s\n", version)
			fmt.Println("You should switch to another version after this operation.")
		}

		// 检查是否是默认版本
		if cfg.Settings.Default == version {
			cfg.Settings.Default = ""
			fmt.Println("Note: Removed version was the default version.")
		}

		// 检查是否有别名指向此版本
		for alias, target := range cfg.Aliases {
			if target == version {
				delete(cfg.Aliases, alias)
				fmt.Printf("Removed alias '%s' that pointed to version %s\n", alias, version)
			}
		}

		// 删除版本
		if err := config.RemoveVersion(version); err != nil {
			return fmt.Errorf("failed to remove version: %w", err)
		}

		fmt.Printf("Successfully removed JDK version %s\n", version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
