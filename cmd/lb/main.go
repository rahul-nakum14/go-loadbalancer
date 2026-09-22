package main

import (
    "fmt"
    "net"
    "log"
    "io"
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
        go handleConnection(conn) // one go routine per connection
    }
}

func handleConnection(clientConn net.Conn) {
    defer clientConn.Close()

    // open a new connecction to the backend server
    backendConn, err := net.Dial("tcp", "localhost:9001")
    if err != nil {
        log.Println("Backend error:", err)
        return
    }
    defer backendConn.Close()

    go io.Copy(backendConn, clientConn)
    io.Copy(clientConn, backendConn)
}