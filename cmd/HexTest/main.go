package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GuPoroca/HexTest/front"
	"github.com/GuPoroca/HexTest/internal/mockclient"
	"github.com/GuPoroca/HexTest/internal/mockserver"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
)

// go run . server
// or
// fo run . client
func main() {
	switch os.Args[1] {
	case "authtest":
		auth := typeDefines.NewoAuth2("client_credentials")
		str, err := auth.Authenticate()
		if err != nil {
			log.Fatalf("error in authentication, %v", err)
		} else {
			fmt.Print(str)
		}
	case "server":
		go mockserver.OpenServer()
	case "client":
		mockclient.MakeAllRequests()
	default:
		front.Run_Front()
	}

	select {}
}
