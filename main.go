package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GuPoroca/HexTest/internal/mockclient"
	"github.com/GuPoroca/HexTest/internal/mockserver"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"github.com/GuPoroca/HexTest/server"
)

// go run . server
// or
// fo run . client
func main() {
	if len(os.Args) > 1 {
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
			server.Run()
		}
	} else {
		server.Run()

	}

}
