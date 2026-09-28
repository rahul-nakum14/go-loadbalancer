package ratelimiter

import (
    "sync"
    "golang.org/x/time/rate"
)

type IPRateLimiter struct {
    mu       sync.Mutex
    limiters map[string]*rate.Limiter
}

func New() *IPRateLimiter {
    return &IPRateLimiter{
        limiters: make(map[string]*rate.Limiter),
    }
}

func (i *IPRateLimiter) Allow(ip string) bool {
    i.mu.Lock()
    defer i.mu.Unlock()

    if _, exists := i.limiters[ip]; !exists {
        i.limiters[ip] = rate.NewLimiter(10, 20) 
    }

    return i.limiters[ip].Allow() 
}