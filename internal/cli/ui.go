package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// ColorProfile defines color settings for different output types
type ColorProfile struct {
	Success   *color.Color
	Error     *color.Color
	Warning   *color.Color
	Info      *color.Color
	Progress  *color.Color
	Highlight *color.Color
	Muted     *color.Color
}

// DefaultColorProfile returns the default color scheme
func DefaultColorProfile() *ColorProfile {
	return &ColorProfile{
		Success:   color.New(color.FgGreen, color.Bold),
		Error:     color.New(color.FgRed, color.Bold),
		Warning:   color.New(color.FgYellow, color.Bold),
		Info:      color.New(color.FgCyan),
		Progress:  color.New(color.FgBlue, color.Bold),
		Highlight: color.New(color.FgMagenta, color.Bold),
		Muted:     color.New(color.FgHiBlack),
	}
}

// UIManager handles console output formatting and colors
type UIManager struct {
	colors *ColorProfile
	quiet  bool
	verbose bool
}

// NewUIManager creates a new UI manager
func NewUIManager(quiet, verbose bool) *UIManager {
	colors := DefaultColorProfile()
	
	// Disable colors if not in a terminal or if NO_COLOR is set
	if !isTerminal() || os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}
	
	return &UIManager{
		colors:  colors,
		quiet:   quiet,
		verbose: verbose,
	}
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	stat, _ := os.Stdout.Stat()
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// Success prints a success message with checkmark
func (ui *UIManager) Success(format string, args ...interface{}) {
	if ui.quiet {
		return
	}
	message := fmt.Sprintf(format, args...)
	ui.colors.Success.Printf("✓ %s\n", message)
}

// Error prints an error message with X mark
func (ui *UIManager) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	ui.colors.Error.Printf("✗ %s\n", message)
}

// Warning prints a warning message with warning symbol
func (ui *UIManager) Warning(format string, args ...interface{}) {
	if ui.quiet {
		return
	}
	message := fmt.Sprintf(format, args...)
	ui.colors.Warning.Printf("⚠ %s\n", message)
}

// Info prints an informational message
func (ui *UIManager) Info(format string, args ...interface{}) {
	if ui.quiet {
		return
	}
	message := fmt.Sprintf(format, args...)
	ui.colors.Info.Printf("ℹ %s\n", message)
}

// Progress prints a progress message with spinner-like indicator
func (ui *UIManager) Progress(format string, args ...interface{}) {
	if ui.quiet {
		return
	}
	message := fmt.Sprintf(format, args...)
	ui.colors.Progress.Printf("⟳ %s\n", message)
}

// Verbose prints a message only in verbose mode
func (ui *UIManager) Verbose(format string, args ...interface{}) {
	if !ui.verbose || ui.quiet {
		return
	}
	message := fmt.Sprintf(format, args...)
	ui.colors.Muted.Printf("  %s\n", message)
}

// Highlight prints a highlighted message
func (ui *UIManager) Highlight(format string, args ...interface{}) {
	if ui.quiet {
		return
	}
	message := fmt.Sprintf(format, args...)
	ui.colors.Highlight.Printf("★ %s\n", message)
}

// Section prints a section header
func (ui *UIManager) Section(title string) {
	if ui.quiet {
		return
	}
	separator := strings.Repeat("─", len(title)+4)
	ui.colors.Highlight.Printf("\n┌%s┐\n", separator)
	ui.colors.Highlight.Printf("│ %s │\n", strings.ToUpper(title))
	ui.colors.Highlight.Printf("└%s┘\n", separator)
}

// ProcessingRepo shows repository processing status
func (ui *UIManager) ProcessingRepo(repo, action string) {
	if ui.quiet {
		return
	}
	ui.colors.Progress.Printf("⟳ %s %s...\n", action, ui.colors.Highlight.Sprint(repo))
}

// RepoResult shows the result of repository processing
func (ui *UIManager) RepoResult(repo string, success bool, duration time.Duration, err error) {
	if ui.quiet && success {
		return
	}
	
	durationStr := ui.colors.Muted.Sprintf("(%v)", duration)
	repoStr := ui.colors.Highlight.Sprint(repo)
	
	if success {
		ui.colors.Success.Printf("✓ %s %s\n", repoStr, durationStr)
	} else {
		ui.colors.Error.Printf("✗ %s %s\n", repoStr, durationStr)
		if err != nil {
			ui.colors.Error.Printf("  Error: %v\n", err)
		}
	}
}

// Summary prints a summary of operations
func (ui *UIManager) Summary(total, successful, failed int, totalDuration time.Duration) {
	if ui.quiet {
		return
	}
	
	fmt.Println()
	ui.Section("Summary")
	
	ui.colors.Info.Printf("Total repositories: %d\n", total)
	
	if successful > 0 {
		ui.colors.Success.Printf("Successful: %d\n", successful)
	}
	
	if failed > 0 {
		ui.colors.Error.Printf("Failed: %d\n", failed)
	}
	
	ui.colors.Info.Printf("Total time: %v\n", totalDuration)
	
	if failed == 0 && total > 0 {
		ui.colors.Success.Printf("\n🎉 All repositories processed successfully!\n")
	}
}

// DryRun prints a dry run message
func (ui *UIManager) DryRun(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	ui.colors.Warning.Printf("[DRY RUN] %s\n", message)
}

// Banner prints the repoll banner
func (ui *UIManager) Banner() {
	if ui.quiet {
		return
	}
	
	banner := `
🚀 repoll - Git Repository Management Tool
`
	ui.colors.Highlight.Print(banner)
}

// ProgressBar represents a simple progress indicator
type ProgressBar struct {
	ui       *UIManager
	total    int
	current  int
	width    int
	prefix   string
}

// NewProgressBar creates a new progress bar
func (ui *UIManager) NewProgressBar(total int, prefix string) *ProgressBar {
	if ui.quiet {
		return &ProgressBar{ui: ui, total: total}
	}
	
	return &ProgressBar{
		ui:      ui,
		total:   total,
		current: 0,
		width:   40,
		prefix:  prefix,
	}
}

// Update updates the progress bar
func (pb *ProgressBar) Update(current int) {
	if pb.ui.quiet {
		return
	}
	
	pb.current = current
	percent := float64(current) / float64(pb.total)
	filled := int(percent * float64(pb.width))
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.width-filled)
	
	fmt.Printf("\r%s [%s] %d/%d (%.1f%%)", 
		pb.prefix, bar, current, pb.total, percent*100)
	
	if current >= pb.total {
		fmt.Println()
	}
}

// Finish completes the progress bar
func (pb *ProgressBar) Finish() {
	if pb.ui.quiet {
		return
	}
	pb.Update(pb.total)
} 