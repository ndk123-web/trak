package helper

import (
	"fmt"
	"strings"

	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/shared"
)

var ValidCategories = map[string]bool{
	"lang":  true,
	"os":    true,
	"cloud": true,
	"db":    true,
	"tool":  true,
}

var nameRegex = shared.TrakRegexes.NameRegex
var userRegex = shared.TrakRegexes.UserRegex
var versionRegex = shared.TrakRegexes.VersionRegex

func ExtractToolAndVersion(rawTool string) (toolName, version string, err error) {
	if strings.Contains(rawTool, "@") {
		parts := strings.SplitN(rawTool, "@", 2)
		toolName = strings.ToLower(parts[0])
		version = parts[1]

		if !nameRegex.MatchString(toolName) {
			return "", "", fmt.Errorf("invalid blueprint name '%s'", toolName)
		}
		if !versionRegex.MatchString(version) {
			return "", "", fmt.Errorf("invalid version tag '%s' (e.g. v1.0.0)", version)
		}
		return toolName, version, nil
	}

	toolName = strings.ToLower(rawTool)
	if !nameRegex.MatchString(toolName) {
		return "", "", fmt.Errorf("invalid blueprint name '%s'", toolName)
	}
	return toolName, "", nil
}

func ParseTemplateString(input string) (*models.ParsedTemplate, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("template identifier cannot be empty")
	}

	input = strings.ReplaceAll(input, "\\", "/")
	parts := strings.Split(input, "/")

	switch len(parts) {
	case 2:
		// Format: <category>/<tool>[@<version>]
		category := strings.ToLower(parts[0])
		if !ValidCategories[category] {
			return nil, fmt.Errorf("invalid category '%s'", category)
		}

		toolName, version, err := ExtractToolAndVersion(parts[1])
		if err != nil {
			return nil, err
		}

		sourceFile := toolName
		if version != "" {
			sourceFile = fmt.Sprintf("%s@%s", toolName, version)
		}

		return &models.ParsedTemplate{
			Author:   "Trak",
			Category: category,
			ToolName: toolName,
			// Version:    version,
			IsOfficial: true,
			SourcePath: fmt.Sprintf("templates/%s/%s.json", category, sourceFile),
			Identifier: fmt.Sprintf("%s/%s", category, sourceFile),
		}, nil

	case 3:
		// Format: <author>/<category>/<tool>[@<version>]
		author := parts[0]
		category := strings.ToLower(parts[1])

		if !ValidCategories[category] {
			return nil, fmt.Errorf("invalid category '%s'", category)
		}

		if !userRegex.MatchString(author) {
			return nil, fmt.Errorf("invalid username '%s'", author)
		}

		toolName, version, err := ExtractToolAndVersion(parts[2])
		if err != nil {
			return nil, err
		}

		sourceFile := toolName
		if version != "" {
			sourceFile = fmt.Sprintf("%s@%s", toolName, version)
		}

		isOfficial := strings.EqualFold(author, "trak") || strings.EqualFold(author, "templates")
		sourcePath := fmt.Sprintf("users/%s/%s/%s.json", author, category, sourceFile)
		if isOfficial {
			sourcePath = fmt.Sprintf("templates/%s/%s.json", category, sourceFile)
		}

		return &models.ParsedTemplate{
			Author:   author,
			Category: category,
			ToolName: toolName,
			// Version:    version,
			IsOfficial: isOfficial,
			SourcePath: sourcePath,
			Identifier: fmt.Sprintf("%s/%s/%s", author, category, sourceFile),
		}, nil

	default:
		return nil, fmt.Errorf("invalid blueprint format '%s'", input)
	}
}
