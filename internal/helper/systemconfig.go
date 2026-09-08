package helper

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/ndk123-web/trak/internal/models"
)

var configMutex sync.Mutex

func matchEmailRegex(email string) bool {
	pattern, err := regexp.Compile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if err != nil {
		return false
	}
	return pattern.MatchString(email)
}

func getTrakSystemConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New(err.Error())
	}

	dir := filepath.Join(homeDir, ".trak")

	if err = os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	dir = filepath.Join(dir, "trak-config.json")
	return dir, nil
}

func UpdateUserConfig(userConfig *models.UserConfig) (bool, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	systemConfigPath, err := getTrakSystemConfigPath()
	if err != nil {
		return false, err
	}

	var trakSystemConfig models.TrakUserConfig
	dataBytes, err := os.ReadFile(systemConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
	} else if len(dataBytes) > 0 {
		_ = json.Unmarshal(dataBytes, &trakSystemConfig)
	}

	if trakSystemConfig.Workspaces == nil {
		trakSystemConfig.Workspaces = []models.SystemConfigWorkspaceModel{}
	}

	if userConfig.Email != "" && !matchEmailRegex(userConfig.Email) {
		return false, errors.New("Error: Email is Invalid")
	}

	_ = trakSystemConfig.SetEmail(userConfig.Email).SetPassword(userConfig.Password).SetUsername(userConfig.Username)

	newDataBytes, err := json.MarshalIndent(trakSystemConfig, "", "  ")
	if err != nil {
		return false, err
	}

	if err = os.WriteFile(systemConfigPath, newDataBytes, 0644); err != nil {
		return false, err
	}

	return true, nil
}
