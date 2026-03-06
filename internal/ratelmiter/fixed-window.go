package ratelimiter

import (
	"sync"
	"time"
)

const (
	cleanupEveryNRequests = 256
	cleanupMinEntries     = 64
)

type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients  map[string]*clientWindow
	limit    int
	window   time.Duration
	requests uint64
}

type clientWindow struct {
	count       int
	windowStart time.Time
}

func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]*clientWindow),
		limit:   limit,
		window:  window,
	}
}

func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	now := time.Now()

	rl.Lock()
	rl.requests++
	if rl.shouldCleanup() {
		rl.cleanupExpired(now)
	}

	client, exists := rl.clients[ip]
	defer rl.Unlock()

	if rl.limit <= 0 {
		return false, rl.window
	}

	if !exists {
		rl.clients[ip] = &clientWindow{
			count:       1,
			windowStart: now,
		}
		return true, 0
	}

	if now.Sub(client.windowStart) >= rl.window {
		client.count = 1
		client.windowStart = now
		return true, 0
	}

	if client.count < rl.limit {
		client.count++
		return true, 0
	}

	retryAfter := client.windowStart.Add(rl.window).Sub(now)
	if retryAfter < 0 {
		retryAfter = 0
	}

	return false, retryAfter
}

func (rl *FixedWindowRateLimiter) shouldCleanup() bool {
	return len(rl.clients) >= cleanupMinEntries && rl.requests%cleanupEveryNRequests == 0
}

func (rl *FixedWindowRateLimiter) cleanupExpired(now time.Time) {
	for ip, client := range rl.clients {
		if now.Sub(client.windowStart) >= rl.window {
			delete(rl.clients, ip)
		}
	}
}
