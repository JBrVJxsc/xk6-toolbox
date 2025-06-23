import toolbox from 'k6/x/toolbox';
import { check, sleep } from 'k6';

export const options = {
  vus: 3,
  duration: '2m',
  thresholds: {
    checks: ['rate>0.8'],
  },
};

// Custom metrics for resource monitoring
export const resourceMetrics = {
  cpu_usage: new Trend('cpu_usage_percent'),
  memory_usage: new Trend('memory_usage_percent'),
  cpu_limit: new Gauge('cpu_limit_cores'),
  memory_limit: new Gauge('memory_limit_mb'),
};

export default function () {
  // Get system information
  const info = toolbox.getSystemInfo();
  
  // Record metrics
  resourceMetrics.cpu_usage.add(info.cpu.usage_percent);
  resourceMetrics.memory_usage.add(info.memory.usage_percent);
  resourceMetrics.cpu_limit.add(info.cpu.limit_cores);
  resourceMetrics.memory_limit.add(info.memory.limit_mb);
  
  // Check resource thresholds
  check(info, {
    'CPU usage is below 80%': (data) => data.cpu.usage_percent < 80,
    'Memory usage is below 85%': (data) => data.memory.usage_percent < 85,
    'CPU limit is reasonable': (data) => data.cpu.limit_cores >= 0.5,
    'Memory limit is reasonable': (data) => data.memory.limit_mb >= 100,
    'Available CPU is positive': (data) => data.cpu.available_cores > 0,
    'Available memory is positive': (data) => data.memory.available_bytes > 0,
  });
  
  // Simulate different workloads based on VU number
  const workload = __VU % 3;
  
  switch (workload) {
    case 0:
      // Light workload - just sleep
      sleep(2);
      break;
    case 1:
      // Medium workload - some CPU work
      let sum = 0;
      for (let i = 0; i < 100000; i++) {
        sum += Math.sqrt(i);
      }
      sleep(1);
      break;
    case 2:
      // Heavy workload - more CPU work
      let heavySum = 0;
      for (let i = 0; i < 500000; i++) {
        heavySum += Math.sqrt(i) * Math.sin(i);
      }
      sleep(0.5);
      break;
  }
  
  // Log resource usage every 10 iterations
  if (__ITER % 10 === 0) {
    console.log(`VU ${__VU}, Iteration ${__ITER}:`);
    console.log(`  CPU: ${info.cpu.usage_percent.toFixed(1)}% (${info.cpu.used_cores.toFixed(2)}/${info.cpu.limit_cores.toFixed(2)} cores)`);
    console.log(`  Memory: ${info.memory.usage_percent.toFixed(1)}% (${info.memory.usage_mb.toFixed(1)}/${info.memory.limit_mb.toFixed(1)} MB)`);
  }
}

export function handleSummary(data) {
  console.log('=== Resource Threshold Test Summary ===');
  
  // Calculate resource statistics
  const cpuStats = data.metrics.cpu_usage_percent;
  const memoryStats = data.metrics.memory_usage_percent;
  
  console.log('CPU Usage Statistics:');
  console.log(`  Average: ${cpuStats.values.avg.toFixed(2)}%`);
  console.log(`  Min: ${cpuStats.values.min.toFixed(2)}%`);
  console.log(`  Max: ${cpuStats.values.max.toFixed(2)}%`);
  console.log(`  95th percentile: ${cpuStats.values['p(95)'].toFixed(2)}%`);
  
  console.log('Memory Usage Statistics:');
  console.log(`  Average: ${memoryStats.values.avg.toFixed(2)}%`);
  console.log(`  Min: ${memoryStats.values.min.toFixed(2)}%`);
  console.log(`  Max: ${memoryStats.values.max.toFixed(2)}%`);
  console.log(`  95th percentile: ${memoryStats.values['p(95)'].toFixed(2)}%`);
  
  // Check for resource violations
  const cpuViolations = data.metrics.cpu_usage_percent.values.count - data.metrics.cpu_usage_percent.values.passes;
  const memoryViolations = data.metrics.memory_usage_percent.values.count - data.metrics.memory_usage_percent.values.passes;
  
  if (cpuViolations > 0) {
    console.warn(`⚠ CPU threshold violations: ${cpuViolations}`);
  }
  
  if (memoryViolations > 0) {
    console.warn(`⚠ Memory threshold violations: ${memoryViolations}`);
  }
  
  return {
    'resource-threshold-results.json': JSON.stringify(data, null, 2),
  };
} 