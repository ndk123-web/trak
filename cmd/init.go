package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ndk123-web/trak/internal/generator"
	"github.com/ndk123-web/trak/internal/helper"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var targetPath string

var initCmd = cobra.Command{
	Use:   "init [template]",
	Short: "Initialize the specified learning template into the workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template := args[0]
		category, toolName, err := helper.ParseTemplateString(template)
		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		s := spinner.New(spinner.CharSets[24], 100*time.Millisecond)
		ui.Info("Fetching Registry...")
		s.Start()

		// 1. Fetch Registry
		tmpl, err := helper.FetchRegistryAndCheck(category, toolName)
		s.Stop()

		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		ui.Success(fmt.Sprintf("Found Template: %s (v%s)", tmpl.Name, tmpl.Version))

		// 2. Fetch Template Blueprint
		ui.Info("Fetching Template Blueprint...")
		s.Start()

		toolTemplate, err := helper.FetchTemplate(category, toolName, tmpl.Source)
		s.Stop()

		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		// 3. Resolve Target Workspace Directory
		var finalPath string
		if targetPath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				ui.Error(fmt.Sprintf("Error determining home directory: %v", err.Error()))
				return
			}
			finalPath = filepath.Join(homeDir, fmt.Sprintf("trak-learn-%s", toolName))
		} else {
			absPath, err := filepath.Abs(targetPath)
			if err != nil {
				ui.Error(fmt.Sprintf("Error resolving target path: %v", err.Error()))
				return
			}
			finalPath = absPath
		}

		ui.Info(fmt.Sprintf("Target Workspace: %s", finalPath))

		// 4. Generate Files & Folders
		s.Start()
		_, err = generator.GenerateDirectories(toolTemplate, finalPath)
		s.Stop()

		if err != nil {
			return
		}

		fmt.Println()
		ui.Success(fmt.Sprintf("Successfully initialized %s workspace at %s! 🎉", tmpl.Name, finalPath))
	},
}

func init() {
	initCmd.Flags().StringVarP(&targetPath, "path", "p", "", "Destination directory path for the learning workspace")
	rootCmd.AddCommand(&initCmd)
}
