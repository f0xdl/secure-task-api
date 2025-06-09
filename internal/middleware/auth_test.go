package middleware

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthMock struct {
	users map[string]string
}

func (auth *AuthMock) Allow(u, p string) bool {
	pass, ok := auth.users[u]
	return ok && pass == p
}

func TestBasicAuth_EmptyArgs(t *testing.T) {
	_, err := NewBasicAuth(nil)
	assert.Error(t, err)
}

func TestBasicAuth_FillArgs(t *testing.T) {
	a := &AuthMock{}
	_, err := NewBasicAuth(a)
	assert.Nil(t, err)
}

func TestBasicAuth_WrongUsers(t *testing.T) {
	correctUsers := map[string]string{"user1": "pass1", "user2": "pass2"}
	wrongUsers := map[string]string{"": "", "user1": "", "user2": "PASS", "user3": "pass3", "user4": "pass4"}

	a := &AuthMock{users: correctUsers}
	mw, err := NewBasicAuth(a)
	assert.Nil(t, err)
	h := mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	for user, pass := range wrongUsers {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth(user, pass)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		assert.Equal(t, 401, resp.Code)
	}
}

func TestBasicAuth_CorrectUsers(t *testing.T) {
	users := map[string]string{"user1": "pass1", "user2": "pass2"}

	a := &AuthMock{users: users}
	mw, err := NewBasicAuth(a)
	assert.Nil(t, err)
	h := mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Error(t, errors.New("not allowed"))
	}))

	for user, pass := range users {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth(user, pass)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	}
}
