package main

import (
	"fmt"
	"os"

	"github.com/GuPoroca/HexTest/internal/mockclient"
	"github.com/GuPoroca/HexTest/internal/mockserver"
	"github.com/GuPoroca/HexTest/pkg/jsonOperations"
	"github.com/GuPoroca/HexTest/server"
)

// go run . server
// or
// fo run . client
func main() {

	switch os.Args[1] {
	case "execute":
		if len(os.Args) <= 2 {
			fmt.Printf("\nadd the project.json file after the execute flag\n")
		} else {
			path := os.Args[2]
			projeto := jsonOperations.ReadJSON(path)
			projeto.ExecuteProject()
		}
	case "server":
		mockserver.OpenServer()
	case "client":
		mockclient.MakeAllRequests()
	case "front":
		//frontend
		server.Run()
	}
}
