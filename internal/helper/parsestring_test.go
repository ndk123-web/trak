package helper

import (
	"testing"
)

func TestParseTemplateString(t *testing.T) {
	tests := []struct {
		input       string
		wantErr     bool
		author      string
		category    string
		toolName    string
		isOfficial  bool
		sourcePath  string
		identifier  string
	}{
		// 1. Short Official
		{
			input:      "lang/go",
			wantErr:    false,
			author:     "Trak",
			category:   "lang",
			toolName:   "go",
			isOfficial: true,
			sourcePath: "templates/lang/go.json",
			identifier: "lang/go",
		},
		// 2. Explicit Official (with trak prefix)
		{
			input:      "trak/lang/go",
			wantErr:    false,
			author:     "Trak",
			category:   "lang",
			toolName:   "go",
			isOfficial: true,
			sourcePath: "templates/lang/go.json",
			identifier: "lang/go",
		},
		{
			input:      "trak/tool/docker",
			wantErr:    false,
			author:     "Trak",
			category:   "tool",
			toolName:   "docker",
			isOfficial: true,
			sourcePath: "templates/tool/docker.json",
			identifier: "tool/docker",
		},
		// 3. Community Namespaces
		{
			input:      "vishal-12/lang/go",
			wantErr:    false,
			author:     "vishal-12",
			category:   "lang",
			toolName:   "go",
			isOfficial: false,
			sourcePath: "users/vishal-12/lang/go.json",
			identifier: "vishal-12/lang/go",
		},
		{
			input:      "alex_dev/db/postgres",
			wantErr:    false,
			author:     "alex_dev",
			category:   "db",
			toolName:   "postgres",
			isOfficial: false,
			sourcePath: "users/alex_dev/db/postgres.json",
			identifier: "alex_dev/db/postgres",
		},
		{
			input:      "Ndk18-wesd/db/postgres@v1.0.0",
			wantErr:    false,
			author:     "Ndk18-wesd",
			category:   "db",
			toolName:   "postgres@v1.0.0",
			isOfficial: false,
			sourcePath: "users/Ndk18-wesd/db/postgres@v1.0.0.json",
			identifier: "Ndk18-wesd/db/postgres@v1.0.0",
		},
		{
			input:      "lang/go@v1.2.0",
			wantErr:    false,
			author:     "Trak",
			category:   "lang",
			toolName:   "go@v1.2.0",
			isOfficial: true,
			sourcePath: "templates/lang/go@v1.2.0.json",
			identifier: "lang/go@v1.2.0",
		},
		// 4. Invalid Cases
		{
			input:   "",
			wantErr: true,
		},
		{
			input:   "just-one-word",
			wantErr: true,
		},
		{
			input:   "invalid_category/go",
			wantErr: true,
		},
		{
			input:   "too/many/nested/parts/deep",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			res, err := ParseTemplateString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseTemplateString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				if res.Author != tt.author {
					t.Errorf("Author = %v, want %v", res.Author, tt.author)
				}
				if res.Category != tt.category {
					t.Errorf("Category = %v, want %v", res.Category, tt.category)
				}
				if res.ToolName != tt.toolName {
					t.Errorf("ToolName = %v, want %v", res.ToolName, tt.toolName)
				}
				if res.IsOfficial != tt.isOfficial {
					t.Errorf("IsOfficial = %v, want %v", res.IsOfficial, tt.isOfficial)
				}
				if res.SourcePath != tt.sourcePath {
					t.Errorf("SourcePath = %v, want %v", res.SourcePath, tt.sourcePath)
				}
				if res.Identifier != tt.identifier {
					t.Errorf("Identifier = %v, want %v", res.Identifier, tt.identifier)
				}
			}
		})
	}
}
