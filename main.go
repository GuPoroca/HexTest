package main

import (
	"fmt"
	"github.com/GuPoroca/HexTest/internal/exampleserver"
	"github.com/GuPoroca/HexTest/pkg/jsonOperations"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	"github.com/GuPoroca/HexTest/server"
	"github.com/GuPoroca/HexTest/tui"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

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

	case "bubble":
		cwd, _ := os.Getwd()
		picker := tui.NewEmojiPicker(cwd, ".json", 0, 0)

		root := tui.RootModel{Current: picker}

		p := tea.NewProgram(root, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			os.Exit(1)
		}
	}
}
