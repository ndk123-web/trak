package models

type Node struct {
	Name     string `json:"name"`
	Type     string `json:"type"` // "directory" or "file"
	Content  string `json:"content,omitempty"`
	Children []Node `json:"children,omitempty"`
}

type ToolTemplateModel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Root        Node   `json:"root"`
}

// ParsedTemplate represents a normalized blueprint request (official or community)
type ParsedTemplate struct {
	Author     string `json:"author"`      // "Trak" or GitHub username e.g. "vishal-12"
	Category   string `json:"category"`    // "lang", "os", "cloud", "db", "tool"
	ToolName   string `json:"tool_name"`   // "go", "rust", "k8s", "docker"
	IsOfficial bool   `json:"is_official"` // true if official (templates/)
	SourcePath string `json:"source_path"` // "templates/lang/go.json" or "users/vishal-12/lang/go.json"
	Identifier string `json:"identifier"`  // "lang/go" or "vishal-12/lang/go"
}
