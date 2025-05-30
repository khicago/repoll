// Package main contains tests for reporting functionality.
// This file provides comprehensive test coverage for report generation
// and formatting functions.
package main

import (
	"strings"
	"testing"
	"time"
)

func TestMakeReport_Report(t *testing.T) {
	// Create a report with multiple actions
	report := MakeReport{
		Actions: []*MakeAction{
			{
				Time:       time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Repository: "test-repo1",
				Duration:   5 * time.Second,
				Success:    true,
				Error:      "",
				Memo:       "Clone completed",
			},
			{
				Time:       time.Date(2023, 1, 1, 12, 1, 0, 0, time.UTC),
				Repository: "test-repo2",
				Duration:   3 * time.Second,
				Success:    false,
				Error:      "failed to clone",
				Memo:       "Clone failed",
			},
		},
	}

	reportStr := report.Report()

	// Check that the report contains expected elements
	expectedElements := []string{
		"test-repo1",
		"test-repo2",
		"Clone completed",
		"Clone failed",
		"failed to clone",
		"5s",
		"3s",
		"Success: true",
		"Success: false",
	}

	for _, element := range expectedElements {
		if !strings.Contains(reportStr, element) {
			t.Errorf("Expected report to contain '%s', but it didn't.\nReport:\n%s", element, reportStr)
		}
	}
}

func TestMakeReport_EmptyReport(t *testing.T) {
	report := MakeReport{Actions: []*MakeAction{}}

	reportStr := report.Report()

	if !strings.Contains(reportStr, "Make Report") {
		t.Error("Expected report to contain header even when empty")
	}
}

func TestMkconfReport_Report(t *testing.T) {
	// Create a mkconf report with multiple actions
	report := MkconfReport{
		Actions: []*MkconfAction{
			{
				Time:        time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Path:        "/home/user/projects/repo1",
				Origin:      "https://github.com/user/repo1.git",
				HasOrigin:   true,
				Uncommitted: false,
				Unmerged:    false,
			},
			{
				Time:        time.Date(2023, 1, 1, 12, 1, 0, 0, time.UTC),
				Path:        "/home/user/projects/repo2",
				Origin:      "git@gitlab.com:org/repo2.git",
				HasOrigin:   true,
				Uncommitted: true,
				Unmerged:    false,
			},
			{
				Time:      time.Date(2023, 1, 1, 12, 2, 0, 0, time.UTC),
				Path:      "/home/user/projects/local-repo",
				HasOrigin: false,
			},
		},
	}

	reportStr := report.Report()

	// Check that the report contains expected elements
	expectedElements := []string{
		"Mkconf Report",
		"/home/user/projects/repo1",
		"/home/user/projects/repo2",
		"/home/user/projects/local-repo",
		"https://github.com/user/repo1.git",
		"git@gitlab.com:org/repo2.git",
		"HasOrigin: true",
		"HasOrigin: false",
		"Uncommitted: true",
		"Uncommitted: false",
		"Unmerged: false",
	}

	for _, element := range expectedElements {
		if !strings.Contains(reportStr, element) {
			t.Errorf("Expected report to contain '%s', but it didn't.\nReport:\n%s", element, reportStr)
		}
	}
}

func TestMkconfReport_EmptyReport(t *testing.T) {
	report := MkconfReport{Actions: []*MkconfAction{}}

	reportStr := report.Report()

	if !strings.Contains(reportStr, "Mkconf Report") {
		t.Error("Expected report to contain header even when empty")
	}
}

func TestMakeAction_String(t *testing.T) {
	action := &MakeAction{
		Time:       time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Repository: "test-repo",
		Duration:   5 * time.Second,
		Success:    true,
		Error:      "",
		Memo:       "Operation completed",
	}

	str := action.String()

	expectedElements := []string{
		"test-repo",
		"5s",
		"Success: true",
		"Operation completed",
	}

	for _, element := range expectedElements {
		if !strings.Contains(str, element) {
			t.Errorf("Expected action string to contain '%s', but it didn't.\nString:\n%s", element, str)
		}
	}
}

func TestMakeAction_StringWithError(t *testing.T) {
	action := &MakeAction{
		Time:       time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Repository: "failed-repo",
		Duration:   2 * time.Second,
		Success:    false,
		Error:      "clone failed: repository not found",
		Memo:       "Failed operation",
	}

	str := action.String()

	expectedElements := []string{
		"failed-repo",
		"2s",
		"Success: false",
		"clone failed: repository not found",
		"Failed operation",
	}

	for _, element := range expectedElements {
		if !strings.Contains(str, element) {
			t.Errorf("Expected action string to contain '%s', but it didn't.\nString:\n%s", element, str)
		}
	}
}

func TestMkconfAction_String(t *testing.T) {
	action := &MkconfAction{
		Time:        time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Path:        "/home/user/projects/test-repo",
		Origin:      "https://github.com/user/test-repo.git",
		HasOrigin:   true,
		Uncommitted: true,
		Unmerged:    false,
	}

	str := action.String()

	expectedElements := []string{
		"/home/user/projects/test-repo",
		"https://github.com/user/test-repo.git",
		"HasOrigin: true",
		"Uncommitted: true",
		"Unmerged: false",
	}

	for _, element := range expectedElements {
		if !strings.Contains(str, element) {
			t.Errorf("Expected action string to contain '%s', but it didn't.\nString:\n%s", element, str)
		}
	}
}

func TestMkconfAction_StringNoOrigin(t *testing.T) {
	action := &MkconfAction{
		Time:        time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Path:        "/home/user/projects/local-repo",
		Origin:      "",
		HasOrigin:   false,
		Uncommitted: false,
		Unmerged:    false,
	}

	str := action.String()

	expectedElements := []string{
		"/home/user/projects/local-repo",
		"HasOrigin: false",
		"Uncommitted: false",
		"Unmerged: false",
	}

	for _, element := range expectedElements {
		if !strings.Contains(str, element) {
			t.Errorf("Expected action string to contain '%s', but it didn't.\nString:\n%s", element, str)
		}
	}

	// Should not contain origin URL when HasOrigin is false
	if strings.Contains(str, "Origin:") && action.Origin == "" {
		t.Error("Expected action string not to contain 'Origin:' when HasOrigin is false and Origin is empty")
	}
}
