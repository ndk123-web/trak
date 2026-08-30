package helper

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func ParseTemplateString(template string) (string, string, error) {
	template = strings.TrimSpace(template)

	if strings.EqualFold(template, "") {
		return "", "", errors.New("template String is Empty")
	}

	var category string
	var toolName string

	pattern, err := regexp.Compile("^[a-zA-z]+/[a-zA-z]+$")
	if err != nil {
		return "", "", err
	}

	// use first
	if valid := pattern.MatchString(template); !valid {
		return "", "", err
	}

	splits := strings.Split(template, "/")

	category = splits[0]
	toolName = splits[1]

	fmt.Printf("Category: %v\n", category)
	fmt.Printf("ToolName: %v\n", toolName)

	return category, toolName, nil
}
