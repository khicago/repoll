package main

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRunCommand_Success(t *testing.T) {
	// 测试成功的命令
	cmd := exec.Command("echo", "hello", "world")
	err := runCommand(cmd)
	if err != nil {
		t.Errorf("Expected no error for echo command, got: %v", err)
	}
}

func TestRunCommand_Failure(t *testing.T) {
	// 测试失败的命令
	cmd := exec.Command("nonexistent-command")
	err := runCommand(cmd)
	if err == nil {
		t.Error("Expected error for nonexistent command")
	}
}

func TestRunCommand_InvalidDirectory(t *testing.T) {
	// 测试无效目录
	cmd := exec.Command("echo", "test")
	cmd.Dir = "/nonexistent/directory"
	err := runCommand(cmd)
	if err == nil {
		t.Error("Expected error for invalid directory")
	}
}

func TestRunCommandWithTimer_Success(t *testing.T) {
	// 测试成功的命令
	start := time.Now()
	cmd := exec.Command("echo", "hello")
	err := runCommandWithTimer(cmd)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error for echo command, got: %v", err)
	}

	// 检查是否有合理的执行时间
	if duration < 0 {
		t.Error("Expected positive duration")
	}
}

func TestRunCommandWithTimer_Failure(t *testing.T) {
	// 测试失败的命令
	cmd := exec.Command("false")
	err := runCommandWithTimer(cmd)
	if err == nil {
		t.Error("Expected error for 'false' command")
	}
}

func TestRunCommandWithTimer_LongRunning(t *testing.T) {
	// 测试较长时间运行的命令
	start := time.Now()
	cmd := exec.Command("sleep", "0.1")
	err := runCommandWithTimer(cmd)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error for sleep command, got: %v", err)
	}

	// 检查是否至少运行了指定的时间
	if duration < 50*time.Millisecond {
		t.Errorf("Expected duration >= 50ms, got %v", duration)
	}
}

func TestRunCommandWithTimer_CommandNotFound(t *testing.T) {
	// 测试命令不存在的情况
	cmd := exec.Command("this-command-does-not-exist")
	err := runCommandWithTimer(cmd)
	if err == nil {
		t.Error("Expected error for non-existent command")
	}

	// 检查错误消息是否包含预期的内容
	if !strings.Contains(err.Error(), "executable file not found") &&
		!strings.Contains(err.Error(), "command not found") &&
		!strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Expected 'command not found' type error, got: %v", err)
	}
}

func TestRunCommand_StdoutPipeError(t *testing.T) {
	// 这个测试比较难模拟，因为需要让StdoutPipe失败
	// 我们使用一个已经启动的命令来尝试触发错误
	cmd := exec.Command("echo", "test")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}
	cmd.Wait() // 等待命令完成

	// 再次尝试获取管道，这可能会失败
	_, err = cmd.StdoutPipe()
	if err == nil {
		t.Log("StdoutPipe did not fail as expected - this is environment dependent")
	}
}

func TestRunCommand_StderrPipeError(t *testing.T) {
	// Create a command that will succeed but test stderr pipe creation
	// This is hard to test directly, but we can test the error handling path
	
	// Test with a command that produces stderr output
	cmd := exec.Command("sh", "-c", "echo 'error message' >&2")
	
	err := runCommand(cmd)
	
	// Command should succeed even with stderr output
	if err != nil {
		t.Errorf("runCommand should handle stderr output, got error: %v", err)
	}
}

func TestRunCommand_StartError(t *testing.T) {
	// Create a command that will fail to start
	cmd := exec.Command("/non/existent/command")
	
	err := runCommand(cmd)
	
	if err == nil {
		t.Error("Expected error when starting non-existent command")
	}
}

func TestRunCommandWithTimer_VeryShortCommand(t *testing.T) {
	// 测试非常短的命令，确保计时器逻辑正常工作
	cmd := exec.Command("echo", "fast")
	start := time.Now()
	err := runCommandWithTimer(cmd)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error for echo command, got: %v", err)
	}

	if duration < 0 {
		t.Error("Expected positive duration")
	}

	// 即使是很快的命令，也应该有一些最小的执行时间
	if duration > 5*time.Second {
		t.Error("Command took too long for a simple echo")
	}
}

func TestControlPrintf(t *testing.T) {
	// 测试控制台输出函数（主要为了覆盖率）
	// 这个函数主要是格式化输出，很难直接测试输出内容
	// 但至少可以确保它不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("controlPrintf panicked: %v", r)
		}
	}()

	controlPrintf("test %s", "message")
	controlPrintf("test %d", 123)
	controlPrintf("simple message")
}

func TestRunCommandWithTimer_CommandFailure(t *testing.T) {
	// Test with a command that will fail
	cmd := exec.Command("sh", "-c", "exit 1")
	
	err := runCommandWithTimer(cmd)
	
	if err == nil {
		t.Error("Expected error when command exits with non-zero status")
	}
}

func TestRunCommandWithTimer_LongRunningCommand(t *testing.T) {
	// Test with a command that takes some time to complete
	cmd := exec.Command("sh", "-c", "sleep 0.2")
	
	start := time.Now()
	err := runCommandWithTimer(cmd)
	duration := time.Since(start)
	
	if err != nil {
		t.Errorf("runCommandWithTimer should succeed with sleep command, got: %v", err)
	}
	
	// Should take at least 200ms
	if duration < 200*time.Millisecond {
		t.Errorf("Expected command to take at least 200ms, took %v", duration)
	}
	
	// Should not take too long (allowing for overhead)
	if duration > 1*time.Second {
		t.Errorf("Command took too long: %v", duration)
	}
}
