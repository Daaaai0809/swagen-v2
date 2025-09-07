/*
Copyright © 2025 NAME HERE dai.tsuruga0809@gmail.com
*/
package cmd

import (
	"github.com/Daaaai0809/swagen-v2/libs/model"
	"github.com/Daaaai0809/swagen-v2/utils"
	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate model schema",
	Run: func(cmd *cobra.Command, args []string) {
		input := utils.NewInputMethods()
		modelHandler := model.NewModelHandler(input)
		if err := modelHandler.HandleGenerateModelCommand(); err != nil {
			cmd.PrintErrf("[ERROR] Generating model schema: %v\n", err)
			return
		}
		cmd.Println("[INFO] Model schema generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(modelCmd)
}
