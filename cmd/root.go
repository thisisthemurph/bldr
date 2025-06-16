package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = cobra.Command{
	Use:   "bldr",
	Short: "Generate Go builder patterns for structs",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateSingleCmd)
	rootCmd.AddCommand(generateFromYamlCmd)
}
