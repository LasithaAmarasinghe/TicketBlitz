import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 10,        // Let's start with 10. This is enough to cause a race condition.
  duration: '5s',
};

export default function () {
  // Use explicit IPv4 IP to avoid "localhost" confusion
  const res = http.post('http://127.0.0.1:9090/buy');
  
  // If it fails, print the error to the terminal so we can debug
  if (res.status !== 200) {
     console.warn(`Error: ${res.status} ${res.body}`);
  }

  check(res, { 'is status 200': (r) => r.status === 200 });
}