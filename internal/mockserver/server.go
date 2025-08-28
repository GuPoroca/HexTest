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
	switch r.Method {
	case "GET":
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "GET called"}`))
		log.Printf("Received a GET request from: %s\n", r.RemoteAddr)
	case "POST":
		time.Sleep(250 * time.Millisecond)
		w.WriteHeader(http.StatusCreated)
		w.Write(([]byte(`{"message": "POST called"}`)))
		io.Copy(res_body, r.Body)
		w.Write([]byte(res_body.String()))
		log.Printf("Received a POST request from: %s\n", r.RemoteAddr)
	case "PUT":
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "PUT called"}`))
		io.Copy(res_body, r.Body)
		w.Write([]byte(res_body.String()))
		log.Printf("Received a PUT request from: %s\n", r.RemoteAddr)
	case "DELETE":
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "DELETE called"}`))
		log.Printf("Received a DELETE request from: %s\n", r.RemoteAddr)
	default:
		time.Sleep(15 * time.Millisecond)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Bro, use real HTTP methods, who tf uses PATCH or OPTIONS, be forreal`))
		log.Printf("Receved a Delusional request from: %s\n", r.RemoteAddr)
	}
}

func OpenServer() {
	http.Handle("/base", new(server))
	log.Printf("Server running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
