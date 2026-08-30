package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ndk123-web/trak/internal/helper"
	"github.com/ndk123-web/trak/internal/ui"
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

		// fmt.Println("Category: ", category)
		// fmt.Println("ToolName: ", toolName)

		s := spinner.New(spinner.CharSets[24], 100*time.Millisecond)
		ui.Info("Fetching Registry")
		s.Start()

		// fetch Registry
		tmpl, err := helper.FetchRegistryAndCheck(category, toolName)

		s.Stop()

		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v\n", err.Error()))
			return
		}

		ui.Success(fmt.Sprintf("Found Template: %s (v%s)\n", tmpl.Name, tmpl.Version))
		// fetch the Template from the
		// https://github.com/%s/%s/blob/%s/templates/%s/%s.json
	},
}

// init runs before main() , speciality of golang
func init() {
	initCmd.Flags().StringVarP(&path, "path", "p", "", "Path where all resources will be going to put")

	rootCmd.AddCommand(&initCmd)
}
