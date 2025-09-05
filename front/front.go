package front

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func Run_Front() {
	component := Hello("Poroca")

	http.Handle("/", templ.Handler(component))
	fmt.Println("listening on :7307")
	http.ListenAndServe(":7307", nil)
}
