# working-memory - In-Memory Caching Extension for k6

`working-memory` is an in-memory caching extension for [k6](https://k6.io/), a modern load testing tool. This extension allows k6 scripts to store, retrieve, and manage temporary data across virtual users (VUs), enabling scenarios where data persistence and caching are required.

## Installation

To use the `working-memory` extension, you need to build k6 with this extension.

```bash
xk6 build --with github.com/gera-cl/xk6-working-memory@latest
```

> Note: Replace `github.com/gera-cl/xk6-working-memory` with the actual repository URL for this package.

## Usage

After building k6 with the `working-memory` extension, you can import and use it in your k6 scripts.

### Importing the Extension

```javascript
import memory from 'k6/x/working-memory';
```

### API

The `working-memory` extension exposes the following methods:

- `init(defaultExpiration, cleanupInterval)`: Initializes the cache with a default expiration time and cleanup interval in seconds.
- `set(id, value, expiration)`: Stores a value in the cache with an optional expiration time in seconds.
- `get(id)`: Retrieves a value from the cache by its identifier.

### Example Script

Here's an example script demonstrating how to use the `memory` extension in k6.

```javascript
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
```

### Explanation

1. **Setup**: The `memory.init()` function is called in the `setup` stage to initialize the cache with a 60-second default expiration and a 120-second cleanup interval.
2. **Check Cache**: In the `default` function, the script checks if the URL's response is already cached using `memory.get()`. If cached, it logs the cached response.
3. **Fetch and Store in Cache**: If the response is not in the cache, an HTTP request is made, and the response is stored in the cache with a 30-second expiration.

### Running the Example

To run the example script:

```bash
k6 run script.js
```

Replace `script.js` with the name of your k6 script file.

## License

MIT License
