package exampleserver

import (
	"encoding/json"
	"net/http"
	"time"
)

// Hardcoded credentials the demo OAuth2 token endpoint accepts. They match the
// values in example.env so the CLI's `auth` command can be tried end to end.
const (
	clientID     = "client-id-example"
	clientSecret = "client-secret-example"
)

// authHandler is a minimal OAuth2 client-credentials token endpoint. It accepts
// application/x-www-form-urlencoded requests (as sent by
// golang.org/x/oauth2/clientcredentials) and returns a throwaway bearer token.
func authHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	id := r.FormValue("client_id")
	secret := r.FormValue("client_secret")

	if id != clientID || secret != clientSecret {
		http.Error(w, "invalid client credentials", http.StatusUnauthorized)
		return
	}

	resp := map[string]any{
		"access_token": "demo-token-" + time.Now().Format("150405"),
		"token_type":   "bearer",
		"expires_in":   3600,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
