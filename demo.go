package main

import (
	"axe-backend/service/ratelimiter"
	"fmt"
	"time"
)

func main() {
	// 创建令牌桶：每秒生成 5 个令牌，最大容量 10
	tb := ratelimiter.NewTokenBucket(10, 5)

	for i := 0; i < 20; i++ {
		if tb.Allow() {
			fmt.Println(i, "=> pass")
		} else {
			fmt.Println(i, "=> limit")
		}
		time.Sleep(20 * time.Millisecond)
	}
}
