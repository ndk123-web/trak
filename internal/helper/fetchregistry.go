package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// "net/http"
	"time"

	"github.com/ndk123-web/trak/internal/config"
	"github.com/ndk123-web/trak/internal/models"
)

func FetchRegistryAndCheck(category string, toolName string) (*models.TemplateModel, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	registryURL := fmt.Sprintf("%sregistry.json", config.TrakConfig.RawBaseUrl)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		registryURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned HTTP %d (%s)", resp.StatusCode, resp.Status)
	}

	var registry models.RegistryModel
	err = json.NewDecoder(resp.Body).Decode(&registry)
	if err != nil {
		return nil, fmt.Errorf("failed to parse registry json: %w", err)
	}

	cat, exists := registry.Categories[category]
	if !exists {
		return nil, fmt.Errorf("category '%s' not found in registry", category)
	}

	template, exists := cat.Templates[toolName]
	if !exists {
		return nil, fmt.Errorf("template '%s' not found in category '%s'", toolName, category)
	}

	return &template, nil
}
