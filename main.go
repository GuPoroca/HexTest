package main

import (
	"fmt"
	"github.com/GuPoroca/HexTest/internal/exampleserver"
	"github.com/GuPoroca/HexTest/pkg/jsonOperations"
	"github.com/GuPoroca/HexTest/server"
	"os"
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
	case "example_server":
		exampleserver.RunExample()
	case "front":
		//frontend
		server.Run()
	}
}
