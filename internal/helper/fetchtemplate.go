package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ndk123-web/trak/internal/config"
	"github.com/ndk123-web/trak/internal/models"
)

// FetchTemplate downloads the AST JSON blueprint from the registry source path
func FetchTemplate(source string) (*models.ToolTemplateModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	candidates := []string{source}
	if strings.HasSuffix(source, ".json") {
		candidates = append(candidates, strings.TrimSuffix(source, ".json"))
	} else {
		candidates = append(candidates, source+".json")
	}

	var lastErr error
	client := http.Client{}

	for _, cand := range candidates {
		url := fmt.Sprintf("%v%v", config.TrakConfig.RawBaseUrl, cand)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to connect to registry: %w", err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var toolTemplate *models.ToolTemplateModel
			if err := json.NewDecoder(resp.Body).Decode(&toolTemplate); err != nil {
				return nil, fmt.Errorf("failed to parse template AST schema: %w", err)
			}
			return toolTemplate, nil
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			lastErr = fmt.Errorf("blueprint not found in registry at '%s'. Check identifier spelling", source)
		} else {
			lastErr = fmt.Errorf("registry returned HTTP %d (%s)", resp.StatusCode, resp.Status)
		}
	}

	return nil, lastErr
}
