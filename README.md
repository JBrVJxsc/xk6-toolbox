# xk6-toolbox

A k6 extension for monitoring system resources in containerized environments. This extension provides real-time access to CPU and memory usage information from cgroups, making it ideal for performance testing and resource monitoring in containerized applications.

## Features

- **CPU Monitoring**: Get current CPU usage, limits, and available cores
- **Memory Monitoring**: Get current memory usage, limits, and available memory
- **Container-Aware**: Automatically detects and reads from cgroup v1 and v2
- **Fallback Support**: Falls back to system-level metrics when cgroups are not available
- **Comprehensive API**: Both individual metrics and complete system information

## Installation

### Prerequisites

- Go 1.19 or later
- xk6 v0.20.1 or later

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
  // Get comprehensive system information
  const systemInfo = toolbox.getSystemInfo();
  console.log('System Info:', JSON.stringify(systemInfo, null, 2));
  
  // Get individual metrics
  const cpuUsage = toolbox.getCPUUsage();
  const memoryUsage = toolbox.getMemoryUsage();
  
  console.log(`CPU Usage: ${cpuUsage}%`);
  console.log(`Memory Usage: ${memoryUsage} bytes`);
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
  });
  
  // Log detailed information
  console.log(`CPU: ${info.cpu.used_cores.toFixed(2)}/${info.cpu.limit_cores.toFixed(2)} cores (${info.cpu.usage_percent.toFixed(1)}%)`);
  console.log(`Memory: ${info.memory.usage_mb.toFixed(1)}/${info.memory.limit_mb.toFixed(1)} MB (${info.memory.usage_percent.toFixed(1)}%)`);
  
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
  const memoryPercent = toolbox.getMemoryUsagePercent();
  
  // Log resource usage
  console.log(`VU ${__VU}: CPU=${cpuUsage.toFixed(1)}%, Memory=${memoryPercent.toFixed(1)}%`);
  
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
Returns comprehensive system resource information.

**Returns:** `SystemInfo` object with the following structure:
```javascript
{
  cpu: {
    usage_percent: number,    // CPU usage as percentage
    limit_cores: number,      // CPU limit in cores
    used_cores: number,       // Currently used CPU cores
    available_cores: number   // Available CPU cores
  },
  memory: {
    usage_bytes: number,      // Memory usage in bytes
    limit_bytes: number,      // Memory limit in bytes
    available_bytes: number,  // Available memory in bytes
    usage_percent: number,    // Memory usage as percentage
    usage_mb: number,         // Memory usage in MB
    limit_mb: number,         // Memory limit in MB
    available_mb: number      // Available memory in MB
  }
}
```

#### `getCPUUsage()`
Returns current CPU usage as a percentage.

**Returns:** `number` - CPU usage percentage (0-100)

#### `getCPULimit()`
Returns the CPU limit in cores.

**Returns:** `number` - CPU limit in cores

#### `getMemoryUsage()`
Returns current memory usage in bytes.

**Returns:** `number` - Memory usage in bytes

#### `getMemoryLimit()`
Returns the memory limit in bytes.

**Returns:** `number` - Memory limit in bytes

#### `getMemoryUsagePercent()`
Returns memory usage as a percentage.

**Returns:** `number` - Memory usage percentage (0-100)

#### `getAvailableMemory()`
Returns available memory in bytes.

**Returns:** `number` - Available memory in bytes

#### `getAvailableCPU()`
Returns available CPU cores.

**Returns:** `number` - Available CPU cores

## Container Support

This extension is designed to work in containerized environments and automatically detects:

- **cgroup v2**: Modern container runtimes (Docker 20.10+, containerd 1.4+)
- **cgroup v1**: Legacy container runtimes
- **System fallback**: When cgroups are not available

### Supported Container Runtimes

- Docker
- containerd
- CRI-O
- Podman
- Kubernetes (with any supported runtime)

## Error Handling

The extension gracefully handles errors and provides meaningful error messages:

- Missing cgroup files (expected in non-container environments)
- Permission issues
- Invalid cgroup data
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

1. **Permission denied errors**: Ensure the container has access to `/sys/fs/cgroup`
2. **No cgroup information**: The extension will fall back to system metrics
3. **High CPU usage**: Consider reducing the frequency of resource checks

### Debug Mode

Enable debug logging by setting the `K6_DEBUG` environment variable:

```bash
K6_DEBUG=true k6 run your-test.js
```

## Examples

See the `examples/` directory for additional usage examples and test scenarios.