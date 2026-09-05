// Command exampleserver runs the bundled demo REST API on its own, without the
// rest of the HexTest CLI. It is handy for pointing an external client or the
// web UI at a predictable server. The same server is also reachable through
// `hextest example_server`.
package main

import (
	"log"

	"github.com/GuPoroca/HexTest/internal/exampleserver"
)

func main() {
	log.Fatal(exampleserver.RunExample())
}
