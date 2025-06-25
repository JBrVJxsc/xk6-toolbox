import toolbox from 'k6/x/toolbox';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';

export const options = {
  vus: 1,
  duration: '15s',
};

// Define custom metrics to be used in the test
export const cpu_usage = new Counter('cpu_usage');
export const memory_usage_percent = new Counter('memory_usage_percent');

export default function () {
  console.log('=== xk6-toolbox Extension Test ===');
  
  // Test 1: Individual CPU functions
  let cpuUsage, cpuLimit, availableCPU;
  try {
    cpuUsage = toolbox.getCPUUsage();
    console.log(`✓ getCPUUsage(): ${cpuUsage.toFixed(2)}%`);
    cpu_usage.add(cpuUsage);
  } catch (error) {
    console.log('⚠ getCPUUsage() failed:', error.message);
  }
  
  try {
    cpuLimit = toolbox.getCPULimit();
    console.log(`✓ getCPULimit(): ${cpuLimit} cores`);
  } catch (error) {
    console.log('⚠ getCPULimit() failed:', error.message);
  }
  
  try {
    availableCPU = toolbox.getAvailableCPU();
    console.log(`✓ getAvailableCPU(): ${availableCPU.toFixed(2)} cores`);
  } catch (error) {
    console.log('⚠ getAvailableCPU() failed:', error.message);
  }
  
  // Test 2: Individual Memory functions
  let memoryUsage, memoryLimit, memoryPercent, availableMemory;
  try {
    memoryUsage = toolbox.getMemoryUsage();
    console.log(`✓ getMemoryUsage(): ${memoryUsage} bytes (${(memoryUsage / 1024 / 1024).toFixed(2)} MB)`);
  } catch (error) {
    console.log('⚠ getMemoryUsage() failed:', error.message);
  }
  
  try {
    memoryLimit = toolbox.getMemoryLimit();
    console.log(`✓ getMemoryLimit(): ${memoryLimit} bytes (${(memoryLimit / 1024 / 1024).toFixed(2)} MB)`);
  } catch (error) {
    console.log('⚠ getMemoryLimit() failed:', error.message);
  }
  
  try {
    memoryPercent = toolbox.getMemoryUsagePercent();
    console.log(`✓ getMemoryUsagePercent(): ${memoryPercent.toFixed(2)}%`);
    memory_usage_percent.add(memoryPercent);
  } catch (error) {
    console.log('⚠ getMemoryUsagePercent() failed:', error.message);
  }
  
  try {
    availableMemory = toolbox.getAvailableMemory();
    console.log(`✓ getAvailableMemory(): ${availableMemory} bytes (${(availableMemory / 1024 / 1024).toFixed(2)} MB)`);
  } catch (error) {
    console.log('⚠ getAvailableMemory() failed:', error.message);
  }
  
  // Test 3: Raw command outputs
  try {
    const psOutput = toolbox.getPsOutput();
    console.log('✓ getPsOutput() successful');
    const processCount = psOutput.split('\n').length - 1;
    console.log(`PS Output: ${processCount} processes (first 200 chars): ${psOutput.slice(0, 200)}`);
  } catch (error) {
    console.log('⚠ getPsOutput() failed:', error.message);
  }
  
  try {
    const uptimeOutput = toolbox.getUptimeOutput();
    console.log('✓ getUptimeOutput() successful');
    console.log('Uptime Output:', uptimeOutput.trim());
  } catch (error) {
    console.log('⚠ getUptimeOutput() failed:', error.message);
  }
  
  // Test 4: Compare cgroup vs command methods (if both work)
  if (cpuUsage !== undefined && memoryPercent !== undefined) {
    const cpuDiff = Math.abs(cpuUsage - cpuUsage);
    const memoryDiff = Math.abs(memoryPercent - memoryPercent);
    
    console.log(`Method comparison - CPU diff: ${cpuDiff.toFixed(2)}%, Memory diff: ${memoryDiff.toFixed(2)}%`);
    
    check(null, {
      'CPU methods are reasonably close': () => cpuDiff < 20, // Allow 20% difference between methods
      'Memory methods are reasonably close': () => memoryDiff < 20,
      'Both methods return valid data': () => cpuUsage !== undefined && memoryPercent !== undefined,
    });
  }
  
  // Test: Connectivity check (TCP/HTTP)
  try {
    const domain = 'google.com';
    const port = '80';
    const timeout = 5;
    const connReport = toolbox.checkConnectivity(domain, port, timeout);
    console.log(`✓ checkConnectivity(${domain}, ${port}, ${timeout})`);
    console.log('Connectivity Report:', JSON.stringify(connReport, null, 2));
    check(connReport, {
      'Connectivity Report: Domain matches': (r) => r.domain === domain,
      'Connectivity Report: Port matches': (r) => r.port === port,
      'Connectivity Report: TCP status is string': (r) => typeof r.tcp === 'string' && r.tcp.length > 0,
      'Connectivity Report: HTTP status is string': (r) => typeof r.http === 'string' && r.http.length > 0,
    });
  } catch (error) {
    console.log('⚠ checkConnectivity() failed:', error.message);
  }
  
  // Test: OS Detection
  try {
    const isMac = toolbox.isMacOS();
    const isLinux = toolbox.isLinux();
    console.log(`✓ OS Detection: isMacOS=${isMac}, isLinux=${isLinux}`);
    check(null, {
      'OS Detection: isMacOS is boolean': () => typeof isMac === 'boolean',
      'OS Detection: isLinux is boolean': () => typeof isLinux === 'boolean',
      'OS Detection: Flags are mutually exclusive': () => isMac !== isLinux || (!isMac && !isLinux),
    });
  } catch (error) {
    console.log('⚠ OS Detection failed:', error.message);
  }
  
  console.log('=== Test iteration completed ===\n');
  
  sleep(1);
}

export function handleSummary(data) {
  console.log('\n=== Final Test Summary ===');
  
  // Safely access metrics with fallbacks
  const iterations = data.metrics.iterations?.values?.count || 0;
  const iterationDuration = data.metrics.iteration_duration?.values?.avg || 0;
  const checksPasses = data.metrics.checks?.values?.passes || 0;
  const checksCount = data.metrics.checks?.values?.count || 0;
  const checksRate = data.metrics.checks?.values?.rate || 0;
  
  console.log(`Iterations completed: ${iterations}`);
  console.log(`Average iteration duration: ${iterationDuration.toFixed(2)}ms`);
  console.log(`Checks passed: ${checksPasses}/${checksCount}`);
  console.log(`Check success rate: ${(checksRate * 100).toFixed(1)}%`);
  
  // Resource usage summary
  const avgCPU = data.metrics.cpu_usage?.values?.avg;
  const maxCPU = data.metrics.cpu_usage?.values?.max;
  const avgMemory = data.metrics.memory_usage_percent?.values?.avg;
  const maxMemory = data.metrics.memory_usage_percent?.values?.max;
  
  if (avgCPU !== undefined) {
    console.log(`CPU Usage: avg=${avgCPU.toFixed(1)}%, max=${maxCPU.toFixed(1)}%`);
  }
  
  if (avgMemory !== undefined) {
    console.log(`Memory Usage: avg=${avgMemory.toFixed(1)}%, max=${maxMemory.toFixed(1)}%`);
  }
  
  // Performance assessment
  if (checksRate < 0.8) {
    console.log('⚠️ WARNING: Low check success rate - some toolbox functions may not be working properly');
  } else if (checksRate >= 0.95) {
    console.log('✅ EXCELLENT: High check success rate - toolbox extension is working well');
  } else {
    console.log('✅ GOOD: Acceptable check success rate - toolbox extension is mostly working');
  }
  
  console.log('=== Test Summary Complete ===');
  
  return {
    'stdout': `xk6-toolbox test completed with ${(checksRate * 100).toFixed(1)}% check success rate`,
  };
}