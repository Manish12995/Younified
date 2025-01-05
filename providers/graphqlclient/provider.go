package graphqlclient

import (
	"net/http"
)

// Provider handles authentication for GraphQL requests
type Provider interface {
	Authenticate(req *http.Request)
}

// BasicAuthProvider implements basic authentication
type BasicAuthProvider struct {
	Username string
	Password string
}

// Authenticate adds basic authentication to the request
func (b *BasicAuthProvider) Authenticate(req *http.Request) {
	req.SetBasicAuth(b.Username, b.Password)
}

// TokenAuthProvider implements token-based authentication
type TokenAuthProvider struct {
	Token string
}

// Authenticate adds token authentication to the request
func (t *TokenAuthProvider) Authenticate(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+t.Token)
}
