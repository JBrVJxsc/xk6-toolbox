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
  
  // Test 1: Get comprehensive system information (auto-fallback)
  let systemInfo;
  try {
    systemInfo = toolbox.getSystemInfo();
    console.log('✓ getSystemInfo() successful');
    console.log('System Info:', JSON.stringify(systemInfo, null, 2));
    console.log(`Method: ${systemInfo.method}, Fallback: ${systemInfo.fallback}`);
    console.log(`Load Average: ${systemInfo.cpu.load_average}`);
    console.log(`Free: ${systemInfo.memory.free_bytes}, Buffer: ${systemInfo.memory.buffer_bytes}, Cached: ${systemInfo.memory.cached_bytes}`);
  } catch (error) {
    console.log('⚠ getSystemInfo() failed (expected in test environment):', error.message);
  }
  
  // Test 2: Get command-based system information
  let systemInfoCmd;
  try {
    systemInfoCmd = toolbox.getSystemInfoCommand();
    console.log('✓ getSystemInfoCommand() successful');
    console.log('System Info (command):', JSON.stringify(systemInfoCmd, null, 2));
    console.log(`Method: ${systemInfoCmd.method}, Fallback: ${systemInfoCmd.fallback}`);
    console.log(`Load Average: ${systemInfoCmd.cpu.load_average}`);
    console.log(`Free: ${systemInfoCmd.memory.free_bytes}, Buffer: ${systemInfoCmd.memory.buffer_bytes}, Cached: ${systemInfoCmd.memory.cached_bytes}`);
  } catch (error) {
    console.log('⚠ getSystemInfoCommand() failed:', error.message);
  }
  
  // Test 3: Individual CPU and Memory functions
  let cpuUsage, memoryUsage;
  try {
    cpuUsage = toolbox.getCPUUsage();
    console.log(`✓ getCPUUsage(): ${cpuUsage.toFixed(2)}%`);
  } catch (error) {
    console.log('⚠ getCPUUsage() failed:', error.message);
  }
  
  try {
    memoryUsage = toolbox.getMemoryUsage();
    console.log(`✓ getMemoryUsage(): ${memoryUsage} bytes (${(memoryUsage / 1024 / 1024).toFixed(2)} MB)`);
  } catch (error) {
    console.log('⚠ getMemoryUsage() failed:', error.message);
  }
  
  // Test 4: Raw command outputs
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
    console.log('Free Output:', freeOutput);
  } catch (error) {
    console.log('⚠ getFreeOutput() failed:', error.message);
  }
  
  try {
    const psOutput = toolbox.getPsOutput();
    console.log('✓ getPsOutput() successful');
    console.log('PS Output (first 200 chars):', psOutput.slice(0, 200));
  } catch (error) {
    console.log('⚠ getPsOutput() failed:', error.message);
  }
  
  try {
    const uptimeOutput = toolbox.getUptimeOutput();
    console.log('✓ getUptimeOutput() successful');
    console.log('Uptime Output:', uptimeOutput);
  } catch (error) {
    console.log('⚠ getUptimeOutput() failed:', error.message);
  }
  
  // Test 5: Validation checks
  if (systemInfo) {
    check(systemInfo, {
      'System info has CPU data': (info) => info.cpu !== undefined,
      'System info has Memory data': (info) => info.memory !== undefined,
      'CPU usage is reasonable': (info) => info.cpu.usage_percent >= 0 && info.cpu.usage_percent <= 100,
      'Memory usage is reasonable': (info) => info.memory.usage_percent === undefined || (info.memory.usage_percent >= 0 && info.memory.usage_percent <= 100),
      'CPU limit is positive': (info) => info.cpu.limit_cores > 0,
      'Memory limit is positive': (info) => info.memory.limit_bytes > 0,
      'Load average is string': (info) => typeof info.cpu.load_average === 'string',
      'Free bytes is number': (info) => typeof info.memory.free_bytes === 'number',
      'Buffer bytes is number': (info) => typeof info.memory.buffer_bytes === 'number',
      'Cached bytes is number': (info) => typeof info.memory.cached_bytes === 'number',
      'Method is string': (info) => typeof info.method === 'string',
      'Fallback is boolean': (info) => typeof info.fallback === 'boolean',
    });
  }
  
  // Test 6: Simulate some CPU work to see usage change
  console.log('Simulating CPU work...');
  const startTime = Date.now();
  let sum = 0;
  for (let i = 0; i < 1000000; i++) {
    sum += Math.sqrt(i);
  }
  const workTime = Date.now() - startTime;
  console.log(`CPU work completed in ${workTime}ms (sum: ${sum.toFixed(2)})`);
  
  // Test 7: Check resource usage after work
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