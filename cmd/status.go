package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var statusCmd = cobra.Command{
	Use:   "status",
	Short: "Display workspace progress, active track details, and module status",
	Long: `Inspects the current workspace for trak.json, calculates module completion metrics,
and renders a visual progress dashboard.`,
	Example: `  trak status`,
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
			fmt.Printf("  %s💡 Tip:%s Navigate into your learning track folder (e.g. %scd ./learn-rust%s)\n", ui.White+ui.Bold, ui.Reset, ui.Green, ui.Reset)
			fmt.Printf("     or run %strak init <category>/<tool>%s to materialize one.\n\n", ui.Green, ui.Reset)
			return
		}

		s := spinner.New(spinner.CharSets[14], 80*time.Millisecond)
		s.Color("green")
		s.Suffix = fmt.Sprintf(" %sAnalyzing workspace & curriculum progress...%s", ui.Gray, ui.Reset)
		s.Start()

		dataBytes, err := os.ReadFile(trakPath)
		if err != nil {
			s.Stop()
			ui.Error(fmt.Sprintf("Failed to read trak.json: %v", err))
			return
		}

		statusModel := &models.StatusModel{}
		err = json.NewDecoder(bytes.NewReader(dataBytes)).Decode(statusModel)
		if err != nil {
			s.Stop()
			ui.Error(fmt.Sprintf("Failed to parse trak.json manifest: %v", err))
			return
		}

		// Small delay for smooth CLI UX
		time.Sleep(700 * time.Millisecond)
		s.Stop()

		// Parse and format CreatedAt timestamp
		createdDisplay := statusModel.CreatedAt
		if t, err := time.Parse(time.RFC3339Nano, statusModel.CreatedAt); err == nil {
			createdDisplay = t.Format("02 Jan 2006, 15:04 UTC")
		} else if t, err := time.Parse(time.RFC3339, statusModel.CreatedAt); err == nil {
			createdDisplay = t.Format("02 Jan 2006, 15:04 UTC")
		}

		trackTitle := statusModel.Name
		if trackTitle == "" {
			trackTitle = statusModel.Id
		}

		// Calculate module progress
		totalModules := len(statusModel.ModuleBreakdown)
		completedCount := 0

		// Sort module keys alphabetically/numerically
		var moduleKeys []string
		for k, isDone := range statusModel.ModuleBreakdown {
			moduleKeys = append(moduleKeys, k)
			if isDone {
				completedCount++
			}
		}
		sort.Strings(moduleKeys)

		var progressPercent float64 = 0
		if totalModules > 0 {
			progressPercent = (float64(completedCount) / float64(totalModules)) * 100.0
		} else if statusModel.Progress > 0 {
			progressPercent = float64(statusModel.Progress)
		}

		// ── Render Professional Terminal Dashboard ──
		fmt.Println()
		fmt.Printf("  %s%s⚡ Trak Workspace Status%s\n", ui.Bold, ui.Green, ui.Reset)
		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)

		fmt.Printf("  %s• Track Name   :%s  %s%s%s\n", ui.Gray, ui.Reset, ui.White+ui.Bold, trackTitle, ui.Reset)
		fmt.Printf("  %s• Blueprint ID :%s  %s%s%s\n", ui.Gray, ui.Reset, ui.Green, statusModel.Id, ui.Reset)
		fmt.Printf("  %s• Version      :%s  v%s\n", ui.Gray, ui.Reset, statusModel.Version)
		if statusModel.Author != "" {
			fmt.Printf("  %s• Author       :%s  %s\n", ui.Gray, ui.Reset, statusModel.Author)
		}
		fmt.Printf("  %s• Initialized  :%s  %s\n", ui.Gray, ui.Reset, createdDisplay)

		fmt.Println()
		fmt.Printf("  %sProgress Overview:%s\n", ui.White+ui.Bold, ui.Reset)
		progressBar := makeProgressBar(progressPercent, 24)
		fmt.Printf("  [%s%s%s] %s%.1f%%%s  %s(%d of %d modules completed)%s\n\n",
			ui.Green, progressBar, ui.Reset,
			ui.Green+ui.Bold, progressPercent, ui.Reset,
			ui.Gray, completedCount, totalModules, ui.Reset,
		)

		if len(moduleKeys) > 0 {
			fmt.Printf("  %sModule Breakdown:%s\n", ui.White+ui.Bold, ui.Reset)
			for _, key := range moduleKeys {
				isDone := statusModel.ModuleBreakdown[key]
				if isDone {
					fmt.Printf("    %s✔%s  %sModule %s%s  %s[COMPLETED]%s\n",
						ui.Green, ui.Reset,
						ui.White, key, ui.Reset,
						ui.Green+ui.Bold, ui.Reset,
					)
				} else {
					fmt.Printf("    %s○%s  %sModule %s%s  %s[PENDING]%s\n",
						ui.Gray, ui.Reset,
						ui.LightGray, key, ui.Reset,
						ui.Gray, ui.Reset,
					)
				}
			}
		}

		fmt.Printf("  %s────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)
		if completedCount < totalModules {
			// Find first pending module
			nextMod := ""
			for _, k := range moduleKeys {
				if !statusModel.ModuleBreakdown[k] {
					nextMod = k
					break
				}
			}
			if nextMod != "" {
				fmt.Printf("  %s💡 Next Step:%s Work on %sModule %s%s and read its %sREADME.md%s\n",
					ui.Green+ui.Bold, ui.Reset, ui.Green+ui.Bold, nextMod, ui.Reset, ui.White+ui.Bold, ui.Reset)
			}
		} else if totalModules > 0 && completedCount == totalModules {
			fmt.Printf("  %s🎉 Congratulations! You have completed all modules in this track!%s\n", ui.Green+ui.Bold, ui.Reset)
		}
		fmt.Println()
	},
}

func makeProgressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}

	var sb strings.Builder
	for i := 0; i < filled; i++ {
		sb.WriteString("█")
	}
	for i := filled; i < width; i++ {
		sb.WriteString("░")
	}
	return sb.String()
}

func init() {
	rootCmd.AddCommand(&statusCmd)
}
