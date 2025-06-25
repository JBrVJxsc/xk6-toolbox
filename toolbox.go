package toolbox

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.k6.io/k6/js/modules"
)

// Error messages
const (
	ErrReadingFile     = "failed to read file"
	ErrParsingValue    = "failed to parse value"
	ErrCgroupNotFound  = "cgroup information not found"
	ErrMemoryNotFound  = "memory information not found"
	ErrCPUNotFound     = "CPU information not found"
	ErrInvalidCgroupV  = "unsupported cgroup version"
	ErrCommandFailed   = "command execution failed"
	ErrCommandNotFound = "command not found"
)

// SystemInfo represents the current system resource information
type SystemInfo struct {
	CPU      CPUInfo    `json:"cpu"`
	Memory   MemoryInfo `json:"memory"`
	Method   string     `json:"method"`   // How the data was collected
	Fallback bool       `json:"fallback"` // Whether fallback methods were used
}

// CPUInfo contains CPU usage and limit information
type CPUInfo struct {
	UsagePercent float64 `json:"usage_percent"`
	LimitCores   float64 `json:"limit_cores"`
	UsedCores    float64 `json:"used_cores"`
	Available    float64 `json:"available_cores"`
	LoadAverage  string  `json:"load_average"`
}

// MemoryInfo contains memory usage and limit information
type MemoryInfo struct {
	UsageBytes     int64   `json:"usage_bytes"`
	LimitBytes     int64   `json:"limit_bytes"`
	AvailableBytes int64   `json:"available_bytes"`
	UsagePercent   float64 `json:"usage_percent"`
	UsageMB        float64 `json:"usage_mb"`
	LimitMB        float64 `json:"limit_mb"`
	AvailableMB    float64 `json:"available_mb"`
	FreeBytes      int64   `json:"free_bytes"`
	BufferBytes    int64   `json:"buffer_bytes"`
	CachedBytes    int64   `json:"cached_bytes"`
}

// ConnectivityReport represents the result of connectivity checks at different layers
type ConnectivityReport struct {
	Domain         string `json:"domain"`
	Port           string `json:"port"`
	TimeoutSeconds int    `json:"timeout_seconds"`
	TCP            string `json:"tcp"`  // e.g. "success" or error message
	HTTP           string `json:"http"` // e.g. "success" or error message
}

func init() {
	modules.Register("k6/x/toolbox", new(Toolbox))
}

// Toolbox is the main module exposed to k6 JavaScript.
// It provides functions for monitoring system resources in containerized environments.
type Toolbox struct{}

// GetSystemInfo returns comprehensive system resource information
// Uses multiple methods with automatic fallback
func (Toolbox) GetSystemInfo() (SystemInfo, error) {
	var info SystemInfo
	var err error

	// Try cgroup method first
	info.CPU, err = getCPUInfoCgroup()
	if err != nil {
		// Fall back to command-based method
		info.CPU, err = getCPUInfoCommand()
		if err != nil {
			return info, fmt.Errorf("failed to get CPU info: %w", err)
		}
		info.Fallback = true
		info.Method = "command"
	} else {
		info.Method = "cgroup"
	}

	info.Memory, err = getMemoryInfoCgroup()
	if err != nil {
		// Fall back to command-based method
		info.Memory, err = getMemoryInfoCommand()
		if err != nil {
			return info, fmt.Errorf("failed to get memory info: %w", err)
		}
		info.Fallback = true
		if info.Method == "cgroup" {
			info.Method = "mixed"
		} else {
			info.Method = "command"
		}
	}

	return info, nil
}

// GetSystemInfoCommand forces using command-based monitoring only
func (Toolbox) GetSystemInfoCommand() (SystemInfo, error) {
	var info SystemInfo
	var err error

	info.CPU, err = getCPUInfoCommand()
	if err != nil {
		return info, fmt.Errorf("failed to get CPU info via commands: %w", err)
	}

	info.Memory, err = getMemoryInfoCommand()
	if err != nil {
		return info, fmt.Errorf("failed to get memory info via commands: %w", err)
	}

	info.Method = "command"
	info.Fallback = false
	return info, nil
}

// GetTopOutput returns raw output from the `top` command
func (Toolbox) GetTopOutput() (string, error) {
	cmd := exec.Command("top", "-b", "-n", "1")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrCommandFailed, err)
	}
	return string(output), nil
}

// GetFreeOutput returns raw output from the `free` command
func (Toolbox) GetFreeOutput() (string, error) {
	cmd := exec.Command("free", "-b") // Output in bytes
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrCommandFailed, err)
	}
	return string(output), nil
}

// GetPsOutput returns raw output from the `ps` command
func (Toolbox) GetPsOutput() (string, error) {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrCommandFailed, err)
	}
	return string(output), nil
}

// GetUptimeOutput returns raw output from the `uptime` command
func (Toolbox) GetUptimeOutput() (string, error) {
	cmd := exec.Command("uptime")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrCommandFailed, err)
	}
	return string(output), nil
}

// GetCPUUsage returns current CPU usage percentage
func (Toolbox) GetCPUUsage() (float64, error) {
	cpuInfo, err := getCPUInfoCgroup()
	if err != nil {
		cpuInfo, err = getCPUInfoCommand()
		if err != nil {
			return 0, err
		}
	}
	return cpuInfo.UsagePercent, nil
}

// GetCPULimit returns the CPU limit in cores
func (Toolbox) GetCPULimit() (float64, error) {
	return getCPULimit()
}

// GetMemoryUsage returns current memory usage in bytes
func (Toolbox) GetMemoryUsage() (int64, error) {
	memInfo, err := getMemoryInfoCgroup()
	if err != nil {
		memInfo, err = getMemoryInfoCommand()
		if err != nil {
			return 0, err
		}
	}
	return memInfo.UsageBytes, nil
}

// GetMemoryLimit returns the memory limit in bytes
func (Toolbox) GetMemoryLimit() (int64, error) {
	return getMemoryLimit()
}

// GetMemoryUsagePercent returns memory usage as a percentage
func (Toolbox) GetMemoryUsagePercent() (float64, error) {
	memInfo, err := getMemoryInfoCgroup()
	if err != nil {
		memInfo, err = getMemoryInfoCommand()
		if err != nil {
			return 0, err
		}
	}
	return memInfo.UsagePercent, nil
}

// GetAvailableMemory returns available memory in bytes
func (Toolbox) GetAvailableMemory() (int64, error) {
	memInfo, err := getMemoryInfoCgroup()
	if err != nil {
		memInfo, err = getMemoryInfoCommand()
		if err != nil {
			return 0, err
		}
	}
	return memInfo.AvailableBytes, nil
}

// GetAvailableCPU returns available CPU cores
func (Toolbox) GetAvailableCPU() (float64, error) {
	cpuInfo, err := getCPUInfoCgroup()
	if err != nil {
		cpuInfo, err = getCPUInfoCommand()
		if err != nil {
			return 0, err
		}
	}
	return cpuInfo.Available, nil
}

// Command-based implementations

// getCPUInfoCommand gets CPU info using system commands
func getCPUInfoCommand() (CPUInfo, error) {
	var info CPUInfo

	// Get CPU count
	cores, err := getCPUCoresCommand()
	if err != nil {
		return info, err
	}
	info.LimitCores = cores

	// Get CPU usage from top
	usage, err := getCPUUsageFromTop()
	if err != nil {
		return info, err
	}
	info.UsagePercent = usage
	info.UsedCores = (usage / 100.0) * cores
	info.Available = cores - info.UsedCores

	// Get load average
	loadAvg, err := getLoadAverage()
	if err == nil {
		info.LoadAverage = loadAvg
	}

	return info, nil
}

// getMemoryInfoCommand gets memory info using system commands
func getMemoryInfoCommand() (MemoryInfo, error) {
	var info MemoryInfo

	// Use free command to get memory info
	output, err := exec.Command("free", "-b").Output()
	if err != nil {
		return info, fmt.Errorf("%s: %w", ErrCommandFailed, err)
	}

	return parseFreeCmdOutput(string(output))
}

// getCPUCoresCommand gets number of CPU cores
func getCPUCoresCommand() (float64, error) {
	output, err := exec.Command("nproc").Output()
	if err != nil {
		// Fallback to parsing /proc/cpuinfo
		return getCPUCoresFromProcInfo()
	}

	cores, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrParsingValue, err)
	}

	return cores, nil
}

// getCPUUsageFromTop parses CPU usage from top command
func getCPUUsageFromTop() (float64, error) {
	cmd := exec.Command("top", "-b", "-n", "1")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrCommandFailed, err)
	}

	return parseTopCPUUsage(string(output))
}

// parseTopCPUUsage extracts CPU usage from top output
func parseTopCPUUsage(output string) (float64, error) {
	lines := strings.Split(output, "\n")

	// Look for CPU usage line (varies by top version)
	cpuRegex := regexp.MustCompile(`%Cpu\(s\):\s*([0-9.]+)\s*us,\s*([0-9.]+)\s*sy,.*?([0-9.]+)\s*id`)

	for _, line := range lines {
		if strings.Contains(line, "Cpu(s)") || strings.Contains(line, "%Cpu") {
			matches := cpuRegex.FindStringSubmatch(line)
			if len(matches) >= 4 {
				idle, err := strconv.ParseFloat(matches[3], 64)
				if err != nil {
					continue
				}
				return 100 - idle, nil
			}
		}

		// Alternative parsing for different top formats
		if strings.Contains(line, "CPU usage:") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "idle" && i > 0 {
					idleStr := strings.TrimSuffix(parts[i-1], "%")
					idle, err := strconv.ParseFloat(idleStr, 64)
					if err == nil {
						return 100 - idle, nil
					}
				}
			}
		}
	}

	return 0, errors.New("could not parse CPU usage from top output")
}

// parseFreeCmdOutput parses the output of the free command
func parseFreeCmdOutput(output string) (MemoryInfo, error) {
	var info MemoryInfo

	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return info, errors.New("invalid free command output")
	}

	// Parse the "Mem:" line
	for _, line := range lines {
		if strings.HasPrefix(line, "Mem:") {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				return info, errors.New("invalid memory line format")
			}

			total, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return info, fmt.Errorf("failed to parse total memory: %w", err)
			}

			used, err := strconv.ParseInt(fields[2], 10, 64)
			if err != nil {
				return info, fmt.Errorf("failed to parse used memory: %w", err)
			}

			free, err := strconv.ParseInt(fields[3], 10, 64)
			if err != nil {
				return info, fmt.Errorf("failed to parse free memory: %w", err)
			}

			info.LimitBytes = total
			info.UsageBytes = used
			info.FreeBytes = free
			info.AvailableBytes = free

			var buffers, cached int64

			// Parse buffers and cached if available
			if len(fields) >= 6 {
				if buf, err := strconv.ParseInt(fields[5], 10, 64); err == nil {
					buffers = buf
					info.BufferBytes = buffers
				}
			}
			if len(fields) >= 7 {
				if cach, err := strconv.ParseInt(fields[6], 10, 64); err == nil {
					cached = cach
					info.CachedBytes = cached
				}
			}

			// Available memory includes buffers and cache
			info.AvailableBytes = free + buffers + cached

			info.UsagePercent = (float64(used) / float64(total)) * 100
			info.UsageMB = float64(used) / (1024 * 1024)
			info.LimitMB = float64(total) / (1024 * 1024)
			info.AvailableMB = float64(info.AvailableBytes) / (1024 * 1024)

			return info, nil
		}
	}

	return info, errors.New("memory information not found in free output")
}

// getLoadAverage gets system load average
func getLoadAverage() (string, error) {
	output, err := exec.Command("uptime").Output()
	if err != nil {
		return "", err
	}

	// Extract load average from uptime output
	uptimeStr := string(output)
	loadIdx := strings.Index(uptimeStr, "load average:")
	if loadIdx == -1 {
		return "", errors.New("load average not found")
	}

	return strings.TrimSpace(uptimeStr[loadIdx+13:]), nil
}

// getCPUCoresFromProcInfo gets CPU cores from /proc/cpuinfo
func getCPUCoresFromProcInfo() (float64, error) {
	content, err := readFile("/proc/cpuinfo")
	if err != nil {
		return 0, err
	}

	count := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "processor") {
			count++
		}
	}

	if count == 0 {
		return 0, errors.New("no processors found in /proc/cpuinfo")
	}

	return float64(count), nil
}

// Original cgroup-based implementations (keeping for primary method)

// getCPUInfoCgroup retrieves CPU usage and limit information from cgroup
func getCPUInfoCgroup() (CPUInfo, error) {
	var info CPUInfo

	// Get CPU limit from cgroup
	limit, err := getCPULimit()
	if err != nil {
		return info, err
	}
	info.LimitCores = limit

	// Get CPU usage
	usage, err := getCPUUsage()
	if err != nil {
		return info, err
	}
	info.UsedCores = usage
	info.UsagePercent = (usage / limit) * 100
	info.Available = limit - usage

	return info, nil
}

// getMemoryInfoCgroup retrieves memory usage and limit information from cgroup
func getMemoryInfoCgroup() (MemoryInfo, error) {
	var info MemoryInfo

	// Get memory limit from cgroup
	limit, err := getMemoryLimit()
	if err != nil {
		return info, err
	}
	info.LimitBytes = limit

	// Get memory usage from cgroup
	usage, err := getMemoryUsage()
	if err != nil {
		return info, err
	}
	info.UsageBytes = usage
	info.AvailableBytes = limit - usage
	info.UsagePercent = (float64(usage) / float64(limit)) * 100

	// Convert to MB for convenience
	info.UsageMB = float64(usage) / (1024 * 1024)
	info.LimitMB = float64(limit) / (1024 * 1024)
	info.AvailableMB = float64(info.AvailableBytes) / (1024 * 1024)

	return info, nil
}

// getCPULimit reads CPU limit from cgroup
func getCPULimit() (float64, error) {
	// Try cgroup v2 first
	if limit, err := readCgroupV2CPULimit(); err == nil {
		return limit, nil
	}

	// Fall back to cgroup v1
	return readCgroupV1CPULimit()
}

// getCPUUsage calculates current CPU usage
func getCPUUsage() (float64, error) {
	// Read CPU usage from cgroup
	if usage, err := readCgroupCPUUsage(); err == nil {
		return usage, nil
	}

	// Fall back to /proc/stat method
	return readProcStatCPUUsage()
}

// getMemoryLimit reads memory limit from cgroup
func getMemoryLimit() (int64, error) {
	// Try cgroup v2 first
	if limit, err := readCgroupV2MemoryLimit(); err == nil {
		return limit, nil
	}

	// Fall back to cgroup v1
	return readCgroupV1MemoryLimit()
}

// getMemoryUsage reads current memory usage from cgroup
func getMemoryUsage() (int64, error) {
	// Try cgroup v2 first
	if usage, err := readCgroupV2MemoryUsage(); err == nil {
		return usage, nil
	}

	// Fall back to cgroup v1
	return readCgroupV1MemoryUsage()
}

// readCgroupV2CPULimit reads CPU limit from cgroup v2
func readCgroupV2CPULimit() (float64, error) {
	content, err := readFile("/sys/fs/cgroup/cpu.max")
	if err != nil {
		return 0, err
	}

	parts := strings.Fields(strings.TrimSpace(content))
	if len(parts) != 2 {
		return 0, errors.New("invalid cpu.max format")
	}

	if parts[0] == "max" {
		// No CPU limit set, use number of CPUs
		return getNumCPUs()
	}

	quota, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}

	period, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, err
	}

	return quota / period, nil
}

// readCgroupV1CPULimit reads CPU limit from cgroup v1
func readCgroupV1CPULimit() (float64, error) {
	quotaContent, err := readFile("/sys/fs/cgroup/cpu,cpuacct/cpu.cfs_quota_us")
	if err != nil {
		return 0, err
	}

	periodContent, err := readFile("/sys/fs/cgroup/cpu,cpuacct/cpu.cfs_period_us")
	if err != nil {
		return 0, err
	}

	quota, err := strconv.ParseFloat(strings.TrimSpace(quotaContent), 64)
	if err != nil {
		return 0, err
	}

	if quota == -1 {
		// No CPU limit set, use number of CPUs
		return getNumCPUs()
	}

	period, err := strconv.ParseFloat(strings.TrimSpace(periodContent), 64)
	if err != nil {
		return 0, err
	}

	return quota / period, nil
}

// readCgroupCPUUsage reads CPU usage from cgroup
func readCgroupCPUUsage() (float64, error) {
	// This is a simplified implementation
	// In practice, we'd need to calculate usage over time
	content, err := readFile("/sys/fs/cgroup/cpuacct/cpuacct.usage")
	if err != nil {
		// Try cgroup v2
		content, err = readFile("/sys/fs/cgroup/cpu.stat")
		if err != nil {
			return 0, err
		}
		return parseCgroupV2CPUUsage(content)
	}

	nanoseconds, err := strconv.ParseFloat(strings.TrimSpace(content), 64)
	if err != nil {
		return 0, err
	}

	// Convert to cores (this is cumulative, so we'd need to track over time)
	// For now, return a placeholder that would need proper implementation
	return nanoseconds / 1e9 / 100, nil // Rough approximation
}

// readProcStatCPUUsage reads CPU usage from /proc/stat
func readProcStatCPUUsage() (float64, error) {
	content, err := readFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return 0, errors.New("empty /proc/stat")
	}

	cpuLine := lines[0]
	if !strings.HasPrefix(cpuLine, "cpu ") {
		return 0, errors.New("invalid /proc/stat format")
	}

	fields := strings.Fields(cpuLine)
	if len(fields) < 8 {
		return 0, errors.New("insufficient CPU fields in /proc/stat")
	}

	// This is a simplified calculation - proper implementation would track over time
	user, _ := strconv.ParseFloat(fields[1], 64)
	system, _ := strconv.ParseFloat(fields[3], 64)
	idle, _ := strconv.ParseFloat(fields[4], 64)

	total := user + system + idle
	used := total - idle

	numCPUs, err := getNumCPUs()
	if err != nil {
		return 0, err
	}

	return (used / total) * numCPUs, nil
}

// readCgroupV2MemoryLimit reads memory limit from cgroup v2
func readCgroupV2MemoryLimit() (int64, error) {
	content, err := readFile("/sys/fs/cgroup/memory.max")
	if err != nil {
		return 0, err
	}

	limitStr := strings.TrimSpace(content)
	if limitStr == "max" {
		// No memory limit, read from /proc/meminfo
		return getSystemMemory()
	}

	return strconv.ParseInt(limitStr, 10, 64)
}

// readCgroupV1MemoryLimit reads memory limit from cgroup v1
func readCgroupV1MemoryLimit() (int64, error) {
	content, err := readFile("/sys/fs/cgroup/memory/memory.limit_in_bytes")
	if err != nil {
		return 0, err
	}

	limit, err := strconv.ParseInt(strings.TrimSpace(content), 10, 64)
	if err != nil {
		return 0, err
	}

	// Check if limit is set to a very large value (indicating no limit)
	if limit > 9223372036854775807/2 { // Very large number indicating no limit
		return getSystemMemory()
	}

	return limit, nil
}

// readCgroupV2MemoryUsage reads memory usage from cgroup v2
func readCgroupV2MemoryUsage() (int64, error) {
	content, err := readFile("/sys/fs/cgroup/memory.current")
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(strings.TrimSpace(content), 10, 64)
}

// readCgroupV1MemoryUsage reads memory usage from cgroup v1
func readCgroupV1MemoryUsage() (int64, error) {
	content, err := readFile("/sys/fs/cgroup/memory/memory.usage_in_bytes")
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(strings.TrimSpace(content), 10, 64)
}

// parseCgroupV2CPUUsage parses CPU usage from cgroup v2 cpu.stat
func parseCgroupV2CPUUsage(content string) (float64, error) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "usage_usec ") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				microseconds, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					return 0, err
				}
				// Convert to cores (simplified calculation)
				return microseconds / 1e6 / 100, nil
			}
		}
	}
	return 0, errors.New("usage_usec not found in cpu.stat")
}

// getNumCPUs returns the number of CPUs available to the container
func getNumCPUs() (float64, error) {
	content, err := readFile("/proc/cpuinfo")
	if err != nil {
		return 0, err
	}

	count := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "processor") {
			count++
		}
	}

	if count == 0 {
		return 0, errors.New("no processors found in /proc/cpuinfo")
	}

	return float64(count), nil
}

// getSystemMemory returns total system memory from /proc/meminfo
func getSystemMemory() (int64, error) {
	content, err := readFile("/proc/meminfo")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				value, err := strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					return 0, err
				}
				// Convert from KB to bytes
				return value * 1024, nil
			}
		}
	}

	return 0, errors.New("MemTotal not found in /proc/meminfo")
}

// readFile reads the contents of a file
func readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrReadingFile, err)
	}
	return string(content), nil
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// CheckConnectivity checks connectivity to a domain at multiple layers (TCP, HTTP)
// timeoutSeconds: timeout for each check in seconds (default 5 if <=0)
// port: port to check (default "80" if empty)
func CheckConnectivity(domain, port string, timeoutSeconds int) ConnectivityReport {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 5
	}
	if port == "" {
		port = "80"
	}
	address := net.JoinHostPort(domain, port)
	report := ConnectivityReport{
		Domain:         domain,
		Port:           port,
		TimeoutSeconds: timeoutSeconds,
	}

	// TCP check
	dialer := net.Dialer{Timeout: time.Duration(timeoutSeconds) * time.Second}
	tcpConn, err := dialer.Dial("tcp", address)
	if err != nil {
		report.TCP = err.Error()
	} else {
		report.TCP = "success"
		tcpConn.Close()
	}

	// HTTP check (only if TCP succeeded)
	if report.TCP == "success" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
		defer cancel()
		url := "http://" + address
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			report.HTTP = err.Error()
		} else {
			client := &http.Client{
				Timeout: time.Duration(timeoutSeconds) * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				report.HTTP = err.Error()
			} else {
				report.HTTP = resp.Status
				resp.Body.Close()
			}
		}
	} else {
		report.HTTP = "skipped (TCP failed)"
	}

	return report
}

// CheckConnectivity exposes CheckConnectivity to k6 JavaScript
func (Toolbox) CheckConnectivity(domain string, port string, timeoutSeconds int) ConnectivityReport {
	return CheckConnectivity(domain, port, timeoutSeconds)
}
