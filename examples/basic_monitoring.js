import { sleep } from 'k6';

export const options = {
  vus: 1,
  duration: '30s',
};

export default function () {
  // Log current resource usage
  console.log(`[${new Date().toISOString()}] Resource Usage:`);
  console.log(`  CPU: ${info.cpu.used_cores.toFixed(2)}/${info.cpu.limit_cores.toFixed(2)} cores (${info.cpu.usage_percent.toFixed(1)}%)`);
  console.log(`  Memory: ${info.memory.usage_mb.toFixed(1)}/${info.memory.limit_mb.toFixed(1)} MB (${info.memory.usage_percent.toFixed(1)}%)`);
  
  // Check for resource warnings
  if (info.cpu.usage_percent > 80) {
    console.warn(`⚠ High CPU usage: ${info.cpu.usage_percent.toFixed(1)}%`);
  }
  
  if (info.memory.usage_percent > 85) {
    console.warn(`⚠ High memory usage: ${info.memory.usage_percent.toFixed(1)}%`);
  }
  
  sleep(5);
} 