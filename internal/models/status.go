package models

type StatusModel struct {
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	Version         string          `json:"version"`
	TemplateVersion string          `json:"template_version"`
	CreatedAt       string          `json:"created_at"`
	Template        string          `json:"template"`
	Progress        float32         `json:"progress"`
	Author          string          `json:"author"`
	Source          string          `json:"source"`
	Repository      string          `json:"repository"`
	ModuleBreakdown map[string]bool `json:"module_breakdown"`
}
