# Proxy

**Proxy** is a simple caching proxy created in my free time. It focuses on speed and simplicity.

## Features

- 🔨 **Logger**: Uses [zap](https://github.com/uber-go/zap) logger for structured, high-performance logging
- 🏎️ **Cache**: Redis to cache responses
- 🌐 **Proxying**: Forward requests seamlessly to the origin server

## Installation

```bash
git clone https://github.com/pyr33x/proxy.git
cd proxy
go build -o proxy
```

## Usage

You can set environment variables to configure the proxy before running it.

### Environment Variables

| Variable         | Description                         | Default Value           |
|-----------------|-------------------------------------|------------------------|
| `ORIGIN_URL`     | The origin server URL to proxy       | `http://dummyjson.com` |
| `PROXY_PORT`     | Port for the proxy server to listen | `1337`                 |
| `ZAP_ENV`        | Logging environment (dev/prod)     | `dev`                  |
| `REDIS_HOST`     | Redis server host                  | `127.0.0.1`            |
| `REDIS_PORT`     | Redis server port                  | `6379`                 |
| `REDIS_USERNAME` | Redis username                     | `proxy`                |
| `REDIS_PASSWORD` | Redis password                     | *(empty)*              |
| `REDIS_DATABASE` | Redis database number              | `0`                    |

## Run Proxy

```bash
# Run with default environment variables
./proxy

# Or set custom variables inline
ORIGIN_URL="http://example.com" PROXY_PORT="8080" ./proxy
```
Cached responses will be stored in Redis and served quickly on repeated requests.

## Docker
You can also run Proxy using Docker.

### Build and Run
```bash
# Build the Docker image
docker build -t proxy .

# Run the container
docker run -d \
  -e ORIGIN_URL="http://example.com" \
  -e PROXY_PORT="8080" \
  -e REDIS_HOST="redis" \
  -p 8080:8080 \
  proxy
```

## Pull from Docker Hub
```bash
docker pull pyr33x/proxy:latest
docker run -d -p 8080:8080 pyr33x/proxy:latest
```

## Docker Compose
If you have Docker Compose, you can spin up Proxy with Redis easily.
Run with:
```bash
docker compose up -d --build
```

## Contributing
Feel free to open issues or submit pull requests. Any help is appreciated!

## License
MIT @ Mehdi Parandak
