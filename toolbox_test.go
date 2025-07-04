package toolbox

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetCPUUsage(t *testing.T) {
	toolbox := Toolbox{}
	usage, err := toolbox.GetCPUUsage()

	if err != nil {
		t.Logf("GetCPUUsage failed (expected in test environment): %v", err)
		return
	}

	if usage < 0 || usage > 100 {
		t.Errorf("Expected CPU usage between 0-100, got %f", usage)
	}

	t.Logf("CPU Usage: %.2f%%", usage)
}

func TestGetCPULimit(t *testing.T) {
	toolbox := Toolbox{}
	limit, err := toolbox.GetCPULimit()

	if err != nil {
		t.Logf("GetCPULimit failed (expected in test environment): %v", err)
		return
	}

	if limit <= 0 {
		t.Errorf("Expected CPU limit > 0, got %f", limit)
	}

	t.Logf("CPU Limit: %.2f cores", limit)
}

func TestGetAvailableCPU(t *testing.T) {
	toolbox := Toolbox{}
	available, err := toolbox.GetAvailableCPU()

	if err != nil {
		t.Logf("GetAvailableCPU failed (expected in test environment): %v", err)
		return
	}

	if available < 0 {
		t.Errorf("Expected available CPU >= 0, got %f", available)
	}

	t.Logf("Available CPU: %.2f cores", available)
}

func TestGetMemoryUsage(t *testing.T) {
	toolbox := Toolbox{}
	usage, err := toolbox.GetMemoryUsage()

	if err != nil {
		t.Logf("GetMemoryUsage failed (expected in test environment): %v", err)
		return
	}

	if usage < 0 {
		t.Errorf("Expected memory usage >= 0, got %d", usage)
	}

	t.Logf("Memory Usage: %d bytes (%.2f MB)", usage, float64(usage)/(1024*1024))
}

func TestGetMemoryLimit(t *testing.T) {
	toolbox := Toolbox{}
	limit, err := toolbox.GetMemoryLimit()

	if err != nil {
		t.Logf("GetMemoryLimit failed (expected in test environment): %v", err)
		return
	}

	if limit <= 0 {
		t.Errorf("Expected memory limit > 0, got %d", limit)
	}

	t.Logf("Memory Limit: %d bytes (%.2f MB)", limit, float64(limit)/(1024*1024))
}

func TestGetMemoryUsagePercent(t *testing.T) {
	toolbox := Toolbox{}
	percent, err := toolbox.GetMemoryUsagePercent()

	if err != nil {
		t.Logf("GetMemoryUsagePercent failed (expected in test environment): %v", err)
		return
	}

	if percent < 0 || percent > 100 {
		t.Errorf("Expected memory usage percent between 0-100, got %f", percent)
	}

	t.Logf("Memory Usage Percent: %.2f%%", percent)
}

func TestGetAvailableMemory(t *testing.T) {
	toolbox := Toolbox{}
	available, err := toolbox.GetAvailableMemory()

	if err != nil {
		t.Logf("GetAvailableMemory failed (expected in test environment): %v", err)
		return
	}

	if available < 0 {
		t.Errorf("Expected available memory >= 0, got %d", available)
	}

	t.Logf("Available Memory: %d bytes (%.2f MB)", available, float64(available)/(1024*1024))
}

func TestGetPsOutput(t *testing.T) {
	toolbox := Toolbox{}
	output, err := toolbox.GetPsOutput()

	if err != nil {
		t.Logf("GetPsOutput failed (ps command may not be available): %v", err)
		return
	}

	if output == "" {
		t.Error("Expected non-empty ps output")
	}

	// Check if output contains process information
	if !strings.Contains(output, "PID") && !strings.Contains(output, "USER") {
		t.Error("Expected ps output to contain process information")
	}

	// Count number of processes
	lines := strings.Split(output, "\n")
	t.Logf("PS output: %d lines", len(lines))
}

func TestGetUptimeOutput(t *testing.T) {
	toolbox := Toolbox{}
	output, err := toolbox.GetUptimeOutput()

	if err != nil {
		t.Logf("GetUptimeOutput failed (uptime command may not be available): %v", err)
		return
	}

	if output == "" {
		t.Error("Expected non-empty uptime output")
	}

	// Check if output contains uptime information
	if !strings.Contains(output, "up") && !strings.Contains(output, "load average") {
		t.Error("Expected uptime output to contain uptime or load average information")
	}

	t.Logf("Uptime: %s", strings.TrimSpace(output))
}

func TestReadFile(t *testing.T) {
	// Test reading a file that exists
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	content, err := readFile(testFile)
	if err != nil {
		t.Errorf("readFile failed: %v", err)
	}

	if content != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, content)
	}

	// Test reading a file that doesn't exist
	_, err = readFile("/nonexistent/file")
	if err == nil {
		t.Error("Expected error when reading nonexistent file")
	}
}

func TestFileExists(t *testing.T) {
	// Test file that exists
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if !fileExists(testFile) {
		t.Error("Expected fileExists to return true for existing file")
	}

	// Test file that doesn't exist
	if fileExists("/nonexistent/file") {
		t.Error("Expected fileExists to return false for nonexistent file")
	}
}

func TestParseCgroupV2CPUUsage(t *testing.T) {
	// Test valid input
	content := "usage_usec 123456789"
	usage, err := parseCgroupV2CPUUsage(content)
	if err != nil {
		t.Errorf("parseCgroupV2CPUUsage failed: %v", err)
	}
	if usage <= 0 {
		t.Errorf("Expected usage > 0, got %f", usage)
	}

	// Test invalid input
	_, err = parseCgroupV2CPUUsage("invalid content")
	if err == nil {
		t.Error("Expected error for invalid content")
	}

	// Test empty content
	_, err = parseCgroupV2CPUUsage("")
	if err == nil {
		t.Error("Expected error for empty content")
	}
}

func TestParseTopCPUUsage(t *testing.T) {
	// Test standard top output format
	output := `top - 10:30:00 up 2 days, 20:45,  1 user,  load average: 0.52, 0.58, 0.59
Tasks: 123 total,   1 running, 122 sleeping,   0 stopped,   0 zombie
%Cpu(s):  5.2 us,  2.1 sy,  0.0 ni, 92.7 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
MiB Mem :  16384.0 total,   8192.0 free,   4096.0 used,   4096.0 buff/cache`

	usage, err := parseTopCPUUsage(output)
	if err != nil {
		t.Errorf("parseTopCPUUsage failed: %v", err)
	}

	expected := 100.0 - 92.7 // 100 - idle
	epsilon := 0.001
	if usage < expected-epsilon || usage > expected+epsilon {
		t.Errorf("Expected CPU usage %f, got %f", expected, usage)
	}

	// Test alternative format
	output2 := `CPU usage: 15.2% user, 8.1% system, 76.7% idle`
	usage2, err := parseTopCPUUsage(output2)
	if err != nil {
		t.Errorf("parseTopCPUUsage failed on alternative format: %v", err)
	}

	expected2 := 100.0 - 76.7 // 100 - idle
	if usage2 < expected2-epsilon || usage2 > expected2+epsilon {
		t.Errorf("Expected CPU usage %f, got %f", expected2, usage2)
	}

	// Test invalid format
	_, err = parseTopCPUUsage("invalid output")
	if err == nil {
		t.Error("Expected error for invalid top output")
	}
}

func TestParseFreeCmdOutput(t *testing.T) {
	// Test standard free output format
	output := `              total        used        free      shared  buff/cache   available
Mem:       16777216     8388608     4194304          0     4194304     8388608
Swap:      16777216            0    16777216`

	info, err := parseFreeCmdOutput(output)
	if err != nil {
		t.Errorf("parseFreeCmdOutput failed: %v", err)
	}

	// Validate parsed values
	if info.LimitBytes != 16777216 {
		t.Errorf("Expected total memory 16777216, got %d", info.LimitBytes)
	}
	if info.UsageBytes != 8388608 {
		t.Errorf("Expected used memory 8388608, got %d", info.UsageBytes)
	}
	if info.FreeBytes != 4194304 {
		t.Errorf("Expected free memory 4194304, got %d", info.FreeBytes)
	}
	if info.BufferBytes != 4194304 {
		t.Errorf("Expected buffer memory 4194304, got %d", info.BufferBytes)
	}

	// Test invalid format
	_, err = parseFreeCmdOutput("invalid output")
	if err == nil {
		t.Error("Expected error for invalid free output")
	}
}

func TestGetLoadAverage(t *testing.T) {
	loadAvg, err := getLoadAverage()
	if err != nil {
		t.Logf("getLoadAverage failed (uptime command may not be available): %v", err)
		return
	}

	if loadAvg == "" {
		t.Error("Expected non-empty load average")
	}

	// Load average should contain numbers
	if !strings.Contains(loadAvg, ",") && !strings.Contains(loadAvg, ".") {
		t.Errorf("Expected load average to contain numbers, got %s", loadAvg)
	}

	t.Logf("Load average: %s", loadAvg)
}

func TestGetNumCPUs(t *testing.T) {
	cpus, err := getNumCPUs()
	if err != nil {
		t.Logf("getNumCPUs failed (expected in test environment): %v", err)
		return
	}

	if cpus <= 0 {
		t.Errorf("Expected number of CPUs > 0, got %f", cpus)
	}

	t.Logf("Number of CPUs: %.0f", cpus)
}

func TestGetSystemMemory(t *testing.T) {
	memory, err := getSystemMemory()
	if err != nil {
		t.Logf("getSystemMemory failed (expected in test environment): %v", err)
		return
	}

	if memory <= 0 {
		t.Errorf("Expected system memory > 0, got %d", memory)
	}

	t.Logf("System memory: %d bytes (%.2f GB)", memory, float64(memory)/(1024*1024*1024))
}

func TestCheckConnectivity(t *testing.T) {
	// This is a basic test that requires network access
	report := CheckConnectivity("google.com", "80", 5)

	if report.Domain != "google.com" {
		t.Errorf("Expected domain 'google.com', got '%s'", report.Domain)
	}

	if report.TCP != "success" && !strings.Contains(report.TCP, "refused") {
		// Allow for connection refused in sandboxed environments
		t.Logf("TCP check did not succeed (as expected in some environments): %s", report.TCP)
	}
}

func TestOSDetection(t *testing.T) {
	toolbox := Toolbox{}
	isMac := toolbox.IsMacOS()
	isLin := toolbox.IsLinux()

	if isMac && isLin {
		t.Error("isMacOS and isLinux cannot both be true")
	}

	switch runtime.GOOS {
	case "darwin":
		if !isMac {
			t.Error("Expected isMacOS to be true on darwin")
		}
		if isLin {
			t.Error("Expected isLinux to be false on darwin")
		}
	case "linux":
		if isMac {
			t.Error("Expected isMacOS to be false on linux")
		}
		if !isLin {
			t.Error("Expected isLinux to be true on linux")
		}
	default:
		if isMac || isLin {
			t.Errorf("Expected both isMacOS and isLinux to be false on %s", runtime.GOOS)
		}
	}
	t.Logf("OS detection: GOOS=%s, isMacOS=%v, isLinux=%v", runtime.GOOS, isMac, isLin)
}
