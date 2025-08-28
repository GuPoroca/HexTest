package main

import (
	"os"

	"github.com/GuPoroca/HexTest/internal/mockclient"
	"github.com/GuPoroca/HexTest/internal/mockserver"
)

// go run . server
// or
// fo run . client
func main() {
	switch os.Args[1] {
	case "server":
		go mockserver.OpenServer()
	case "client":
		mockclient.MakeAllRequests()
	}

	select {}
}
