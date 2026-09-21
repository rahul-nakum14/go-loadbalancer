package main

import (
    "fmt"
    "net"
    "log"
)

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()
    fmt.Println("Load balancer listening on :8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Accept error:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Fprintf(conn, "Hello from handle connection\n")
    fmt.Printf("New connection from %s\n", conn.RemoteAddr())
}