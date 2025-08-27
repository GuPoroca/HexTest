package mockclient

import "github.com/GuPoroca/HexTest/pkg/typedefines"

func MakeAllRequests() {
	server_url := "http://localhost:8080/base"

	teste_get := typedefines.NewTest("Teste_GET", "GET", "", "{aaa}", "200 OK")
	teste_get.Execute(server_url, nil)

	teste_post := typedefines.NewTest("Teste_POST", "POST", `{"title":"post","body":"request","userId":1}`, "{aaa}", "201 Created")
	teste_post.Execute(server_url, nil)

	teste_put := typedefines.NewTest("Teste_PUT", "PUT", `{"title":"put","body":"request","userId":2}`, "{aaa}", "202 Accepted")
	teste_put.Execute(server_url, nil)

	teste_delete := typedefines.NewTest("Teste_DELETE", "DELETE", "", "{aaa}", "200 OK")
	teste_delete.Execute(server_url, nil)

	teste_options := typedefines.NewTest("Teste_OPTIONS", "OPTIONS", "", "{aaa}", "404 Not Found")
	teste_options.Execute(server_url, nil)

}
