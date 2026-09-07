package studio

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ndk123-web/trak/internal/models"
)

// GetTrakConfigFile returns the absolute path to the global trak-config.json file.
func GetTrakConfigFile() string {
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

// LoadTrakConfig loads the user's global workspace history and preferences.
func LoadTrakConfig() TrakGlobalConfig {
	cfgFile := GetTrakConfigFile()
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

// SaveTrakConfig writes the global config to disk.
func SaveTrakConfig(cfg TrakGlobalConfig) error {
	cfgFile := GetTrakConfigFile()
	_ = os.MkdirAll(filepath.Dir(cfgFile), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfgFile, data, 0644)
}

// AddWorkspaceToTrakConfig registers or updates a workspace in the global history.
func AddWorkspaceToTrakConfig(wsPath string) {
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

	cfg := LoadTrakConfig()
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
	_ = SaveTrakConfig(cfg)
}

// RemoveWorkspaceFromTrakConfig deletes a workspace entry from global history.
func RemoveWorkspaceFromTrakConfig(wsPath string) {
	clean := filepath.Clean(filepath.FromSlash(strings.TrimSpace(wsPath)))
	cfg := LoadTrakConfig()
	var updated []WorkspaceHistoryItem
	for _, item := range cfg.Workspaces {
		if strings.EqualFold(filepath.Clean(item.Path), clean) {
			continue
		}
		updated = append(updated, item)
	}
	cfg.Workspaces = updated
	_ = SaveTrakConfig(cfg)
}
