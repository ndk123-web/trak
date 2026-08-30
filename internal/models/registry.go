package models

type TemplateModel struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Source      string   `json:"source"`
	Tags        []string `json:"tags"`
}

type Category struct {
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Templates   map[string]TemplateModel `json:"templates"`
}

type RegistryModel struct {
	SchemaVersion string              `json:"schema_version"`
	UpdatedAt     string              `json:"updated_at"`
	Categories    map[string]Category `json:"categories"`
}
