package cmd

import (
	"fmt"

	"github.com/ndk123-web/trak/internal/config"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "trak",
	Short: "Trak - Developer Learning Workspace Generator",
	Long: fmt.Sprintf(`%s%s  _______ _____            _  __
 |__   __|  __ \     /\   | |/ /
    | |  | |__) |   /  \  | ' / 
    | |  |  _  /   / /\ \ |  <  
    | |  | | \ \  / ____ \| . \ 
    |_|  |_|  \_\/_/    \_\_|\_\%s

%s⚡ Trak%s — Local-first developer learning workspace generator.

Scaffolds structured, multi-module project folders directly onto your machine
complete with hands-on runnable code, exercises, and architectural notes.

Explore 20+ production-grade curricula across Languages, Operating Systems,
Cloud Platforms, Databases, and DevOps Tools.`, ui.Bold, ui.Green, ui.Reset, ui.Bold+ui.Green, ui.Reset),
	Example: `  # Discover all available curriculum tracks:
  trak list

  # Filter by category:
  trak list lang
  trak list db
  trak list os
  trak list cloud
  trak list tool

  # Initialize a workspace in the current directory:
  trak init lang/go
  trak init db/postgres
  trak init tool/docker --path ./my-docker-lab

  # Check workspace progress:
  trak status

  # Check CLI version:
  trak version`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n  %s%s⚡ Trak - Developer Learning Workspace Generator%s\n", ui.Bold, ui.Green, ui.Reset)
		fmt.Printf("  %sLocal-first curriculum materialization CLI (%s)%s\n\n", ui.Gray, config.TrakConfig.Version, ui.Reset)

		fmt.Printf("  %sAvailable Commands:%s\n", ui.Bold, ui.Reset)
		fmt.Printf("    %sdone%s       Mark a curriculum module as completed\n", ui.Green, ui.Reset)
		fmt.Printf("    %sinit%s       Initialize a hands-on learning workspace\n", ui.Green, ui.Reset)
		fmt.Printf("    %slist%s       Explore all available tracks in an interactive tree\n", ui.Green, ui.Reset)
		fmt.Printf("    %sstatus%s     Display workspace progress and module status\n", ui.Green, ui.Reset)
		fmt.Printf("    %sundo%s       Reset or unmark a module back to pending\n", ui.Green, ui.Reset)
		fmt.Printf("    %sversion%s    Display CLI version and build details\n\n", ui.Green, ui.Reset)

		fmt.Printf("  %sQuick Start:%s\n", ui.Bold, ui.Reset)
		fmt.Printf("    %strak list%s                   # Browse all tracks\n", ui.Green, ui.Reset)
		fmt.Printf("    %strak init lang/go%s           # Generate Go learning workspace\n", ui.Green, ui.Reset)
		fmt.Printf("    %strak status%s                 # Check current workspace progress\n", ui.Green, ui.Reset)
		fmt.Printf("    %strak done 00%s                # Mark Module 00 as completed\n", ui.Green, ui.Reset)
		fmt.Printf("    %strak undo 00%s                # Reset Module 00 back to pending\n", ui.Green, ui.Reset)
		fmt.Printf("    %strak init --help%s            # View init options & examples\n\n", ui.Green, ui.Reset)
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
