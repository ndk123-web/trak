package cmd

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/ndk123-web/trak/internal/studio"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

//go:embed all:dist
var studioDist embed.FS

var targetPort string

var studioCmd = cobra.Command{
	Use:   "studio",
	Short: "Launch the local Trak Studio Web Dashboard",
	Long: fmt.Sprintf(`%s%s⚡ Trak Studio%s — Local-First Visual IDE & Workspace Dashboard.

Starts an embedded HTTP server and mounts a browser-based Single Page Application
connected to your active curriculum workspace via a high-performance local bridge.

Features include:
  • Interactive Curriculum Roadmap & Module Overview
  • Integrated Monaco Code Editor with live disk synchronization
  • Integrated Native Test Runner with real-time pass/fail assertions
  • Multi-Workspace Hub for managing and switching between local tracks
  • Interactive Track Documentation viewer with built-in cheat sheets
  • Live trak.json manifest editor and workspace configuration`, ui.Bold, ui.Green, ui.Reset),
	Example: `  # Launch Studio for the current workspace:
  trak studio

  # Launch on a specific custom port:
  trak studio --port 8500
  trak studio -p 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		distFS, err := fs.Sub(studioDist, "dist")
		if err != nil {
			fmt.Printf("%sError:%s Failed to mount embedded studio UI: %v\n", ui.Red, ui.Reset, err)
			return
		}

		studio.Run(studio.ServerOptions{
			Port:   targetPort,
			DistFS: distFS,
		})
	},
}

func init() {
	studioCmd.Flags().StringVarP(&targetPort, "port", "p", "", "Custom port for Trak Studio")
	rootCmd.AddCommand(&studioCmd)
}
