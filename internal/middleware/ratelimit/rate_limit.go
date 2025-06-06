package ratelimit

import (
	"golang.org/x/time/rate"
	"net/http"
	"strings"
	"sync"
)

type RateIpLimit struct {
	clients   sync.Map // map[string]*rate.Limiter
	maxTokens int
	rts       float64
}

func NewRateIpLimit(maxTokens int, rts float64) *RateIpLimit {
	return &RateIpLimit{
		maxTokens: maxTokens,
		rts:       rts,
		clients:   sync.Map{},
	}
}

func (limit *RateIpLimit) getLimiter(ip string) *rate.Limiter {
	l, ok := limit.clients.Load(ip)
	if !ok {
		l = rate.NewLimiter(rate.Limit(limit.rts), limit.maxTokens)
		limit.clients.Store(ip, l)
	}
	return l.(*rate.Limiter)
}

func getIp(req *http.Request) string {
	fwdAddr := req.Header.Get("X-Forwarded-For")
	if fwdAddr == "" {
		fwdAddr = strings.Split(req.RemoteAddr, ":")[0]
	}
	return fwdAddr
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
