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

// createNode recursively creates directories and files based on the Node structure
func createNode(basePath string, node models.Node) error {
	targetPath := filepath.Join(basePath, node.Name)

	switch node.Type {
	case "directory":
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
		}
		ui.Success(fmt.Sprintf("Created Directory: %s", targetPath))

		for _, child := range node.Children {
			if err := createNode(targetPath, child); err != nil {
				return err
			}
		}

	case "file":
		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent dir for file %s: %w", targetPath, err)
		}

		if err := os.WriteFile(targetPath, []byte(node.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
		ui.Success(fmt.Sprintf("Created File: %s", targetPath))

	default:
		return fmt.Errorf("unknown node type '%s' for node '%s'", node.Type, node.Name)
	}

	return nil
}

// GenerateDirectories materializes the template workspace and stamps trak.json
func GenerateDirectories(toolTemplate *models.ToolTemplateModel, targetPath string) (bool, error) {
	// 1. Ensure target root directory exists
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		ui.Error(fmt.Sprintf("Failed to create workspace directory: %v", err))
		return false, err
	}

	// 2. Materialize all child nodes recursively
	for _, child := range toolTemplate.Root.Children {
		if err := createNode(targetPath, child); err != nil {
			ui.Error(fmt.Sprintf("Error creating workspace item: %v", err))
			return false, err
		}
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
		ui.Success(fmt.Sprintf("Stamped Metadata: %s", trakJSONPath))
	}

	return true, nil
}
