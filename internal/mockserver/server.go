package mockserver

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res_body := new(strings.Builder)
	switch r.RequestURI {
	case "/base":
		if r.Header.Get("Authorization") != "MTQ0NjJkZmQ5OTM2NDE1ZTZjNGZmZjI3" {
			log.Printf("User not authorized")
			return
		}
		switch r.Method {
		case "GET":
			time.Sleep(500 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"testing": "is true", "message": [{"test": "deep"}], "field 1": "a", "field 2": "b", "field 4": "c"}`))
			log.Printf("Received a GET request from: %s\n", r.RemoteAddr)
		case "POST":
			time.Sleep(250 * time.Millisecond)
			w.WriteHeader(http.StatusCreated)
			w.Write(([]byte(`{"message": "POST called"}`)))
			io.Copy(res_body, r.Body)
			log.Printf("Received a POST request from: %s\n", r.RemoteAddr)
		case "PUT":
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(`{"message": "PUT called"}`))
			io.Copy(res_body, r.Body)
			log.Printf("Received a PUT request from: %s\n", r.RemoteAddr)
		case "DELETE":
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "DELETE called"}`))
			log.Printf("Received a DELETE request from: %s\n", r.RemoteAddr)
		default:
			time.Sleep(15 * time.Millisecond)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Bro, use real HTTP methods, who tf uses PATCH or OPTIONS, be forreal"}`))
			log.Printf("Receved a Delusional request from: %s\n", r.RemoteAddr)
		}
	case "/auth":
		if r.FormValue("client_id") == "abc" && r.FormValue("client_secret") == "123" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
 			"access_token":"MTQ0NjJkZmQ5OTM2NDE1ZTZjNGZmZjI3",
			"token_type":"Bearer",
			"expires_in":3600,
			"refresh_token":"IwOGYzYTlmM2YxOTQ5MGE3YmNmMDFkNTVk",
			"scope":"create"
		}`))
			log.Printf("Received a sucessful auth request from: %s\n", r.RemoteAddr)
		} else {
			log.Printf("Received a failed auth request from: %s\n", r.RemoteAddr)
		}
	}
}

func OpenServer() {
	http.Handle("/base", new(server))
	http.Handle("/auth", new(server))
	log.Printf("Server running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
