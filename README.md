# K6 Toolbox Extension

A comprehensive system monitoring extension for k6 that provides CPU and memory metrics in containerized environments (like Docker and EKS) and on macOS.

## Features

- 📊 **Comprehensive Metrics**: CPU usage, memory usage, limits, and availability.
- 🐳 **Container-Aware**: Reads actual container limits from cgroup v1/v2 on Linux.
- 🍏 **macOS Support**: Uses native system commands for resource metrics on macOS.
- 🛡️ **Alpine Compatible**: Works with BusyBox commands in Alpine containers.
- ⚡ **Performance**: Optimized for minimal overhead.
- 📈 **Real-time**: Kernel-level accuracy for container resource monitoring.

## Installation

Build the k6 binary with the extension:

```bash
# Build k6 with the extension
xk6 build --with github.com/JBrVJxsc/xk6-toolbox@latest
```

## Quick Start

```javascript
import toolbox from 'k6/x/toolbox';
import { sleep } from 'k6';

export default function() {
    // Get individual CPU and memory usage stats
    const cpuUsage = toolbox.getCPUUsage();
    const memoryUsagePercent = toolbox.getMemoryUsagePercent();

    console.log(`CPU: ${cpuUsage.toFixed(1)}%, Memory: ${memoryUsagePercent.toFixed(1)}%`);

    // Alert if memory usage is high
    if (memoryUsagePercent > 80) {
        const memoryUsageBytes = toolbox.getMemoryUsage();
        const memoryLimitBytes = toolbox.getMemoryLimit();
        console.warn(`High memory usage: ${(memoryUsageBytes / 1024 / 1024).toFixed(1)}MB / ${(memoryLimitBytes / 1024 / 1024).toFixed(1)}MB`);
    }

    sleep(1);
}
```

## API Reference

### CPU Metrics

| Method | Return Type | Description |
|--------|-------------|-------------|
| `getCPUUsage()` | `float64` | Current CPU usage percentage (0-100). |
| `getCPULimit()` | `float64` | CPU limit in cores. |
| `getAvailableCPU()` | `float64` | Available CPU cores (limit - usage). |

### Memory Metrics

| Method | Return Type | Description |
|--------|-------------|-------------|
| `getMemoryUsage()` | `int64` | Current memory usage in bytes. |
| `getMemoryLimit()` | `int64` | Memory limit in bytes. |
| `getMemoryUsagePercent()` | `float64` | Memory usage percentage (0-100). |
| `getAvailableMemory()` | `int64` | Available memory in bytes. |

### Raw Command Output

| Method | Return Type | Description |
|--------|-------------|-------------|
| `getPsOutput()` | `string` | Raw `ps aux` output. |
| `getUptimeOutput()` | `string` | Raw `uptime` output. |

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
    //   "domain": "google.com",
    //   "port": "80",
    //   "timeout_seconds": 5,
    //   "tcp": "success",
    //   "http": "200 OK"
    // }
}
```

#### ConnectivityReport Structure

```javascript
{
  "domain": "string",           // The domain checked
  "port": "string",             // The port checked
  "timeout_seconds": number,    // Timeout used for each check
  "tcp": "string",              // 'success' or error message
  "http": "string"              // HTTP status or error/skipped message
}
```

## Use Cases

### Resource Monitoring During Load Tests

```javascript
import toolbox from 'k6/x/toolbox';
import { sleep } from 'k6';

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
    sleep(1);
}
```

### Memory Usage Tracking

```javascript
import toolbox from 'k6/x/toolbox';
import { sleep } from 'k6';

export function setup() {
    const initialMemory = toolbox.getMemoryUsage();
    console.log(`Test starting with ${(initialMemory / 1024 / 1024).toFixed(0)}MB used memory`);
    return { initialMemory };
}

export default function(data) {
    // Track memory growth
    const currentMemory = toolbox.getMemoryUsage();
    const memoryGrowth = currentMemory - data.initialMemory;

    if (memoryGrowth > 100 * 1024 * 1024) { // 100MB growth
        console.warn(`Memory increased by ${(memoryGrowth / 1024 / 1024).toFixed(1)}MB`);
    }
    sleep(1);
}
```

### Container Limits Validation

```javascript
import toolbox from 'k6/x/toolbox';

export function setup() {
    const cpuLimit = toolbox.getCPULimit();
    const memoryLimit = toolbox.getMemoryLimit();

    // Validate container has sufficient resources
    const minCPU = 2.0;
    const minMemoryGB = 4.0;

    if (cpuLimit < minCPU) {
        throw new Error(`Insufficient CPU: ${cpuLimit} cores (need ${minCPU})`);
    }

    if (memoryLimit < (minMemoryGB * 1024 * 1024 * 1024)) {
        throw new Error(`Insufficient memory: ${(memoryLimit / 1024 / 1024).toFixed(0)}MB (need ${minMemoryGB * 1024}MB)`);
    }

    console.log(`✅ Resource validation passed: ${cpuLimit} cores, ${(memoryLimit / 1024 / 1024).toFixed(0)}MB`);
}

export default function() {
    // test logic
}
```

### Performance Regression Detection

```javascript
import toolbox from 'k6/x/toolbox';
import { sleep } from 'k6';

let previousCPU = 0;
let cpuSpikes = 0;

export function setup() {
    previousCPU = toolbox.getCPUUsage();
}

export default function() {
    const currentCPU = toolbox.getCPUUsage();

    // Detect CPU spikes
    if (currentCPU > previousCPU + 20) {
        cpuSpikes++;
        console.warn(`CPU spike detected: ${previousCPU.toFixed(1)}% → ${currentCPU.toFixed(1)}%`);
    }

    previousCPU = currentCPU;
    sleep(1);
}

export function teardown() {
    console.log(`Test completed. CPU spikes detected: ${cpuSpikes}`);
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