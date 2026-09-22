package backend

import "sync/atomic" 

type Backend struct {
	Address string
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
        p.backends = append(p.backends, &Backend{Address: addr})
    }
    return p
}

// It will basically picks the next backend in a round-robin rotation :)
func (p *Pool) Next() *Backend {
    idx := p.counter.Add(1) - 1
    return p.backends[idx%uint64(len(p.backends))]
}