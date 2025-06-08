package ratelimit

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func rateLimiter() *redis_rate.Limiter {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{"server0": ":6379"},
	})
	if err := ring.FlushDB(context.TODO()).Err(); err != nil {
		panic(err)
	}
	return redis_rate.NewLimiter(ring)
}

func TestRedisRts_Handle(t *testing.T) {
	rts := RedisRts{limiter: rateLimiter(), limit: redis_rate.Limit{Rate: 1, Burst: 1, Period: time.Millisecond}}
	h := rts.Handle(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
	}))
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(resp, req)
	assert.Equal(t, 429, resp.Code)
}

func TestRedisRts_Handle_Period(t *testing.T) {
	rts := RedisRts{limiter: rateLimiter(), limit: redis_rate.Limit{Rate: 1, Burst: 1, Period: time.Millisecond}}
	h := rts.Handle(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)
	}))
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	time.Sleep(time.Millisecond)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
}
