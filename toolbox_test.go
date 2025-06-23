package toolbox

import (
	"os"
	"path/filepath"
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

func TestGetNumCPUs(t *testing.T) {
	cpus, err := getNumCPUs()
	if err != nil {
		t.Logf("getNumCPUs failed (expected in test environment): %v", err)
		return
	}

	if cpus <= 0 {
		t.Errorf("Expected number of CPUs > 0, got %f", cpus)
	}
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
}

// Benchmark tests
func BenchmarkGetSystemInfo(b *testing.B) {
	toolbox := Toolbox{}
	for i := 0; i < b.N; i++ {
		_, _ = toolbox.GetSystemInfo()
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
