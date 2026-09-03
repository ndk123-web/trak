package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/ui"
)

// createNode recursively creates directories and files based on the Node structure and returns the count of items created
func createNode(basePath string, node models.Node) (int, error) {
	targetPath := filepath.Join(basePath, node.Name)
	count := 1 // Count this node (file or directory)

	switch node.Type {
	case "directory":
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			return 0, fmt.Errorf("failed to create directory %s: %w", targetPath, err)
		}
		ui.CreatedDir(targetPath)

		for _, child := range node.Children {
			childCount, err := createNode(targetPath, child)
			if err != nil {
				return 0, err
			}
			count += childCount
		}

	case "file":
		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return 0, fmt.Errorf("failed to create parent dir for file %s: %w", targetPath, err)
		}

		if err := os.WriteFile(targetPath, []byte(node.Content), 0644); err != nil {
			return 0, fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
		ui.CreatedFile(targetPath)

	default:
		return 0, fmt.Errorf("unknown node type '%s' for node '%s'", node.Type, node.Name)
	}

	return count, nil
}

// GenerateDirectories materializes the template workspace and stamps trak.json with provenance metadata
func GenerateDirectories(toolTemplate *models.ToolTemplateModel, targetPath string, parsed *models.ParsedTemplate) (int, error) {
	totalResources := 1 // Root directory

	// 1. Ensure target root directory exists
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		ui.Error(fmt.Sprintf("Failed to create workspace directory: %v", err))
		return 0, err
	}

	// 2. Materialize all child nodes recursively
	for _, child := range toolTemplate.Root.Children {
		count, err := createNode(targetPath, child)
		if err != nil {
			ui.Error(fmt.Sprintf("Error creating workspace item: %v", err))
			return 0, err
		}
		totalResources += count
	}

	// 3. Resolve metadata fields based on parsed blueprint origin
	author := "Trak"
	source := fmt.Sprintf("templates/%s.json", toolTemplate.Id)
	id := toolTemplate.Id

	if parsed != nil {
		author = parsed.Author
		source = parsed.SourcePath
		id = parsed.Identifier
	}

	var moduleBreakdown map[string]bool = make(map[string]bool)
	
	for _, node := range toolTemplate.Root.Children {
		if node.Type == "directory" {
			moduleBreakdown[node.Name] = false
		}
	}

	metadata := models.WorkspaceMetadata{
		Id:              id,
		Name:            toolTemplate.Name,
		Version:         toolTemplate.Version,
		Template:        id,
		TemplateVersion: toolTemplate.Version,
		Author:          author,
		Source:          source,
		Repository:      "https://github.com/ndk123-web/trak-registry",
		CreatedAt:       time.Now().UTC(),
		ModuleBreakdown: moduleBreakdown,
	}

	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err == nil {
		trakJSONPath := filepath.Join(targetPath, "trak.json")
		_ = os.WriteFile(trakJSONPath, metadataBytes, 0644)
		ui.StampedMetadata(trakJSONPath)
		totalResources++ // Count trak.json file
	}

	return totalResources, nil
}
