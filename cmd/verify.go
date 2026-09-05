package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ndk123-web/trak/internal/helper"
	"github.com/ndk123-web/trak/internal/shared"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var (
	verifyAll  bool
	verifyList bool
	detail     bool
)

var verifyCmd = cobra.Command{
	Use:     "verify [module] [flags]",
	Aliases: []string{"test", "check"},
	Short:   "Run automated tests to verify your exercise implementation",
	Long: `Executes native language test suites against your exercise code.
If all tests pass, the module is automatically marked complete in trak.json
and your curriculum progress is updated.

You can specify:
  • No arguments:  Verifies the current pending module
  • Module query:  trak verify 00, trak verify 02, trak verify escape
  • All modules:   trak verify --all (-a)
  • Supported:     trak verify --list (-l) or trak verify allowlists`,
	Example: `  # Verify current pending exercise:
  trak verify

  # Verify specific module by number:
  trak verify 00
  trak verify 2

  # Verify all modules across the workspace:
  trak verify --all
  trak verify -a

  # List supported language runtimes:
  trak verify --list
  trak verify allowlists`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if verifyList {
			renderSupportedRuntimes()
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to get current directory: %v", err))
			return
		}

		trakStruct, err := helper.FetchTrakJsonWorkspace()
		if err != nil {
			fmt.Println()
			fmt.Printf("  %s✖ Not a Trak workspace%s\n", ui.Red+ui.Bold, ui.Reset)
			fmt.Printf("  %sCould not find 'trak.json' in: %s%s\n\n", ui.Gray, cwd, ui.Reset)
			fmt.Printf("  %sTip:%s Navigate into your learning track folder (e.g. %scd ./learn-go%s)\n", ui.White+ui.Bold, ui.Reset, ui.Green, ui.Reset)
			fmt.Printf("     or run %strak init <category>/<tool>%s to materialize one.\n\n", ui.Green, ui.Reset)
			return
		}

		if len(trakStruct.ModuleBreakdown) == 0 {
			ui.Error("This workspace does not have any modules registered in trak.json.")
			return
		}

		// Parse template identifier to detect track category & tool
		parsed, err := helper.ParseTemplateString(trakStruct.Id)
		if err != nil {
			// Fallback parsing from template field
			parsed, err = helper.ParseTemplateString(trakStruct.Template)
			if err != nil {
				ui.Error(fmt.Sprintf("Failed to parse track template '%s': %v", trakStruct.Id, err))
				return
			}
		}

		// Only "lang" tracks have automated compilers/test runners
		if parsed.Category != "lang" {
			fmt.Println()
			fmt.Printf("  %s[i] Hands-On Laboratory Track%s\n", ui.Green+ui.Bold, ui.Reset)
			fmt.Printf("  %sThis track (%s/%s) is an architectural laboratory without automated compiler tests.%s\n\n", ui.Gray, parsed.Category, parsed.ToolName, ui.Reset)
			fmt.Printf("  %sUse %strak done <module>%s to record your completed exercises.\n\n", ui.Gray, ui.Green, ui.Reset)
			return
		}

		// Resolve toolchain
		resolvedBin, runtimeCfg, err := shared.ResolveToolchain(parsed.ToolName)
		if err != nil {
			fmt.Println()
			fmt.Printf("  %s* Toolchain Not Found%s\n", ui.Red+ui.Bold, ui.Reset)
			fmt.Printf("  %s%s%s\n\n", ui.Gray, err.Error(), ui.Reset)
			fmt.Printf("  %sNote:%s Install %s to run verification tests locally.\n\n", ui.White+ui.Bold, ui.Reset, parsed.ToolName)
			return
		}

		// Sort module keys alphabetically/numerically
		var allKeys []string
		for k := range trakStruct.ModuleBreakdown {
			allKeys = append(allKeys, k)
		}
		sort.Strings(allKeys)

		// Determine target modules to verify
		var targets []string

		if verifyAll {
			targets = allKeys
		} else if len(args) == 1 {
			query := strings.TrimSpace(args[0])

			// 1. Direct match
			if _, exists := trakStruct.ModuleBreakdown[query]; exists {
				targets = append(targets, query)
			}

			// 2. Prefix match (e.g. "00", "0", "1")
			if len(targets) == 0 {
				prefix := query + "-"
				paddedPrefix := fmt.Sprintf("%02s-", query)
				for _, k := range allKeys {
					if strings.HasPrefix(k, prefix) || strings.HasPrefix(k, paddedPrefix) {
						targets = append(targets, k)
					}
				}
			}

			// 3. Keyword match
			if len(targets) == 0 {
				queryLower := strings.ToLower(query)
				for _, k := range allKeys {
					if strings.Contains(strings.ToLower(k), queryLower) {
						targets = append(targets, k)
					}
				}
			}

			if len(targets) == 0 {
				fmt.Printf("\n  %s✖ Module not found: '%s'%s\n", ui.Red+ui.Bold, query, ui.Reset)
				fmt.Printf("  %sAvailable modules in this workspace:%s\n", ui.Gray, ui.Reset)
				for _, k := range allKeys {
					fmt.Printf("    • %s%s%s\n", ui.Green, k, ui.Reset)
				}
				fmt.Println()
				return
			}

			if len(targets) > 1 {
				fmt.Printf("\n  %s✖ Ambiguous module query '%s'%s\n", ui.Red+ui.Bold, query, ui.Reset)
				fmt.Printf("  %sMatches multiple modules:%s\n", ui.Gray, ui.Reset)
				for _, m := range targets {
					fmt.Printf("    • %s%s%s\n", ui.Green, m, ui.Reset)
				}
				fmt.Printf("\n  %sPlease specify the exact module name or prefix.%s\n\n", ui.Gray, ui.Reset)
				return
			}
		} else {
			// No args: Find current pending module
			for _, k := range allKeys {
				if !trakStruct.ModuleBreakdown[k] {
					targets = append(targets, k)
					break
				}
			}

			if len(targets) == 0 {
				fmt.Println()
				fmt.Printf("  %sAll Modules Already Completed!%s\n", ui.Green+ui.Bold, ui.Reset)
				fmt.Printf("  %sRun %strak verify --all%s to re-test the entire curriculum.\n\n", ui.Gray, ui.Green, ui.Reset)
				return
			}
		}

		// Execute Verification
		fmt.Println()
		fmt.Printf("  %s>> Trak Automated Test Runner (%s)%s\n", ui.Bold+ui.Green, runtimeCfg.Name, ui.Reset)
		fmt.Printf("  %sTrack: %s%s\n", ui.Gray, trakStruct.Name, ui.Reset)
		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)

		passCount := 0
		failCount := 0
		hasStateChange := false

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Color("green")

		for _, mod := range targets {
			modDir := filepath.Join(cwd, mod)
			if stat, err := os.Stat(modDir); os.IsNotExist(err) || !stat.IsDir() {
				failCount++
				fmt.Printf("  %s✖ FAIL%s  %s%s%s\n", ui.Red+ui.Bold, ui.Reset, ui.White+ui.Bold, mod, ui.Reset)
				fmt.Printf("  %s┌── Directory Not Found ──────────────────────────────────────%s\n", ui.Red, ui.Reset)
				fmt.Printf("  %s│ Module directory '%s' was not found in: %s\n", ui.Gray, mod, cwd)
				fmt.Printf("  %s│ Make sure you run trak verify inside your initialized track folder.%s\n", ui.Gray, ui.Reset)
				fmt.Printf("  %s└── Check workspace folder and try again ─────────────────────%s\n\n", ui.Red, ui.Reset)
				continue
			}

			s.Suffix = fmt.Sprintf(" Testing %s...", mod)
			s.Start()

			bin, cmdArgs := runtimeCfg.BuildCommand(resolvedBin, mod)
			testCmd := exec.Command(bin, cmdArgs...)
			testCmd.Dir = cwd

			output, testErr := testCmd.CombinedOutput()
			s.Stop()
			fmt.Print("\r\033[K")

			if testErr == nil {
				passCount++
				fmt.Printf("  %s✔ PASS%s  %s%s%s\n", ui.Green+ui.Bold, ui.Reset, ui.White+ui.Bold, mod, ui.Reset)

				if !trakStruct.ModuleBreakdown[mod] {
					trakStruct.ModuleBreakdown[mod] = true
					hasStateChange = true
				}
			} else {
				failCount++
				fmt.Printf("  %s✖ FAIL%s  %s%s%s\n", ui.Red+ui.Bold, ui.Reset, ui.White+ui.Bold, mod, ui.Reset)

				// Print trimmed test failure output
				outStr := strings.TrimSpace(string(output))
				if outStr != "" && detail {
					lines := strings.Split(outStr, "\n")
					limit := 10
					if len(lines) < limit {
						limit = len(lines)
					}
					fmt.Printf("  %s┌── Test Failure Details ──────────────────────────────────────%s\n", ui.Gray, ui.Reset)
					for i := 0; i < limit; i++ {
						fmt.Printf("  %s│%s %s\n", ui.Gray, ui.Reset, lines[i])
					}
					if len(lines) > limit {
						fmt.Printf("  %s│%s %s... (%d more lines)%s\n", ui.Gray, ui.Reset, ui.Gray, len(lines)-limit, ui.Reset)
					}
					fmt.Printf("  %s└── Check exercise code & hints to fix ────────────────────────%s\n\n", ui.Gray, ui.Reset)
				}
			}
		}

		// Save state if any module status changed to complete
		if hasStateChange {
			trakPath := filepath.Join(cwd, "trak.json")
			if updatedBytes, err := json.MarshalIndent(trakStruct, "", "  "); err == nil {
				_ = os.WriteFile(trakPath, updatedBytes, 0644)
			}
		}

		// Calculate updated metrics
		totalModules := len(allKeys)
		completedCount := 0
		for _, k := range allKeys {
			if trakStruct.ModuleBreakdown[k] {
				completedCount++
			}
		}
		progressPercent := (float64(completedCount) / float64(totalModules)) * 100.0

		// Summary Banner
		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)
		fmt.Printf("  %s• Result        :%s  %d passed, %d failed\n", ui.Gray, ui.Reset, passCount, failCount)
		fmt.Printf("  %s• Track Progress:%s  [%s%s%s] %s%.1f%%%s %s(%d of %d completed)%s\n",
			ui.Gray, ui.Reset,
			ui.Green, ui.ProgressBar(progressPercent, 20), ui.Reset,
			ui.Green+ui.Bold, progressPercent, ui.Reset,
			ui.Gray, completedCount, totalModules, ui.Reset,
		)
		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n\n", ui.Gray, ui.Reset)

		if hasStateChange {
			fmt.Printf("  %s✔ Marked completed module(s) in trak.json%s\n", ui.Green, ui.Reset)
			fmt.Printf("  %sRun %strak next%s to jump into your next exercise.\n\n", ui.Gray, ui.Green, ui.Reset)
		}
	},
}

func init() {
	verifyCmd.Flags().BoolVarP(&detail, "detail", "d", false, "In Detail Testing Results")
	verifyCmd.Flags().BoolVarP(&verifyAll, "all", "a", false, "Verify all modules across the workspace")
	verifyCmd.Flags().BoolVarP(&verifyList, "list", "l", false, "List all supported verification runtimes and local toolchain status")
	rootCmd.AddCommand(&verifyCmd)
}
