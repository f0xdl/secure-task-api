package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type serveMock struct {
	racePanic bool
}

func (mock *serveMock) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if mock.racePanic {
		panic("panic")
	}
	resp.WriteHeader(http.StatusOK)
}

func TestRecover_PanicWithoutMiddleware(t *testing.T) {
	h := &serveMock{racePanic: true}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	assert.Panics(t, func() { h.ServeHTTP(resp, req) }, "should panic")
}
func TestRecover_PanicWithMiddleware(t *testing.T) {
	h := &serveMock{racePanic: true}
	hWithRecover := Recover(h)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	assert.NotPanics(t, func() { hWithRecover.ServeHTTP(resp, req) }, "should panic")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestRecover_NoPanicWithMiddleware(t *testing.T) {
	h := &serveMock{racePanic: false}
	hWithRecover := Recover(h)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	assert.NotPanics(t, func() { hWithRecover.ServeHTTP(resp, req) }, "should panic")
	assert.Equal(t, http.StatusOK, resp.Code)
}
