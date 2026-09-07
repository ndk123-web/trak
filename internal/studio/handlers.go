package studio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ndk123-web/trak/internal/helper"
	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/shared"
)

// Handler holds the dependencies for Trak Studio HTTP API controllers.
type Handler struct {
	wm *WorkspaceManager
}

// NewHandler creates a new API controller Handler with the provided WorkspaceManager.
func NewHandler(wm *WorkspaceManager) *Handler {
	return &Handler{
		wm: wm,
	}
}

// EnableCORS sets permissive cross-origin headers for local studio development.
func EnableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	}
}

// HandleBrowse opens a native OS folder picker dialog and returns the selected path.
func (h *Handler) HandleBrowse(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}
	path, err := OpenNativeFolderDialog()
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

// HandleWorkspacesHistory manages recently opened workspaces stored in ~/.trak/trak-config.json.
func (h *Handler) HandleWorkspacesHistory(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		cfg := LoadTrakConfig()
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
				AddWorkspaceToTrakConfig(clean)
				cfg := LoadTrakConfig()
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
			RemoveWorkspaceFromTrakConfig(targetPath)
			cfg := LoadTrakConfig()
			_ = json.NewEncoder(w).Encode(cfg.Workspaces)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Missing path parameter"})
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

// HandleWorkspace gets or switches the current active workspace directory.
func (h *Handler) HandleWorkspace(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
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

			if err := ValidateTrackWorkspace(cleanPath); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}

			h.wm.SetActiveDir(cleanPath)
			AddWorkspaceToTrakConfig(cleanPath)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid workspace payload",
			})
			return
		}
	}

	activeDir := h.wm.GetActiveDir()
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

// HandleStatus returns the active module completion breakdown and track progress.
func (h *Handler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	activeDir := h.wm.GetActiveDir()
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

// HandleTree returns the filtered file and folder tree of the active workspace.
func (h *Handler) HandleTree(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	activeDir := h.wm.GetActiveDir()
	tree := BuildFileTree(activeDir, "")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tree)
}

// HandleFile reads or saves file contents within the active workspace.
func (h *Handler) HandleFile(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method == http.MethodOptions {
		return
	}

	activeDir := h.wm.GetActiveDir()

	if r.Method == http.MethodGet {
		relPath := r.URL.Query().Get("path")
		if strings.TrimSpace(relPath) == "" {
			http.Error(w, "Missing path parameter", http.StatusBadRequest)
			return
		}
		targetPath, err := ResolveSafeWorkspacePath(activeDir, relPath)
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

		targetPath, err := ResolveSafeWorkspacePath(activeDir, req.Path)
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

// HandleVerify runs the native test runner for a module and records the result.
func (h *Handler) HandleVerify(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
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

	activeDir := h.wm.GetActiveDir()
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

// HandleDone manually toggles the completion status of a module.
func (h *Handler) HandleDone(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
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

	activeDir := h.wm.GetActiveDir()
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

// HandleCreateItem creates a new file or directory inside the active workspace.
func (h *Handler) HandleCreateItem(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
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

	activeDir := h.wm.GetActiveDir()
	targetPath, err := ResolveSafeWorkspacePath(activeDir, req.Path)
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

// HandleDeleteItem deletes a file or directory inside the active workspace.
func (h *Handler) HandleDeleteItem(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
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

	activeDir := h.wm.GetActiveDir()
	targetPath, err := ResolveSafeWorkspacePath(activeDir, req.Path)
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
