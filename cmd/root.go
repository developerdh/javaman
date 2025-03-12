package cmd

import (
	"fmt"
	"os"

	"javaman/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "javaman",
	Short: "Java version manager",
	Long: `Javaman is a version manager for JDK.
It helps you easily switch between different versions of JDK installed on your system.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
