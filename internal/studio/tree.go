package studio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BuildFileTree recursively constructs the filtered file and folder hierarchy for the Studio UI.
func BuildFileTree(rootPath, relPath string) []*StudioFileNode {
	dirPath := filepath.Join(rootPath, relPath)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil
	}

	var dirs []*StudioFileNode
	var files []*StudioFileNode

	for _, entry := range entries {
		name := entry.Name()
		childRel := filepath.Join(relPath, name)
		slashPath := filepath.ToSlash(childRel)

		if entry.IsDir() {
			if name == ".git" || name == "node_modules" || name == "dist" ||
				name == "build" || name == ".vs" || name == ".vscode" ||
				name == ".idea" || name == "__pycache__" || name == ".gemini" {
				continue
			}

			children := BuildFileTree(rootPath, childRel)
			dirs = append(dirs, &StudioFileNode{
				Name:     name,
				Path:     slashPath,
				IsDir:    true,
				Children: children,
			})
		} else {
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".exe" || ext == ".dll" || ext == ".so" || ext == ".dylib" ||
				ext == ".o" || ext == ".obj" || ext == ".a" || ext == ".bin" ||
				name == ".DS_Store" || strings.HasPrefix(name, "_trak_test") {
				continue
			}

			info, _ := entry.Info()
			var size int64
			if info != nil {
				size = info.Size()
			}

			files = append(files, &StudioFileNode{
				Name:  name,
				Path:  slashPath,
				IsDir: false,
				Size:  size,
			})
		}
	}

	return append(dirs, files...)
}

// ResolveSafeWorkspacePath ensures the requested path stays strictly jailed within the active workspace.
func ResolveSafeWorkspacePath(activeDir, relPath string) (string, error) {
	trimmed := strings.TrimSpace(relPath)
	trimmed = strings.TrimPrefix(filepath.ToSlash(trimmed), "/")
	if trimmed == "" || trimmed == "." {
		return "", fmt.Errorf("invalid path")
	}

	cleaned := filepath.Clean(filepath.FromSlash(trimmed))
	target := filepath.Join(activeDir, cleaned)

	rel, err := filepath.Rel(activeDir, target)
	if err != nil || strings.HasPrefix(rel, "..") || strings.HasPrefix(rel, "/") || strings.HasPrefix(rel, "\\") {
		return "", fmt.Errorf("access denied: path outside workspace")
	}
	return target, nil
}
