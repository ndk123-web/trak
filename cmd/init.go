package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ndk123-web/trak/internal/generator"
	"github.com/ndk123-web/trak/internal/helper"
	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var targetPath string

var initCmd = cobra.Command{
	Use:   "init <blueprint>",
	Short: "Initialize a hands-on learning workspace from the registry",
	Long: `Download and materialize a comprehensive, multi-module learning workspace
on your local filesystem based on blueprints in Trak Registry.

Supports both official and community creator blueprints:
  • Official Short:     lang/go, db/postgres, tool/docker, os/linux
  • Official Explicit:  trak/lang/go, trak/tool/docker
  • Community Creator:  <username>/<category>/<tool> (e.g. vishal-12/lang/go)

If --path (-p) is not specified, Trak will automatically create a './learn-<toolName>'
directory in your current working directory.`,
	Example: `  # Initialize official Go workspace in ./learn-go (current directory):
  trak init lang/go

  # Initialize explicit official track:
  trak init trak/lang/rust --path ./my-rust-track

  # Initialize community track by creator:
  trak init vishal-12/lang/go -p ./advanced-go-microservices

  # Initialize Operating Systems or DevOps tools:
  trak init os/linux
  trak init cloud/aws
  trak init tool/docker --path ./docker-mastery`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateArg := args[0]
		parsed, err := helper.ParseTemplateString(templateArg)
		if err != nil {
			ui.Error(err.Error())
			return
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Color("cyan")

		var toolTemplate *models.ToolTemplateModel
		displayName := parsed.ToolName
		_ = displayName

		if parsed.IsOfficial {
			// 1. Fetch Official Registry Catalog
			s.Suffix = " Querying Trak official catalog..."
			s.Start()
			tmpl, err := helper.FetchRegistryAndCheck(parsed.Category, parsed.ToolName)
			s.Stop()
			fmt.Print("\r\033[K")

			if err != nil {
				ui.Error(err.Error())
				return
			}

			displayName = tmpl.Name
			ui.FoundTemplate(tmpl.Name, tmpl.Version)

			// 2. Fetch Official Template Blueprint
			s.Suffix = fmt.Sprintf(" Downloading %s official blueprint...", tmpl.Name)
			s.Start()
			toolTemplate, err = helper.FetchTemplate(parsed.SourcePath)
			s.Stop()
			fmt.Print("\r\033[K")

			if err != nil {
				ui.Error(fmt.Sprintf("Failed to download template: %v", err))
				return
			}
		} else {
			// Community Creator Track
			s.Suffix = fmt.Sprintf(" Downloading community blueprint from @%s (%s)...", parsed.Author, parsed.SourcePath)
			s.Start()
			toolTemplate, err = helper.FetchTemplate(parsed.SourcePath)
			s.Stop()
			fmt.Print("\r\033[K")

			if err != nil {
				ui.Error(fmt.Sprintf("Failed to download community blueprint: %v", err))
				return
			}

			displayName = fmt.Sprintf("%s (by @%s)", toolTemplate.Name, parsed.Author)
			ui.FoundTemplate(displayName, toolTemplate.Version)
		}

		// 3. Resolve Target Workspace Directory
		var finalPath string
		if targetPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				homeDir, _ := os.UserHomeDir()
				finalPath = filepath.Join(homeDir, fmt.Sprintf("learn-%s", parsed.ToolName))
			} else {
				finalPath = filepath.Join(cwd, fmt.Sprintf("learn-%s", parsed.ToolName))
			}
		} else {
			absPath, err := filepath.Abs(targetPath)
			if err != nil {
				ui.Error(fmt.Sprintf("Error resolving target path: %v", err))
				return
			}
			finalPath = absPath
		}

		ui.TargetWorkspace(finalPath)

		// 4. Materialize Files & Folders
		createdCount, err := generator.GenerateDirectories(toolTemplate, finalPath)
		if err != nil {
			ui.Error(fmt.Sprintf("Generation failed: %v", err))
			return
		}

		ui.CompletedBanner(displayName, createdCount)
		ui.NextSteps(finalPath)
	},
}

func init() {
	initCmd.Flags().StringVarP(&targetPath, "path", "p", "", "Destination directory path (default: ./learn-<template>)")
	rootCmd.AddCommand(&initCmd)
}
