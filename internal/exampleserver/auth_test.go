package exampleserver

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	form := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"client_credentials"},
	}
	rec := do(t, http.MethodPost, "/auth", form.Encode(),
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if rec.Code != http.StatusOK {
		t.Fatalf("valid credentials: got %d, want %d", rec.Code, http.StatusOK)
	}
	body := decode(t, rec)
	if token, _ := body["access_token"].(string); !strings.HasPrefix(token, "demo-token-") {
		t.Fatalf("unexpected access_token %v", body["access_token"])
	}
	if body["token_type"] != "bearer" {
		t.Fatalf("token_type: got %v, want bearer", body["token_type"])
	}
}

func TestAuthHandlerRejectsBadCredentials(t *testing.T) {
	form := url.Values{
		"client_id":     {"wrong"},
		"client_secret": {"nope"},
	}
	rec := do(t, http.MethodPost, "/auth", form.Encode(),
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("bad credentials: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
