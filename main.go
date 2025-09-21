package main

import (
	"fmt"
	"os"

	"github.com/GuPoroca/HexTest/internal/exampleserver"
	"github.com/GuPoroca/HexTest/pkg/jsonOperations"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"github.com/GuPoroca/HexTest/server"
)

// go run . server
// or
// fo run . client
func main() {

	switch os.Args[1] {
	case "auth":
		auth := typeDefines.NewoAuth2("grant_type")
		token := auth.Authenticate()
		fmt.Printf("\n\n%s\n\n", token)
	case "execute":
		if len(os.Args) <= 2 {
			fmt.Printf("\nadd the project.json file after the execute flag\n")
		} else {
			path := os.Args[2]
			projeto := jsonOperations.ReadJSON(path)
			projeto.ExecuteProject()
			projeto.PrintResults()
		}
	case "example_server":
		exampleserver.RunExample()
	case "front":
		//frontend
		server.Run()
	}

}
