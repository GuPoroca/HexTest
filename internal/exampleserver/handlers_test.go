package exampleserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// resetState clears the in-memory "database" so tests do not leak into each
// other regardless of the order the Go test runner picks.
func resetState(t *testing.T) {
	t.Helper()
	mu.Lock()
	defer mu.Unlock()
	users = make(map[string]User)
	sessions = make(map[string]string)
}

// do sends a request through the real mux and returns the recorded response.
func do(t *testing.T, method, target, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	Handler().ServeHTTP(rec, req)
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("response body is not a JSON object: %v (%q)", err, rec.Body.String())
	}
	return out
}

func TestCreateAccount(t *testing.T) {
	resetState(t)

	rec := do(t, http.MethodPost, "/account/create",
		`{"username":"alice","password":"pw","data":"secret"}`, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("first create: got %d, want %d", rec.Code, http.StatusCreated)
	}

	rec = do(t, http.MethodPost, "/account/create",
		`{"username":"alice","password":"pw"}`, nil)
	if rec.Code != http.StatusConflict {
		t.Fatalf("duplicate create: got %d, want %d", rec.Code, http.StatusConflict)
	}

	rec = do(t, http.MethodPost, "/account/create", `not json`, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid json: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestLoginAndGetData(t *testing.T) {
	resetState(t)
	do(t, http.MethodPost, "/account/create",
		`{"username":"bob","password":"hunter2","data":"bob-data"}`, nil)

	rec := do(t, http.MethodPost, "/account/login",
		`{"username":"bob","password":"wrong"}`, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("bad password: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	rec = do(t, http.MethodPost, "/account/login",
		`{"username":"bob","password":"hunter2"}`, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("login: got %d, want %d", rec.Code, http.StatusOK)
	}
	token, _ := decode(t, rec)["token"].(string)
	if token == "" {
		t.Fatal("login did not return a token")
	}

	rec = do(t, http.MethodGet, "/account/data", "", nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("data without token: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	rec = do(t, http.MethodGet, "/account/data", "", map[string]string{"Authorization": token})
	if rec.Code != http.StatusOK {
		t.Fatalf("data with token: got %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decode(t, rec)["data"]; got != "bob-data" {
		t.Fatalf("data payload: got %v, want %q", got, "bob-data")
	}
}

func TestGetAllUsersRequiresAdmin(t *testing.T) {
	resetState(t)

	rec := do(t, http.MethodGet, "/admin/users", "", nil)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("no token: got %d, want %d", rec.Code, http.StatusForbidden)
	}

	rec = do(t, http.MethodGet, "/admin/users", "",
		map[string]string{"Authorization": "token-admin-123456"})
	if rec.Code != http.StatusOK {
		t.Fatalf("admin token: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestSchemaEndpointIsStable(t *testing.T) {
	rec := do(t, http.MethodGet, "/test/schema", "", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusOK)
	}
	body := decode(t, rec)
	for _, key := range []string{"id", "name", "email"} {
		if _, ok := body[key]; !ok {
			t.Errorf("schema response is missing key %q", key)
		}
	}
}

func TestRegexEndpointShape(t *testing.T) {
	rec := do(t, http.MethodGet, "/test/regex", "", nil)
	body := decode(t, rec)
	if email, _ := body["email"].(string); !strings.Contains(email, "@") {
		t.Errorf("email %q does not look like an email", email)
	}
}

func TestHeaderEchoReflectsRequestHeaders(t *testing.T) {
	rec := do(t, http.MethodGet, "/test/headers", "",
		map[string]string{"X-Custom-Test": "hello"})
	if got := rec.Header().Get("X-Custom-Test"); got != "hello" {
		t.Fatalf("echoed header: got %q, want %q", got, "hello")
	}
	received, _ := decode(t, rec)["received_headers"].(map[string]any)
	if received["X-Custom-Test"] != "hello" {
		t.Fatalf("received_headers did not include the custom header: %v", received)
	}
}

func TestLargePayloadCount(t *testing.T) {
	rec := do(t, http.MethodGet, "/test/large", "", nil)
	if got := decode(t, rec)["count"]; got != float64(1000) {
		t.Fatalf("count: got %v, want 1000", got)
	}
}

func TestEmptyEndpointIsAlwaysEmpty(t *testing.T) {
	// The endpoint randomly returns an empty array or an empty object; both
	// must be empty every time.
	for i := 0; i < 20; i++ {
		rec := do(t, http.MethodGet, "/test/empty", "", nil)
		trimmed := strings.TrimSpace(rec.Body.String())
		if trimmed != "[]" && trimmed != "{}" {
			t.Fatalf("iteration %d: got %q, want [] or {}", i, trimmed)
		}
	}
}

func TestUnknownRouteIs404(t *testing.T) {
	rec := do(t, http.MethodGet, "/does/not/exist", "", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusNotFound)
	}
}
