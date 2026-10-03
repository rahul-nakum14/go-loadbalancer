# Go TCP Load Balancer

A lightweight Layer 4 TCP load balancer built from scratch in Go. 

I’ve been working as a backend engineer primarily with Node.js and NestJS, and I built this project to dive deep into Go's concurrency model, networking internals, and systems programming—building something from the ground up rather than just configuring Nginx.

---

## What It Does

```
Client (curl / Postman) 
       │
       ▼ (:8080)
┌──────────────────────────────────────────────┐
│             Go Load Balancer                 │
│                                              │
│  1. Rate Limiting (Token Bucket per IP)      │
│  2. Round-Robin Selection (Atomic counter)   │
│  3. Background Health Check (Ping backends)  │
└──────────────────────────────────────────────┘
       │                      │
       ▼ (:9001)              ▼ (:9002)
   Backend 1              Backend 2
```

1. **TCP Proxying**: Listens on `:8080` and opens a direct TCP tunnel to backend servers, forwarding incoming and outgoing bytes concurrently using `io.Copy`.
2. **Round-Robin Routing**: Distributes requests evenly across healthy backend instances using a thread-safe `atomic.Uint64` counter (no lock contention).
3. **Active Health Checks**: A background goroutine periodically pings each backend every few seconds. If a backend goes down, traffic automatically skips it until it recovers.
4. **Rate Limiting**: Protects backends from floods using a token-bucket algorithm per client IP (`golang.org/x/time/rate`).
5. **Graceful Shutdown**: Listens for `SIGINT` / `SIGTERM` (`Ctrl+C`), stops accepting new connections, and waits for all active connections to finish via `sync.WaitGroup` before exiting.
6. **YAML Config**: Reads port, health check intervals, and upstream backend addresses from `config.yaml`.

---

## Project Structure

```
├── cmd/
│   ├── lb/
│   │   └── main.go           # Entry point, TCP listener & shutdown logic
│   └── tester/
│       └── main.go           # Concurrency stress tester using worker pools
├── internal/
│   ├── backend/
│   │   └── backend.go        # Backend pool, round-robin & health checking
│   ├── config/
│   │   └── config.go         # YAML config parsing
│   └── ratelimiter/
│       └── limiter.go        # Per-IP token bucket rate limiting
├── config.yaml               # Server & backend configuration
└── demoServer.go             # Simple HTTP server to simulate backends
```

---

## How to Run It

### 1. Start two test backend servers
Open two terminal windows:

```bash
# Terminal 1
go run demoServer.go 9001

# Terminal 2
go run demoServer.go 9002
```

### 2. Start the Load Balancer
```bash
# Terminal 3
go run cmd/lb/main.go
```

Now send requests to the load balancer:
```bash
curl http://localhost:8080
```
Hit it a few times — you will see it alternate between `9001` and `9002`. If you kill one backend (`Ctrl+C`), the load balancer will detect it and send 100% of the traffic to the surviving one!

### 3. Run the stress test
```bash
# Terminal 4
go run cmd/tester/main.go
```
This spawns 50 concurrent worker goroutines pushing 1,000 requests to test throughput and verify the rate limiter.

