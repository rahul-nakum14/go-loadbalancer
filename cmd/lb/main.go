package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "net"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "github.com/rahul-nakum14/go-loadbalancer/internal/backend"
    "github.com/rahul-nakum14/go-loadbalancer/internal/config"
)
var pool *backend.Pool
var wg sync.WaitGroup


func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	addresses := []string{}
	for _, backend := range cfg.Backends {
		addresses = append(addresses, backend.Address)
	}

	pool = backend.NewPool(addresses)

	go pool.HealthCheck(time.Duration(cfg.HealthCheckInterval) * time.Second)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Load balancer listening on :8080")

	// 3. Accept incoming client connections forever

    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                select {
                case <-ctx.Done(): 
                    return
                default:
                    log.Println("Accept error:", err)
                    continue
                }
            }
            wg.Add(1)
            go handleConnection(conn)
        }
    }()

    <-ctx.Done()
    log.Println("Shutting down... waiting for active connections to finish") // gaceful shutdown
    listener.Close()
    wg.Wait()      
}

func handleConnection(clientConn net.Conn) {
	defer wg.Done()
	defer clientConn.Close()

	target := pool.Next()
	if target == nil {
		log.Println("No alive backends!")
		clientConn.Close()
		return
	}
	
	log.Printf("Routing connection to: %s", target.Address)

	backendConn, err := net.Dial("tcp", target.Address)
	if err != nil {
        log.Println("Backend error:", err)
		return
	}
	defer backendConn.Close()

	go io.Copy(backendConn, clientConn)
	io.Copy(clientConn, backendConn)
}