# K6 Toolbox Extension

A comprehensive system monitoring extension for k6 that provides CPU and memory metrics in containerized environments, specifically designed for EKS/Docker deployments.

## Features

- 🔄 **Automatic Fallback**: Tries cgroup first, falls back to system commands
- 📊 **Comprehensive Metrics**: CPU usage, memory usage, limits, and availability
- 🐳 **Container-Aware**: Reads actual container limits from cgroup v1/v2
- 🛡️ **Alpine Compatible**: Works with BusyBox commands in Alpine containers
- ⚡ **Performance**: Optimized for minimal overhead
- 📈 **Real-time**: Kernel-level accuracy for container resource monitoring

## Installation

Build the k6 binary with the toolbox extension:

```bash
# Build k6 with toolbox extension
xk6 build --with github.com/your-org/k6-toolbox-extension@latest
```

## Quick Start

```javascript
import toolbox from 'k6/x/toolbox';

export default function() {
    // Get complete system information
    const info = toolbox.getSystemInfo();
    console.log(`CPU: ${info.cpu.usage_percent}%, Memory: ${info.memory.usage_percent}%`);
    
    // Check if resources are available
    if (info.memory.usage_percent > 80) {
        console.warn(`High memory usage: ${info.memory.usage_mb}MB/${info.memory.limit_mb}MB`);
    }
}
```

## API Reference

### System Overview

#### `getSystemInfo()`
Returns comprehensive system information with automatic fallback.

```javascript
const info = toolbox.getSystemInfo();
// Returns: SystemInfo object with CPU, Memory, Method, Fallback fields
```

#### `getSystemInfoCommand()`
Forces command-based monitoring (no cgroup files).

```javascript
const info = toolbox.getSystemInfoCommand();
// Returns: SystemInfo object using only system commands
```

### CPU Metrics

| Method | Return Type | Description |
|--------|-------------|-------------|
| `getCPUUsage()` | `float64` | Current CPU usage percentage (0-100) |
| `getCPULimit()` | `float64` | CPU limit in cores |
| `getAvailableCPU()` | `float64` | Available CPU cores (limit - usage) |

### Memory Metrics

| Method | Return Type | Description |
|--------|-------------|-------------|
| `getMemoryUsage()` | `int64` | Current memory usage in bytes |
| `getMemoryLimit()` | `int64` | Memory limit in bytes |
| `getMemoryUsagePercent()` | `float64` | Memory usage percentage (0-100) |
| `getAvailableMemory()` | `int64` | Available memory in bytes |

### Raw Command Output

| Method | Return Type | Description |
|--------|-------------|-------------|
| `getTopOutput()` | `string` | Raw `top -b -n 1` output |
| `getFreeOutput()` | `string` | Raw `free -b` output |
| `getPsOutput()` | `string` | Raw `ps aux` output |
| `getUptimeOutput()` | `string` | Raw `uptime` output |

### Connectivity Check

| Method | Return Type | Description |
|--------|-------------|-------------|
| `checkConnectivity(domain, port, timeout)` | `ConnectivityReport` | Checks TCP and HTTP connectivity to the given domain and port, with a configurable timeout (seconds, default 5). |

#### Example

```javascript
import toolbox from 'k6/x/toolbox';

export default function() {
    const report = toolbox.checkConnectivity('google.com', '80', 5);
    console.log('Connectivity Report:', JSON.stringify(report, null, 2));
    // Example output:
    // {
    //   domain: 'google.com',
    //   port: '80',
    //   timeout_seconds: 5,
    //   tcp: 'success',
    //   http: '200 OK'
    // }
}
```

#### ConnectivityReport Structure

```javascript
{
  domain: 'string',           // The domain checked
  port: 'string',             // The port checked
  timeout_seconds: number,    // Timeout used for each check
  tcp: 'success' | string,    // 'success' or error message
  http: 'success' | string    // HTTP status or error/skipped message
}
```

## Data Structures

### SystemInfo

```javascript
{
    cpu: {
        usage_percent: 45.2,      // Current CPU usage %
        limit_cores: 4.0,         // CPU limit in cores
        used_cores: 1.8,          // Currently used cores
        available_cores: 2.2,     // Available cores
        load_average: "1.2, 1.5, 1.8"  // System load average
    },
    memory: {
        usage_bytes: 1073741824,   // Memory usage in bytes
        limit_bytes: 2147483648,   // Memory limit in bytes
        available_bytes: 1073741824,  // Available memory
        usage_percent: 50.0,       // Memory usage %
        usage_mb: 1024.0,         // Usage in MB
        limit_mb: 2048.0,         // Limit in MB
        available_mb: 1024.0,     // Available in MB
        free_bytes: 536870912,    // Free memory (from free command)
        buffer_bytes: 268435456,  // Buffer memory
        cached_bytes: 268435456   // Cached memory
    },
    method: "cgroup",             // "cgroup", "command", or "mixed"
    fallback: false               // Whether fallback was used
}
```

## Use Cases

### Resource Monitoring During Load Tests

```javascript
import toolbox from 'k6/x/toolbox';

export default function() {
    // Monitor resources every 10 iterations
    if (__ITER % 10 === 0) {
        const cpu = toolbox.getCPUUsage();
        const memory = toolbox.getMemoryUsagePercent();
        
        console.log(`Iteration ${__ITER}: CPU=${cpu.toFixed(1)}%, Memory=${memory.toFixed(1)}%`);
        
        // Alert if resources are constrained
        if (cpu > 90 || memory > 85) {
            console.error(`⚠️  Resource constraint detected!`);
        }
    }
    
    // Your test logic here
    // ...
}
```

### Memory Usage Tracking

```javascript
import toolbox from 'k6/x/toolbox';

export function setup() {
    const info = toolbox.getSystemInfo();
    console.log(`Test starting with ${info.memory.available_mb.toFixed(0)}MB available memory`);
    return { initialMemory: info.memory.usage_bytes };
}

export default function(data) {
    // Track memory growth
    const currentMemory = toolbox.getMemoryUsage();
    const memoryGrowth = currentMemory - data.initialMemory;
    
    if (memoryGrowth > 100 * 1024 * 1024) { // 100MB growth
        console.warn(`Memory increased by ${(memoryGrowth / 1024 / 1024).toFixed(1)}MB`);
    }
}
```

### Container Limits Validation

```javascript
import toolbox from 'k6/x/toolbox';

export function setup() {
    const info = toolbox.getSystemInfo();
    
    // Validate container has sufficient resources
    const minCPU = 2.0;
    const minMemoryGB = 4.0;
    
    if (info.cpu.limit_cores < minCPU) {
        throw new Error(`Insufficient CPU: ${info.cpu.limit_cores} cores (need ${minCPU})`);
    }
    
    if (info.memory.limit_mb < minMemoryGB * 1024) {
        throw new Error(`Insufficient memory: ${info.memory.limit_mb}MB (need ${minMemoryGB * 1024}MB)`);
    }
    
    console.log(`✅ Resource validation passed: ${info.cpu.limit_cores} cores, ${info.memory.limit_mb}MB`);
    console.log(`📊 Detection method: ${info.method} ${info.fallback ? '(with fallback)' : ''}`);
}
```

### Performance Regression Detection

```javascript
import toolbox from 'k6/x/toolbox';

let previousCPU = 0;
let cpuSpikes = 0;

export default function() {
    const currentCPU = toolbox.getCPUUsage();
    
    // Detect CPU spikes
    if (currentCPU > previousCPU + 20) {
        cpuSpikes++;
        console.warn(`CPU spike detected: ${previousCPU.toFixed(1)}% → ${currentCPU.toFixed(1)}%`);
    }
    
    previousCPU = currentCPU;
}

export function teardown() {
    const finalInfo = toolbox.getSystemInfo();
    console.log(`Test completed. CPU spikes: ${cpuSpikes}`);
    console.log(`Final resource usage: CPU=${finalInfo.cpu.usage_percent.toFixed(1)}%, Memory=${finalInfo.memory.usage_percent.toFixed(1)}%`);
}
```

## Environment Compatibility

### Container Environments ✅
- **EKS/Kubernetes**: Full cgroup v1/v2 support
- **Docker**: Container-aware limits and usage
- **Alpine Linux**: BusyBox command compatibility
- **Ubuntu/Debian**: Full feature support

### Fallback Chain
1. **Primary**: cgroup v2 files (`/sys/fs/cgroup/memory.current`, etc.)
2. **Secondary**: cgroup v1 files (`/sys/fs/cgroup/memory/memory.usage_in_bytes`, etc.)
3. **Fallback**: System commands (`top`, `free`, `nproc`, `uptime`)

### Required Permissions
- ✅ Standard container permissions (no root required)
- ✅ Read access to `/proc/` and `/sys/fs/cgroup/`
- ✅ Basic command execution (`top`, `free`, `ps`, `uptime`)

## Error Handling

The extension gracefully handles missing commands or files:

```javascript
// Always check for errors in production
const usage = toolbox.getCPUUsage();
if (usage === undefined) {
    console.log("CPU monitoring not available in this environment");
}

// Or use the comprehensive method with built-in fallback
const info = toolbox.getSystemInfo();
console.log(`Monitoring method: ${info.method}`);
if (info.fallback) {
    console.log("Using fallback monitoring - some metrics may be less accurate");
}
```

## Troubleshooting

### Common Issues

**"Object has no member 'getMemoryLimit'"**
- Ensure you're using the latest version with all public methods
- Rebuild the k6 binary with the updated extension

**"failed to read cgroup files"**
- Normal in non-containerized environments
- Extension will automatically fall back to system commands

**"command not found: nproc"**
- Common in Alpine/BusyBox environments
- Extension falls back to parsing `/proc/cpuinfo`

### Debug Information

```javascript
// Get detailed information about monitoring method
const info = toolbox.getSystemInfo();
console.log(`Method: ${info.method}, Fallback: ${info.fallback}`);

// Get raw command output for debugging
if (info.fallback) {
    console.log("Raw top output:", toolbox.getTopOutput());
    console.log("Raw free output:", toolbox.getFreeOutput());
}
```

## Performance Notes

- **Minimal Overhead**: Designed for production load testing
- **File I/O**: ~1-2ms per metric read from cgroup files
- **Command Execution**: ~10-50ms per system command (fallback only)
- **Memory Impact**: Negligible - reads system metrics, doesn't store data

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.