package studio

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ndk123-web/trak/internal/models"
)

// WorkspaceManager provides thread-safe access and validation for the active Trak Studio workspace.
type WorkspaceManager struct {
	mu        sync.RWMutex
	activeDir string
}

// NewWorkspaceManager creates a WorkspaceManager initialized with initialDir or current working directory.
func NewWorkspaceManager(initialDir string) *WorkspaceManager {
	if initialDir == "" {
		initialDir, _ = os.Getwd()
	}
	initialDir = filepath.Clean(initialDir)

	// Check if trak.json is in initialDir; if not, check parent directory
	trakPath := filepath.Join(initialDir, "trak.json")
	if _, err := os.Stat(trakPath); os.IsNotExist(err) {
		parentDir := filepath.Dir(initialDir)
		if _, err := os.Stat(filepath.Join(parentDir, "trak.json")); err == nil {
			initialDir = parentDir
		}
	}

	return &WorkspaceManager{
		activeDir: initialDir,
	}
}

// GetActiveDir returns the current active workspace directory path safely.
func (wm *WorkspaceManager) GetActiveDir() string {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.activeDir
}

// SetActiveDir updates the active workspace directory path safely.
func (wm *WorkspaceManager) SetActiveDir(dir string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.activeDir = filepath.Clean(dir)
}

// InspectActiveWorkspace inspects trak.json and module folders in the current active workspace.
func (wm *WorkspaceManager) InspectActiveWorkspace() (hasTrak bool, trackName string, hasModules bool, missingList []string) {
	dir := wm.GetActiveDir()
	trakPath := filepath.Join(dir, "trak.json")

	var sm models.StatusModel
	dataBytes, err := os.ReadFile(trakPath)
	if err != nil {
		return false, "", false, nil
	}

	if err := json.Unmarshal(dataBytes, &sm); err != nil {
		return false, "", false, nil
	}

	hasTrak = true
	trackName = sm.Name
	if trackName == "" {
		trackName = sm.Template
	}

	if len(sm.ModuleBreakdown) > 0 {
		existingCount := 0
		for modName := range sm.ModuleBreakdown {
			modPath := filepath.Join(dir, modName)
			if stat, err := os.Stat(modPath); err == nil && stat.IsDir() {
				existingCount++
			} else {
				missingList = append(missingList, modName)
			}
		}
		if existingCount > 0 {
			hasModules = true
		}
	}

	return hasTrak, trackName, hasModules, missingList
}

// ValidateTrackWorkspace validates that a candidate directory exists and if trak.json is present, module directories exist.
func ValidateTrackWorkspace(cleanPath string) error {
	cleanPath = strings.Trim(strings.TrimSpace(cleanPath), `"'`)
	cleanPath = filepath.Clean(cleanPath)
	stat, err := os.Stat(cleanPath)
	if err != nil || !stat.IsDir() {
		return fmt.Errorf("directory does not exist: %s", cleanPath)
	}

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
				return fmt.Errorf("invalid track workspace: Found 'trak.json', but none of the %d module folders exist in '%s'", len(sm.ModuleBreakdown), cleanPath)
			}
		}
	}

	return nil
}
