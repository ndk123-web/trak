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

// GenerateDirectories materializes the template workspace and stamps trak.json, returning the total resource count
func GenerateDirectories(toolTemplate *models.ToolTemplateModel, targetPath string) (int, error) {
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

	// 3. Stamp trak.json metadata in the workspace root
	metadata := models.WorkspaceMetadata{
		Template:        toolTemplate.Id,
		TemplateVersion: toolTemplate.Version,
		CreatedAt:       time.Now().UTC(),
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
