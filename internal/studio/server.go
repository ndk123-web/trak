package studio

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ndk123-web/trak/internal/ui"
)

// Run initializes the workspace, mounts routes and static SPA assets, and starts the Trak Studio HTTP server.
func Run(opts ServerOptions) {
	port := opts.Port
	if port == "" {
		port = "8200"
	}

	wm := NewWorkspaceManager("")

	_, trackName, hasModules, missingList := wm.InspectActiveWorkspace()
	if hasModules {
		AddWorkspaceToTrakConfig(wm.GetActiveDir())
		if len(missingList) > 0 {
			fmt.Printf("\n  %s! Warning: %d module folders are missing on disk in this workspace.%s\n",
				ui.Yellow+ui.Bold, len(missingList), ui.Reset)
		}
	}

	h := NewHandler(wm)
	mux := http.NewServeMux()

	// Register API Handlers
	mux.HandleFunc("/api/workspace", h.HandleWorkspace)
	mux.HandleFunc("/api/workspaces", h.HandleWorkspacesHistory)
	mux.HandleFunc("/api/status", h.HandleStatus)
	mux.HandleFunc("/api/tree", h.HandleTree)
	mux.HandleFunc("/api/file", h.HandleFile)
	mux.HandleFunc("/api/verify", h.HandleVerify)
	mux.HandleFunc("/api/done", h.HandleDone)
	mux.HandleFunc("/api/item", h.HandleCreateItem)
	mux.HandleFunc("/api/item/delete", h.HandleDeleteItem)
	mux.HandleFunc("/api/browse", h.HandleBrowse)

	// Static UI handler with SPA routing fallback
	if opts.DistFS != nil {
		fileServer := http.FileServer(http.FS(opts.DistFS))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				return
			}
			path := strings.TrimPrefix(r.URL.Path, "/")
			if path != "" {
				if f, err := opts.DistFS.Open(path); err == nil {
					_ = f.Close()
					fileServer.ServeHTTP(w, r)
					return
				}
			}
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
		})
	}

	routeSuffix := "/#/workspaces"
	url := fmt.Sprintf("http://localhost:%s%s", port, routeSuffix)
	netUrl := fmt.Sprintf("http://127.0.0.1:%s%s", port, routeSuffix)

	OpenBrowser(url)

	fmt.Println()
	fmt.Printf("  %sTrak Studio%s running at:\n\n", ui.Green+ui.Bold, ui.Reset)
	fmt.Printf("  > %-10s %s%s%s\n", "Local:", ui.Cyan+ui.Bold, url, ui.Reset)
	fmt.Printf("  > %-10s %s%s%s\n", "Network:", ui.Gray, netUrl, ui.Reset)
	if hasModules {
		fmt.Printf("  > %-10s %s\n", "Track:", trackName)
	} else {
		fmt.Printf("  > %-10s %s\n", "Mode:", "Workspace Hub (No active modules)")
		fmt.Printf("  > %-10s %s\n", "Workspace:", wm.GetActiveDir())
	}
	fmt.Printf("\n  Ready. Press %sCtrl+C%s to stop.\n\n", ui.White+ui.Bold, ui.Reset)

	// Start HTTP Server (blocking)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("%sServer failed:%s %v\n", ui.Red, ui.Reset, err)
	}
}
