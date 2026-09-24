package backend

import (
    "sync/atomic"
    "sync"
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