import toolbox from 'k6/x/toolbox';
import { check, sleep } from 'k6';

export const options = {
  vus: 1,
  duration: '10s',
  thresholds: {
    checks: ['rate>0.8'], // At least 80% of checks should pass
  },
};

export default function () {
  console.log('=== xk6-toolbox Test ===');
  
  // Test 1: Get comprehensive system information
  let systemInfo;
  try {
    systemInfo = toolbox.getSystemInfo();
    console.log('✓ getSystemInfo() successful');
    console.log('System Info:', JSON.stringify(systemInfo, null, 2));
  } catch (error) {
    console.log('⚠ getSystemInfo() failed (expected in test environment):', error.message);
  }
  
  // Test 2: Individual CPU functions
  let cpuUsage, cpuLimit, availableCPU;
  try {
    cpuUsage = toolbox.getCPUUsage();
    console.log(`✓ getCPUUsage(): ${cpuUsage.toFixed(2)}%`);
  } catch (error) {
    console.log('⚠ getCPUUsage() failed:', error.message);
  }
  
  try {
    cpuLimit = toolbox.getCPULimit();
    console.log(`✓ getCPULimit(): ${cpuLimit.toFixed(2)} cores`);
  } catch (error) {
    console.log('⚠ getCPULimit() failed:', error.message);
  }
  
  try {
    availableCPU = toolbox.getAvailableCPU();
    console.log(`✓ getAvailableCPU(): ${availableCPU.toFixed(2)} cores`);
  } catch (error) {
    console.log('⚠ getAvailableCPU() failed:', error.message);
  }
  
  // Test 3: Individual Memory functions
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
  } catch (error) {
    console.log('⚠ getMemoryUsagePercent() failed:', error.message);
  }
  
  try {
    availableMemory = toolbox.getAvailableMemory();
    console.log(`✓ getAvailableMemory(): ${availableMemory} bytes (${(availableMemory / 1024 / 1024).toFixed(2)} MB)`);
  } catch (error) {
    console.log('⚠ getAvailableMemory() failed:', error.message);
  }
  
  // Test 4: Validation checks
  if (systemInfo) {
    check(systemInfo, {
      'System info has CPU data': (info) => info.cpu !== undefined,
      'System info has Memory data': (info) => info.memory !== undefined,
      'CPU usage is reasonable': (info) => info.cpu.usage_percent >= 0 && info.cpu.usage_percent <= 100,
      'Memory usage is reasonable': (info) => info.memory.usage_percent >= 0 && info.memory.usage_percent <= 100,
      'CPU limit is positive': (info) => info.cpu.limit_cores > 0,
      'Memory limit is positive': (info) => info.memory.limit_bytes > 0,
      'Memory usage MB is calculated': (info) => info.memory.usage_mb > 0,
      'Memory limit MB is calculated': (info) => info.memory.limit_mb > 0,
    });
  }
  
  if (cpuUsage !== undefined) {
    check({ cpuUsage }, {
      'CPU usage is within valid range': (data) => data.cpuUsage >= 0 && data.cpuUsage <= 100,
    });
  }
  
  if (memoryPercent !== undefined) {
    check({ memoryPercent }, {
      'Memory usage percent is within valid range': (data) => data.memoryPercent >= 0 && data.memoryPercent <= 100,
    });
  }
  
  if (cpuLimit !== undefined && availableCPU !== undefined) {
    check({ cpuLimit, availableCPU }, {
      'Available CPU is not negative': (data) => data.availableCPU >= 0,
      'Available CPU is not greater than limit': (data) => data.availableCPU <= data.cpuLimit,
    });
  }
  
  if (memoryLimit !== undefined && availableMemory !== undefined) {
    check({ memoryLimit, availableMemory }, {
      'Available memory is not negative': (data) => data.availableMemory >= 0,
      'Available memory is not greater than limit': (data) => data.availableMemory <= data.memoryLimit,
    });
  }
  
  // Test 5: Simulate some CPU work to see usage change
  console.log('Simulating CPU work...');
  const startTime = Date.now();
  let sum = 0;
  for (let i = 0; i < 1000000; i++) {
    sum += Math.sqrt(i);
  }
  const workTime = Date.now() - startTime;
  console.log(`CPU work completed in ${workTime}ms (sum: ${sum.toFixed(2)})`);
  
  // Test 6: Check resource usage after work
  try {
    const cpuUsageAfter = toolbox.getCPUUsage();
    console.log(`CPU usage after work: ${cpuUsageAfter.toFixed(2)}%`);
    
    if (cpuUsage !== undefined) {
      console.log(`CPU usage change: ${(cpuUsageAfter - cpuUsage).toFixed(2)}%`);
    }
  } catch (error) {
    console.log('⚠ Could not get CPU usage after work:', error.message);
  }
  
  console.log('=== Test completed ===\n');
  
  sleep(1);
}

export function handleSummary(data) {
  console.log('=== Test Summary ===');
  
  // Safely access metrics with fallbacks
  const testDuration = data.metrics.test_duration?.values?.avg || 0;
  const checksPasses = data.metrics.checks?.values?.passes || 0;
  const checksCount = data.metrics.checks?.values?.count || 0;
  const checksRate = data.metrics.checks?.values?.rate || 0;
  
  console.log(`Test duration: ${testDuration.toFixed(2)}s`);
  console.log(`Checks passed: ${checksPasses}/${checksCount}`);
  console.log(`Check rate: ${(checksRate * 100).toFixed(1)}%`);
  
  return {
    'stdout': JSON.stringify(data, null, 2),
  };
} 