package auth

import (
	"net/http"
)

// Mutual TLS authentication strategy
//
// Mutual TLS authentication is a secure method of authentication.
// Unlike a traditional password or token-based method, mutual TLS
// does not exchange a secret value during the authentication process.
// The client and server each present their certificate, which contains
// a public key, during the handshake.
type MutualTLSStrategy struct {
	Certificate string
}

// Do sends requests with an http.Client modified to verify client certificates
func (a *MutualTLSStrategy) Do(*http.Request) (*http.Response, error) {
	panic("Not implemented")
}

var _ Strategy = new(MutualTLSStrategy)
