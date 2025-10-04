package cmd

import (
	"fmt"

	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/handler/api"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/validator"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "path",
	Short: "Generate an Path file",
	Long:  `Interactively generate an API file for your endpoints.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		isAddMode, err := cmd.Flags().GetBool("add")
		if err != nil {
			return err
		}

		inputMethods := input.NewInputMethods()
		inputValidator := validator.NewInputValidator()
		propsValidator := validator.NewPropsValidator()
		directoryFetcher := fetcher.NewDirectoryFetcher(inputMethods, inputValidator)
		apiHandler := api.NewAPIHandler(inputMethods, inputValidator, fetcher.NewFileFetcher(), directoryFetcher)

		if err := propsValidator.Validate_Environment_Props(); err != nil {
			return err
		}

		switch {
		case isAddMode:
			if err := apiHandler.HandleAddToAPICommand(); err != nil {
				return fmt.Errorf("adding to API: %w", err)
			}
			cmd.Println("[INFO] Added to API successfully.")
			return nil
		default:
			if err := apiHandler.HandleGenerateAPICommand(); err != nil {
				return fmt.Errorf("generating API: %w", err)
			}
			cmd.Println("[INFO] API generated successfully.")
			return nil
		}
	},
}

func init() {
	apiCmd.Flags().Bool("add", false, "Add to existing API file if it exists")

	rootCmd.AddCommand(apiCmd)
}
