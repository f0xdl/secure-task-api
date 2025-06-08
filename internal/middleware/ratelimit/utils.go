package ratelimit

import (
	"net/http"
	"strings"
)

type LimitHandler interface {
	Handle(next http.Handler) http.Handler
}

func getIp(req *http.Request) string {
	fwdAddr := req.Header.Get("X-Forwarded-For")
	if fwdAddr == "" {
		fwdAddr = strings.Split(req.RemoteAddr, ":")[0]
	}
	return fwdAddr
}
