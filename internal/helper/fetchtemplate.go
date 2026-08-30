package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ndk123-web/trak/internal/config"
	"github.com/ndk123-web/trak/internal/models"
)

func FetchTemplate(category string, toolName string, source string) (*models.ToolTemplateModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%v%v", config.TrakConfig.RawBaseUrl, source)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned HTTP %d (%s)", resp.StatusCode, resp.Status)
	}

	var toolTemplate *models.ToolTemplateModel

	if err := json.NewDecoder(resp.Body).Decode(&toolTemplate); err != nil {
		return nil, err
	}

	return toolTemplate, nil
}
