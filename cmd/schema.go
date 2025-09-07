package cmd

import (
	"github.com/Daaaai0809/swagen-v2/libs/schema"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Generate a Request/Response Schema file",
	Long:  `Interactively generate a schema file for your models.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputMethods := utils.NewInputMethods()
		schemaHandler := schema.NewSchemaHandler(inputMethods)
		if err := schemaHandler.HandleGenerateSchemaCommand(); err != nil {
			cmd.PrintErrf("[ERROR] Generating schema: %v\n", err)
			return err
		}
		cmd.Println("[INFO] Schema generated successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}
