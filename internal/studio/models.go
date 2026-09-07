package studio

import (
	"io/fs"
)

// StudioFileNode represents a node in the directory tree for Trak Studio.
type StudioFileNode struct {
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	IsDir    bool              `json:"isDir"`
	Size     int64             `json:"size,omitempty"`
	Children []*StudioFileNode `json:"children,omitempty"`
}

// WorkspaceHistoryItem records a recently accessed workspace in global config.
type WorkspaceHistoryItem struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	LastOpened string `json:"lastOpened"`
	TrackId    string `json:"trackId,omitempty"`
}

// TrakGlobalConfig represents the persistent user configuration stored in ~/.trak/trak-config.json.
type TrakGlobalConfig struct {
	Username   string                 `json:"username,omitempty"`
	Email      string                 `json:"email,omitempty"`
	Password   string                 `json:"password,omitempty"`
	Workspaces []WorkspaceHistoryItem `json:"workspaces"`
}

// ServerOptions contains configuration required to start the Trak Studio web server.
type ServerOptions struct {
	Port   string
	DistFS fs.FS
}
