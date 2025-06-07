package reporter

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestMakeReport_NewReport(t *testing.T) {
	report := &MakeReport{Actions: make([]*MakeAction, 0)}
	
	if report.Actions == nil {
		t.Error("Actions slice should be initialized")
	}
	
	if len(report.Actions) != 0 {
		t.Error("New report should have empty actions")
	}
}

func TestMakeReport_AddAction(t *testing.T) {
	report := &MakeReport{Actions: make([]*MakeAction, 0)}
	
	action := &MakeAction{
		Repository: "test/repo",
		Success:    true,
		Time:       time.Now(),
		Duration:   time.Second,
		Error:      "",
		Memo:       "Test memo",
	}
	
	report.Actions = append(report.Actions, action)
	
	if len(report.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(report.Actions))
	}
	
	if report.Actions[0].Repository != "test/repo" {
		t.Errorf("Expected repository 'test/repo', got %s", report.Actions[0].Repository)
	}
}

func TestMakeReport_Report_NoActions(t *testing.T) {
	report := &MakeReport{Actions: make([]*MakeAction, 0)}
	
	output := report.Report()
	if !strings.Contains(output, "No actions performed") {
		t.Error("Expected 'No actions performed' message for empty report")
	}
}

func TestMakeReport_Report_WithActions(t *testing.T) {
	report := &MakeReport{Actions: make([]*MakeAction, 0)}
	
	// 添加成功的操作
	successAction := &MakeAction{
		Repository: "test/success",
		Success:    true,
		Time:       time.Now(),
		Duration:   time.Millisecond * 100,
		Error:      "",
		Memo:       "Test success",
	}
	
	// 添加失败的操作
	failAction := &MakeAction{
		Repository: "test/fail",
		Success:    false,
		Time:       time.Now(),
		Duration:   time.Millisecond * 50,
		Error:      "Failed to process",
		Memo:       "Test fail",
	}
	
	report.Actions = append(report.Actions, successAction, failAction)
	
	output := report.Report()
	
	// 检查总结信息
	if !strings.Contains(output, "Repository Processing Report") {
		t.Error("Expected 'Repository Processing Report' header")
	}
	
	if !strings.Contains(output, "Total repositories: 2") {
		t.Error("Expected total count of 2")
	}
	
	if !strings.Contains(output, "Successful: 1") {
		t.Error("Expected success count of 1")
	}
	
	if !strings.Contains(output, "Failed: 1") {
		t.Error("Expected failure count of 1")
	}
	
	// 检查具体操作信息
	if !strings.Contains(output, "test/success") {
		t.Error("Expected successful repository name")
	}
	
	if !strings.Contains(output, "test/fail") {
		t.Error("Expected failed repository name")
	}
	
	if !strings.Contains(output, "SUCCESS") {
		t.Error("Expected success indicator")
	}
	
	if !strings.Contains(output, "FAILED") {
		t.Error("Expected failure indicator")
	}
}

func TestMakeReport_Report_TimeFormatting(t *testing.T) {
	report := &MakeReport{Actions: make([]*MakeAction, 0)}
	
	action := &MakeAction{
		Repository: "test/time",
		Success:    true,
		Time:       time.Now(),
		Duration:   time.Second + time.Millisecond*500,
		Error:      "",
		Memo:       "Test time formatting",
	}
	
	report.Actions = append(report.Actions, action)
	
	output := report.Report()
	
	// 检查持续时间格式
	if !strings.Contains(output, "1.5s") {
		t.Error("Expected duration '1.5s' in output")
	}
}

func TestMakeAction_Fields(t *testing.T) {
	now := time.Now()
	duration := time.Second * 5
	
	action := &MakeAction{
		Repository: "github.com/test/repo",
		Success:    true,
		Time:       now,
		Duration:   duration,
		Error:      "",
		Memo:       "Important repository for project X",
	}
	
	// 验证所有字段都正确设置
	if action.Repository != "github.com/test/repo" {
		t.Errorf("Repository field mismatch")
	}
	
	if !action.Success {
		t.Errorf("Success field should be true")
	}
	
	if !action.Time.Equal(now) {
		t.Errorf("Time field mismatch")
	}
	
	if action.Duration != duration {
		t.Errorf("Duration field mismatch")
	}
	
	if action.Error != "" {
		t.Errorf("Error field should be empty")
	}
	
	if action.Memo != "Important repository for project X" {
		t.Errorf("Memo field mismatch")
	}
}

func TestMakeAction_String(t *testing.T) {
	action := &MakeAction{
		Repository: "test/repo",
		Success:    true,
		Time:       time.Now(),
		Duration:   time.Second,
		Error:      "",
		Memo:       "Test memo",
	}
	
	str := action.String()
	if !strings.Contains(str, "SUCCESS") {
		t.Error("Expected 'SUCCESS' in string representation")
	}
	
	if !strings.Contains(str, "test/repo") {
		t.Error("Expected repository name in string representation")
	}
	
	if !strings.Contains(str, "Test memo") {
		t.Error("Expected memo in string representation")
	}
}

func TestMakeAction_String_Failed(t *testing.T) {
	action := &MakeAction{
		Repository: "test/repo",
		Success:    false,
		Time:       time.Now(),
		Duration:   time.Second,
		Error:      "Some error",
		Memo:       "Test memo",
	}
	
	str := action.String()
	if !strings.Contains(str, "FAILED") {
		t.Error("Expected 'FAILED' in string representation")
	}
}

func TestMkconfReport_Report_NoActions(t *testing.T) {
	report := &MkconfReport{Actions: make([]*MkconfAction, 0)}
	
	output := report.Report()
	if !strings.Contains(output, "No repositories discovered") {
		t.Error("Expected 'No repositories discovered' message for empty report")
	}
}

func TestMkconfReport_Report_WithActions(t *testing.T) {
	report := &MkconfReport{Actions: make([]*MkconfAction, 0)}
	
	action := &MkconfAction{
		Time:        time.Now(),
		Path:        "/path/to/repo",
		Origin:      "https://github.com/user/repo.git",
		HasOrigin:   true,
		Uncommitted: false,
		Unmerged:    false,
	}
	
	report.Actions = append(report.Actions, action)
	
	output := report.Report()
	
	if !strings.Contains(output, "Repository Discovery Report") {
		t.Error("Expected 'Repository Discovery Report' header")
	}
	
	if !strings.Contains(output, "/path/to/repo") {
		t.Error("Expected repository path")
	}
	
	if !strings.Contains(output, "https://github.com/user/repo.git") {
		t.Error("Expected origin URL")
	}
}

func TestMkconfAction_String(t *testing.T) {
	now := time.Now()
	action := &MkconfAction{
		Time:        now,
		Path:        "/path/to/repo",
		Origin:      "https://github.com/user/repo.git",
		HasOrigin:   true,
		Uncommitted: false,
		Unmerged:    false,
	}
	
	str := action.String()
	if !strings.Contains(str, "/path/to/repo") {
		t.Error("Expected path in string representation")
	}
	
	if !strings.Contains(str, "origin: true") {
		t.Error("Expected origin status in string representation")
	}
}

func TestMakeReport_Report_LargeOutput(t *testing.T) {
	report := &MakeReport{Actions: make([]*MakeAction, 0)}
	
	// 添加大量操作来测试输出格式的稳定性
	for i := 0; i < 100; i++ {
		action := &MakeAction{
			Repository: fmt.Sprintf("test/repo%d", i),
			Success:    i%2 == 0, // 交替成功/失败
			Time:       time.Now().Add(time.Duration(i) * time.Second),
			Duration:   time.Duration(i+1) * time.Millisecond,
			Error:      "",
			Memo:       fmt.Sprintf("Memo %d", i),
		}
		report.Actions = append(report.Actions, action)
	}
	
	output := report.Report()
	
	// 检查总结统计
	if !strings.Contains(output, "Total repositories: 100") {
		t.Error("Expected total count of 100")
	}
	
	if !strings.Contains(output, "Successful: 50") {
		t.Error("Expected success count of 50")
	}
	
	if !strings.Contains(output, "Failed: 50") {
		t.Error("Expected failure count of 50")
	}
	
	// 确保输出包含所有仓库
	for i := 0; i < 100; i++ {
		repoName := fmt.Sprintf("test/repo%d", i)
		if !strings.Contains(output, repoName) {
			t.Errorf("Expected repository %s in output", repoName)
		}
	}
} 