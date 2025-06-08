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

func (mock *serveMock) ServeHTTP(http.ResponseWriter, *http.Request) {
	if mock.racePanic {
		panic("panic")
	}
}

func TestRecover(t *testing.T) {
	h := &serveMock{racePanic: true}
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	assert.Panics(t, func() { h.ServeHTTP(resp, req) }, "should panic")
	hWithRecover := Recover(h)
	resp = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	assert.NotPanics(t, func() { hWithRecover.ServeHTTP(resp, req) }, "should panic")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
