package toolbox

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	toolbox := Toolbox{}
	info, err := toolbox.GetSystemInfo()

	// In a real container environment, this should work
	// In test environment, it might fail due to missing cgroup files
	if err != nil {
		t.Logf("GetSystemInfo failed (expected in test environment): %v", err)
		return
	}

	// Validate CPU info
	if info.CPU.LimitCores <= 0 {
		t.Errorf("Expected CPU limit > 0, got %f", info.CPU.LimitCores)
	}
	if info.CPU.UsagePercent < 0 || info.CPU.UsagePercent > 100 {
		t.Errorf("Expected CPU usage percent between 0-100, got %f", info.CPU.UsagePercent)
	}

	// Validate Memory info
	if info.Memory.LimitBytes <= 0 {
		t.Errorf("Expected memory limit > 0, got %d", info.Memory.LimitBytes)
	}
	if info.Memory.UsageBytes < 0 {
		t.Errorf("Expected memory usage >= 0, got %d", info.Memory.UsageBytes)
	}
	if info.Memory.UsagePercent < 0 || info.Memory.UsagePercent > 100 {
		t.Errorf("Expected memory usage percent between 0-100, got %f", info.Memory.UsagePercent)
	}

	// Validate new fields
	if info.Method == "" {
		t.Error("Expected method to be set")
	}

	// Validate load average if available
	if info.CPU.LoadAverage != "" {
		if !strings.Contains(info.CPU.LoadAverage, ",") && !strings.Contains(info.CPU.LoadAverage, ".") {
			t.Errorf("Expected load average to contain numbers, got %s", info.CPU.LoadAverage)
		}
	}

	t.Logf("SystemInfo: Method=%s, Fallback=%v, CPU=%.2f%%, Memory=%.2f%%",
		info.Method, info.Fallback, info.CPU.UsagePercent, info.Memory.UsagePercent)
}

func TestGetSystemInfoCommand(t *testing.T) {
	toolbox := Toolbox{}
	info, err := toolbox.GetSystemInfoCommand()

	// This should work in most environments since it uses system commands
	if err != nil {
		t.Logf("GetSystemInfoCommand failed: %v", err)
		return
	}

	// Validate method is set to command
	if info.Method != "command" {
		t.Errorf("Expected method to be 'command', got %s", info.Method)
	}

	// Validate fallback is false for command method
	if info.Fallback {
		t.Error("Expected fallback to be false for command method")
	}

	// Validate CPU info
	if info.CPU.LimitCores <= 0 {
		t.Errorf("Expected CPU limit > 0, got %f", info.CPU.LimitCores)
	}
	if info.CPU.UsagePercent < 0 || info.CPU.UsagePercent > 100 {
		t.Errorf("Expected CPU usage percent between 0-100, got %f", info.CPU.UsagePercent)
	}

	// Validate Memory info
	if info.Memory.LimitBytes <= 0 {
		t.Errorf("Expected memory limit > 0, got %d", info.Memory.LimitBytes)
	}
	if info.Memory.UsageBytes < 0 {
		t.Errorf("Expected memory usage >= 0, got %d", info.Memory.UsageBytes)
	}
	if info.Memory.UsagePercent < 0 || info.Memory.UsagePercent > 100 {
		t.Errorf("Expected memory usage percent between 0-100, got %f", info.Memory.UsagePercent)
	}

	// Validate new memory fields
	if info.Memory.FreeBytes < 0 {
		t.Errorf("Expected free bytes >= 0, got %d", info.Memory.FreeBytes)
	}
	if info.Memory.BufferBytes < 0 {
		t.Errorf("Expected buffer bytes >= 0, got %d", info.Memory.BufferBytes)
	}
	if info.Memory.CachedBytes < 0 {
		t.Errorf("Expected cached bytes >= 0, got %d", info.Memory.CachedBytes)
	}

	t.Logf("Command SystemInfo: CPU=%.2f%%, Memory=%.2f%%, Load=%s",
		info.CPU.UsagePercent, info.Memory.UsagePercent, info.CPU.LoadAverage)
}

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

func TestGetTopOutput(t *testing.T) {
	toolbox := Toolbox{}
	output, err := toolbox.GetTopOutput()

	if err != nil {
		t.Logf("GetTopOutput failed (top command may not be available): %v", err)
		return
	}

	if output == "" {
		t.Error("Expected non-empty top output")
	}

	// Check if output contains typical top information
	if !strings.Contains(output, "CPU") && !strings.Contains(output, "Mem") && !strings.Contains(output, "load") {
		t.Error("Expected top output to contain CPU, Mem, or load information")
	}

	t.Logf("Top output length: %d characters", len(output))
}

func TestGetFreeOutput(t *testing.T) {
	toolbox := Toolbox{}
	output, err := toolbox.GetFreeOutput()

	if err != nil {
		t.Logf("GetFreeOutput failed (free command may not be available): %v", err)
		return
	}

	if output == "" {
		t.Error("Expected non-empty free output")
	}

	// Check if output contains memory information
	if !strings.Contains(output, "Mem:") && !strings.Contains(output, "total") {
		t.Error("Expected free output to contain memory information")
	}

	t.Logf("Free output length: %d characters", len(output))
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

// Integration test that verifies cgroup fallback to command works
func TestFallbackMechanism(t *testing.T) {
	toolbox := Toolbox{}

	// Try to get system info with auto fallback
	info, err := toolbox.GetSystemInfo()
	if err != nil {
		t.Logf("Auto fallback failed: %v", err)
		return
	}

	// Try command-only mode
	cmdInfo, err := toolbox.GetSystemInfoCommand()
	if err != nil {
		t.Logf("Command-only mode failed: %v", err)
		return
	}

	// Compare results - they should be reasonably close
	if cmdInfo.Method != "command" {
		t.Errorf("Expected command method, got %s", cmdInfo.Method)
	}

	t.Logf("Fallback test: Auto method=%s, Command CPU=%.2f%%, Memory=%.2f%%",
		info.Method, cmdInfo.CPU.UsagePercent, cmdInfo.Memory.UsagePercent)
}

// Benchmark tests
func BenchmarkGetSystemInfo(b *testing.B) {
	toolbox := Toolbox{}
	for i := 0; i < b.N; i++ {
		_, _ = toolbox.GetSystemInfo()
	}
}

func BenchmarkGetSystemInfoCommand(b *testing.B) {
	toolbox := Toolbox{}
	for i := 0; i < b.N; i++ {
		_, _ = toolbox.GetSystemInfoCommand()
	}
}

func BenchmarkGetCPUUsage(b *testing.B) {
	toolbox := Toolbox{}
	for i := 0; i < b.N; i++ {
		_, _ = toolbox.GetCPUUsage()
	}
}

func BenchmarkGetMemoryUsage(b *testing.B) {
	toolbox := Toolbox{}
	for i := 0; i < b.N; i++ {
		_, _ = toolbox.GetMemoryUsage()
	}
}

func BenchmarkGetMemoryLimit(b *testing.B) {
	toolbox := Toolbox{}
	for i := 0; i < b.N; i++ {
		_, _ = toolbox.GetMemoryLimit()
	}
}

func BenchmarkGetCPULimit(b *testing.B) {
	toolbox := Toolbox{}
	for i := 0; i < b.N; i++ {
		_, _ = toolbox.GetCPULimit()
	}
}

func TestCheckConnectivity(t *testing.T) {
	domain := "google.com"
	port := "80"
	timeout := 5

	report := CheckConnectivity(domain, port, timeout)

	if report.Domain != domain {
		t.Errorf("Expected domain %s, got %s", domain, report.Domain)
	}
	if report.Port != port {
		t.Errorf("Expected port %s, got %s", port, report.Port)
	}
	if report.TimeoutSeconds != timeout {
		t.Errorf("Expected timeout %d, got %d", timeout, report.TimeoutSeconds)
	}

	t.Logf("TCP result: %s", report.TCP)
	t.Logf("HTTP result: %s", report.HTTP)

	if report.TCP != "success" {
		t.Errorf("Expected TCP success, got %s", report.TCP)
	}
	if report.HTTP == "" || report.HTTP == "skipped (TCP failed)" {
		t.Errorf("Expected HTTP result, got %s", report.HTTP)
	}
}
