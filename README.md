# xk6-toolbox

A k6 extension for monitoring system resources in containerized and non-containerized environments. This extension provides real-time access to CPU and memory usage information using multiple methods with automatic fallback, making it ideal for performance testing and resource monitoring in various environments.

## Features

- **Multi-Method Monitoring**: Automatically uses cgroup-based monitoring in containers and falls back to command-based monitoring
- **CPU Monitoring**: Get current CPU usage, limits, available cores, and load average
- **Memory Monitoring**: Get current memory usage, limits, available memory, buffers, and cache information
- **Container-Aware**: Automatically detects and reads from cgroup v1 and v2
- **Command Fallback**: Falls back to system commands (`top`, `free`, `ps`, `uptime`) when cgroups are not available
- **Raw Command Output**: Access raw output from system monitoring commands
- **Comprehensive API**: Both individual metrics and complete system information

## Installation

### Prerequisites

- Go 1.19 or later
- xk6 v0.20.1 or later
- System commands: `top`, `free`, `ps`, `uptime`, `nproc` (for command-based monitoring)

### Building k6 with the extension

```bash
# Clone the repository
git clone https://github.com/your-username/xk6-toolbox.git
cd xk6-toolbox

# Build k6 with the toolbox extension
xk6 build v1.0.0 --with github.com/your-username/xk6-toolbox=./
```

## Usage

### Basic Example

```javascript
import toolbox from 'k6/x/toolbox';

export default function () {
  // Get comprehensive system information (auto-fallback)
  const systemInfo = toolbox.getSystemInfo();
  console.log('System Info:', JSON.stringify(systemInfo, null, 2));
  
  // Get individual metrics
  const cpuUsage = toolbox.getCPUUsage();
  const memoryUsage = toolbox.getMemoryUsage();
  
  console.log(`CPU Usage: ${cpuUsage}%`);
  console.log(`Memory Usage: ${memoryUsage} bytes`);
  console.log(`Collection Method: ${systemInfo.method}`);
  console.log(`Used Fallback: ${systemInfo.fallback}`);
}
```

### Command-Based Monitoring

```javascript
import toolbox from 'k6/x/toolbox';

export default function () {
  // Force command-based monitoring only
  const info = toolbox.getSystemInfoCommand();
  console.log('Command-based monitoring:', JSON.stringify(info, null, 2));
  
  // Get raw command outputs
  const topOutput = toolbox.getTopOutput();
  const freeOutput = toolbox.getFreeOutput();
  const psOutput = toolbox.getPsOutput();
  const uptimeOutput = toolbox.getUptimeOutput();
  
  console.log('Top output:', topOutput);
  console.log('Free output:', freeOutput);
  console.log('PS output:', psOutput);
  console.log('Uptime output:', uptimeOutput);
}
```

### Advanced Example

```javascript
import toolbox from 'k6/x/toolbox';
import { check } from 'k6';

export default function () {
  // Get system information
  const info = toolbox.getSystemInfo();
  
  // Check resource usage
  check(info, {
    'CPU usage is reasonable': (info) => info.cpu.usage_percent < 90,
    'Memory usage is reasonable': (info) => info.memory.usage_percent < 85,
    'CPU limit is set': (info) => info.cpu.limit_cores > 0,
    'Memory limit is set': (info) => info.memory.limit_bytes > 0,
    'Load average is available': (info) => info.cpu.load_average !== undefined,
  });
  
  // Log detailed information
  console.log(`CPU: ${info.cpu.used_cores.toFixed(2)}/${info.cpu.limit_cores.toFixed(2)} cores (${info.cpu.usage_percent.toFixed(1)}%)`);
  console.log(`Memory: ${info.memory.usage_mb.toFixed(1)}/${info.memory.limit_mb.toFixed(1)} MB (${info.memory.usage_percent.toFixed(1)}%)`);
  console.log(`Load Average: ${info.cpu.load_average}`);
  console.log(`Free Memory: ${info.memory.free_bytes} bytes`);
  console.log(`Buffer Memory: ${info.memory.buffer_bytes} bytes`);
  console.log(`Cached Memory: ${info.memory.cached_bytes} bytes`);
  console.log(`Collection Method: ${info.method}`);
  
  // Simulate some work
  const start = Date.now();
  while (Date.now() - start < 1000) {
    // Do some CPU-intensive work
  }
}
```

### Resource Monitoring Test

```javascript
import toolbox from 'k6/x/toolbox';
import { sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  // Monitor resources during test execution
  const cpuUsage = toolbox.getCPUUsage();
  const memoryUsage = toolbox.getMemoryUsage();
  
  // Get detailed info
  const info = toolbox.getSystemInfo();
  const memoryPercent = info.memory.usage_percent;
  
  // Log resource usage
  console.log(`VU ${__VU}: CPU=${cpuUsage.toFixed(1)}%, Memory=${memoryPercent.toFixed(1)}%`);
  console.log(`Load: ${info.cpu.load_average}, Method: ${info.method}`);
  
  // Check for resource constraints
  if (cpuUsage > 95) {
    console.warn(`High CPU usage detected: ${cpuUsage}%`);
  }
  
  if (memoryPercent > 90) {
    console.warn(`High memory usage detected: ${memoryPercent}%`);
  }
  
  sleep(1);
}
```

## API Reference

### Functions

#### `getSystemInfo()`
Returns comprehensive system resource information with automatic fallback.

**Returns:** `SystemInfo` object with the following structure:
```javascript
{
  cpu: {
    usage_percent: number,    // CPU usage as percentage
    limit_cores: number,      // CPU limit in cores
    used_cores: number,       // Currently used CPU cores
    available_cores: number,  // Available CPU cores
    load_average: string      // System load average
  },
  memory: {
    usage_bytes: number,      // Memory usage in bytes
    limit_bytes: number,      // Memory limit in bytes
    available_bytes: number,  // Available memory in bytes
    usage_percent: number,    // Memory usage as percentage
    usage_mb: number,         // Memory usage in MB
    limit_mb: number,         // Memory limit in MB
    available_mb: number,     // Available memory in MB
    free_bytes: number,       // Free memory in bytes
    buffer_bytes: number,     // Buffer memory in bytes
    cached_bytes: number      // Cached memory in bytes
  },
  method: string,             // How data was collected ("cgroup", "command", "mixed")
  fallback: boolean           // Whether fallback methods were used
}
```

#### `getSystemInfoCommand()`
Forces using command-based monitoring only.

**Returns:** `SystemInfo` object (same structure as above, but `method` will be "command" and `fallback` will be false)

#### `getCPUUsage()`
Returns current CPU usage as a percentage.

**Returns:** `number` - CPU usage percentage (0-100)

#### `getMemoryUsage()`
Returns current memory usage in bytes.

**Returns:** `number` - Memory usage in bytes

#### `getTopOutput()`
Returns raw output from the `top` command.

**Returns:** `string` - Raw top command output

#### `getFreeOutput()`
Returns raw output from the `free` command.

**Returns:** `string` - Raw free command output

#### `getPsOutput()`
Returns raw output from the `ps aux` command.

**Returns:** `string` - Raw ps command output

#### `getUptimeOutput()`
Returns raw output from the `uptime` command.

**Returns:** `string` - Raw uptime command output

## Monitoring Methods

### Cgroup-Based Monitoring (Primary)
- **CPU**: Reads from `/sys/fs/cgroup/cpu.max` (v2) or `/sys/fs/cgroup/cpu,cpuacct/cpu.cfs_quota_us` (v1)
- **Memory**: Reads from `/sys/fs/cgroup/memory.current` (v2) or `/sys/fs/cgroup/memory/memory.usage_in_bytes` (v1)
- **Used in**: Container environments (Docker, Kubernetes, etc.)

### Command-Based Monitoring (Fallback)
- **CPU**: Uses `top` command to parse CPU usage and `nproc` for core count
- **Memory**: Uses `free` command to get memory statistics
- **Load Average**: Uses `uptime` command
- **Used in**: Non-container environments or when cgroups are not available

### Automatic Fallback
The extension automatically tries cgroup-based monitoring first, then falls back to command-based monitoring if cgroup files are not accessible.

## Container Support

This extension is designed to work in containerized environments and automatically detects:

- **cgroup v2**: Modern container runtimes (Docker 20.10+, containerd 1.4+)
- **cgroup v1**: Legacy container runtimes
- **Command fallback**: When cgroups are not available

### Supported Container Runtimes

- Docker
- containerd
- CRI-O
- Podman
- Kubernetes (with any supported runtime)

### Supported Commands

For command-based monitoring, the following system commands are used:
- `top` - CPU usage information
- `free` - Memory usage information
- `ps` - Process information
- `uptime` - System load average
- `nproc` - Number of CPU cores

## Error Handling

The extension gracefully handles errors and provides meaningful error messages:

- Missing cgroup files (expected in non-container environments)
- Missing system commands (fallback to available commands)
- Permission issues
- Invalid data formats
- System resource access failures

## Development

### Running Tests

```bash
# Run Go unit tests
make test-go

# Run k6 JavaScript tests
make test-k6

# Run all tests
make test
```

### Building

```bash
# Build k6 with extension
make build
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## Troubleshooting

### Common Issues

1. **Permission denied errors**: Ensure the container has access to `/sys/fs/cgroup` or system commands
2. **No cgroup information**: The extension will automatically fall back to command-based monitoring
3. **Missing system commands**: Install required commands (`top`, `free`, `ps`, `uptime`, `nproc`)
4. **High CPU usage**: Consider reducing the frequency of resource checks

### Debug Mode

Enable debug logging by setting the `K6_DEBUG` environment variable:

```bash
K6_DEBUG=true k6 run your-test.js
```

### Method Detection

You can check which monitoring method is being used:

```javascript
const info = toolbox.getSystemInfo();
console.log(`Using method: ${info.method}`);
console.log(`Fallback used: ${info.fallback}`);
```

## Examples

See the `examples/` directory for additional usage examples and test scenarios.