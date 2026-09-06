package cmd

import (
	"embed"
	// "encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	// "os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

//go:embed all:dist
var studioDist embed.FS

var (
	targetPort string
)

var studioCmd = cobra.Command{
	Use:   "studio",
	Short: "Launch the local Trak Studio Web Dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		port := targetPort
		if port == "" {
			port = "8200"
		}

		mux := http.NewServeMux()

		// 1. API Handlers
		// mux.HandleFunc("/api/workspace", handleWorkspace)
		// mux.HandleFunc("/api/status", handleStatus)
		// mux.HandleFunc("/api/tree", handleTree)
		// mux.HandleFunc("/api/file", handleFileContent)
		// mux.HandleFunc("/api/verify", handleVerify)

		// 2. Embedded Static UI Handler
		distFS, err := fs.Sub(studioDist, "dist")
		if err != nil {
			fmt.Printf("Error loading studio UI: %v\n", err)
			return
		}
		mux.Handle("/", http.FileServer(http.FS(distFS)))

		url := fmt.Sprintf("http://localhost:%s", port)
		fmt.Printf("\n🚀 Trak Studio running at %s\n", url)
		fmt.Println("Press Ctrl+C to stop.")

		// 3. Auto-open default browser
		openBrowser(url)

		// 4. Start HTTP Server
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			fmt.Printf("Server failed: %v\n", err)
		}
	},
}

// Browser auto-opener utility
func openBrowser(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

func init() {
	studioCmd.Flags().StringVarP(&targetPort, "port", "p", "8200", "Custom port for Trak Studio")
	rootCmd.AddCommand(&studioCmd)
}
