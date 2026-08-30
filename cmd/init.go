package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ndk123-web/trak/internal/helper"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var targetPath string

var initCmd = cobra.Command{
	Use:   "init [template]",
	Short: "Initialize the Specified Template into the Given Path, Default Current Directory / User Directory",
	Run: func(cmd *cobra.Command, args []string) {
		template := args[0]
		category, toolName, err := helper.ParseTemplateString(template)

		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v\n", err.Error()))
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

		ui.Success(fmt.Sprintf("Found Template: %s (v%s)", tmpl.Name, tmpl.Version))

		// fetch the Template from the
		// https://github.com/%s/%s/blob/%s/templates/%s/%s.json

		ui.Success("Fetching Template")
		s.Start()

		_, err = helper.FetchTemplate(category, toolName, tmpl.Source)
		s.Stop()

		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		var finalPath string

		// if targetPath empty
		if targetPath == "" {
			absPath, err := os.UserHomeDir()
			if err != nil {
				ui.Error(fmt.Sprintf("Error: %v", err.Error()))
				return
			}

			finalPath = filepath.Join(absPath, fmt.Sprintf("trak-learn-%v", toolName))
		} else {
			absPath, err := filepath.Abs(targetPath)
			if err != nil {
				ui.Error(fmt.Sprintf("Error: %v", err.Error()))
				return
			}

			finalPath = absPath
		}

		ui.Info(fmt.Sprintf("Using Directory: %v", finalPath))
	},
}

// init runs before main() , speciality of golang
func init() {
	initCmd.Flags().StringVarP(&targetPath, "path", "p", "", "Path where all resources will be going to put")

	rootCmd.AddCommand(&initCmd)
}
