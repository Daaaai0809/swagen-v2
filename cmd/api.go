package cmd

import (
	"github.com/Daaaai0809/swagen-v2/handler/api"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "path",
	Short: "Generate an Path file",
	Long:  `Interactively generate an API file for your endpoints.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputMethods := input.NewInputMethods()
		apiHandler := api.NewAPIHandler(inputMethods)
		if err := apiHandler.HandleGenerateAPICommand(); err != nil {
			cmd.PrintErrf("[ERROR] Generating API: %v\n", err)
			return err
		}
		cmd.Println("[INFO] API generated successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
