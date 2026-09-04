package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

type SupportedRuntimeEntry struct {
	Name         string
	TrackPattern string
	Executables  []string
	TestRunner   string
}

var supportedRuntimesList = []SupportedRuntimeEntry{
	{
		Name:         "Go",
		TrackPattern: "lang/go",
		Executables:  []string{"go"},
		TestRunner:   "go test -v ./...",
	},
	{
		Name:         "Python",
		TrackPattern: "lang/python",
		Executables:  []string{"python", "python3"},
		TestRunner:   "unittest discover",
	},
	{
		Name:         "Rust",
		TrackPattern: "lang/rust",
		Executables:  []string{"cargo"},
		TestRunner:   "cargo test",
	},
	{
		Name:         "JavaScript",
		TrackPattern: "lang/js",
		Executables:  []string{"node", "bun"},
		TestRunner:   "node --test",
	},
	{
		Name:         "TypeScript",
		TrackPattern: "lang/ts",
		Executables:  []string{"node", "bun"},
		TestRunner:   "node --test",
	},
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func renderSupportedRuntimes() {
	fmt.Println()
	fmt.Printf("  %s>> Supported Language Verification Runtimes%s\n", ui.Bold+ui.Green, ui.Reset)
	fmt.Printf("  %sList of programming languages supported by 'trak verify' for automated testing.%s\n\n", ui.Gray, ui.Reset)

	const (
		colLangWidth   = 14
		colTrackWidth  = 16
		colRunnerWidth = 22
	)

	fmt.Printf("  %s%s %s %s %s%s\n",
		ui.White+ui.Bold,
		padRight("LANGUAGE", colLangWidth),
		padRight("TRACK PATTERN", colTrackWidth),
		padRight("TEST RUNNER", colRunnerWidth),
		"LOCAL TOOLCHAIN",
		ui.Reset,
	)
	fmt.Printf("  %s────────────────────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)

	installedCount := 0

	for _, entry := range supportedRuntimesList {
		// Detect local toolchain availability
		var detectedBin string
		for _, bin := range entry.Executables {
			if path, err := exec.LookPath(bin); err == nil && path != "" {
				detectedBin = bin
				break
			}
		}

		var statusStr string
		if detectedBin != "" {
			installedCount++
			statusStr = fmt.Sprintf("%s✔ Installed%s  %s(%s)%s", ui.Green, ui.Reset, ui.Gray, detectedBin, ui.Reset)
		} else {
			statusStr = fmt.Sprintf("%s✘ Not in PATH%s %s(requires %s)%s", ui.Red, ui.Reset, ui.Gray, strings.Join(entry.Executables, " / "), ui.Reset)
		}

		fmt.Printf("  %s%s%s %s%s%s %s%s%s %s\n",
			ui.White+ui.Bold, padRight(entry.Name, colLangWidth), ui.Reset,
			ui.Green, padRight(entry.TrackPattern, colTrackWidth), ui.Reset,
			ui.LightGray, padRight(entry.TestRunner, colRunnerWidth), ui.Reset,
			statusStr,
		)
	}

	fmt.Printf("  %s────────────────────────────────────────────────────────────────────────────%s\n", ui.Gray, ui.Reset)
	fmt.Printf("  %s• Total Supported: %d runtimes  • Ready on this machine: %d%s\n\n",
		ui.Gray,
		len(supportedRuntimesList),
		installedCount,
		ui.Reset,
	)

	fmt.Printf("  %sNote:%s Hands-on laboratory tracks (e.g. %stool/docker%s, %sdb/postgres%s, %scloud/aws%s)\n",
		ui.White+ui.Bold, ui.Reset,
		ui.Green, ui.Reset,
		ui.Green, ui.Reset,
		ui.Green, ui.Reset,
	)
	fmt.Printf("        do not use compilers. Verify them manually via %strak done <module>%s.\n\n",
		ui.Green, ui.Reset,
	)
}

var allowedListsCmd = cobra.Command{
	Use:     "allowlists",
	Aliases: []string{"allowlist", "list", "supported", "runtimes"},
	Short:   "List all supported programming language runtimes for automated verification",
	Long: `Display the programming language tracks supported by 'trak verify',
along with the required compiler/test toolchains and their installation status on your system.`,
	Example: `  trak verify allowlists
  trak verify list
  trak verify --list`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		renderSupportedRuntimes()
	},
}

func init() {
	verifyCmd.AddCommand(&allowedListsCmd)
}
