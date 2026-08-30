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
	Use:   "init <category>/<template>",
	Short: "Initialize a hands-on learning workspace from the registry",
	Long: `Download and materialize a comprehensive, multi-module learning workspace
on your local filesystem based on the template registered in Trak Registry.

The template argument must follow the '<category>/<template>' format:
  • lang/<language>   (e.g. lang/go, lang/rust, lang/typescript, lang/python, lang/cpp)
  • os/<os-name>      (e.g. os/linux, os/windows, os/macos)
  • cloud/<provider>  (e.g. cloud/aws)
  • db/<database>     (e.g. db/postgres, db/redis, db/sql)
  • tool/<tool-name>  (e.g. tool/docker, tool/k8s, tool/git, tool/terraform, tool/ansible)

If --path (-p) is not specified, Trak will automatically create a './learn-<toolName>'
directory in your current working directory.`,
	Example: `  # Initialize Go workspace in ./learn-go (current directory):
  trak init lang/go

  # Initialize in a custom path:
  trak init lang/rust --path ./my-rust-track
  trak init db/postgres -p ./postgres-labs

  # Initialize Operating Systems or DevOps tools:
  trak init os/linux
  trak init cloud/aws
  trak init tool/docker --path ./docker-mastery`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateArg := args[0]
		category, toolName, err := helper.ParseTemplateString(templateArg)
		if err != nil {
			ui.Error(err.Error())
			return
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Color("cyan")

		// 1. Fetch Registry Catalog
		s.Suffix = " Querying Trak registry catalog..."
		s.Start()
		tmpl, err := helper.FetchRegistryAndCheck(category, toolName)
		s.Stop()

		if err != nil {
			ui.Error(err.Error())
			return
		}

		ui.Success(fmt.Sprintf("Found Template: %s%s%s (v%s)", ui.Bold, tmpl.Name, ui.Reset, tmpl.Version))

		// 2. Fetch Template Blueprint
		s.Suffix = fmt.Sprintf(" Downloading %s curriculum blueprint...", tmpl.Name)
		s.Start()
		toolTemplate, err := helper.FetchTemplate(category, toolName, tmpl.Source)
		s.Stop()

		if err != nil {
			ui.Error(fmt.Sprintf("Failed to download template: %v", err))
			return
		}

		// 3. Resolve Target Workspace Directory
		var finalPath string
		if targetPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				homeDir, _ := os.UserHomeDir()
				finalPath = filepath.Join(homeDir, fmt.Sprintf("learn-%s", toolName))
			} else {
				finalPath = filepath.Join(cwd, fmt.Sprintf("learn-%s", toolName))
			}
		} else {
			absPath, err := filepath.Abs(targetPath)
			if err != nil {
				ui.Error(fmt.Sprintf("Error resolving target path: %v", err))
				return
			}
			finalPath = absPath
		}

		ui.Info(fmt.Sprintf("Target Workspace: %s%s%s", ui.Cyan, finalPath, ui.Reset))

		// 4. Materialize Files & Folders
		s.Suffix = " Materializing directories and curriculum files on disk..."
		s.Start()
		createdCount, err := generator.GenerateDirectories(toolTemplate, finalPath)
		s.Stop()

		if err != nil {
			ui.Error(fmt.Sprintf("Generation failed: %v", err))
			return
		}

		fmt.Println()
		ui.Success(fmt.Sprintf("Successfully initialized %s%s%s workspace with %d resources! 🎉",
			ui.Bold, tmpl.Name, ui.Reset, createdCount))

		fmt.Println()
		fmt.Printf("  %sNext steps:%s\n", ui.Bold, ui.Reset)
		fmt.Printf("    1. cd %s\n", finalPath)
		fmt.Printf("    2. Open in your editor (%scode .%s)\n", ui.Cyan, ui.Reset)
		fmt.Printf("    3. Read %sREADME.md%s and start Module 00!\n\n", ui.Cyan, ui.Reset)
	},
}

func init() {
	initCmd.Flags().StringVarP(&targetPath, "path", "p", "", "Destination directory path (default: ./learn-<template>)")
	rootCmd.AddCommand(&initCmd)
}
