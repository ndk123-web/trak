package models

import "time"

type WorkspaceMetadata struct {
	Id              string          `json:"id"`
	Template        string          `json:"template"`
	TemplateVersion string          `json:"template_version"`
	CreatedAt       time.Time       `json:"created_at"`
	Author          string          `json:"author"`
	Source          string          `json:"source"`
	Version         string          `json:"version"`
	Repository      string          `json:"repository"`
	Name            string          `json:"name"`
	ModuleBreakdown map[string]bool `json:"module_breakdown"`
}
