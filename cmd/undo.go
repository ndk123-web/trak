package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var undoCmd = cobra.Command{
	Use:     "undo <module>",
	Aliases: []string{"reset", "unmark"},
	Short:   "Reset or unmark a curriculum module back to pending",
	Long: `Revert a previously completed module back to pending in trak.json.
Recalculates your completion percentage and updates your workspace progress.

You can specify the module using:
  • Full module name:  00-setup-and-prerequisites
  • Number prefix:     00, 01, 1, 12
  • Short keyword:     setup, escape, goroutines`,
	Example: `  # Reset by module number:
  trak undo 00
  trak undo 1

  # Reset by full folder name:
  trak undo 00-setup-and-prerequisites

  # Using aliases:
  trak reset 01
  trak unmark 02`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.TrimSpace(args[0])
		if query == "" {
			ui.Error("Please specify a module name or number (e.g. 'trak undo 00')")
			return
		}

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
			ui.Error("This workspace does not have any modules registered in trak.json.")
			return
		}

		// Sort module keys alphabetically/numerically for consistent ordering
		var allKeys []string
		for k := range trakStruct.ModuleBreakdown {
			allKeys = append(allKeys, k)
		}
		sort.Strings(allKeys)

		// ── Smart Module Matching ──
		resolvedKey := ""

		// 1. Direct exact match
		if _, exists := trakStruct.ModuleBreakdown[query]; exists {
			resolvedKey = query
		}

		// 2. Prefix match e.g. "00" matches "00-setup-and-prerequisites"
		if resolvedKey == "" {
			prefix := query + "-"
			paddedPrefix := fmt.Sprintf("%02s-", query)

			var matches []string
			for _, k := range allKeys {
				if strings.HasPrefix(k, prefix) || strings.HasPrefix(k, paddedPrefix) {
					matches = append(matches, k)
				}
			}

			if len(matches) == 1 {
				resolvedKey = matches[0]
			} else if len(matches) > 1 {
				fmt.Printf("\n  %s✖ Ambiguous module query '%s'%s\n", ui.Red+ui.Bold, query, ui.Reset)
				fmt.Printf("  %sMatches multiple modules:%s\n", ui.Gray, ui.Reset)
				for _, m := range matches {
					fmt.Printf("    • %s%s%s\n", ui.Green, m, ui.Reset)
				}
				fmt.Printf("\n  %sPlease specify the exact module name.%s\n\n", ui.Gray, ui.Reset)
				return
			}
		}

		// 3. Substring match fallback e.g. "setup" or "escape"
		if resolvedKey == "" {
			var matches []string
			queryLower := strings.ToLower(query)
			for _, k := range allKeys {
				if strings.Contains(strings.ToLower(k), queryLower) {
					matches = append(matches, k)
				}
			}

			if len(matches) == 1 {
				resolvedKey = matches[0]
			} else if len(matches) > 1 {
				fmt.Printf("\n  %s✖ Ambiguous module query '%s'%s\n", ui.Red+ui.Bold, query, ui.Reset)
				fmt.Printf("  %sMatches multiple modules:%s\n", ui.Gray, ui.Reset)
				for _, m := range matches {
					fmt.Printf("    • %s%s%s\n", ui.Green, m, ui.Reset)
				}
				fmt.Printf("\n  %sPlease specify a more specific name.%s\n\n", ui.Gray, ui.Reset)
				return
			}
		}

		// If no matches found, show available modules
		if resolvedKey == "" {
			fmt.Printf("\n  %s✖ Module not found: '%s'%s\n", ui.Red+ui.Bold, query, ui.Reset)
			fmt.Printf("  %sAvailable modules in this workspace:%s\n", ui.Gray, ui.Reset)
			for _, k := range allKeys {
				statusTag := fmt.Sprintf("%s[PENDING]%s", ui.Gray, ui.Reset)
				if trakStruct.ModuleBreakdown[k] {
					statusTag = fmt.Sprintf("%s[COMPLETED]%s", ui.Green+ui.Bold, ui.Reset)
				}
				fmt.Printf("    • %s%-36s%s %s\n", ui.White, k, ui.Reset, statusTag)
			}
			fmt.Println()
			return
		}

		// ── Check if already pending ──
		wasAlreadyPending := !trakStruct.ModuleBreakdown[resolvedKey]
		trakStruct.ModuleBreakdown[resolvedKey] = false

		// Calculate updated metrics
		totalModules := len(allKeys)
		completedCount := 0

		for _, k := range allKeys {
			if trakStruct.ModuleBreakdown[k] {
				completedCount++
			}
		}

		progressPercent := (float64(completedCount) / float64(totalModules)) * 100.0

		// Write back formatted JSON
		updatedBytes, err := json.MarshalIndent(trakStruct, "", "  ")
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to encode metadata: %v", err))
			return
		}

		if err := os.WriteFile(trakPath, updatedBytes, 0644); err != nil {
			ui.Error(fmt.Sprintf("Failed to update trak.json: %v", err))
			return
		}

		// ── Render Clean Terminal Feedback ──
		progressBar := ui.ProgressBar(progressPercent, 24)

		fmt.Println()
		if wasAlreadyPending {
			fmt.Printf("  %sℹ Module is already pending:%s %s%s%s\n",
				ui.Gray, ui.Reset, ui.White+ui.Bold, resolvedKey, ui.Reset)
		} else {
			fmt.Printf("  %s➜ Module Reset:%s %s%s%s\n",
				ui.Green, ui.Reset, ui.White+ui.Bold, resolvedKey, ui.Reset)
		}
		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)

		fmt.Printf("  %s• Status        :%s  %s[PENDING]%s\n", ui.Gray, ui.Reset, ui.LightGray, ui.Reset)
		fmt.Printf("  %s• Track Name    :%s  %s%s%s\n", ui.Gray, ui.Reset, ui.White, trakStruct.Name, ui.Reset)
		fmt.Printf("  %s• Track Progress:%s  [%s%s%s] %s%.1f%%%s %s(%d of %d modules)%s\n",
			ui.Gray, ui.Reset,
			ui.Green, progressBar, ui.Reset,
			ui.Green+ui.Bold, progressPercent, ui.Reset,
			ui.Gray, completedCount, totalModules, ui.Reset,
		)
		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n\n", ui.Gray, ui.Reset)
	},
}

func init() {
	rootCmd.AddCommand(&undoCmd)
}
