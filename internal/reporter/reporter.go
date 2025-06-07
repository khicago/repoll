// report.go
package reporter

import (
	"fmt"
	"strings"
	"time"
)

// MakeReport represents a report for repository processing operations
type MakeReport struct {
	Actions []*MakeAction
}

// MakeAction represents a single repository processing action
type MakeAction struct {
	Time       time.Time
	Repository string
	Duration   time.Duration
	Success    bool
	Error      string
	Memo       string
}

// MkconfReport represents a report for configuration generation operations
type MkconfReport struct {
	Actions []*MkconfAction
}

// MkconfAction represents a single repository discovery action
type MkconfAction struct {
	Time        time.Time
	Path        string
	Origin      string
	HasOrigin   bool
	Uncommitted bool
	Unmerged    bool
}

// Report generates a formatted string report of repository operations
func (mr *MakeReport) Report() string {
	if len(mr.Actions) == 0 {
		return "No actions performed."
	}

	var report strings.Builder
	report.WriteString("=== Repository Processing Report ===\n\n")

	successCount := 0
	var totalDuration time.Duration

	for _, action := range mr.Actions {
		totalDuration += action.Duration
		
		status := "✅ SUCCESS"
		if !action.Success {
			status = "❌ FAILED"
		} else {
			successCount++
		}

		report.WriteString(fmt.Sprintf("%s %s (took %v)\n", status, action.Repository, action.Duration))
		
		if action.Memo != "" {
			report.WriteString(fmt.Sprintf("   📝 %s\n", action.Memo))
		}
		
		if !action.Success && action.Error != "" {
			report.WriteString(fmt.Sprintf("   ❗ Error: %s\n", action.Error))
		}
		
		report.WriteString("\n")
	}

	report.WriteString("=== Summary ===\n")
	report.WriteString(fmt.Sprintf("Total repositories: %d\n", len(mr.Actions)))
	report.WriteString(fmt.Sprintf("Successful: %d\n", successCount))
	report.WriteString(fmt.Sprintf("Failed: %d\n", len(mr.Actions)-successCount))
	report.WriteString(fmt.Sprintf("Total time: %v\n", totalDuration))

	return report.String()
}

// Report generates a formatted string report of repository discovery
func (mr *MkconfReport) Report() string {
	if len(mr.Actions) == 0 {
		return "No repositories discovered."
	}

	var report strings.Builder
	report.WriteString("=== Repository Discovery Report ===\n\n")

	originCount := 0
	uncommittedCount := 0
	unmergedCount := 0

	for _, action := range mr.Actions {
		status := "📁"
		warnings := []string{}

		if action.HasOrigin {
			status = "📦"
			originCount++
		} else {
			warnings = append(warnings, "no origin")
		}

		if action.Uncommitted {
			warnings = append(warnings, "uncommitted changes")
			uncommittedCount++
		}

		if action.Unmerged {
			warnings = append(warnings, "unmerged changes")
			unmergedCount++
		}

		report.WriteString(fmt.Sprintf("%s %s\n", status, action.Path))
		
		if action.HasOrigin && action.Origin != "" {
			report.WriteString(fmt.Sprintf("   🔗 %s\n", action.Origin))
		}

		if len(warnings) > 0 {
			report.WriteString(fmt.Sprintf("   ⚠️  %s\n", strings.Join(warnings, ", ")))
		}

		report.WriteString("\n")
	}

	report.WriteString("=== Summary ===\n")
	report.WriteString(fmt.Sprintf("Total repositories: %d\n", len(mr.Actions)))
	report.WriteString(fmt.Sprintf("With origin: %d\n", originCount))
	if uncommittedCount > 0 {
		report.WriteString(fmt.Sprintf("With uncommitted changes: %d\n", uncommittedCount))
	}
	if unmergedCount > 0 {
		report.WriteString(fmt.Sprintf("With unmerged changes: %d\n", unmergedCount))
	}

	return report.String()
}

// String provides a simple string representation of MakeAction
func (ma *MakeAction) String() string {
	status := "SUCCESS"
	if !ma.Success {
		status = "FAILED"
	}
	return fmt.Sprintf("[%s] %s (%v) - %s", status, ma.Repository, ma.Duration, ma.Memo)
}

// String provides a simple string representation of MkconfAction
func (ma *MkconfAction) String() string {
	return fmt.Sprintf("[%s] %s (origin: %v, uncommitted: %v)", 
		ma.Time.Format("15:04:05"), ma.Path, ma.HasOrigin, ma.Uncommitted)
}
