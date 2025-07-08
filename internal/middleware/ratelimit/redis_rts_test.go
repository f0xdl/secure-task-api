package ratelimit

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func rateLimiter(t *testing.T) *redis_rate.Limiter {
	mr := miniredis.RunT(t)

	rc := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	if err := rc.FlushDB(context.TODO()).Err(); err != nil {
		panic(err)
	}
	return redis_rate.NewLimiter(rc)
}

func TestRedisRts_Handle(t *testing.T) {
	rts := RedisRts{limiter: rateLimiter(t), limit: redis_rate.Limit{Rate: 1, Burst: 1, Period: time.Millisecond}}
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
	rts := RedisRts{limiter: rateLimiter(t), limit: redis_rate.Limit{Rate: 1, Burst: 1, Period: time.Millisecond}}
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

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	burst := 100
	clients := []string{"1.2.3.4", "127.0.0.1", "192.168.0.1"}

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		assert.Fail(t, err.Error())
	}
	defer rdb.FlushAll(context.Background())
	limiter, err := NewRedisRts(rdb, 10, burst, time.Minute)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	if limiter == nil {
		assert.Fail(t, "limiter is nil")
		return
	}
	handler := limiter.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	var wg sync.WaitGroup
	var mu sync.Mutex
	successes := 0
	failures := 0

	for _, addr := range clients {
		wg.Add(burst * 10)
		go func(ip string) {
			for i := 0; i < burst*10; i++ {
				go func() {
					defer wg.Done()
					req := httptest.NewRequest(http.MethodGet, "/", nil)
					req.RemoteAddr = net.JoinHostPort(ip, "12345")
					rec := httptest.NewRecorder()
					handler.ServeHTTP(rec, req)
					mu.Lock()
					defer mu.Unlock()
					if rec.Code == http.StatusOK {
						successes++
					} else if rec.Code == http.StatusTooManyRequests {
						failures++
					}
				}()
			}
		}(addr)
	}

	wg.Wait()
	assert.Equal(t, burst*len(clients), successes)
	assert.Equal(t, burst*len(clients)*10-burst*len(clients), failures)
}
