package api

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tamnd/tama/pkg/store"
)

// rateLimiter is an in-memory token bucket per client IP: burst capacity of
// 10, refilling at 10 tokens per minute. Register and login share one
// limiter so a bot cannot alternate endpoints.
type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64 // tokens per second
	burst   float64
}

type bucket struct {
	tokens float64
	last   time.Time
}

func newRateLimiter(perMinute int) *rateLimiter {
	return &rateLimiter{
		buckets: make(map[string]*bucket),
		rate:    float64(perMinute) / 60,
		burst:   float64(perMinute),
	}
}

// allow spends one token for ip, or reports how long until one is available.
func (l *rateLimiter) allow(ip string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := store.Now()
	b, ok := l.buckets[ip]
	if !ok {
		b = &bucket{tokens: l.burst, last: now}
		l.buckets[ip] = b
		l.maybeSweep(now)
	}
	b.tokens = min(l.burst, b.tokens+now.Sub(b.last).Seconds()*l.rate)
	b.last = now

	if b.tokens >= 1 {
		b.tokens--
		return true, 0
	}
	wait := time.Duration((1 - b.tokens) / l.rate * float64(time.Second))
	return false, wait
}

// maybeSweep drops buckets idle long enough to be full again, keeping the
// map from growing with one entry per IP ever seen. Runs on the locked path.
func (l *rateLimiter) maybeSweep(now time.Time) {
	if len(l.buckets) < 1024 {
		return
	}
	idle := time.Duration(l.burst/l.rate) * time.Second
	for ip, b := range l.buckets {
		if now.Sub(b.last) > idle {
			delete(l.buckets, ip)
		}
	}
}

// limit wraps a handler with the token bucket; over-limit requests get a
// rate_limited envelope and a Retry-After header.
func (l *rateLimiter) limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		ok, wait := l.allow(ip)
		if !ok {
			w.Header().Set("Retry-After", itoaCeil(wait))
			WriteError(w, CodeRateLimited, "too many attempts, slow down")
			return
		}
		next(w, r)
	}
}

func itoaCeil(d time.Duration) string {
	secs := int(math.Ceil(d.Seconds()))
	if secs < 1 {
		secs = 1
	}
	return strconv.Itoa(secs)
}
