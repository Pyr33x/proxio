# Proxio

**Proxio** is a simple caching proxy created in my free time. It focuses on speed and simplicity.

## Features

- 🔨 **Logger**: Uses [zap](https://github.com/uber-go/zap) logger for structured, high-performance logging
- 🏎️ **Cache**: Redis to cache responses
- 🌐 **Proxying**: Forward requests seamlessly to the origin server

## Installation

```bash
git clone https://github.com/pyr33x/proxio.git
cd proxio
go build -o proxio
```

## Usage

You can set environment variables to configure the proxy before running it.

### Environment Variables

| Variable         | Description                         | Default Value           |
|-----------------|-------------------------------------|------------------------|
| `ORIGIN_URL`     | The origin server URL to proxy       | `http://dummyjson.com` |
| `PROXY_PORT`     | Port for the proxy server to listen | `1337`                 |
| `PROXY_CACHE`    | Cache storage (redis/memory) | `memory` |
| `PROXY_TTL`      | Cache TTL (seconds) | `60` |
| `ZAP_ENV`        | Logging environment (dev/prod)     | `dev`                  |
| `REDIS_HOST`     | Redis server host (used only if cache = redis)               | `redis_proxio`            |
| `REDIS_PORT`     | Redis server port                  | `6379`                 |
| `REDIS_USERNAME` | Redis username                     | `proxio`                |
| `REDIS_PASSWORD` | Redis password                     | *(empty)*              |
| `REDIS_DATABASE` | Redis database number              | `0`                    |

### Notes on Caching
- **Default cache storage is `memory`**
- When using memory cache, **you do not need a Redis instance**
- Switch to Redis by setting:
```env
PROXY_CACHE="redis"
```

## Run Proxio

```bash
# Run with default environment variables
./proxio

# Or set custom variables inline
ORIGIN_URL="http://example.com" \
PROXY_PORT="8080" \
PROXY_CACHE="redis" \
PROXY_TTL="120" \
./proxio
```
Cached responses will be stored in Redis and served quickly on repeated requests.

## Docker
You can also run Proxio using Docker.

### Build and Run
```bash
# Build the Docker image
docker build -t proxio .

# Run the container
docker run -d \
  -e ORIGIN_URL="http://example.com" \
  -e PROXY_PORT="8080" \
  -p 8080:8080 \
  proxio
```

If using redis:
```bash
docker run -d \
  -e PROXY_CACHE="redis" \
  -e REDIS_HOST="redis" \
  -p 8080:8080 \
  proxio
```

## Using Container Registries
```bash
# Pull from Docker Hub
docker pull me3di/proxio:latest

# Pull from GitHub Container Registry (GHCR)
docker pull ghcr.io/pyr33x/proxio:latest

# Or build locally
docker build -t proxio .
docker run -d -p 8080:8080 ghcr.io/pyr33x/proxio:latest
```

## Docker Compose
If you have Docker Compose, you can spin up Proxio with Redis easily.
Run with:
```bash
docker compose up -d --build
```
If `PROXY_CACHE=memory`, Redis won't be used even if it’s running.

## Contributing
Feel free to open issues or submit pull requests. Any help is appreciated!

## License
MIT @ Mehdi Parandak
