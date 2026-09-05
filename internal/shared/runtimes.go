package shared

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type LanguageRuntimeConfig struct {
	Name         string
	Executables  []string // List of binary candidates to search in PATH
	BuildCommand func(binary string, moduleDir string) (bin string, args []string)
}

func findCFiles(moduleDir string) []string {
	testFile := filepath.Join(moduleDir, "exercise_test.c")
	exFile := filepath.Join(moduleDir, "exercise.c")
	starterFile := filepath.Join(moduleDir, "exercise", "starter.c")
	mainFile := filepath.Join(moduleDir, "main.c")

	if _, err := os.Stat(testFile); err == nil {
		var files []string
		if _, err := os.Stat(exFile); err == nil {
			files = append(files, exFile)
		}
		files = append(files, testFile)
		return files
	}

	matches, _ := filepath.Glob(filepath.Join(moduleDir, "*test*.c"))
	if len(matches) > 0 {
		var files []string
		if _, err := os.Stat(exFile); err == nil {
			files = append(files, exFile)
		}
		files = append(files, matches[0])
		return files
	}

	if _, err := os.Stat(starterFile); err == nil {
		return []string{starterFile}
	}

	if _, err := os.Stat(mainFile); err == nil {
		return []string{mainFile}
	}

	if _, err := os.Stat(exFile); err == nil {
		return []string{exFile}
	}

	allC, _ := filepath.Glob(filepath.Join(moduleDir, "*.c"))
	var nonExample []string
	for _, f := range allC {
		base := filepath.Base(f)
		if base != "example.c" && base != "_trak_test.c" {
			nonExample = append(nonExample, f)
		}
	}
	if len(nonExample) > 0 {
		return nonExample
	}

	return []string{filepath.Join(moduleDir, "main.c")}
}

func findCppFiles(moduleDir string) []string {
	testFile := filepath.Join(moduleDir, "exercise_test.cpp")
	exFile := filepath.Join(moduleDir, "exercise.cpp")
	starterFile := filepath.Join(moduleDir, "exercise", "starter.cpp")
	mainFile := filepath.Join(moduleDir, "main.cpp")

	if _, err := os.Stat(testFile); err == nil {
		var files []string
		if _, err := os.Stat(exFile); err == nil {
			files = append(files, exFile)
		}
		files = append(files, testFile)
		return files
	}

	matches, _ := filepath.Glob(filepath.Join(moduleDir, "*test*.cpp"))
	if len(matches) > 0 {
		var files []string
		if _, err := os.Stat(exFile); err == nil {
			files = append(files, exFile)
		}
		files = append(files, matches[0])
		return files
	}

	if _, err := os.Stat(starterFile); err == nil {
		return []string{starterFile}
	}

	if _, err := os.Stat(mainFile); err == nil {
		return []string{mainFile}
	}

	if _, err := os.Stat(exFile); err == nil {
		return []string{exFile}
	}

	allCpp, _ := filepath.Glob(filepath.Join(moduleDir, "*.cpp"))
	var nonExample []string
	for _, f := range allCpp {
		base := filepath.Base(f)
		if base != "example.cpp" && base != "_trak_test.cpp" {
			nonExample = append(nonExample, f)
		}
	}
	if len(nonExample) > 0 {
		return nonExample
	}

	return []string{filepath.Join(moduleDir, "main.cpp")}
}

func buildCCommand(binary string, moduleDir string) (string, []string) {
	outName := "_trak_test"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outExe := filepath.Join(moduleDir, outName)

	files := findCFiles(moduleDir)
	compileArgs := []string{"-Wall", "-Wextra", "-std=c11", "-I" + moduleDir, "-o", outExe}
	compileArgs = append(compileArgs, files...)
	if runtime.GOOS != "windows" {
		compileArgs = append(compileArgs, "-lm")
	}

	compileCmd := exec.Command(binary, compileArgs...)
	if _, err := compileCmd.CombinedOutput(); err != nil {
		return binary, compileArgs
	}

	absExe, err := filepath.Abs(outExe)
	if err == nil {
		return absExe, []string{}
	}
	return outExe, []string{}
}

func buildCppCommand(binary string, moduleDir string) (string, []string) {
	outName := "_trak_test"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outExe := filepath.Join(moduleDir, outName)

	files := findCppFiles(moduleDir)
	compileArgs := []string{"-Wall", "-Wextra", "-std=c++17", "-I" + moduleDir, "-o", outExe}
	compileArgs = append(compileArgs, files...)

	compileCmd := exec.Command(binary, compileArgs...)
	if _, err := compileCmd.CombinedOutput(); err != nil {
		return binary, compileArgs
	}

	absExe, err := filepath.Abs(outExe)
	if err == nil {
		return absExe, []string{}
	}
	return outExe, []string{}
}

var Runtimes = map[string]LanguageRuntimeConfig{
	"go": {
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
	"c": {
		Name:         "C",
		Executables:  []string{"gcc", "clang"},
		BuildCommand: buildCCommand,
	},
	"cpp": {
		Name:         "C++",
		Executables:  []string{"g++", "clang++"},
		BuildCommand: buildCppCommand,
	},
	"c++": {
		Name:         "C++",
		Executables:  []string{"g++", "clang++"},
		BuildCommand: buildCppCommand,
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
