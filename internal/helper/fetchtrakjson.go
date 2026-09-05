package helper

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ndk123-web/trak/internal/models"
)

func FetchTrakJsonWorkspace() (*models.WorkspaceMetadata, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	trakPath := filepath.Join(cwd, "trak.json")

	if _, err = os.Stat(trakPath); err != nil {
		return nil, err
	}

	dataBytes, err := os.ReadFile(trakPath)
	if err != nil {
		return nil, err
	}

	var trakStruct models.WorkspaceMetadata
	if err = json.NewDecoder(bytes.NewReader(dataBytes)).Decode(&trakStruct); err != nil {
		return nil, err
	}

	return &trakStruct, err
}
