package ratelimit

import (
	"net/http"
	"sync"
	"time"
)

type CustomTokenBucketLimit struct {
	mx        sync.Mutex
	tokens    uint
	maxTokens uint
	rts       uint

	refilledAt time.Time
}

// NewCustomTokenBucketLimit create a new Token Bucket instance
// maxTokens - max available tokens
// rts - refill rate per second
func NewCustomTokenBucketLimit(maxTokens uint, rts uint) *CustomTokenBucketLimit {
	s := &CustomTokenBucketLimit{
		maxTokens:  maxTokens,
		tokens:     maxTokens,
		rts:        rts,
		refilledAt: time.Now(),
	}
	return s
}

func (limit *CustomTokenBucketLimit) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if !limit.available() {
			http.Error(resp, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(resp, req)
	})
}

func (limit *CustomTokenBucketLimit) available() bool {
	limit.mx.Lock()
	limit.refill()
	if limit.tokens != 0 {
		limit.tokens--
		limit.mx.Unlock()
		return true
	} else {
		limit.mx.Unlock()
		return false
	}
}

func (limit *CustomTokenBucketLimit) refill() {
	now := time.Now()
	addedTokens := uint(now.Sub(limit.refilledAt).Seconds()) * limit.rts
	if addedTokens > 0 {
		limit.tokens += addedTokens
		if limit.tokens > limit.maxTokens {
			limit.tokens = limit.maxTokens
		}
		limit.refilledAt = now
	}
}
