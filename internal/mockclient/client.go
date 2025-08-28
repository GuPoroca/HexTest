package mockclient

import "github.com/GuPoroca/HexTest/pkg/jsonOperations"

func MakeAllRequests() {
	projeto := jsonOperations.ReadJSON()

	projeto.ExecuteProject()

	// teste_get := typeDefines.NewTest("Teste_GET", "GET", "", "{aaa}", "200 OK", "/base")
	// teste_get.Execute(server_url, auth)
	//
	// teste_post := typeDefines.NewTest("Teste_POST", "POST", `{"title":"post","body":"request","userId":1}`, "{aaa}", "201 Created", "/base")
	// teste_post.Execute(server_url, auth)
	//
	// teste_put := typeDefines.NewTest("Teste_PUT", "PUT", `{"title":"put","body":"request","userId":2}`, "{aaa}", "202 Accepted", "/base")
	// teste_put.Execute(server_url, auth)
	//
	// teste_delete := typeDefines.NewTest("Teste_DELETE", "DELETE", "", "{aaa}", "200 OK", "/base")
	// teste_delete.Execute(server_url, auth)
	//
	// teste_options := typeDefines.NewTest("Teste_OPTIONS", "OPTIONS", "", "{aaa}", "404 Not Found", "/base")
	// teste_options.Execute(server_url, auth)

}
