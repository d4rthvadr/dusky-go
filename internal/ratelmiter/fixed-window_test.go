package ratelimiter

import (
	"fmt"
	"testing"
	"time"
)

func TestFixedWindowRateLimiter_AllowsUpToLimitThenDenies(t *testing.T) {
	rl := NewFixedWindowRateLimiter(2, time.Second)

	allowed, retryAfter := rl.Allow("127.0.0.1")
	if !allowed || retryAfter != 0 {
		t.Fatalf("first request should be allowed with zero retryAfter, got allowed=%v retryAfter=%v", allowed, retryAfter)
	}

	allowed, retryAfter = rl.Allow("127.0.0.1")
	if !allowed || retryAfter != 0 {
		t.Fatalf("second request should be allowed with zero retryAfter, got allowed=%v retryAfter=%v", allowed, retryAfter)
	}

	allowed, retryAfter = rl.Allow("127.0.0.1")
	if allowed {
		t.Fatalf("third request should be denied")
	}
	if retryAfter <= 0 || retryAfter > time.Second {
		t.Fatalf("retryAfter should be within (0, window], got %v", retryAfter)
	}
}

func TestFixedWindowRateLimiter_ResetsAfterWindow(t *testing.T) {
	window := 40 * time.Millisecond
	rl := NewFixedWindowRateLimiter(1, window)

	allowed, _ := rl.Allow("127.0.0.1")
	if !allowed {
		t.Fatalf("first request should be allowed")
	}

	allowed, _ = rl.Allow("127.0.0.1")
	if allowed {
		t.Fatalf("second request should be denied before window expiry")
	}

	time.Sleep(window + 20*time.Millisecond)

	allowed, retryAfter := rl.Allow("127.0.0.1")
	if !allowed {
		t.Fatalf("request should be allowed after window reset, retryAfter=%v", retryAfter)
	}
}

func TestFixedWindowRateLimiter_DifferentIPsAreIsolated(t *testing.T) {
	rl := NewFixedWindowRateLimiter(1, time.Second)

	allowedA, _ := rl.Allow("10.0.0.1")
	allowedB, _ := rl.Allow("10.0.0.2")
	if !allowedA || !allowedB {
		t.Fatalf("first request from each IP should be allowed")
	}

	allowedA, _ = rl.Allow("10.0.0.1")
	allowedB, _ = rl.Allow("10.0.0.2")
	if allowedA || allowedB {
		t.Fatalf("second request in same window should be denied independently per IP")
	}
}

func TestFixedWindowRateLimiter_CleansUpExpiredClients(t *testing.T) {
	window := 20 * time.Millisecond
	rl := NewFixedWindowRateLimiter(1, window)

	// Create more than cleanupMinEntries unique clients.
	for i := 0; i < cleanupMinEntries+8; i++ {
		ip := fmt.Sprintf("10.0.0.%d", i)
		allowed, _ := rl.Allow(ip)
		if !allowed {
			t.Fatalf("initial request for %s should be allowed", ip)
		}
	}

	time.Sleep(window + 10*time.Millisecond)

	// Trigger periodic cleanup by reaching the cleanup cadence.
	for i := 0; i < cleanupEveryNRequests; i++ {
		ip := fmt.Sprintf("172.16.0.%d", i)
		rl.Allow(ip)
	}

	rl.RLock()
	defer rl.RUnlock()
	if len(rl.clients) > cleanupEveryNRequests {
		t.Fatalf("expected expired clients to be cleaned up, current size=%d", len(rl.clients))
	}
}

func TestFixedWindowRateLimiter_NonPositiveLimitAlwaysDenies(t *testing.T) {
	rl := NewFixedWindowRateLimiter(0, 100*time.Millisecond)

	allowed, retryAfter := rl.Allow("127.0.0.1")
	if allowed {
		t.Fatalf("request should be denied when limit is non-positive")
	}
	if retryAfter != 100*time.Millisecond {
		t.Fatalf("retryAfter should equal configured window, got %v", retryAfter)
	}
}

func TestFixedWindowRateLimiter_DoesNotCleanupBelowMinEntries(t *testing.T) {
	window := 20 * time.Millisecond
	rl := NewFixedWindowRateLimiter(1, window)

	for i := 0; i < cleanupMinEntries-1; i++ {
		ip := fmt.Sprintf("192.168.1.%d", i)
		allowed, _ := rl.Allow(ip)
		if !allowed {
			t.Fatalf("initial request for %s should be allowed", ip)
		}
	}

	time.Sleep(window + 10*time.Millisecond)

	// Force the request counter close to cleanup cadence and make one more call.
	rl.Lock()
	rl.requests = cleanupEveryNRequests - 1
	rl.Unlock()

	allowed, _ := rl.Allow("203.0.113.10")
	if !allowed {
		t.Fatalf("request should be allowed")
	}

	rl.RLock()
	defer rl.RUnlock()
	if len(rl.clients) != cleanupMinEntries {
		t.Fatalf("expected no cleanup below minimum entries, current size=%d", len(rl.clients))
	}
}
