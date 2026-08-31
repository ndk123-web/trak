package ui

import "fmt"

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[90m"
	LightGray = "\033[37m"
	White     = "\033[97m"
)

// Success prints a green checkmark with bold message
func Success(msg string) {
	fmt.Printf("%s✔%s %s\n", Green, Reset, msg)
}

// Error prints a red cross with red message
func Error(msg string) {
	fmt.Printf("%s✘%s %s%s%s\n", Red, Reset, Red, msg, Reset)
}

// Info prints a cyan arrow with info message
func Info(msg string) {
	fmt.Printf("%s➜%s %s\n", Cyan, Reset, msg)
}

// FoundTemplate prints the brand-themed template discovery line matching Trak Green
func FoundTemplate(name, version string) {
	fmt.Printf("%s✔%s %sFound Template:%s  %s%s%s %s(v%s)%s\n",
		Green, Reset, Gray, Reset, Green+Bold, name, Reset, Cyan, version, Reset)
}

// TargetWorkspace prints the resolved target workspace path
func TargetWorkspace(path string) {
	fmt.Printf("%s➜%s %sTarget Workspace:%s %s%s%s\n",
		Cyan, Reset, Gray, Reset, White+Bold, path, Reset)
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

// StampedMetadata prints the stamped trak.json metadata manifest line with cyan accent
func StampedMetadata(path string) {
	fmt.Printf("  %s✔%s %sStamped Metadata: %s%s%s%s\n",
		Cyan, Reset, Gray, Reset, Cyan+Bold, path, Reset)
}

// CompletedBanner prints the final success message with highlighted resource count
func CompletedBanner(name string, resourceCount int) {
	fmt.Printf("\n%s✔%s %sSuccessfully initialized %s%s%s workspace with %s%d%s resources! 🎉\n",
		Green, Reset, Green+Bold, White+Bold, name, Green+Bold, Yellow+Bold, resourceCount, Green+Bold+Reset)
}

// NextSteps prints the formatted next steps guidance
func NextSteps(path string) {
	fmt.Printf("\n  %sNext steps:%s\n", White+Bold, Reset)
	fmt.Printf("    1. %scd %s%s\n", Cyan, path, Reset)
	fmt.Printf("    2. %scode .%s %s(or open in your preferred editor)%s\n", Green, Reset, Gray, Reset)
	fmt.Printf("    3. Read %sREADME.md%s and start Module 00!\n\n", Yellow, Reset)
}
