package ratelimiter

import (
    "sync"
    "time"
)

type TokenBucket struct {
    capacity float64
    fillRate float64
    tokens   float64
    last     time.Time
    mu       sync.Mutex
}

func NewTokenBucket(capacity int, rate int) *TokenBucket {
    return &TokenBucket{
        capacity: float64(capacity),
        fillRate: float64(rate),
        tokens:   float64(capacity),
        last:     time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.last).Seconds()
    if elapsed > 0 {
        tb.tokens += elapsed * tb.fillRate
        if tb.tokens > tb.capacity {
            tb.tokens = tb.capacity
        }
        tb.last = now
    }

    if tb.tokens >= 1 {
        tb.tokens -= 1
        return true
    }
    return false
}