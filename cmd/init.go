package cmd

import (
	"fmt"

	"github.com/ndk123-web/trak/internal/helper"
	"github.com/spf13/cobra"
)

var path string

var initCmd = cobra.Command{
	Use:   "init [template]",
	Short: "Initialize the Specified Template into the Given Path, Default Current Directory / User Directory",
	Run: func(cmd *cobra.Command, args []string) {
		template := args[0]
		category, toolName, err := helper.ParseTemplateString(template)

		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return
		}

		fmt.Println("All is Right!!")
		fmt.Println(category)
		fmt.Println(toolName)

		if path != "" {
			fmt.Println("Path: ", path)
		}
	},
}

func init() {
	initCmd.Flags().StringVarP(&path, "path", "p", "", "Path where all resources will be going to put")
	rootCmd.AddCommand(&initCmd)
}
