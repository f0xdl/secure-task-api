package ratelimit

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type RateIpLimit struct {
	clients   sync.Map // map[string]*rate.Limiter
	maxTokens int
	rts       float64
}

func NewRateIpLimit(maxTokens int, rts float64) (*RateIpLimit, error) {
	if maxTokens <= 0 || rts <= 0 {
		return nil, fmt.Errorf("maxTokens and rts must be positive")
	}
	return &RateIpLimit{
		maxTokens: maxTokens,
		rts:       rts,
		clients:   sync.Map{},
	}, nil
}

func (limit *RateIpLimit) getLimiter(ip string) *rate.Limiter {
	l, ok := limit.clients.Load(ip)
	if !ok {
		l = rate.NewLimiter(rate.Limit(limit.rts), limit.maxTokens)
		limit.clients.Store(ip, l)
	}
	return l.(*rate.Limiter)
}

func (limit *RateIpLimit) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if !limit.getLimiter(getIp(req)).Allow() {
			http.Error(resp, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(resp, req)
	})
}
