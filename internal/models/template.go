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
