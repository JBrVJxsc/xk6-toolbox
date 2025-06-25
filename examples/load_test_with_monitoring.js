import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 5 },  // Ramp up to 5 users
    { duration: '1m', target: 5 },   // Stay at 5 users
    { duration: '30s', target: 0 },  // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    checks: ['rate>0.9'],            // At least 90% of checks should pass
  },
};

export default function () {
  // Monitor resources before making requests
  // const resourcesBefore = toolbox.getSystemInfo();
  
  // Make HTTP request
  const response = http.get('https://httpbin.org/delay/1');
  
  // Monitor resources after making requests
  // const resourcesAfter = toolbox.getSystemInfo();
  
  // Check response
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 2000ms': (r) => r.timings.duration < 2000,
  });
  
  // Log resource usage if it changed significantly
  // const cpuChange = resourcesAfter.cpu.usage_percent - resourcesBefore.cpu.usage_percent;
  // const memoryChange = resourcesAfter.memory.usage_percent - resourcesBefore.memory.usage_percent;
  
  // if (Math.abs(cpuChange) > 5 || Math.abs(memoryChange) > 2) {
  //   console.log(`VU ${__VU}: CPU change: ${cpuChange.toFixed(1)}%, Memory change: ${memoryChange.toFixed(1)}%`);
  // }
  
  // Check for resource constraints
  // if (resourcesAfter.cpu.usage_percent > 90) {
  //   console.warn(`VU ${__VU}: High CPU usage detected: ${resourcesAfter.cpu.usage_percent.toFixed(1)}%`);
  // }
  
  // if (resourcesAfter.memory.usage_percent > 90) {
  //   console.warn(`VU ${__VU}: High memory usage detected: ${resourcesAfter.memory.usage_percent.toFixed(1)}%`);
  // }
  
  sleep(1);
}

export function handleSummary(data) {
  console.log('=== Load Test Summary ===');
  console.log(`Total requests: ${data.metrics.http_reqs.values.count}`);
  console.log(`Failed requests: ${data.metrics.http_req_failed.values.passes}`);
  console.log(`Average response time: ${data.metrics.http_req_duration.values.avg.toFixed(2)}ms`);
  console.log(`95th percentile: ${data.metrics.http_req_duration.values['p(95)'].toFixed(2)}ms`);
  
  return {
    'load-test-results.json': JSON.stringify(data, null, 2),
  };
} 