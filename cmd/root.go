/*
Copyright © 2025 NAME HERE dai.tsuruga0809@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "swagen",
	Long: `
This is a CLI application that helps your OpenAPI schema definition.
You can generate API endpoint schemas, models, and other related files.
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
