package exampleserver

import (
	"log"
	"net/http"
)

// DefaultAddr is the address the demo server listens on when started through
// RunExample.
const DefaultAddr = ":3443"

// Handler builds the demo REST API used to exercise HexTest against a real,
// predictable server. It is exposed (instead of registering on
// http.DefaultServeMux) so it can be reused by the standalone
// cmd/exampleserver binary and driven directly from tests with httptest.
func Handler() http.Handler {
	mux := http.NewServeMux()

	// Core account/auth flow
	mux.HandleFunc("/account/create", createAccount)
	mux.HandleFunc("/account/login", login)
	mux.HandleFunc("/account/data", getData)
	mux.HandleFunc("/admin/users", getAllUsers)

	// Endpoints that stress specific assertion features
	mux.HandleFunc("/test/slow", slowEndpoint)
	mux.HandleFunc("/test/random-error", randomErrorEndpoint)
	mux.HandleFunc("/test/weird-schema", weirdSchemaEndpoint)
	mux.HandleFunc("/test/headers", headerEchoEndpoint)
	mux.HandleFunc("/test/large", largePayloadEndpoint)

	mux.HandleFunc("/test/bodytype", bodyTypeEndpoint)
	mux.HandleFunc("/test/schema", schemaEndpoint)
	mux.HandleFunc("/test/regex", regexEndpoint)
	mux.HandleFunc("/test/empty", emptyEndpoint)

	// OAuth2 client-credentials token endpoint
	mux.HandleFunc("/auth", authHandler)

	return mux
}

// RunExample starts the demo server and blocks until it stops. It returns any
// error from the underlying HTTP server so callers can decide how to handle it.
func RunExample() error {
	log.Printf("Example Server running on %s", DefaultAddr)
	return http.ListenAndServe(DefaultAddr, Handler())
}
