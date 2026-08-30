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

// SearchRegistry renders a beautiful ASCII tree visualization of available templates
func SearchRegistry(category string, all bool) error {
	registry, err := FetchRegistry()
	if err != nil {
		return err
	}

	category = strings.TrimSpace(strings.ToLower(category))

	// If no category specified or "all" passed, render all categories
	if all || category == "" || category == "all" {
		renderAllCategoriesTree(registry)
		return nil
	}

	// Render single category tree
	cat, exists := registry.Categories[category]
	if !exists {
		var available []string
		for k := range registry.Categories {
			available = append(available, k)
		}
		sort.Strings(available)
		return fmt.Errorf("category '%s' not found. Available categories: %s", category, strings.Join(available, ", "))
	}

	renderSingleCategoryTree(category, cat)
	return nil
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

// renderAllCategoriesTree renders full catalog tree
func renderAllCategoriesTree(registry *models.RegistryModel) {
	fmt.Printf("\n%s%s🌳 Trak Learning Catalog%s %s(v%s)%s\n\n", ui.Bold, ui.Green, ui.Reset, ui.Gray, registry.SchemaVersion, ui.Reset)

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

	for i, catKey := range categoryKeys {
		cat := registry.Categories[catKey]
		isLastCat := i == len(categoryKeys)-1

		catBranch := "├──"
		subPrefix := "│  "
		if isLastCat {
			catBranch = "└──"
			subPrefix = "   "
		}

		icon := getCategoryIcon(catKey)
		fmt.Printf("%s %s %s%s%s %s(%s)%s\n", catBranch, icon, ui.Bold, cat.Title, ui.Reset, ui.Cyan, catKey, ui.Reset)

		// Sort templates alphabetically
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

		for j, tKey := range tplKeys {
			tpl := cat.Templates[tKey]
			isLastTpl := j == len(tplKeys)-1

			tplBranch := "├──"
			if isLastTpl {
				tplBranch = "└──"
			}

			padding := strings.Repeat(" ", maxKeyLen-len(tKey)+2)
			fmt.Printf("%s %s %s%s%s%s%s%s%s\n",
				subPrefix,
				tplBranch,
				ui.Cyan, tKey, ui.Reset,
				padding,
				ui.Gray, tpl.Description, ui.Reset,
			)
		}

		if !isLastCat {
			fmt.Println("│")
		}
	}

	fmt.Printf("\n%s💡 Tip:%s Initialize any workspace with: %s%strak init <category>/<template> --path ./my-workspace%s\n\n",
		ui.Yellow, ui.Reset, ui.Green, ui.Bold, ui.Reset)
}

// renderSingleCategoryTree renders a single category in focused view
func renderSingleCategoryTree(catKey string, cat models.Category) {
	icon := getCategoryIcon(catKey)
	fmt.Printf("\n%s %s%s%s %s(%s)%s\n", icon, ui.Bold, cat.Title, ui.Reset, ui.Cyan, catKey, ui.Reset)
	if cat.Description != "" {
		fmt.Printf("%s%s%s\n", ui.Gray, cat.Description, ui.Reset)
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

	for j, tKey := range tplKeys {
		tpl := cat.Templates[tKey]
		isLastTpl := j == len(tplKeys)-1

		tplBranch := "├──"
		if isLastTpl {
			tplBranch = "└──"
		}

		padding := strings.Repeat(" ", maxKeyLen-len(tKey)+2)
		fmt.Printf("%s %s%s%s%s%s%s%s\n",
			tplBranch,
			ui.Cyan, tKey, ui.Reset,
			padding,
			ui.Gray, tpl.Description, ui.Reset,
		)
	}

	fmt.Printf("\n%s💡 Tip:%s Run: %s%strak init %s/<template> --path ./learn-%s%s\n\n",
		ui.Yellow, ui.Reset, ui.Green, ui.Bold, catKey, catKey, ui.Reset)
}
