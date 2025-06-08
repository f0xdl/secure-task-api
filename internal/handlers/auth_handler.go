package handlers

type BasicAuthHandler struct {
	store map[string]string //map[username]password
}

func NewBasicAuthHandler() *BasicAuthHandler {
	return &BasicAuthHandler{store: make(map[string]string)}
}

func (auth *BasicAuthHandler) AddCredentials(username string, password string) {
	auth.store[username] = password
}

func (auth *BasicAuthHandler) Allow(name, pass string) bool {
	if len(auth.store) == 0 {
		return false
	}
	storePass, ok := auth.store[name]
	return ok && pass == storePass
}
