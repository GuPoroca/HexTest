package exampleserver

import (
	"log"
	"net/http"
)

func RunExample() {
	http.HandleFunc("/account/create", createAccount)
	http.HandleFunc("/account/login", login)
	http.HandleFunc("/account/data", getData)
	http.HandleFunc("/admin/users", getAllUsers)

	http.HandleFunc("/test/slow", slowEndpoint)
	http.HandleFunc("/test/random-error", randomErrorEndpoint)
	http.HandleFunc("/test/weird-schema", weirdSchemaEndpoint)
	http.HandleFunc("/test/headers", headerEchoEndpoint)
	http.HandleFunc("/test/large", largePayloadEndpoint)

	http.HandleFunc("/test/bodytype", bodyTypeEndpoint)
	http.HandleFunc("/test/schema", schemaEndpoint)
	http.HandleFunc("/test/regex", regexEndpoint)
	http.HandleFunc("/test/empty", emptyEndpoint)

	log.Println("Example Server running on :3443")
	log.Fatal(http.ListenAndServe(":3443", nil))
}
