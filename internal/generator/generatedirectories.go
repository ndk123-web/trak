package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/ui"
)

func GenerateDirectories(toolTemplate *models.ToolTemplateModel, targetPath string) (bool, error) {

	name := toolTemplate.Name
	description := toolTemplate.Description
	root := toolTemplate.Root

	childrens := root.Children

	for _, val := range childrens {

		switch val.Type {
		case "file":
			err := os.WriteFile(filepath.Join(targetPath, val.Name), []byte(val.Content), 0755)
			if err != nil {
				ui.Error(fmt.Sprintf("%v", err.Error()))
				return false, err
			}
			ui.Success(fmt.Sprintf("Created File: %v", filepath.Join(targetPath, val.Name)))

		case "directory":
			parentDir := filepath.Join(targetPath, val.Name)

			err := os.Mkdir(parentDir, 0755)

			if err != nil {
				ui.Error(fmt.Sprintf("%v", err.Error()))
				return false, err
			}
			ui.Success(fmt.Sprintf("Created Directory: %v", fmt.Sprintf("%v/%v", targetPath, val.Name)))

			// create files
			for _, file := range val.Children {
				err := os.WriteFile(filepath.Join(targetPath, file.Name), []byte(file.Content), 0755)
				if err != nil {
					ui.Error(fmt.Sprintf("%v", err.Error()))
					return false, err
				}
				ui.Success(fmt.Sprintf("Created File: %v", filepath.Join(targetPath, file.Name)))
			}
		}
	}

	ui.Success(fmt.Sprintf("Name: %v", name))
	ui.Success(fmt.Sprintf("Description: %v", description))

	return true, nil
}
