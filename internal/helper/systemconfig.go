package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/ndk123-web/trak/internal/models"
)

var configMutex sync.Mutex

func getTrakSystemConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New(err.Error())
	}

	dir := filepath.Join(homeDir, ".trak", "trak-config.json")

	if err = os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return dir, nil
}

func UpdateUserConfig(userConfig *models.UserConfig) (bool, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	systemConfigPath, err := getTrakSystemConfigPath()
	if err != nil {
		return false, err
	}

	dataBytes, err := os.ReadFile(systemConfigPath)
	if err != nil {
		return false, err
	}

	var trakSystemConfig models.TrakUserConfig
	if err = json.NewDecoder(bytes.NewReader(dataBytes)).Decode(&trakSystemConfig); err != nil {
		return false, err
	}

	_ = trakSystemConfig.SetEmail(userConfig.Email).SetPassword(userConfig.Password).SetUsername(userConfig.Username)
	return true, nil
}
