import http from 'k6/http';
import { sleep } from 'k6';
import memory from 'k6/x/working-memory';

export const options = {
    vus: 10,
    duration: '30s',
};

export function setup() {
    // Initialize the cache with a 60-second default expiration and 120-second cleanup interval.
    memory.init(60, 120);
}

export default function () {
    const url = 'https://httpbin.test.k6.io/get';

    // Check if the URL is already cached
    let cachedResponse = memory.get(url);
    if (cachedResponse) {
        console.log(`Cache hit: ${cachedResponse}`);
    } else {
        console.log('Cache miss, making HTTP request...');
        const response = http.get(url);
        
        // Cache the response for 30 seconds
        memory.set(url, response.body, 30);
    }

    // Simulate test execution delay
    sleep(1);
}