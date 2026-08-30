package models

import "time"

type WorkspaceMetadata struct {
	Template        string    `json:"template"`
	TemplateVersion string    `json:"template_version"`
	CreatedAt       time.Time `json:"created_at"`
}
