package backend

import (
    "sync/atomic"
    "sync"
    "time"
    "log"
    "net"
)

type Backend struct {
    alive   bool 
	Address string
    mu      sync.RWMutex
}

type Pool struct {
    backends []*Backend
    counter  atomic.Uint64
}

// NewPool initializes a pool with backend addresses
// It will take a list of backe	d Addresses and create a pool of backends and return a pointer.......
func NewPool(addresses []string) *Pool {
    p := &Pool{}
    for _, addr := range addresses {
        b := &Backend{Address: addr}
        b.SetAlive(true)
        p.backends = append(p.backends, b)
    }
    return p
}

// It will basically picks the next backend in a round-robin rotation and alive which backend server is alive else it 
// will return the nill :)
func (p *Pool) Next() *Backend {
    total := uint64(len(p.backends))
    
    for i := uint64(0); i < total; i++ {
        idx := p.counter.Add(1) - 1
        backend := p.backends[idx%total]
        if backend.Alive() {
            return backend
        }
    }
    return nil // all backends dead
}

//it will set the status backend alive or dead based on Health checkss.
func (b *Backend) SetAlive(alive bool) {
    b.mu.Lock()
    b.alive = alive
    b.mu.Unlock()
}

func (b *Backend) Alive() bool {
    b.mu.RLock()
    defer b.mu.RUnlock()
    return b.alive
}


func (p *Pool) HealthCheck(interval time.Duration) {
    ticker := time.NewTicker(interval) // here ticker is like a cron which runs like given interval
    defer ticker.Stop()
    for range ticker.C { // here ticker.C is a channel coz ticker return us a channel
        for _, b := range p.backends {
            conn, err := net.DialTimeout("tcp", b.Address, 2*time.Second)
            if err != nil {
                b.SetAlive(false)
                log.Printf("Backend %s is DOWN", b.Address)
            } else {
                b.SetAlive(true)
                conn.Close()
                log.Printf("Backend %s is UP", b.Address)
            }
        }
    }
}