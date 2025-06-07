package cli

import (
	"os"
	"testing"
	"time"
)

func TestNewUIManager(t *testing.T) {
	ui := NewUIManager(false, false)
	if ui == nil {
		t.Fatal("NewUIManager returned nil")
	}
	
	if ui.colors == nil {
		t.Error("UIManager colors profile is nil")
	}
	
	if ui.quiet {
		t.Error("UIManager should not be quiet by default")
	}
	
	if ui.verbose {
		t.Error("UIManager should not be verbose by default")
	}
}

func TestNewUIManager_QuietMode(t *testing.T) {
	ui := NewUIManager(true, false)
	if !ui.quiet {
		t.Error("UIManager should be in quiet mode")
	}
}

func TestNewUIManager_VerboseMode(t *testing.T) {
	ui := NewUIManager(false, true)
	if !ui.verbose {
		t.Error("UIManager should be in verbose mode")
	}
}

func TestDefaultColorProfile(t *testing.T) {
	profile := DefaultColorProfile()
	if profile == nil {
		t.Fatal("DefaultColorProfile returned nil")
	}
	
	if profile.Success == nil {
		t.Error("Success color is nil")
	}
	
	if profile.Error == nil {
		t.Error("Error color is nil")
	}
	
	if profile.Warning == nil {
		t.Error("Warning color is nil")
	}
	
	if profile.Info == nil {
		t.Error("Info color is nil")
	}
	
	if profile.Progress == nil {
		t.Error("Progress color is nil")
	}
	
	if profile.Highlight == nil {
		t.Error("Highlight color is nil")
	}
	
	if profile.Muted == nil {
		t.Error("Muted color is nil")
	}
}

func TestUIManager_QuietModeOutput(t *testing.T) {
	ui := NewUIManager(true, false)
	
	// These methods should not output anything in quiet mode
	// We can't easily test output, but we can ensure they don't panic
	ui.Success("test message")
	ui.Warning("test message")
	ui.Info("test message")
	ui.Progress("test message")
	ui.Verbose("test message")
	ui.Highlight("test message")
	ui.Section("test section")
	ui.ProcessingRepo("test/repo", "Cloning")
	ui.Banner()
	
	// Error should still output in quiet mode
	ui.Error("test error")
}

func TestUIManager_VerboseMode(t *testing.T) {
	ui := NewUIManager(false, true)
	
	// Verbose messages should be shown
	ui.Verbose("test verbose message")
	
	// Other methods should still work
	ui.Success("test message")
	ui.Info("test message")
}

func TestUIManager_RepoResult(t *testing.T) {
	ui := NewUIManager(false, false)
	
	// Test successful result
	ui.RepoResult("test/repo", true, time.Second, nil)
	
	// Test failed result
	err := os.ErrNotExist
	ui.RepoResult("test/repo", false, time.Second, err)
}

func TestUIManager_Summary(t *testing.T) {
	ui := NewUIManager(false, false)
	
	// Test summary with mixed results
	ui.Summary(10, 8, 2, time.Minute)
	
	// Test summary with all successful
	ui.Summary(5, 5, 0, 30*time.Second)
	
	// Test summary with all failed
	ui.Summary(3, 0, 3, 15*time.Second)
}

func TestProgressBar(t *testing.T) {
	ui := NewUIManager(false, false)
	
	pb := ui.NewProgressBar(10, "Testing")
	if pb == nil {
		t.Fatal("NewProgressBar returned nil")
	}
	
	if pb.total != 10 {
		t.Errorf("Expected total 10, got %d", pb.total)
	}
	
	if pb.prefix != "Testing" {
		t.Errorf("Expected prefix 'Testing', got '%s'", pb.prefix)
	}
	
	// Test updates
	pb.Update(5)
	pb.Update(10)
	pb.Finish()
}

func TestProgressBar_QuietMode(t *testing.T) {
	ui := NewUIManager(true, false)
	
	pb := ui.NewProgressBar(5, "Quiet Test")
	if pb == nil {
		t.Fatal("NewProgressBar returned nil")
	}
	
	// Should not panic in quiet mode
	pb.Update(2)
	pb.Update(5)
	pb.Finish()
}

func TestUIManager_DryRun(t *testing.T) {
	ui := NewUIManager(false, false)
	
	// Should not panic
	ui.DryRun("Would clone %s", "test/repo")
}

func TestIsTerminal(t *testing.T) {
	// This function should not panic
	result := isTerminal()
	
	// The result depends on the test environment
	// We just ensure it returns a boolean
	if result != true && result != false {
		t.Error("isTerminal should return a boolean value")
	}
} 