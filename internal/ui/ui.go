package ui

import (
	"fmt"
	"strings"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"

	// Trak Web Exact Signature Palette (24-bit TrueColor)
	Emerald   = "\033[38;2;16;185;129m"  // #10b981 (Brand Emerald Green)
	Green     = "\033[38;2;16;185;129m"  // #10b981
	White     = "\033[38;2;245;244;239m" // #f5f4ef (Crisp Warm White)
	LightGray = "\033[38;2;203;213;225m" // #cbd5e1 (Slate 300)
	Gray      = "\033[38;2;148;163;184m" // #94a3b8 (Slate 400)
	Red       = "\033[38;2;239;68;68m"   // #ef4444 (Clean Error Red)
	Yellow    = "\033[38;2;245;158;11m"  // #f59e0b (Amber Accent)
	Cyan      = "\033[38;2;16;185;129m"  // Mapped to Brand Emerald Green
	Blue      = "\033[38;2;16;185;129m"  // Mapped to Brand Emerald Green
	Magenta   = "\033[38;2;16;185;129m"
)

// Success prints a green checkmark with message
func Success(msg string) {
	fmt.Printf("%s✔%s %s\n", Green, Reset, msg)
}

// Error prints a red cross with red message
func Error(msg string) {
	fmt.Printf("%s✘%s %s%s%s\n", Red, Reset, Red, msg, Reset)
}

// Info prints a brand green arrow with info message
func Info(msg string) {
	fmt.Printf("%s➜%s %s\n", Green, Reset, msg)
}

// FoundTemplate prints the brand-themed template discovery line
func FoundTemplate(name, version string) {
	fmt.Printf("%s✔%s %sFound Template:%s  %s%s%s %s(v%s)%s\n",
		Green, Reset, Gray, Reset, White+Bold, name, Reset, Green, version, Reset)
}

// TargetWorkspace prints the resolved target workspace path
func TargetWorkspace(path string) {
	fmt.Printf("%s➜%s %sTarget Workspace:%s %s%s%s\n",
		Green, Reset, Gray, Reset, White+Bold, path, Reset)
}

// CreatedFile prints a clean file creation line with muted gray label and crisp white path
func CreatedFile(path string) {
	fmt.Printf("  %s✔%s %sCreated File:     %s%s%s%s\n",
		Green, Reset, Gray, Reset, White, path, Reset)
}

// CreatedDir prints a clean directory creation line with muted gray label and crisp path
func CreatedDir(path string) {
	fmt.Printf("  %s✔%s %sCreated Directory:%s%s%s%s\n",
		Green, Reset, Gray, Reset, LightGray, path, Reset)
}

// StampedMetadata prints the stamped trak.json metadata manifest line with brand green accent
func StampedMetadata(path string) {
	fmt.Printf("  %s✔%s %sStamped Metadata: %s%s%s%s\n",
		Green, Reset, Gray, Reset, Green+Bold, path, Reset)
}

// CompletedBanner prints the final success message with highlighted resource count
func CompletedBanner(name string, resourceCount int) {
	fmt.Printf("\n%s✔%s %sSuccessfully initialized %s%s%s workspace with %s%d%s resources!\n",
		Green, Reset, Green+Bold, White+Bold, name, Green+Bold, White+Bold, resourceCount, Green+Bold+Reset)
}

// NextSteps prints the formatted next steps guidance
func NextSteps(path string) {
	fmt.Printf("\n  %sNext steps:%s\n", White+Bold, Reset)
	fmt.Printf("    1. %scd %s%s\n", Green, path, Reset)
	fmt.Printf("    2. %scode .%s %s(or open in your preferred editor)%s\n", Green, Reset, Gray, Reset)
	fmt.Printf("    3. Read %sREADME.md%s and start Module 00!\n\n", Green+Bold, Reset)
}

// ProgressBar renders a clean block progress bar
func ProgressBar(percent float64, width int) string {
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
