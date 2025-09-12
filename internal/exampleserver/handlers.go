package exampleserver

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Data     string `json:"data"`
}

var (
	users    = make(map[string]User)   // simple in-memory db
	sessions = make(map[string]string) // session token -> username
	mu       sync.Mutex
)

// Utils
func jsonResponse(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func getUsernameFromToken(r *http.Request) (string, bool) {
	token := r.Header.Get("Authorization")
	mu.Lock()
	defer mu.Unlock()
	username, ok := sessions[token]
	return username, ok
}

// Core Handlers
func createAccount(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := users[u.Username]; exists {
		jsonResponse(w, http.StatusConflict, map[string]string{"error": "user exists"})
		return
	}
	users[u.Username] = u
	jsonResponse(w, http.StatusCreated, map[string]string{"message": "account created"})
}

func login(w http.ResponseWriter, r *http.Request) {
	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	mu.Lock()
	user, exists := users[creds.Username]
	mu.Unlock()
	if !exists || user.Password != creds.Password {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	token := fmt.Sprintf("token-%s-%d", creds.Username, time.Now().UnixNano())
	mu.Lock()
	sessions[token] = creds.Username
	mu.Unlock()

	jsonResponse(w, http.StatusOK, map[string]string{"token": token})
}

func getData(w http.ResponseWriter, r *http.Request) {
	username, ok := getUsernameFromToken(r)
	if !ok {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "not logged in"})
		return
	}

	user := users[username]
	jsonResponse(w, http.StatusOK, map[string]string{
		"username": user.Username,
		"data":     user.Data,
	})
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	//Hardcoding admin token by hand
	sessions["token-admin-123456"] = "admin"
	//Don't try this at home
	username, ok := getUsernameFromToken(r)
	if !ok || username != "admin" {
		jsonResponse(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	jsonResponse(w, http.StatusOK, users)
}

// Test Handlers
func slowEndpoint(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	jsonResponse(w, http.StatusOK, map[string]string{"message": "slow response done"})
}

func randomErrorEndpoint(w http.ResponseWriter, _ *http.Request) {
	if rand.Intn(2) == 0 {
		jsonResponse(w, http.StatusOK, map[string]string{"message": "all good"})
	} else {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "random failure"})
	}
}

func weirdSchemaEndpoint(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(2) == 0 {
		jsonResponse(w, http.StatusOK, map[string]any{
			"id":   123,
			"name": "weird object",
		})
	} else {
		jsonResponse(w, http.StatusOK, map[string]any{
			"unexpected_field": "oops",
			"nested": map[string]int{
				"deep": 42,
			},
		})
	}
}

func headerEchoEndpoint(w http.ResponseWriter, r *http.Request) {
	headers := map[string]string{}
	for k, v := range r.Header {
		headers[k] = v[0]
		w.Header().Set(k, v[0])
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"received_headers": headers,
	})
}

func largePayloadEndpoint(w http.ResponseWriter, r *http.Request) {
	// simulate large payload
	large := make([]string, 0, 1000)
	for i := range 1000 {
		large = append(large, fmt.Sprintf("item-%d", i))
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"count": len(large),
		"data":  large,
	})
}

func bodyTypeEndpoint(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(2) == 0 {
		// return array
		jsonResponse(w, http.StatusOK, []string{"alpha", "beta"})
	} else {
		// return object
		jsonResponse(w, http.StatusOK, map[string]string{"key": "value"})
	}
}

func schemaEndpoint(w http.ResponseWriter, r *http.Request) {
	// predictable schema
	jsonResponse(w, http.StatusOK, map[string]any{
		"id":    123,
		"name":  "demo",
		"email": "demo@example.com",
	})
}

func regexEndpoint(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]any{
		"email": "user123@example.com",
		"uuid":  "550e8400-e29b-41d4-a716-446655440000",
	})
}

func emptyEndpoint(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(2) == 0 {
		jsonResponse(w, http.StatusOK, []string{}) // empty array
	} else {
		jsonResponse(w, http.StatusOK, map[string]string{}) // empty object
	}
}
