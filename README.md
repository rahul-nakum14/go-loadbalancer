# go-loadbalancer

A high-performance **Layer 4 TCP Load Balancer** built from scratch in Go. Designed to demonstrate production-grade systems engineering, concurrent networking, lock-free data structures, active health monitoring, and defensive traffic controls.

---

## Architecture Overview

```
                                  +-----------------------+
                                  |     Clients / Postman |
                                  +-----------+-----------+
                                              |
                                              v  [Port :8080]
                                  +-----------+-----------+
                                  |  go-loadbalancer      |
                                  |                       |
                                  |  - Token Bucket Limit |
                                  |  - Atomic Round-Robin |
                                  |  - Active Health Check|
                                  +-----+-----------+-----+
                                        |           |
                        [TCP Proxy :9001|           |:9002]
                                        v           v
                                 +------+---+   +---+------+
                                 | Backend1 |   | Backend2 |
                                 +----------+   +----------+
```

---

## Key Features

- **Layer 4 TCP Proxying**: Direct bidirectional byte streaming via `io.Copy` with zero-copy overhead.
- **Goroutine Concurrency**: Non-blocking connection handling with lightweight goroutines (~2KB overhead per connection).
- **Lock-Free Round-Robin**: High-throughput target selection powered by `sync/atomic.Uint64` without mutex contention.
- **Active Health Monitoring**: Dedicated background goroutine (`time.Ticker`) continuously probing backend health over TCP with thread-safe state management (`sync.RWMutex`).
- **Token-Bucket Rate Limiter**: Per-IP traffic throttling using `golang.org/x/time/rate` to defend against flooding and denial-of-service attacks.
- **Graceful Shutdown**: Signal interception (`SIGINT`/`SIGTERM`) utilizing `context.NotifyContext` and `sync.WaitGroup` to ensure in-flight connections complete before termination.
- **Dynamic YAML Configuration**: Port, intervals, and backend topologies configured via `config.yaml`.

---

## Go Concurrency & Systems Concepts Implemented

| Concept | Implementation in Project |
|---|---|
| **Atomic Operations** | `atomic.Uint64` counter for lock-free round-robin load distribution. |
| **Worker Goroutines** | Independent health-checking daemon running concurrently with the listener. |
| **RWMutex** | Reader/Writer lock (`sync.RWMutex`) allowing simultaneous health reads with safe write isolation. |
| **Context & Cancellation** | `signal.NotifyContext` for graceful OS interrupt propagation. |
| **WaitGroups** | `sync.WaitGroup` tracking active connections during graceful draining. |
| **Buffered Channels & Worker Pools** | Used in benchmark test suite (`test_load.go`) for parallel load generation. |

---

## Performance & Stress Testing

Benchmarked using an internal concurrent worker pool generator (`test_load.go`):

```bash
======================================
Total Requests: 1000
Concurrency:    50 simultaneous workers
Duration:       ~1.06s
Status:         Defensive rate limiter actively throttling flood
======================================
```

---

## Quick Start

### 1. Prerequisites
- Go 1.22+ installed

### 2. Configuration (`config.yaml`)
```yaml
port: 8080
health_check_interval: 10
backends:
  - address: "localhost:9001"
  - address: "localhost:9002"
```

### 3. Running the Stack

**Terminal 1 — Backend 1:**
```bash
go run demoServer.go 9001
```

**Terminal 2 — Backend 2:**
```bash
go run demoServer.go 9002
```

**Terminal 3 — Load Balancer:**
```bash
go run cmd/lb/main.go
```

**Terminal 4 — Run Concurrent Load Test:**
```bash
go run test_load.go
```

---

## Resume Bullet Points

> - **High-Performance Layer 4 TCP Load Balancer (Go)**: Built a concurrent Layer 4 TCP load balancer proxying raw socket traffic using Go's `net` package and bidirectional `io.Copy` streaming.
> - **Lock-Free Concurrency & Thread-Safety**: Engineered a round-robin routing engine using `sync/atomic` primitives and `sync.RWMutex` to eliminate lock contention under concurrent load.
> - **Fault Tolerance & Reliability**: Implemented background health-check daemon using `time.Ticker` to dynamically isolate degraded nodes and recover healthy upstream servers.
> - **Defensive Traffic Throttling & Graceful Shutdown**: Integrated per-IP token bucket rate limiting (`golang.org/x/time/rate`) and signal-driven graceful draining using `context.NotifyContext` and `sync.WaitGroup`.
