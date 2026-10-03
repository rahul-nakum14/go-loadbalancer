# Go-loadbalancer 🚀

A TCP load balancer built from scratch in Go. similar to what Nginx — Using Go's standard `net` package.

## What it does

- Routes traffic across multiple backend servers using **round-robin**
- Runs **health checks** in the background and skips dead backends automatically
- **Rate limits** requests per IP using a token bucket (10 req/s, burst of 20)
- **Graceful shutdown** — waits for active connections to finish before exiting
- All config (port, backends, intervals) lives in `config.yaml`

## Go Concepts Used

- **Goroutines** — one goroutine per connection + dedicated health check worker running in background
- **Atomic operations** (`sync/atomic`) — lock-free round-robin counter for high-throughput routing
- **RWMutex** — multiple goroutines can safely read backend health status simultaneously, with exclusive lock only on writes
- **Mutex** — protects the per-IP rate limiter map from concurrent writes
- **WaitGroup** — tracks in-flight connections during graceful shutdown
- **Buffered channels** — used in load tester as a work queue distributed across 50 worker goroutines
- **Context & OS signals** — `signal.NotifyContext` listens for `Ctrl+C` and propagates cancellation

## Project Structure

```
cmd/
  lb/main.go           # TCP listener + proxy logic
  tester/main.go       # Concurrent load tester (50 workers, 1000 requests)
internal/
  backend/backend.go   # Backend pool, round-robin, health checks
  config/config.go     # YAML config parser
  ratelimiter/         # Per-IP token bucket rate limiting
config.yaml            # Port + backend addresses
demoServer.go          # Dummy HTTP server for testing
```

## Run it

```bash
# Start two backends
go run demoServer.go 9001
go run demoServer.go 9002

# Start the load balancer
go run cmd/lb/main.go

# Hit it
curl http://localhost:8080

# Stress test
go run cmd/tester/main.go
```

## Config

```yaml
port: 8080
health_check_interval: 10
backends:
  - address: "localhost:9001"
  - address: "localhost:9002"
```
