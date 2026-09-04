package shared

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type LanguageRuntimeConfig struct {
	Name         string
	Executables  []string // List of binary candidates to search in PATH
	BuildCommand func(binary string, moduleDir string) (bin string, args []string)
}

var Runtimes = map[string]LanguageRuntimeConfig{
	"go": {
		Name:        "Go",
		Executables: []string{"go"},
		BuildCommand: func(binary string, moduleDir string) (string, []string) {
			return binary, []string{"test", "-v", "./" + filepath.ToSlash(moduleDir) + "/..."}
		},
	},
	"golang": {
		Name:        "Go",
		Executables: []string{"go"},
		BuildCommand: func(binary string, moduleDir string) (string, []string) {
			return binary, []string{"test", "-v", "./" + filepath.ToSlash(moduleDir) + "/..."}
		},
	},
	"python": {
		Name:        "Python",
		Executables: []string{"python", "python3"},
		BuildCommand: func(binary string, moduleDir string) (string, []string) {
			return binary, []string{"-m", "unittest", "discover", "-s", "./" + filepath.ToSlash(moduleDir), "-p", "*test*.py"}
		},
	},
	"rust": {
		Name:        "Rust",
		Executables: []string{"cargo"},
		BuildCommand: func(binary string, moduleDir string) (string, []string) {
			manifestPath := filepath.Join(moduleDir, "Cargo.toml")
			return binary, []string{"test", "--manifest-path", manifestPath}
		},
	},
	"javascript": {
		Name:        "JavaScript",
		Executables: []string{"node", "bun"},
		BuildCommand: func(binary string, moduleDir string) (string, []string) {
			return binary, []string{"--test", "./" + filepath.ToSlash(moduleDir) + "/*test*.js"}
		},
	},
	"typescript": {
		Name:        "TypeScript",
		Executables: []string{"node", "bun"},
		BuildCommand: func(binary string, moduleDir string) (string, []string) {
			return binary, []string{"--test", "./" + filepath.ToSlash(moduleDir) + "/*test*.ts"}
		},
	},
	"js": {
		Name:        "JavaScript",
		Executables: []string{"node", "bun"},
		BuildCommand: func(binary string, moduleDir string) (string, []string) {
			return binary, []string{"--test", "./" + filepath.ToSlash(moduleDir) + "/*test*.js"}
		},
	},
	"ts": {
		Name:        "TypeScript",
		Executables: []string{"node", "bun"},
		BuildCommand: func(binary string, moduleDir string) (string, []string) {
			return binary, []string{"--test", "./" + filepath.ToSlash(moduleDir) + "/*test*.ts"}
		},
	},
}

// ResolveToolchain looks up whether any valid executable candidate is available in PATH.
func ResolveToolchain(lang string) (string, *LanguageRuntimeConfig, error) {
	cfg, ok := Runtimes[lang]
	if !ok {
		return "", nil, fmt.Errorf("no automated test runner configured for language '%s'", lang)
	}

	for _, bin := range cfg.Executables {
		if path, err := exec.LookPath(bin); err == nil && path != "" {
			return bin, &cfg, nil
		}
	}

	return "", &cfg, fmt.Errorf("toolchain '%s' was not found in your system PATH (checked: %v)", cfg.Name, cfg.Executables)
}
