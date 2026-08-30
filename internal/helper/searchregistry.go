package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ndk123-web/trak/internal/config"
	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/ui"
)

// FetchRegistry downloads and decodes the registry.json catalog
func FetchRegistry() (*models.RegistryModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	registryURL := fmt.Sprintf("%sregistry.json", config.TrakConfig.RawBaseUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, registryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to registry (%s): %w", registryURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned HTTP %d (%s)", resp.StatusCode, resp.Status)
	}

	var registry models.RegistryModel
	if err := json.NewDecoder(resp.Body).Decode(&registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry catalog: %w", err)
	}

	return &registry, nil
}

// RenderRegistry renders the registry catalog after it has been fetched
func RenderRegistry(registry *models.RegistryModel, category string, all bool) error {
	category = strings.TrimSpace(strings.ToLower(category))

	// If no category specified or "all" passed, render all categories
	if all || category == "" || category == "all" {
		renderAllCategories(registry)
		return nil
	}

	// Render single category
	cat, exists := registry.Categories[category]
	if !exists {
		var available []string
		for k := range registry.Categories {
			available = append(available, k)
		}
		sort.Strings(available)
		return fmt.Errorf("category '%s' not found. Available categories: %s", category, strings.Join(available, ", "))
	}

	renderSingleCategory(category, cat)
	return nil
}

// SearchRegistry fetches and renders the registry catalog
func SearchRegistry(category string, all bool) error {
	registry, err := FetchRegistry()
	if err != nil {
		return err
	}
	return RenderRegistry(registry, category, all)
}

// getCategoryIcon returns a fitting emoji for category
func getCategoryIcon(catKey string) string {
	switch catKey {
	case "lang":
		return "📦"
	case "os":
		return "🐧"
	case "cloud":
		return "☁️ "
	case "db":
		return "🗄️ "
	case "tool":
		return "🛠️ "
	default:
		return "📁"
	}
}

// renderAllCategories renders the master catalog in clean, spaced, categorized blocks
func renderAllCategories(registry *models.RegistryModel) {
	totalTracks := 0
	for _, cat := range registry.Categories {
		totalTracks += len(cat.Templates)
	}

	fmt.Printf("\n%s%sTrak Learning Catalog%s %s(v%s • %d Blueprints)%s\n",
		ui.Bold, ui.White, ui.Reset, ui.Gray, registry.SchemaVersion, totalTracks, ui.Reset)

	// Preferred category display order
	preferredOrder := []string{"lang", "os", "cloud", "db", "tool"}
	categoryKeys := make([]string, 0, len(registry.Categories))
	seen := make(map[string]bool)

	for _, k := range preferredOrder {
		if _, ok := registry.Categories[k]; ok {
			categoryKeys = append(categoryKeys, k)
			seen[k] = true
		}
	}
	for k := range registry.Categories {
		if !seen[k] {
			categoryKeys = append(categoryKeys, k)
		}
	}

	for _, catKey := range categoryKeys {
		cat := registry.Categories[catKey]
		icon := getCategoryIcon(catKey)

		// Category Header Badge
		fmt.Printf("\n%s%s%s %s%s %s(%s/)%s\n",
			ui.Bold, icon, ui.Reset,
			ui.Bold, strings.ToUpper(cat.Title),
			ui.Cyan, catKey, ui.Reset,
		)

		// Sort templates alphabetically
		tplKeys := make([]string, 0, len(cat.Templates))
		for tKey := range cat.Templates {
			tplKeys = append(tplKeys, tKey)
		}
		sort.Strings(tplKeys)

		// Calculate max key width for clean column alignment
		maxKeyLen := 0
		for _, tKey := range tplKeys {
			if len(tKey) > maxKeyLen {
				maxKeyLen = len(tKey)
			}
		}

		// Print templates with column alignment and clean spacing
		for _, tKey := range tplKeys {
			tpl := cat.Templates[tKey]
			padding := strings.Repeat(" ", maxKeyLen-len(tKey)+3)

			fmt.Printf("   %s%s%s%s%s%s%s\n",
				ui.Cyan, tKey, ui.Reset,
				padding,
				ui.Gray, tpl.Description, ui.Reset,
			)
		}
	}

	// Clean Quickstart Tip Footer
	fmt.Printf("\n%s──────────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)
	fmt.Printf("%s💡 Initialize any workspace:%s %s%strak init <category>/<template>%s\n",
		ui.Yellow, ui.Reset, ui.Green, ui.Bold, ui.Reset)
	fmt.Printf("   %sexample:%s trak init lang/go --path ./learn-go\n\n", ui.Gray, ui.Reset)
}

// renderSingleCategory renders a focused view of a single category
func renderSingleCategory(catKey string, cat models.Category) {
	icon := getCategoryIcon(catKey)
	fmt.Printf("\n%s%s%s %s%s %s(%s/)%s\n",
		ui.Bold, icon, ui.Reset,
		ui.Bold, strings.ToUpper(cat.Title),
		ui.Cyan, catKey, ui.Reset,
	)
	if cat.Description != "" {
		fmt.Printf("   %s%s%s\n", ui.Gray, cat.Description, ui.Reset)
	}
	fmt.Println()

	tplKeys := make([]string, 0, len(cat.Templates))
	for tKey := range cat.Templates {
		tplKeys = append(tplKeys, tKey)
	}
	sort.Strings(tplKeys)

	maxKeyLen := 0
	for _, tKey := range tplKeys {
		if len(tKey) > maxKeyLen {
			maxKeyLen = len(tKey)
		}
	}

	for _, tKey := range tplKeys {
		tpl := cat.Templates[tKey]
		padding := strings.Repeat(" ", maxKeyLen-len(tKey)+3)
		fmt.Printf("   %s%s%s%s%s%s%s\n",
			ui.Cyan, tKey, ui.Reset,
			padding,
			ui.Gray, tpl.Description, ui.Reset,
		)
	}

	fmt.Printf("\n%s💡 Run:%s %s%strak init %s/<template> --path ./learn-%s%s\n\n",
		ui.Yellow, ui.Reset, ui.Green, ui.Bold, catKey, catKey, ui.Reset)
}
