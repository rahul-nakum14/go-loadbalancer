package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"github.com/rahul-nakum14/go-loadbalancer/internal/backend"
	"github.com/rahul-nakum14/go-loadbalancer/internal/config"
	"time"	
)

var pool *backend.Pool

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	    addresses := []string{}
	for _, backend := range cfg.Backends {
		addresses = append(addresses, backend.Address)
	}
	// 1. Initialize the backend pool with our two dummy servers
	pool = backend.NewPool(addresses)

	go pool.HealthCheck(time.Duration(cfg.HealthCheckInterval) * time.Second)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Load balancer listening on :8080")

	// 3. Accept incoming client connections forever
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
        go handleConnection(conn) // it will be one go routine per connection
	}
}

func handleConnection(clientConn net.Conn) {
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