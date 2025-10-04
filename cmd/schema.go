package cmd

import (
	"fmt"

	"github.com/Daaaai0809/swagen-v2/fetcher"
	"github.com/Daaaai0809/swagen-v2/handler/schema"
	"github.com/Daaaai0809/swagen-v2/input"
	"github.com/Daaaai0809/swagen-v2/validator"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Generate a Request/Response Schema file",
	Long:  `Interactively generate a schema file for your models.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputMethods := input.NewInputMethods()
		validation := validator.NewInputValidator()
		directoryFetcher := fetcher.NewDirectoryFetcher(inputMethods, validation)
		schemaHandler := schema.NewSchemaHandler(inputMethods, validation, fetcher.NewFileFetcher(), directoryFetcher)

		if err := validation.Validate_Environment_Props(); err != nil {
			return fmt.Errorf("environment validation: %w", err)
		}

		if err := schemaHandler.HandleGenerateSchemaCommand(); err != nil {
			return fmt.Errorf("generating schema: %w", err)
		}
		cmd.Println("[INFO] Schema generated successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}
