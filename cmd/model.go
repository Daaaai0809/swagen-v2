/*
Copyright © 2025 NAME HERE dai.tsuruga0809@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/handler/model"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/validator"
	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate model schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		inputMethods := input.NewInputMethods()
		validation := validator.NewInputValidator()
		directoryFetcher := fetcher.NewDirectoryFetcher(inputMethods, validation)
		modelHandler := model.NewModelHandler(inputMethods, validation, directoryFetcher)

		if err := validation.Validate_Environment_Props(); err != nil {
			return fmt.Errorf("environment validation: %w", err)
		}

		if err := modelHandler.HandleGenerateModelCommand(); err != nil {
			return fmt.Errorf("generating model schema: %w", err)
		}
		cmd.Println("[INFO] Model schema generated successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(modelCmd)
}
