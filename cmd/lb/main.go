package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/rahul-nakum14/go-loadbalancer/internal/backend"
)

var pool *backend.Pool

func main() {
	// 1. Initialize the backend pool with our two dummy servers
	pool = backend.NewPool([]string{
		"localhost:9001",
		"localhost:9002",
	})

	listener, err := net.Listen("tcp", ":8080")
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