package helper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ndk123-web/trak/internal/models"
)

// ValidCategories defines the recognized curriculum disciplines in Trak Registry
var ValidCategories = map[string]bool{
	"lang":  true,
	"os":    true,
	"cloud": true,
	"db":    true,
	"tool":  true,
}

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\.@]+$`)
var userRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-]*$`)

// ParseTemplateString parses both short official (lang/go), explicit official (trak/lang/go),
// and community namespaced blueprints (<username>/<category>/<tool>).
func ParseTemplateString(input string) (*models.ParsedTemplate, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return nil, fmt.Errorf("template identifier cannot be empty")
	}

	// Normalize any backslashes to forward slashes
	input = strings.ReplaceAll(input, "\\", "/")
	parts := strings.Split(input, "/")

	switch len(parts) {
	case 2:
		// Format: <category>/<tool> (e.g. lang/go, db/postgres) -> Official Track
		category := strings.ToLower(parts[0])
		toolName := strings.ToLower(parts[1])

		if !ValidCategories[category] {
			return nil, fmt.Errorf("invalid category '%s'. Must be one of: lang, os, cloud, db, tool", category)
		}

		if !nameRegex.MatchString(toolName) {
			return nil, fmt.Errorf("invalid blueprint name '%s'", toolName)
		}

		return &models.ParsedTemplate{
			Author:     "Trak",
			Category:   category,
			ToolName:   toolName,
			IsOfficial: true,
			SourcePath: fmt.Sprintf("templates/%s/%s.json", category, toolName),
			Identifier: fmt.Sprintf("%s/%s", category, toolName),
		}, nil

	case 3:
		// Format: <author>/<category>/<tool> (e.g. trak/lang/go OR <username>/lang/go)
		author := parts[0]
		category := strings.ToLower(parts[1])
		toolName := strings.ToLower(parts[2])

		if !ValidCategories[category] {
			return nil, fmt.Errorf("invalid category '%s'. Must be one of: lang, os, cloud, db, tool", category)
		}

		if !nameRegex.MatchString(toolName) {
			return nil, fmt.Errorf("invalid blueprint name '%s'", toolName)
		}

		if !userRegex.MatchString(author) {
			return nil, fmt.Errorf("invalid username '%s'. Must contain alphanumeric characters and hyphens", author)
		}

		// Check if author is official "trak" or alias "templates"
		if strings.EqualFold(author, "trak") || strings.EqualFold(author, "templates") {
			return &models.ParsedTemplate{
				Author:     "Trak",
				Category:   category,
				ToolName:   toolName,
				IsOfficial: true,
				SourcePath: fmt.Sprintf("templates/%s/%s.json", category, toolName),
				Identifier: fmt.Sprintf("%s/%s", category, toolName),
			}, nil
		}

		// Community blueprint under users/<author>/<category>/<tool>.json
		return &models.ParsedTemplate{
			Author:     author,
			Category:   category,
			ToolName:   toolName,
			IsOfficial: false,
			SourcePath: fmt.Sprintf("users/%s/%s/%s.json", author, category, toolName),
			Identifier: fmt.Sprintf("%s/%s/%s", author, category, toolName),
		}, nil

	default:
		return nil, fmt.Errorf("invalid blueprint format '%s'.\nExpected:\n  • lang/go\n  • trak/lang/go\n  • <username>/<category>/<tool>", input)
	}
}
