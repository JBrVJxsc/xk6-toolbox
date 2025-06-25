import toolbox from 'k6/x/toolbox';
import { check, sleep } from 'k6';

export const options = {
  vus: 1,
  duration: '15s',
};

export default function () {
  console.log('=== xk6-toolbox Extension Test ===');
  
  // Test 1: Get comprehensive system information (auto-fallback)
  let systemInfo;
  try {
    systemInfo = toolbox.getSystemInfo();
    console.log('✓ getSystemInfo() successful');
    console.log(`Method: ${systemInfo.method}, Fallback: ${systemInfo.fallback}`);
    console.log(`CPU: ${systemInfo.cpu.usage_percent.toFixed(2)}% (${systemInfo.cpu.used_cores.toFixed(2)}/${systemInfo.cpu.limit_cores} cores)`);
    console.log(`Memory: ${systemInfo.memory.usage_percent.toFixed(2)}% (${systemInfo.memory.usage_mb.toFixed(0)}/${systemInfo.memory.limit_mb.toFixed(0)} MB)`);
    console.log(`Load Average: ${systemInfo.cpu.load_average || 'N/A'}`);
    
    if (systemInfo.memory.free_bytes !== undefined) {
      console.log(`Memory breakdown: Free=${(systemInfo.memory.free_bytes / 1024 / 1024).toFixed(0)}MB, Buffer=${(systemInfo.memory.buffer_bytes / 1024 / 1024).toFixed(0)}MB, Cached=${(systemInfo.memory.cached_bytes / 1024 / 1024).toFixed(0)}MB`);
    }
  } catch (error) {
    console.log('⚠ getSystemInfo() failed (expected in test environment):', error.message);
  }
  
  // Test 2: Get command-based system information
  let systemInfoCmd;
  try {
    systemInfoCmd = toolbox.getSystemInfoCommand();
    console.log('✓ getSystemInfoCommand() successful');
    console.log(`Command Method: ${systemInfoCmd.method}, Fallback: ${systemInfoCmd.fallback}`);
    console.log(`Command CPU: ${systemInfoCmd.cpu.usage_percent.toFixed(2)}% (${systemInfoCmd.cpu.limit_cores} cores)`);
    console.log(`Command Memory: ${systemInfoCmd.memory.usage_percent.toFixed(2)}% (${systemInfoCmd.memory.usage_mb.toFixed(0)}MB)`);
    console.log(`Command Load: ${systemInfoCmd.cpu.load_average || 'N/A'}`);
  } catch (error) {
    console.log('⚠ getSystemInfoCommand() failed:', error.message);
  }
  
  // Test 3: Individual CPU functions
  let cpuUsage, cpuLimit, availableCPU;
  try {
    cpuUsage = toolbox.getCPUUsage();
    console.log(`✓ getCPUUsage(): ${cpuUsage.toFixed(2)}%`);
    __VU.metrics.counters['cpu_usage'].add(cpuUsage);
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
  
  // Test 4: Individual Memory functions
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
    __VU.metrics.counters['memory_usage_percent'].add(memoryPercent);
  } catch (error) {
    console.log('⚠ getMemoryUsagePercent() failed:', error.message);
  }
  
  try {
    availableMemory = toolbox.getAvailableMemory();
    console.log(`✓ getAvailableMemory(): ${availableMemory} bytes (${(availableMemory / 1024 / 1024).toFixed(2)} MB)`);
  } catch (error) {
    console.log('⚠ getAvailableMemory() failed:', error.message);
  }
  
  // Test 5: Raw command outputs
  try {
    const topOutput = toolbox.getTopOutput();
    console.log('✓ getTopOutput() successful');
    console.log('Top Output (first 200 chars):', topOutput.slice(0, 200));
  } catch (error) {
    console.log('⚠ getTopOutput() failed:', error.message);
  }
  
  try {
    const freeOutput = toolbox.getFreeOutput();
    console.log('✓ getFreeOutput() successful');
    console.log('Free Output:', freeOutput.trim());
  } catch (error) {
    console.log('⚠ getFreeOutput() failed:', error.message);
  }
  
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
  
  // Test 6: Comprehensive validation checks
  if (systemInfo) {
    check(systemInfo, {
      'System info has CPU data': (info) => info.cpu !== undefined,
      'System info has Memory data': (info) => info.memory !== undefined,
      'CPU usage is valid percentage': (info) => info.cpu.usage_percent >= 0 && info.cpu.usage_percent <= 100,
      'Memory usage is valid percentage': (info) => info.memory.usage_percent >= 0 && info.memory.usage_percent <= 100,
      'CPU limit is positive': (info) => info.cpu.limit_cores > 0,
      'Memory limit is positive': (info) => info.memory.limit_bytes > 0,
      'CPU used cores is reasonable': (info) => info.cpu.used_cores >= 0 && info.cpu.used_cores <= info.cpu.limit_cores,
      'CPU available cores is reasonable': (info) => info.cpu.available >= 0,
      'Memory usage bytes is positive': (info) => info.memory.usage_bytes >= 0,
      'Available memory is reasonable': (info) => info.memory.available_bytes >= 0,
      'Memory usage MB matches bytes': (info) => Math.abs(info.memory.usage_mb - (info.memory.usage_bytes / 1024 / 1024)) < 1,
      'Memory limit MB matches bytes': (info) => Math.abs(info.memory.limit_mb - (info.memory.limit_bytes / 1024 / 1024)) < 1,
      'Method is specified': (info) => info.method !== undefined && info.method !== '',
      'Fallback is boolean': (info) => typeof info.fallback === 'boolean',
      'Load average is string or undefined': (info) => info.cpu.load_average === undefined || typeof info.cpu.load_average === 'string',
    });
  }
  
  // Test 7: Individual method validation
  if (cpuUsage !== undefined) {
    check(null, {
      'Individual CPU usage is valid': () => cpuUsage >= 0 && cpuUsage <= 100,
    });
  }
  
  if (cpuLimit !== undefined) {
    check(null, {
      'Individual CPU limit is positive': () => cpuLimit > 0,
    });
  }
  
  if (memoryUsage !== undefined && memoryLimit !== undefined) {
    check(null, {
      'Memory usage does not exceed limit': () => memoryUsage <= memoryLimit,
      'Individual memory usage is positive': () => memoryUsage >= 0,
      'Individual memory limit is positive': () => memoryLimit > 0,
    });
  }
  
  if (memoryPercent !== undefined) {
    check(null, {
      'Individual memory percent is valid': () => memoryPercent >= 0 && memoryPercent <= 100,
    });
  }
  
  if (availableMemory !== undefined) {
    check(null, {
      'Available memory is reasonable': () => availableMemory >= 0,
    });
  }
  
  if (availableCPU !== undefined) {
    check(null, {
      'Available CPU is reasonable': () => availableCPU >= 0,
    });
  }
  
  // Test 8: Consistency checks between methods
  if (systemInfo && cpuUsage !== undefined && memoryPercent !== undefined) {
    check(null, {
      'CPU usage consistency': () => Math.abs(systemInfo.cpu.usage_percent - cpuUsage) < 5, // Allow 5% difference
      'Memory percent consistency': () => Math.abs(systemInfo.memory.usage_percent - memoryPercent) < 5,
    });
  }
  
  // Test 9: Simulate some CPU/memory work to see changes
  console.log('Simulating CPU and memory work...');
  const startTime = Date.now();
  
  // CPU work
  let sum = 0;
  for (let i = 0; i < 500000; i++) {
    sum += Math.sqrt(i * Math.random());
  }
  
  // Memory allocation
  const tempArray = new Array(10000).fill(0).map((_, i) => ({
    id: i,
    data: `test-data-${i}-${Math.random()}`,
    timestamp: Date.now(),
    largeString: 'x'.repeat(100),
  }));
  
  const workTime = Date.now() - startTime;
  console.log(`Work completed in ${workTime}ms (sum: ${sum.toFixed(2)}, array size: ${tempArray.length})`);
  
  // Test 10: Check resource usage after work
  try {
    const cpuUsageAfter = toolbox.getCPUUsage();
    const memoryUsageAfter = toolbox.getMemoryUsage();
    
    console.log(`CPU usage after work: ${cpuUsageAfter.toFixed(2)}%`);
    console.log(`Memory usage after work: ${(memoryUsageAfter / 1024 / 1024).toFixed(2)} MB`);
    
    if (cpuUsage !== undefined) {
      const cpuChange = cpuUsageAfter - cpuUsage;
      console.log(`CPU change: ${cpuChange > 0 ? '+' : ''}${cpuChange.toFixed(2)}%`);
    }
    
    if (memoryUsage !== undefined) {
      const memoryChange = (memoryUsageAfter - memoryUsage) / 1024 / 1024;
      console.log(`Memory change: ${memoryChange > 0 ? '+' : ''}${memoryChange.toFixed(2)} MB`);
    }
    
    // Validate the changes are reasonable
    check(null, {
      'CPU usage after work is still valid': () => cpuUsageAfter >= 0 && cpuUsageAfter <= 100,
      'Memory usage after work is positive': () => memoryUsageAfter >= 0,
    });
    
  } catch (error) {
    console.log('⚠ Could not get resource usage after work:', error.message);
  }
  
  // Test 11: Compare cgroup vs command methods (if both work)
  if (systemInfo && systemInfoCmd) {
    const cpuDiff = Math.abs(systemInfo.cpu.usage_percent - systemInfoCmd.cpu.usage_percent);
    const memoryDiff = Math.abs(systemInfo.memory.usage_percent - systemInfoCmd.memory.usage_percent);
    
    console.log(`Method comparison - CPU diff: ${cpuDiff.toFixed(2)}%, Memory diff: ${memoryDiff.toFixed(2)}%`);
    
    check(null, {
      'CPU methods are reasonably close': () => cpuDiff < 20, // Allow 20% difference between methods
      'Memory methods are reasonably close': () => memoryDiff < 20,
      'Both methods return valid data': () => systemInfo.method !== '' && systemInfoCmd.method === 'command',
    });
  }
  
  // Test: Connectivity check (TCP/HTTP)
  try {
    const domain = 'google.com';
    const port = '80';
    const timeout = 5;
    const connReport = toolbox.checkConnectivity(domain, port, timeout);
    console.log(`\u2713 checkConnectivity(${domain}, ${port}, ${timeout})`);
    console.log('Connectivity Report:', JSON.stringify(connReport, null, 2));
    check(connReport, {
      'Domain matches': (r) => r.domain === domain,
      'Port matches': (r) => r.port === port,
      'TCP success': (r) => r.tcp === 'success',
      'HTTP result present': (r) => typeof r.http === 'string' && r.http.length > 0 && r.http !== 'skipped (TCP failed)',
    });
  } catch (error) {
    console.log('\u26a0 checkConnectivity() failed:', error.message);
  }
  
  // Cleanup large array to free memory
  tempArray.length = 0;
  
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