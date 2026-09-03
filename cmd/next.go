package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var openInEditor bool

var nextCmd = cobra.Command{
	Use:     "next [flags]",
	Aliases: []string{"todo", "current"},
	Short:   "Discover and jump directly to your next pending curriculum exercise",
	Long: `Inspects trak.json, determines the next sequential incomplete module,
and displays its folder path, overview, files, and navigation commands.

Optionally pass --open (-o) to automatically launch the module in your code editor.`,
	Example: `  # Show the next pending module:
  trak next

  # Automatically open the next module in VS Code:
  trak next --open
  trak next -o`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to get current directory: %v", err))
			return
		}

		trakPath := filepath.Join(cwd, "trak.json")
		if _, err := os.Stat(trakPath); os.IsNotExist(err) {
			fmt.Println()
			fmt.Printf("  %s✖ Not a Trak workspace%s\n", ui.Red+ui.Bold, ui.Reset)
			fmt.Printf("  %sCould not find 'trak.json' in: %s%s\n\n", ui.Gray, cwd, ui.Reset)
			fmt.Printf("  %s💡 Tip:%s Navigate into your learning track folder (e.g. %scd ./learn-go%s)\n", ui.White+ui.Bold, ui.Reset, ui.Green, ui.Reset)
			fmt.Printf("     or run %strak init <category>/<tool>%s to materialize one.\n\n", ui.Green, ui.Reset)
			return
		}

		dataBytes, err := os.ReadFile(trakPath)
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to read trak.json: %v", err))
			return
		}

		var trakStruct models.WorkspaceMetadata
		if err := json.NewDecoder(bytes.NewReader(dataBytes)).Decode(&trakStruct); err != nil {
			ui.Error(fmt.Sprintf("Failed to parse trak.json: %v", err))
			return
		}

		if len(trakStruct.ModuleBreakdown) == 0 {
			ui.Error("No curriculum modules found in trak.json.")
			return
		}

		// Sort module keys alphabetically/numerically to guarantee deterministic order
		var sortedKeys []string
		for k := range trakStruct.ModuleBreakdown {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		// Find first incomplete module
		nextModule := ""
		completedCount := 0
		totalModules := len(sortedKeys)

		for _, k := range sortedKeys {
			if trakStruct.ModuleBreakdown[k] {
				completedCount++
			} else if nextModule == "" {
				nextModule = k
			}
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = " Resolving next curriculum module..."
		s.Start()
		time.Sleep(300 * time.Millisecond)
		s.Stop()

		// All modules completed case
		if nextModule == "" {
			fmt.Println()
			fmt.Printf("  %s🎉 Track Fully Completed!%s\n", ui.Green+ui.Bold, ui.Reset)
			fmt.Printf("  %sYou have completed all %d of %d modules in:%s\n", ui.Gray, totalModules, totalModules, ui.Reset)
			fmt.Printf("  %s%s%s\n\n", ui.White+ui.Bold, trakStruct.Name, ui.Reset)
			fmt.Printf("  %sRun %strak status%s to review your complete workspace summary.\n\n", ui.Gray, ui.Green, ui.Reset)
			return
		}

		progressPercent := (float64(completedCount) / float64(totalModules)) * 100.0
		modulePath := filepath.Join(cwd, nextModule)

		// Extract README description if available
		readmeSnippet := ""
		readmePath := filepath.Join(modulePath, "README.md")
		if readmeData, err := os.ReadFile(readmePath); err == nil {
			lines := strings.Split(string(readmeData), "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed == "" || strings.HasPrefix(trimmed, "#") {
					continue
				}
				readmeSnippet = trimmed
				break
			}
		}

		// List files in module
		var exerciseFiles []string
		if entries, err := os.ReadDir(modulePath); err == nil {
			for _, entry := range entries {
				exerciseFiles = append(exerciseFiles, entry.Name())
			}
		}

		// Terminal Output
		fmt.Println()
		fmt.Printf("  %s%s➜ Next Curriculum Module%s\n", ui.Bold, ui.Green, ui.Reset)
		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)
		fmt.Printf("  %s• Module Name   :%s  %s%s%s\n", ui.Gray, ui.Reset, ui.White+ui.Bold, nextModule, ui.Reset)
		fmt.Printf("  %s• Track         :%s  %s%s%s\n", ui.Gray, ui.Reset, ui.White, trakStruct.Name, ui.Reset)
		fmt.Printf("  %s• Progress So Far:%s [%s%s%s] %s%.1f%%%s %s(%d of %d completed)%s\n",
			ui.Gray, ui.Reset,
			ui.Green, ui.ProgressBar(progressPercent, 20), ui.Reset,
			ui.Green+ui.Bold, progressPercent, ui.Reset,
			ui.Gray, completedCount, totalModules, ui.Reset,
		)
		fmt.Printf("  %s• Directory     :%s  %s./%s%s\n", ui.Gray, ui.Reset, ui.Green, nextModule, ui.Reset)

		if readmeSnippet != "" {
			if len(readmeSnippet) > 80 {
				readmeSnippet = readmeSnippet[:77] + "..."
			}
			fmt.Printf("  %s• Objective     :%s  %s%s%s\n", ui.Gray, ui.Reset, ui.LightGray, readmeSnippet, ui.Reset)
		}

		if len(exerciseFiles) > 0 {
			fmt.Printf("  %s• Files         :%s  %s%s%s\n", ui.Gray, ui.Reset, ui.Gray, strings.Join(exerciseFiles, ", "), ui.Reset)
		}

		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)

		// Action Guidance
		fmt.Printf("\n  %sGet Started:%s\n", ui.White+ui.Bold, ui.Reset)
		fmt.Printf("    %scd ./%s%s\n", ui.Green, nextModule, ui.Reset)
		fmt.Printf("    %strak done %s%s     # When completed\n\n", ui.Green, strings.Split(nextModule, "-")[0], ui.Reset)

		// Optional editor auto-open
		if openInEditor {
			openEditor(modulePath)
		}
	},
}

func openEditor(targetPath string) {
	// Try VS Code first
	if _, err := exec.LookPath("code"); err == nil {
		c := exec.Command("code", targetPath)
		if err := c.Start(); err == nil {
			fmt.Printf("  %s✔ Launched VS Code at: %s%s\n\n", ui.Green, targetPath, ui.Reset)
			return
		}
	}

	// OS fallback opener
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", targetPath)
	case "darwin":
		cmd = exec.Command("open", targetPath)
	default:
		cmd = exec.Command("xdg-open", targetPath)
	}

	if err := cmd.Start(); err == nil {
		fmt.Printf("  %s✔ Opened module in default viewer%s\n\n", ui.Green, ui.Reset)
	}
}

func init() {
	nextCmd.Flags().BoolVarP(&openInEditor, "open", "o", false, "Open the next module directory in your code editor")
	rootCmd.AddCommand(&nextCmd)
}
