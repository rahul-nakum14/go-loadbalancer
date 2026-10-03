package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	totalRequests := 1000
	concurrency := 50

	ch := make(chan int, totalRequests)
	for i := 0; i < totalRequests; i++ {
		ch <- i
	}
	close(ch)

	var wg sync.WaitGroup
	start := time.Now()

	successCount := 0
	var mu sync.Mutex

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{Timeout: 2 * time.Second}

			for range ch {
				resp, err := client.Get("http://localhost:8080")
				if err == nil && resp.StatusCode == 200 {
					mu.Lock()
					successCount++
					mu.Unlock()
					resp.Body.Close()
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Println("======================================")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Successful:     %d\n", successCount)
	fmt.Printf("Time Taken:     %v\n", duration)
	fmt.Printf("Requests/sec:   %.2f\n", float64(successCount)/duration.Seconds())
	fmt.Println("======================================")
}