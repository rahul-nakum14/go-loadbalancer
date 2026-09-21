package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    port := os.Args[1] // taking args from command line
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Response from backend on port we got%s\n", port)
    })
    http.ListenAndServe(":"+port, nil)
}