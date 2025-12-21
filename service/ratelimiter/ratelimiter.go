package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity      int           // 桶最大容量
	tokens        int           // 当前令牌数
	rate          int           // 每秒生成的令牌数
	lastTimestamp time.Time     // 上一次补充令牌的时间
	mu            sync.Mutex
}

// NewTokenBucket 创建一个令牌桶
func NewTokenBucket(capacity int, rate int) *TokenBucket {
	return &TokenBucket{
		capacity:      capacity,
		tokens:        capacity,   // 初始令牌满桶
		rate:          rate,
		lastTimestamp: time.Now(),
	}
}

// Allow 尝试获取 1 个令牌，成功返回 true，否则 false
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 补充令牌
	now := time.Now()
	elapsed := now.Sub(tb.lastTimestamp).Seconds()
	if elapsed > 0 {
		// 新增令牌数 = 速率 * 时间
		newTokens := int(float64(tb.rate) * elapsed)
		if newTokens > 0 {
			tb.tokens = min(tb.capacity, tb.tokens+newTokens)
			tb.lastTimestamp = now
		}
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
