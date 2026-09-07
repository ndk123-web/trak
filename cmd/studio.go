package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ndk123-web/trak/internal/helper"
	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/shared"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

//go:embed all:dist
var studioDist embed.FS

var (
	targetPort         string
	studioWorkspaceDir string
	studioWorkspaceMu  sync.RWMutex
)

type StudioFileNode struct {
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	IsDir    bool              `json:"isDir"`
	Size     int64             `json:"size,omitempty"`
	Children []*StudioFileNode `json:"children,omitempty"`
}

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
		port := targetPort
		if port == "" {
			port = "8200"
		}

		// 1. Workspace Discovery
		cwd, _ := os.Getwd()
		studioWorkspaceMu.Lock()
		studioWorkspaceDir = cwd

		// Check if trak.json is in cwd; if not, check parent directory
		trakPath := filepath.Join(studioWorkspaceDir, "trak.json")
		if _, err := os.Stat(trakPath); os.IsNotExist(err) {
			parentDir := filepath.Dir(studioWorkspaceDir)
			if _, err := os.Stat(filepath.Join(parentDir, "trak.json")); err == nil {
				studioWorkspaceDir = parentDir
				trakPath = filepath.Join(studioWorkspaceDir, "trak.json")
			}
		}
		studioWorkspaceMu.Unlock()

		trackName := ""
		hasTrak := false
		var sm models.StatusModel

		if dataBytes, err := os.ReadFile(trakPath); err == nil {
			if err := json.Unmarshal(dataBytes, &sm); err == nil {
				hasTrak = true
				trackName = sm.Name
				if trackName == "" {
					trackName = sm.Template
				}
			}
		}

		// Integrity check: if trak.json exists, verify that module folders exist on disk!
		hasModules := false
		if hasTrak && len(sm.ModuleBreakdown) > 0 {
			existingCount := 0
			var missingList []string
			for modName := range sm.ModuleBreakdown {
				modPath := filepath.Join(studioWorkspaceDir, modName)
				if stat, err := os.Stat(modPath); err == nil && stat.IsDir() {
					existingCount++
				} else {
					missingList = append(missingList, modName)
				}
			}

			if existingCount > 0 {
				hasModules = true
				addWorkspaceToTrakConfig(studioWorkspaceDir)
				if len(missingList) > 0 {
					fmt.Printf("\n  %s! Warning: %d of %d module folders are missing on disk in this workspace.%s\n",
						ui.Yellow+ui.Bold, len(missingList), len(sm.ModuleBreakdown), ui.Reset)
				}
			}
		}

		// Virtual UI filesystem
		distFS, err := fs.Sub(studioDist, "dist")
		if err != nil {
			fmt.Printf("%sError:%s Failed to mount embedded studio UI: %v\n", ui.Red, ui.Reset, err)
			return
		}

		// API & Static routing
		mux := http.NewServeMux()

		// Register API Handlers
		mux.HandleFunc("/api/workspace", handleWorkspace)
		mux.HandleFunc("/api/workspaces", handleWorkspacesHistory)
		mux.HandleFunc("/api/status", handleStatus)
		mux.HandleFunc("/api/tree", handleTree)
		mux.HandleFunc("/api/file", handleFile)
		mux.HandleFunc("/api/verify", handleVerify)
		mux.HandleFunc("/api/done", handleDone)
		mux.HandleFunc("/api/item", handleCreateItem)
		mux.HandleFunc("/api/item/delete", handleDeleteItem)
		mux.HandleFunc("/api/browse", handleBrowse)

		// Static UI handler with SPA routing fallback
		fileServer := http.FileServer(http.FS(distFS))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				return
			}
			path := strings.TrimPrefix(r.URL.Path, "/")
			if path != "" {
				if f, err := distFS.Open(path); err == nil {
					_ = f.Close()
					fileServer.ServeHTTP(w, r)
					return
				}
			}
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
		})

		routeSuffix := "/#/workspaces"

		url := fmt.Sprintf("http://localhost:%s%s", port, routeSuffix)
		netUrl := fmt.Sprintf("http://127.0.0.1:%s%s", port, routeSuffix)

		openBrowser(url)

		fmt.Println()
		fmt.Printf("  %sTrak Studio%s running at:\n\n", ui.Green+ui.Bold, ui.Reset)
		fmt.Printf("  > %-10s %s%s%s\n", "Local:", ui.Cyan+ui.Bold, url, ui.Reset)
		fmt.Printf("  > %-10s %s%s%s\n", "Network:", ui.Gray, netUrl, ui.Reset)
		if hasModules {
			fmt.Printf("  > %-10s %s\n", "Track:", trackName)
		} else {
			fmt.Printf("  > %-10s %s\n", "Mode:", "Workspace Hub (No active modules)")
			fmt.Printf("  > %-10s %s\n", "Workspace:", studioWorkspaceDir)
		}
		fmt.Printf("\n  Ready. Press %sCtrl+C%s to stop.\n\n", ui.White+ui.Bold, ui.Reset)

		// Start HTTP Server (blocking)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			fmt.Printf("%sServer failed:%s %v\n", ui.Red, ui.Reset, err)
		}
	},
}

// CORS helper
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	}
}

func handleBrowse(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}
	path, err := openNativeFolderDialog()
	w.Header().Set("Content-Type", "application/json")
	if err != nil || path == "" {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"path":    "",
		})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"path":    filepath.ToSlash(path),
	})
}

func openNativeFolderDialog() (string, error) {
	if runtime.GOOS == "windows" {
		psScript := `
Add-Type -AssemblyName System.Windows.Forms
$f = New-Object System.Windows.Forms.FolderBrowserDialog
$f.Description = "Select Trak Workspace Directory"
$f.ShowNewFolderButton = $true
if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $f.SelectedPath
}
`
		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("osascript", "-e", `POSIX path of (choose folder with prompt "Select Trak Workspace Directory")`)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	} else {
		cmd := exec.Command("zenity", "--file-selection", "--directory", "--title=Select Trak Workspace Directory")
		out, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(out)), nil
		}
		return "", fmt.Errorf("dialog not available")
	}
}

type WorkspaceHistoryItem struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	LastOpened string `json:"lastOpened"`
	TrackId    string `json:"trackId,omitempty"`
}

type TrakGlobalConfig struct {
	Workspaces []WorkspaceHistoryItem `json:"workspaces"`
}

func getTrakConfigFile() string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		trakDir := filepath.Join(home, ".trak")
		_ = os.MkdirAll(trakDir, 0755)
		return filepath.Join(trakDir, "trak-config.json")
	}
	exe, err := os.Executable()
	if err == nil {
		return filepath.Join(filepath.Dir(exe), "trak-config.json")
	}
	return "trak-config.json"
}

func loadTrakConfig() TrakGlobalConfig {
	cfgFile := getTrakConfigFile()
	var cfg TrakGlobalConfig
	data, err := os.ReadFile(cfgFile)
	if err == nil {
		_ = json.Unmarshal(data, &cfg)
	}
	if cfg.Workspaces == nil {
		cfg.Workspaces = []WorkspaceHistoryItem{}
	}
	return cfg
}

func saveTrakConfig(cfg TrakGlobalConfig) error {
	cfgFile := getTrakConfigFile()
	_ = os.MkdirAll(filepath.Dir(cfgFile), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfgFile, data, 0644)
}

func addWorkspaceToTrakConfig(wsPath string) {
	clean := filepath.Clean(filepath.FromSlash(strings.TrimSpace(wsPath)))
	parts := strings.Split(filepath.ToSlash(clean), "/")
	name := parts[len(parts)-1]
	if name == "" {
		name = "workspace"
	}

	trackId := ""
	trakFile := filepath.Join(clean, "trak.json")
	if tData, err := os.ReadFile(trakFile); err == nil {
		var sm models.StatusModel
		if err := json.Unmarshal(tData, &sm); err == nil {
			trackId = sm.Id
			if sm.Name != "" {
				name = sm.Name
			}
		}
	}

	cfg := loadTrakConfig()
	var updated []WorkspaceHistoryItem
	updated = append(updated, WorkspaceHistoryItem{
		Path:       filepath.ToSlash(clean),
		Name:       name,
		LastOpened: time.Now().UTC().Format(time.RFC3339),
		TrackId:    trackId,
	})
	for _, item := range cfg.Workspaces {
		if strings.EqualFold(filepath.Clean(item.Path), clean) {
			continue
		}
		updated = append(updated, item)
	}
	if len(updated) > 20 {
		updated = updated[:20]
	}
	cfg.Workspaces = updated
	_ = saveTrakConfig(cfg)
}

func removeWorkspaceFromTrakConfig(wsPath string) {
	clean := filepath.Clean(filepath.FromSlash(strings.TrimSpace(wsPath)))
	cfg := loadTrakConfig()
	var updated []WorkspaceHistoryItem
	for _, item := range cfg.Workspaces {
		if strings.EqualFold(filepath.Clean(item.Path), clean) {
			continue
		}
		updated = append(updated, item)
	}
	cfg.Workspaces = updated
	_ = saveTrakConfig(cfg)
}

// 0. /api/workspaces (History persistent in ~/.trak/trak-config.json)
func handleWorkspacesHistory(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		cfg := loadTrakConfig()
		_ = json.NewEncoder(w).Encode(cfg.Workspaces)
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Path string `json:"path"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Path != "" {
			clean := filepath.Clean(strings.TrimSpace(req.Path))
			if stat, err := os.Stat(clean); err == nil && stat.IsDir() {
				addWorkspaceToTrakConfig(clean)
				cfg := loadTrakConfig()
				_ = json.NewEncoder(w).Encode(cfg.Workspaces)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Directory does not exist"})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid path payload"})
		return
	}

	if r.Method == http.MethodDelete {
		targetPath := r.URL.Query().Get("path")
		if targetPath != "" {
			removeWorkspaceFromTrakConfig(targetPath)
			cfg := loadTrakConfig()
			_ = json.NewEncoder(w).Encode(cfg.Workspaces)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Missing path parameter"})
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

// 1. /api/workspace
func handleWorkspace(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Cwd string `json:"cwd"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Cwd != "" {
			cleanPath := strings.Trim(strings.TrimSpace(req.Cwd), `"'`)
			cleanPath = filepath.Clean(cleanPath)
			stat, err := os.Stat(cleanPath)
			if err != nil || !stat.IsDir() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": fmt.Sprintf("Directory does not exist: %s", cleanPath),
				})
				return
			}

			// Check if trak.json is in this directory and verify module folders exist!
			candTrak := filepath.Join(cleanPath, "trak.json")
			if tBytes, err := os.ReadFile(candTrak); err == nil {
				var sm models.StatusModel
				if err := json.Unmarshal(tBytes, &sm); err == nil && len(sm.ModuleBreakdown) > 0 {
					existingMods := 0
					for mod := range sm.ModuleBreakdown {
						if mStat, err := os.Stat(filepath.Join(cleanPath, mod)); err == nil && mStat.IsDir() {
							existingMods++
						}
					}
					if existingMods == 0 {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusBadRequest)
						_ = json.NewEncoder(w).Encode(map[string]string{
							"error": fmt.Sprintf("Invalid track workspace: Found 'trak.json', but none of the %d module folders exist in '%s'.", len(sm.ModuleBreakdown), cleanPath),
						})
						return
					}
				}
			}

			studioWorkspaceMu.Lock()
			studioWorkspaceDir = cleanPath
			studioWorkspaceMu.Unlock()
			addWorkspaceToTrakConfig(cleanPath)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid workspace payload",
			})
			return
		}
	}

	studioWorkspaceMu.RLock()
	activeDir := studioWorkspaceDir
	studioWorkspaceMu.RUnlock()

	trakPath := filepath.Join(activeDir, "trak.json")
	hasTrak := false
	var trakStruct models.StatusModel
	var fileCount int

	if dataBytes, err := os.ReadFile(trakPath); err == nil {
		if err := json.Unmarshal(dataBytes, &trakStruct); err == nil {
			hasTrak = true
		}
	}

	_ = filepath.Walk(activeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := info.Name()
			if base == ".git" || base == "node_modules" || base == "dist" || base == ".vs" || base == ".vscode" {
				return filepath.SkipDir
			}
		} else {
			fileCount++
		}
		return nil
	})

	completedCount := 0
	for _, done := range trakStruct.ModuleBreakdown {
		if done {
			completedCount++
		}
	}

	resp := map[string]interface{}{
		"cwd":              activeDir,
		"hasTrakJson":      hasTrak,
		"activeTrack":      trakStruct.Id,
		"totalFiles":       fileCount,
		"totalModules":     len(trakStruct.ModuleBreakdown),
		"completedModules": completedCount,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// 2. /api/status
func handleStatus(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	studioWorkspaceMu.RLock()
	activeDir := studioWorkspaceDir
	studioWorkspaceMu.RUnlock()

	trakPath := filepath.Join(activeDir, "trak.json")
	dataBytes, err := os.ReadFile(trakPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "trak.json not found in workspace"})
		return
	}

	var status models.StatusModel
	if err := json.Unmarshal(dataBytes, &status); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid trak.json manifest"})
		return
	}

	completed := 0
	total := len(status.ModuleBreakdown)
	for _, done := range status.ModuleBreakdown {
		if done {
			completed++
		}
	}
	if total > 0 {
		status.Progress = float32((float64(completed) / float64(total)) * 100)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

// 3. /api/tree
func handleTree(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	studioWorkspaceMu.RLock()
	activeDir := studioWorkspaceDir
	studioWorkspaceMu.RUnlock()

	tree := buildFileTree(activeDir, "")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tree)
}

// Recursive directory tree builder
func buildFileTree(rootPath, relPath string) []*StudioFileNode {
	dirPath := filepath.Join(rootPath, relPath)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil
	}

	var dirs []*StudioFileNode
	var files []*StudioFileNode

	for _, entry := range entries {
		name := entry.Name()
		childRel := filepath.Join(relPath, name)
		slashPath := filepath.ToSlash(childRel)

		if entry.IsDir() {
			if name == ".git" || name == "node_modules" || name == "dist" ||
				name == "build" || name == ".vs" || name == ".vscode" ||
				name == ".idea" || name == "__pycache__" || name == ".gemini" {
				continue
			}

			children := buildFileTree(rootPath, childRel)
			dirs = append(dirs, &StudioFileNode{
				Name:     name,
				Path:     slashPath,
				IsDir:    true,
				Children: children,
			})
		} else {
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".exe" || ext == ".dll" || ext == ".so" || ext == ".dylib" ||
				ext == ".o" || ext == ".obj" || ext == ".a" || ext == ".bin" ||
				name == ".DS_Store" || strings.HasPrefix(name, "_trak_test") {
				continue
			}

			info, _ := entry.Info()
			var size int64
			if info != nil {
				size = info.Size()
			}

			files = append(files, &StudioFileNode{
				Name:  name,
				Path:  slashPath,
				IsDir: false,
				Size:  size,
			})
		}
	}

	return append(dirs, files...)
}

// Safe path resolution within active workspace directory
func resolveSafeWorkspacePath(activeDir, relPath string) (string, error) {
	trimmed := strings.TrimSpace(relPath)
	trimmed = strings.TrimPrefix(filepath.ToSlash(trimmed), "/")
	if trimmed == "" || trimmed == "." {
		return "", fmt.Errorf("invalid path")
	}

	cleaned := filepath.Clean(filepath.FromSlash(trimmed))
	target := filepath.Join(activeDir, cleaned)

	rel, err := filepath.Rel(activeDir, target)
	if err != nil || strings.HasPrefix(rel, "..") || strings.HasPrefix(rel, "/") || strings.HasPrefix(rel, "\\") {
		return "", fmt.Errorf("access denied: path outside workspace")
	}
	return target, nil
}

// 4. /api/file
func handleFile(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	studioWorkspaceMu.RLock()
	activeDir := studioWorkspaceDir
	studioWorkspaceMu.RUnlock()

	if r.Method == http.MethodGet {
		relPath := r.URL.Query().Get("path")
		if strings.TrimSpace(relPath) == "" {
			http.Error(w, "Missing path parameter", http.StatusBadRequest)
			return
		}
		targetPath, err := resolveSafeWorkspacePath(activeDir, relPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		data, err := os.ReadFile(targetPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("File not found: %v", err), http.StatusNotFound)
			return
		}

		cleanRel, _ := filepath.Rel(activeDir, targetPath)
		ext := strings.TrimPrefix(filepath.Ext(targetPath), ".")
		resp := map[string]interface{}{
			"path":      filepath.ToSlash(cleanRel),
			"content":   string(data),
			"extension": ext,
			"size":      len(data),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Path) == "" {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		targetPath, err := resolveSafeWorkspacePath(activeDir, req.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create parent directory: %v", err), http.StatusInternalServerError)
			return
		}

		if err := os.WriteFile(targetPath, []byte(req.Content), 0644); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
			return
		}

		cleanRel, _ := filepath.Rel(activeDir, targetPath)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"path":    filepath.ToSlash(cleanRel),
			"size":    len(req.Content),
		})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// 5. /api/verify
func handleVerify(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Module string `json:"module"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Module == "" {
		http.Error(w, "Missing module parameter", http.StatusBadRequest)
		return
	}

	studioWorkspaceMu.RLock()
	activeDir := studioWorkspaceDir
	studioWorkspaceMu.RUnlock()

	trakPath := filepath.Join(activeDir, "trak.json")
	dataBytes, err := os.ReadFile(trakPath)
	if err != nil {
		http.Error(w, "trak.json not found", http.StatusNotFound)
		return
	}

	var trakStruct models.StatusModel
	if err := json.Unmarshal(dataBytes, &trakStruct); err != nil {
		http.Error(w, "Failed to parse trak.json", http.StatusInternalServerError)
		return
	}

	templateId := trakStruct.Id
	if templateId == "" {
		templateId = trakStruct.Template
	}
	parsed, err := helper.ParseTemplateString(templateId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template '%s': %v", templateId, err), http.StatusBadRequest)
		return
	}

	start := time.Now()
	var testOutput string
	passed := false

	if parsed.Category != "lang" {
		testOutput = fmt.Sprintf("Notice: '%s' is an architectural laboratory without automated tests.\nMarking module as complete.", req.Module)
		passed = true
	} else {
		resolvedBin, runtimeCfg, err := shared.ResolveToolchain(parsed.ToolName)
		if err != nil {
			testOutput = fmt.Sprintf("Toolchain Error: %v\nPlease make sure %s is installed on your machine and available in PATH.", err, parsed.ToolName)
			passed = false
		} else {
			bin, cmdArgs := runtimeCfg.BuildCommand(resolvedBin, req.Module, activeDir)
			testCmd := exec.Command(bin, cmdArgs...)
			testCmd.Dir = activeDir
			out, testErr := testCmd.CombinedOutput()
			testOutput = string(out)
			passed = (testErr == nil)
			if passed && testOutput == "" {
				testOutput = fmt.Sprintf("=== RUN Module %s\nPASS: All test assertions passed successfully!\nOK", req.Module)
			}
		}
	}

	duration := time.Since(start)

	// If passed, update trak.json on disk
	if passed {
		var rawMap map[string]interface{}
		if err := json.Unmarshal(dataBytes, &rawMap); err == nil {
			mb, ok := rawMap["module_breakdown"].(map[string]interface{})
			if !ok {
				mb = make(map[string]interface{})
				rawMap["module_breakdown"] = mb
			}
			mb[req.Module] = true
			if updatedBytes, err := json.MarshalIndent(rawMap, "", "  "); err == nil {
				_ = os.WriteFile(trakPath, updatedBytes, 0644)
			}
		}
	}

	resp := map[string]interface{}{
		"module":     req.Module,
		"passed":     passed,
		"output":     testOutput,
		"durationMs": int(duration.Milliseconds()),
		"timestamp":  time.Now().Format("15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// 6. /api/done
func handleDone(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Module string `json:"module"`
		Done   bool   `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Module == "" {
		http.Error(w, "Missing module parameter", http.StatusBadRequest)
		return
	}

	studioWorkspaceMu.RLock()
	activeDir := studioWorkspaceDir
	studioWorkspaceMu.RUnlock()

	trakPath := filepath.Join(activeDir, "trak.json")
	dataBytes, err := os.ReadFile(trakPath)
	if err != nil {
		http.Error(w, "trak.json not found", http.StatusNotFound)
		return
	}

	var rawMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &rawMap); err != nil {
		http.Error(w, "Failed to parse trak.json", http.StatusInternalServerError)
		return
	}

	mb, ok := rawMap["module_breakdown"].(map[string]interface{})
	if !ok {
		mb = make(map[string]interface{})
		rawMap["module_breakdown"] = mb
	}
	mb[req.Module] = req.Done

	if updatedBytes, err := json.MarshalIndent(rawMap, "", "  "); err == nil {
		_ = os.WriteFile(trakPath, updatedBytes, 0644)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"module":  req.Module,
		"done":    req.Done,
	})
}

// 7. /api/item (create file / folder)
func handleCreateItem(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path  string `json:"path"`
		IsDir bool   `json:"isDir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Path) == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	studioWorkspaceMu.RLock()
	activeDir := studioWorkspaceDir
	studioWorkspaceMu.RUnlock()

	targetPath, err := resolveSafeWorkspacePath(activeDir, req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if req.IsDir {
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create folder: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create parent directory: %v", err), http.StatusInternalServerError)
			return
		}
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			if err := os.WriteFile(targetPath, []byte(""), 0644); err != nil {
				http.Error(w, fmt.Sprintf("Failed to create file: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	cleanRel, _ := filepath.Rel(activeDir, targetPath)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"path":    filepath.ToSlash(cleanRel),
		"isDir":   req.IsDir,
	})
}

// 8. /api/item/delete
func handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Path) == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	studioWorkspaceMu.RLock()
	activeDir := studioWorkspaceDir
	studioWorkspaceMu.RUnlock()

	targetPath, err := resolveSafeWorkspacePath(activeDir, req.Path)
	if err != nil || targetPath == activeDir {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if err := os.RemoveAll(targetPath); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete item: %v", err), http.StatusInternalServerError)
		return
	}

	cleanRel, _ := filepath.Rel(activeDir, targetPath)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"path":    filepath.ToSlash(cleanRel),
	})
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
	studioCmd.Flags().StringVarP(&targetPort, "port", "p", "", "Custom port for Trak Studio")
	rootCmd.AddCommand(&studioCmd)
}
