package ratelimit

import (
	"fmt"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type RedisRts struct {
	limiter *redis_rate.Limiter
	limit   redis_rate.Limit
}

func NewRedisRts(db *redis.Client, rate, burst int, period time.Duration) (*RedisRts, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	if rate <= 0 || burst <= 0 {
		return nil, fmt.Errorf("rate and burst must be positive")
	}
	return &RedisRts{
		limiter: redis_rate.NewLimiter(db),
		limit: redis_rate.Limit{
			Rate:   rate,
			Burst:  burst,
			Period: period,
		},
	}, nil
}

func (limit *RedisRts) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		ok, err := limit.limiter.Allow(req.Context(), getIp(req), limit.limit)
		if err != nil {
			http.Error(resp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if ok.Allowed <= 0 {
			http.Error(resp, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(resp, req)
	})
}
