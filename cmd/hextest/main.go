package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GuPoroca/HexTest/internal/exampleserver"
	"github.com/GuPoroca/HexTest/pkg/jsonOperations"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"github.com/GuPoroca/HexTest/server"
)

const usage = `HexTest - structured REST API testing

Usage:
  hextest <command> [arguments]

Commands:
  execute <project.json>   Run every suite/test in a project file and print results
  front                    Start the local web UI (templ + HTMX) on :3773
  example_server           Start the bundled demo REST API on :3443
  auth                     Fetch an OAuth2 token using the values in .env
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "auth":
		auth := typeDefines.NewoAuth2("grant_type")
		token := auth.Authenticate()
		fmt.Printf("\n\n%s\n\n", token)

	case "execute":
		if len(os.Args) <= 2 {
			fmt.Println("add the project.json file after the execute flag")
			os.Exit(1)
		}
		projeto := jsonOperations.ReadJSON(os.Args[2])
		projeto.ExecuteProject()
		projeto.PrintResults()

	case "example_server":
		if err := exampleserver.RunExample(); err != nil {
			log.Fatal(err)
		}

	case "front":
		server.Run()

	default:
		fmt.Printf("unknown command %q\n\n", os.Args[1])
		fmt.Print(usage)
		os.Exit(1)
	}
}
