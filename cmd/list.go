package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/ndk123-web/trak/internal/helper"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var (
	categoryFlag string
	fetchAllFlag bool
)

var listCmd = cobra.Command{
	Use:   "list [category]",
	Short: "List all available learning templates in a tree structure",
	Long: `Explore the complete Trak catalog of 20+ production-grade learning blueprints
organized across Programming Languages, Operating Systems, Cloud Providers, Databases,
and DevOps Tools.

You can inspect the entire catalog or drill down into a specific category.`,
	Example: `  # Browse full catalog tree:
  trak list
  trak list --all

  # Browse specific category:
  trak list lang          # Programming languages (go, rust, python, etc.)
  trak list os            # Operating systems (linux, windows, macos)
  trak list cloud         # Cloud platforms (aws)
  trak list db            # Databases (postgres, redis, sql)
  trak list tool          # DevOps tools (docker, k8s, terraform, ansible, git, jenkins)

  # Using flags:
  trak list -c db
  trak list --category os`,
	Run: func(cmd *cobra.Command, args []string) {
		targetCategory := ""

		// Check if positional argument provided e.g. "trak list db"
		if len(args) > 0 {
			targetCategory = args[0]
		} else if cmd.Flags().Changed("category") {
			targetCategory = categoryFlag
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Fetching Trak registry catalog..."
		s.Color("cyan")
		s.Start()

		err := helper.SearchRegistry(targetCategory, fetchAllFlag)
		s.Stop()

		if err != nil {
			ui.Error(err.Error())
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(&categoryFlag, "category", "c", "", "Filter templates by category (lang, os, cloud, db, tool)")
	listCmd.Flags().BoolVarP(&fetchAllFlag, "all", "a", false, "List all categories and templates")

	rootCmd.AddCommand(&listCmd)
}
