package main

import (
	"github.com/GuPoroca/HexTest/internal/mockclient"
	"github.com/GuPoroca/HexTest/internal/mockserver"
	"os"
)

func main() {
	switch os.Args[1] {
	case "server":
		go mockserver.OpenServer()
	case "client":
		mockclient.MakeAllRequests()
	}

	select {}
}
